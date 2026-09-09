package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Addresses struct {
	ID uint

	Local              string
	Remote             string
	RemoteIsAccessible bool
}

func (a *Addresses) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO addresses (game_id, local, remote, remote_is_accessible)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
}

func (a *Addresses) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.Local, a.Remote, a.RemoteIsAccessible}

	err := tx.QueryRowx(data.query, args...).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *Addresses) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, a.Local, a.Remote, a.RemoteIsAccessible}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&a.ID)
	if err != nil {
		return err
	}

	return nil
}
