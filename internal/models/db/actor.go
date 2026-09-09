package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Actor struct {
	ID string

	Img            string
	Name           string
	Type           string
	Folder         string
	Sort           int
	PrototypeToken *Token
	Stats          Stats
	Items          []*Item
	Ownership      []OwnershipString
}

func (a *Actor) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO actor (game_id, id, img, name, type, folder, sort)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (a *Actor) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: a.ID, fieldName: "actor_id"}
	InsertWithCtxParallel(group, tx, a.PrototypeToken, relId)
	InsertWithCtxParallel(group, tx, a.Stats, relId)
	InsertSliceParallel(group, tx, a.Ownership, relId)
	InsertSliceParallel(group, tx, a.Items, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (a *Actor) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.ID, a.Img, a.Name, a.Type, a.Folder, a.Sort}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return a.InsertObjects(tx)
	}
	return nil
}

func (a *Actor) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.ID, a.Img, a.Name, a.Type, a.Folder, a.Sort}

	mutex := GetMutex("actor_insert")
	mutex.Lock()
	defer mutex.Unlock()

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return a.InsertObjects(tx)
	}
	return nil
}
