package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Media struct {
	ID uint `db:"id"`

	Type    string `db:"type"`
	URL     string `db:"url"`
	Caption string `db:"caption"`
}

func (m *Media) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO media (%s, type, url, caption)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (m *Media) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.Type, m.URL, m.Caption}

	err := tx.QueryRowx(data.query, args...).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Media) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.Type, m.URL, m.Caption}

	mutex := GetMutex("media_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

type MediaSlice []*Media

func (m *MediaSlice) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, type, url, caption
		FROM media
		WHERE %s = $1`, data.fieldName)
}

func (m *MediaSlice) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, m, data.query, data.id)
}
