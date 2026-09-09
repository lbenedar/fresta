package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Light struct {
	ID uint

	Color          string
	Priority       int
	Angle          int
	Negative       bool
	Alpha          float64
	Bright         float64
	Coloration     float64
	Dim            float64
	Attenuation    float64
	Luminosity     float64
	Saturation     float64
	Contrast       float64
	Shadows        float64
	LightAnimation LightAnimation
	LightDarkness  LightDarkness
}

func (l *Light) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_light (%s, color, priority, angle, negative, alpha, bright, coloration, dim, 
			attenuation, luminosity, saturation, contrast, shadows)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`, data.fieldName)
}

func (l *Light) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[uint]{id: l.ID, fieldName: "token_light_id"}
	InsertWithCtxParallel(group, tx, l.LightAnimation, relId)
	InsertWithCtxParallel(group, tx, l.LightDarkness, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (l *Light) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Color, l.Priority, l.Angle, l.Negative, l.Alpha, l.Bright, l.Coloration, l.Dim,
		l.Attenuation, l.Luminosity, l.Saturation, l.Contrast, l.Shadows}

	err := tx.QueryRowx(data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l *Light) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Color, l.Priority, l.Angle, l.Negative, l.Alpha, l.Bright, l.Coloration, l.Dim,
		l.Attenuation, l.Luminosity, l.Saturation, l.Contrast, l.Shadows}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

type LightAnimation struct {
	ID uint

	Speed     int
	Intensity int
	Reverse   bool
}

func (l LightAnimation) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_light_animation (%s, speed, intensity, reverse)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (l LightAnimation) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Speed, l.Intensity, l.Reverse}

	err := tx.QueryRowx(data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l LightAnimation) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Speed, l.Intensity, l.Reverse}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

type LightDarkness struct {
	ID uint

	Min float64
	Max float64
}

func (l LightDarkness) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_light_darkness (%s, min, max)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (l LightDarkness) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Min, l.Max}

	err := tx.QueryRowx(data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l LightDarkness) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Min, l.Max}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}
