package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type User struct {
	ID string `db:"id"`

	Name      string `db:"name"`
	Avatar    string `db:"avatar"`
	Character string `db:"character"`
	Color     string `db:"color"`
	Pronouns  string `db:"pronouns"`
	Role      int    `db:"role"`
	Stats     Stats
	Hotbar    UserHotbars
}

func (u *User) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO user (game_id, id, name, avatar, character, color, pronouns, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (u *User) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: u.ID, fieldName: "user_id"}
	InsertWithCtxParallel(group, tx, u.Stats, relId)
	InsertSliceParallel(group, tx, u.Hotbar, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (u *User) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, u.ID, u.Name, u.Avatar, u.Character, u.Color, u.Pronouns, u.Role}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return u.InsertObjects(tx)
	}
	return nil
}

func (u *User) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, u.ID, u.Name, u.Avatar, u.Character, u.Color, u.Pronouns, u.Role}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return u.InsertObjects(tx)
	}
	return nil
}

type Users []*User

func (u *Users) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, name, avatar, character, color, pronouns, role
		FROM user
		WHERE %s = $1`, data.fieldName)
}

func (u *Users) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for _, v := range *u {
		group.Go(func() error {
			relId := InsertId[string]{id: v.ID, fieldName: "user_id"}
			funcGroup, _ := errgroup.WithContext(context.Background())

			GetWithCtxParallel(group, db, &v.Hotbar, relId)

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (u *Users) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, u, data.query, data.id)
	if err != nil {
		return err
	}

	return u.GetObjects(db)
}

type UserHotbar struct {
	ID uint `db:"id"`

	Key   int    `db:"key_"`
	Value string `db:"value"`
}

func (u UserHotbar) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO hotbar (user_id, key_, value)
		VALUES ($1, $2, $3)`
}

func (u UserHotbar) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, u.Key, u.Value}

	err := tx.QueryRowx(data.query, args...).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u UserHotbar) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, u.Key, u.Value}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

type UserHotbars []UserHotbar

func (u *UserHotbars) GetQuery(data *InsertId[string]) {
	data.query = `
		SELECT key_, value
		FROM hotbar
		WHERE user_id = $1`
}

func (u *UserHotbars) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, u, data.query, data.id)
	if err != nil {
		return err
	}

	return nil
}
