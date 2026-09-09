package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Release struct {
	ID uint `db:"id"`

	Generation          int       `db:"generation"`
	Build               int       `db:"build"`
	NodeVersion         int       `db:"node_version"`
	MaxGeneration       int       `db:"max_generation"`
	MaxStableGeneration int       `db:"max_stable_generation"`
	Time                time.Time `db:"time"`
	Channel             string    `db:"channel"`
	Suffix              string    `db:"suffix"`
}

func (r *Release) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO release_ (%s, generation, build, node_version, max_generation, max_stable_generation,
			time, channel, suffix)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`, data.fieldName)
}

func (r *Release) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Generation, r.Build, r.NodeVersion, r.MaxGeneration, r.MaxStableGeneration,
		r.Time, r.Channel, r.Suffix}

	err := tx.QueryRowx(data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Release) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Generation, r.Build, r.NodeVersion, r.MaxGeneration, r.MaxStableGeneration,
		r.Time, r.Channel, r.Suffix}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Release) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, generation, build, node_version, max_generation, max_stable_generation,
			time, channel, suffix 
		FROM release_
		WHERE %s = $1`, data.fieldName)
}

func (r *Release) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, r, data.query, data.id)
}
