package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Files struct {
	ID uint `db:"id"`

	Storages FilesStorage
}

func (f *Files) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO files (%s)
		VALUES ($1)
		RETURNING id`, data.fieldName)
}

func (f *Files) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	err := tx.QueryRowx(data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	InsertSliceParallel(group, tx, f.Storages, InsertId[uint]{id: f.ID})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (f *Files) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	InsertSliceParallel(group, tx, f.Storages, InsertId[uint]{id: f.ID})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (f *Files) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id 
		FROM files
		WHERE %s = $1`, data.fieldName)
}

func (f *Files) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.GetContext(ctx, f, data.query, data.id)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	GetWithCtxParallel(group, db, &f.Storages, InsertId[uint]{id: f.ID, fieldName: "files_id"})

	return group.Wait()
}

type FileStorage struct {
	ID uint `db:"id"`

	Storage string `db:"storage"`
}

func (f FileStorage) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO files_storage (files_id, storage)
		VALUES ($1, $2)
		RETURNING id`
}

func (f FileStorage) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Storage}

	err := tx.QueryRowx(data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f FileStorage) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Storage}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

type FilesStorage []FileStorage

func (f *FilesStorage) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, storage 
		FROM files_storage
		WHERE %s = $1`, data.fieldName)
}

func (f *FilesStorage) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, f, data.query, data.id)
}
