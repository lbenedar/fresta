package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Stats struct {
	ID uint

	CoreVersion    string
	SystemID       string
	SystemVersion  string
	LastModifiedBy string
	ModifiedTime   int64
}

func (s Stats) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO stats (%s, core_version, system_id, system_version, last_modified_by, modified_time)
		VALUES ($1, $2, $3, $4, $5, $6)`, data.fieldName)
}

func (s Stats) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.CoreVersion, s.SystemID, s.SystemVersion, s.LastModifiedBy, s.ModifiedTime}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Stats) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.CoreVersion, s.SystemID, s.SystemVersion, s.LastModifiedBy, s.ModifiedTime}

	mutex := GetMutex("stats_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}
