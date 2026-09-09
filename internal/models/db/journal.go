package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Journal struct {
	ID string

	Name      string
	Sort      int
	Pages     []*JournalPage
	Ownership []OwnershipString
}

func (j *Journal) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO journal (game_id, id, name, sort)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (j *Journal) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: j.ID, fieldName: "journal_id"}
	InsertSliceParallel(group, tx, j.Pages, relId)
	InsertSliceParallel(group, tx, j.Ownership, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (j *Journal) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, j.ID, j.Name, j.Sort}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return j.InsertObjects(tx)
	}
	return nil
}

func (j *Journal) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, j.ID, j.Name, j.Sort}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return j.InsertObjects(tx)
	}
	return nil
}
