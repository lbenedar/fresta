package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Macro struct {
	ID string

	Command   string
	Name      string
	Type      string
	Img       string
	Author    string
	Scope     string
	Folder    string
	Sort      int
	Stats     Stats
	Ownership []OwnershipString
}

func (m *Macro) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO macro (game_id, id, command, name, type, img, author, scope, folder, sort)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (m *Macro) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: m.ID, fieldName: "macro_id"}
	InsertWithCtxParallel(group, tx, m.Stats, relId)
	InsertSliceParallel(group, tx, m.Ownership, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (m *Macro) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID, m.Command, m.Name, m.Type, m.Img, m.Author, m.Scope, m.Folder, m.Sort}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return m.InsertObjects(tx)
	}
	return nil
}

func (m *Macro) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID, m.Command, m.Name, m.Type, m.Img, m.Author, m.Scope, m.Folder, m.Sort}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return m.InsertObjects(tx)
	}
	return nil
}
