package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Compatibility struct {
	ID uint `db:"id"`

	Minimum  string `db:"minimum"`
	Verified string `db:"verified"`
	Maximum  string `db:"maximum"`
}

func (c Compatibility) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO compatibility (%s, minimum, verified, maximum)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (c Compatibility) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.Minimum, c.Verified, c.Maximum}

	err := tx.QueryRowx(data.query, args...).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c Compatibility) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.Minimum, c.Verified, c.Maximum}

	mutex := GetMutex("compatibility_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Compatibility) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, minimum, verified, maximum
		FROM compatibility
		WHERE %s = $1`, data.fieldName)
}

func (c *Compatibility) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, c, data.query, data.id)
}
