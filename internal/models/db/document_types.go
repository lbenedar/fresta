package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type DocumentTypes struct {
	ID uint `db:"id"`

	Data DocumentsTypeData
}

func (d DocumentTypes) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO document_types (%s)
		VALUES ($1)
		RETURNING id`, data.fieldName)
}

func (d DocumentTypes) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	err := tx.QueryRowx(data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertSliceParallel(group, tx, d.Data, InsertId[uint]{id: d.ID})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (d DocumentTypes) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	mutex := GetMutex("document_types_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertSliceParallel(group, tx, d.Data, InsertId[uint]{id: d.ID})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (d *DocumentTypes) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id 
		FROM document_types
		WHERE %s = $1`, data.fieldName)
}

func (d *DocumentTypes) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.GetContext(ctx, d, data.query, data.id)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	GetWithCtxParallel(group, db, &d.Data, InsertId[uint]{id: d.ID})

	return group.Wait()
}

type DocumentTypeData struct {
	ID uint `db:"id"`

	Type       string `db:"type"`
	HtmlFields []string
}

func (d *DocumentTypeData) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO document_types_data (document_types_id, type)
		VALUES ($1, $2)
		RETURNING id`
}

func (d *DocumentTypeData) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, d.Type}

	err := tx.QueryRowx(data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	relData := &InsertId[uint]{id: d.ID, fieldName: "document_types_data_id", tableName: "document_types_data_html"}
	InsertSimpleSliceParallel(group, tx, d.HtmlFields, relData)

	return group.Wait()
}

func (d *DocumentTypeData) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, d.Type}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	relData := &InsertId[uint]{id: d.ID, fieldName: "document_types_data_id", tableName: "document_types_data_html"}
	InsertSimpleSliceParallel(group, tx, d.HtmlFields, relData)

	return group.Wait()
}

type DocumentsTypeData []*DocumentTypeData

func (d *DocumentsTypeData) GetQuery(data *InsertId[uint]) {
	data.query = `
		SELECT id, type 
		FROM document_types_data
		WHERE document_types_id = $1`
}

func (d *DocumentsTypeData) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *d {
		GetSimpleSliceParallel(group, db, &(*d)[i].HtmlFields,
			&InsertId[uint]{id: (*d)[i].ID, fieldName: "document_types_data_id", tableName: "document_types_data_html"})
	}

	return group.Wait()
}

func (d *DocumentsTypeData) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, d, data.query, data.id)
	if err != nil {
		return err
	}

	return d.GetObjects(db)
}
