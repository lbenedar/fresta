package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Item struct {
	ID string

	Img       string
	Name      string
	Type      string
	Folder    string
	Sort      int
	Stats     Stats
	Ownership []OwnershipString
}

func (i *Item) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO item (%[1]s, id, img, name, type, folder, sort)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(id) DO UPDATE SET
			%[1]s = EXCLUDED.%[1]s,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`, data.fieldName)
}

func (i *Item) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: i.ID, fieldName: "item_id"}
	InsertWithCtxParallel(group, tx, i.Stats, relId)
	InsertSliceParallel(group, tx, i.Ownership, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (i *Item) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, i.ID, i.Img, i.Name, i.Type, i.Folder, i.Sort}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return i.InsertObjects(tx)
	}
	return nil
}

func (i *Item) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, i.ID, i.Img, i.Name, i.Type, i.Folder, i.Sort}

	mutex := GetMutex("item_insert")
	mutex.Lock()
	defer mutex.Unlock()

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return i.InsertObjects(tx)
	}
	return nil
}
