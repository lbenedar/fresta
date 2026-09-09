package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Setting struct {
	ID string

	Key   string
	Value string
	Stats Stats
}

func (s *Setting) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO setting (game_id, id, key_, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (s *Setting) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: s.ID, fieldName: "setting_id"}
	InsertWithCtxParallel(group, tx, s.Stats, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (s *Setting) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID, s.Key, s.Value}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return s.InsertObjects(tx)
	}
	return nil
}

func (s *Setting) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID, s.Key, s.Value}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return s.InsertObjects(tx)
	}
	return nil
}
