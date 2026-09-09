package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Author struct {
	ID uint `db:"id"`

	Name    string `db:"name"`
	URL     string `db:"url"`
	Email   string `db:"email"`
	Discord string `db:"discord"`
}

func (a *Author) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO author (%s, name, url, email, discord)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, data.fieldName)
}

func (a *Author) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.Name, a.URL, a.Email, a.Discord}

	err := tx.QueryRowx(data.query, args...).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *Author) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.Name, a.URL, a.Email, a.Discord}

	mutex := GetMutex("author_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

type Authors []*Author

func (a *Authors) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, name, url, email, discord
		FROM author
		WHERE %s = $1`, data.fieldName)
}

func (a *Authors) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, a, data.query, data.id)
}
