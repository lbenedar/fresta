package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CoreUpdate struct {
	ID uint `db:"id"`

	HasUpdate          bool   `db:"has_update"`
	CanUpdate          bool   `db:"can_update"`
	CouldReachWebsite  bool   `db:"could_reach_website"`
	SlowResponse       bool   `db:"slow_response"`
	WillDisableModules bool   `db:"will_disable_modules"`
	Version            string `db:"version"`
	Channel            string `db:"channel"`
}

func (c *CoreUpdate) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO core_update (%s, has_update, can_update, could_reach_website, slow_response, will_disable_modules, version, channel)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`, data.fieldName)
}

func (c *CoreUpdate) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.HasUpdate, c.CanUpdate, c.CouldReachWebsite, c.SlowResponse, c.WillDisableModules, c.Version, c.Channel}

	err := tx.QueryRowx(data.query, args...).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoreUpdate) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.HasUpdate, c.CanUpdate, c.CouldReachWebsite, c.SlowResponse, c.WillDisableModules, c.Version, c.Channel}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoreUpdate) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT has_update, can_update, could_reach_website, slow_response, will_disable_modules, 
			version, channel 
		FROM core_update
		WHERE %s = $1`, data.fieldName)
}

func (c *CoreUpdate) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, c, data.query, data.id)
}

type SystemUpdate struct {
	ID uint

	HasUpdate bool
	Version   string
}

func (s *SystemUpdate) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO system_update (%s, has_update, version)
		VALUES ($1, $2, $3)
		RETURNING id`, data.fieldName)
}

func (s *SystemUpdate) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.HasUpdate, s.Version}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemUpdate) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.HasUpdate, s.Version}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}
