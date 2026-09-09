package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Style struct {
	ID uint `db:"id"`

	Src string `db:"src"`
}

func (s *Style) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO style (%s, src)
		VALUES ($1, $2)
		RETURNING id`, data.fieldName)
}

func (s *Style) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Src}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Style) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Src}

	mutex := GetMutex("style_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

type StyleSlice []*Style

func (s *StyleSlice) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, src
		FROM style
		WHERE %s = $1`, data.fieldName)
}

func (s *StyleSlice) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, s, data.query, data.id)
}
