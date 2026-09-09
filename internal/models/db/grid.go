package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Grid struct {
	ID uint `db:"id"`

	Type      int     `db:"type"`
	Size      int     `db:"size"`
	Distance  int     `db:"distance"`
	Diagonals int     `db:"diagonals"`
	Thickness int     `db:"thickness"`
	Alpha     float64 `db:"alpha"`
	Color     string  `db:"color"`
	Units     string  `db:"units"`
	Style     string  `db:"style"`
}

func (g *Grid) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO grid (system_id, type, size, distance, diagonals, thickness,
							alpha, color, units, style)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
}

func (g *Grid) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, g.Type, g.Size, g.Distance, g.Diagonals, g.Thickness,
		g.Alpha, g.Color, g.Units, g.Style}

	err := tx.QueryRowx(data.query, args...).Scan(&g.ID)
	if err != nil {
		return err
	}

	return nil
}

func (g *Grid) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, g.Type, g.Size, g.Distance, g.Diagonals, g.Thickness,
		g.Alpha, g.Color, g.Units, g.Style}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&g.ID)
	if err != nil {
		return err
	}

	return nil
}

func (g *Grid) GetQuery(data *InsertId[string]) {
	data.query = `
		SELECT id, type, size, distance, diagonals, thickness,
			alpha, color, units, style
		FROM grid
		WHERE system_id = $1`

}

func (g *Grid) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, g, data.query, data.id)
}
