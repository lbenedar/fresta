package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Table struct {
	ID string

	Name        string
	Description string
	Formula     string
	Img         string
	Folder      string
	Sort        int
	Replacement bool
	DisplayRoll bool
	Stats       Stats
	Ownership   []OwnershipString
	Results     []*TableResult
}

func (t *Table) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO table_ (game_id, id, name, description, formula, img, folder, sort, replacement, display_roll)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (t *Table) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: t.ID, fieldName: "table_id"}
	InsertWithCtxParallel(group, tx, t.Stats, relId)
	InsertSliceParallel(group, tx, t.Ownership, relId)
	InsertSliceParallel(group, tx, t.Results, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (t *Table) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.ID, t.Name, t.Description, t.Formula, t.Img, t.Folder, t.Sort, t.Replacement, t.DisplayRoll}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return t.InsertObjects(tx)
	}
	return nil
}

func (t *Table) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.ID, t.Name, t.Description, t.Formula, t.Img, t.Folder, t.Sort, t.Replacement, t.DisplayRoll}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return t.InsertObjects(tx)
	}
	return nil
}

type TableResult struct {
	ID string

	Type        string
	Img         string
	Description string
	Name        string
	Weight      int
	Drawn       bool
	Stats       Stats
	Range       []int
}

func (t *TableResult) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO table_result (table_id, id, type, img, description, name, weight, drawn)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT(id) DO UPDATE SET
			table_id = EXCLUDED.table_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (t *TableResult) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, t.Stats, InsertId[string]{id: t.ID, fieldName: "table_result_id"})
	InsertSimpleSlice(tx, t.Range, &InsertId[string]{id: t.ID, fieldName: "table_result_id", tableName: "table_result_range"})

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (t *TableResult) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.ID, t.Type, t.Img, t.Description, t.Name, t.Weight, t.Drawn}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return t.InsertObjects(tx)
	}
	return nil
}

func (t *TableResult) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.ID, t.Type, t.Img, t.Description, t.Name, t.Weight, t.Drawn}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return t.InsertObjects(tx)
	}
	return nil
}
