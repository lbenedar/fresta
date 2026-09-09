package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Index struct {
	ID string `db:"id"`

	Folder string `db:"folder"`
	Img    string `db:"img"`
	Name   string `db:"name"`
	Type   string `db:"type"`
}

func (i *Index) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO index_ (id, folder, img, name, type)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(id) DO NOTHING`
}

func (i *Index) ConnectGameQuery(data *InsertId[string]) {
	data.query = `
		INSERT INTO pack_to_index (pack_id, index_id)
		VALUES ($1, $2)
		ON CONFLICT(pack_id, index_id) DO NOTHING`
}

func (i *Index) ConnectGame(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, i.ID}

	_, err := tx.Exec(data.query, args...)
	return err
}

func (i *Index) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{i.ID, i.Folder, i.Img, i.Name, i.Type}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	dataCopy := *data
	i.ConnectGameQuery(&dataCopy)

	return i.ConnectGame(tx, &dataCopy)
}

func (i *Index) ConnectGameCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, i.ID}

	_, err := tx.ExecContext(ctx, data.query, args...)
	return err
}

func (i *Index) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{i.ID, i.Folder, i.Img, i.Name, i.Type}

	mutex := GetMutex("index_insert")
	mutex.Lock()
	defer mutex.Unlock()

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	dataCopy := *data
	i.ConnectGameQuery(&dataCopy)

	return i.ConnectGameCtx(ctx, tx, &dataCopy)
}

type IndexSlice []*Index

func (i *IndexSlice) GetQuery(data *InsertId[string]) {
	data.query = `
		SELECT i.id, i.folder, i.img, i.name, i.type
		FROM index_ i
		JOIN pack_to_index AS pi ON pi.index_id = i.id
		WHERE pi.pack_id = $1`
}

func (i *IndexSlice) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, i, data.query, data.id)
}
