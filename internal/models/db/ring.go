package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Ring struct {
	ID uint

	Enabled    bool
	Effects    int
	RingColors RingColors
	Subject    Subject
}

func (r *Ring) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO ring (%s, enabled, effects)
		VALUES ($1, $2, $3)
		RETURNING id`, data.fieldName)
}

func (r *Ring) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[uint]{id: r.ID, fieldName: "ring_id"}
	InsertWithCtxParallel(group, tx, r.RingColors, relId)
	InsertWithCtxParallel(group, tx, r.Subject, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (r *Ring) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Enabled, r.Effects}

	err := tx.QueryRowx(data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Ring) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Enabled, r.Effects}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

type RingColors struct {
	ID uint

	Ring       string
	Background string
}

func (r RingColors) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO ring_colors (%s, ring, background)
		VALUES ($1, $2, $3)
		RETURNING id`, data.fieldName)
}

func (r RingColors) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Ring, r.Background}

	err := tx.QueryRowx(data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r RingColors) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Ring, r.Background}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	return nil
}

type Subject struct {
	ID uint

	Scale   int
	Texture string
}

func (s Subject) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO ring_colors (%s, scale, texture)
		VALUES ($1, $2, $3)
		RETURNING id`, data.fieldName)
}

func (s Subject) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Scale, s.Texture}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Subject) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Scale, s.Texture}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}
