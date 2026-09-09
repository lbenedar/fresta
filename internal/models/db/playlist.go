package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Playlist struct {
	ID string

	Name        string
	Folder      string
	Sorting     string
	Description string
	Channel     string
	Mode        int
	Fade        int
	Seed        int
	Sort        int
	Playing     bool
	Stats       Stats
	Ownership   []OwnershipString
	Sounds      []*Sound
}

func (p *Playlist) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO playlist (game_id, id, name, folder, sorting, description, channel, mode, fade, seed, sort, playing)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (p *Playlist) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: p.ID, fieldName: "playlist_id"}
	InsertWithCtxParallel(group, tx, p.Stats, relId)
	InsertSliceParallel(group, tx, p.Ownership, relId)
	InsertSliceParallel(group, tx, p.Sounds, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (p *Playlist) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID, p.Name, p.Folder, p.Sorting, p.Description, p.Channel, p.Mode, p.Fade, p.Seed, p.Sort, p.Playing}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return p.InsertObjects(tx)
	}
	return nil
}

func (p *Playlist) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID, p.Name, p.Folder, p.Sorting, p.Description, p.Channel, p.Mode, p.Fade, p.Seed, p.Sort, p.Playing}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return p.InsertObjects(tx)
	}
	return nil
}

type Sound struct {
	ID string

	Name        string
	Path        string
	Channel     string
	Description string
	Fade        int
	Sort        int
	Repeat      bool
	Playing     bool
	Volume      float64
	PausedTime  float64
}

func (s *Sound) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO sound (playlist_id, id, name, path, channel, description, fade, sort, repeat, playing, volume, paused_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT(id) DO UPDATE SET
			playlist_id = EXCLUDED.playlist_id,
			updated_at = datetime('now')`
}

func (s *Sound) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID, s.Name, s.Path, s.Channel, s.Description, s.Fade, s.Sort, s.Repeat, s.Playing, s.Volume, s.PausedTime}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sound) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID, s.Name, s.Path, s.Channel, s.Description, s.Fade, s.Sort, s.Repeat, s.Playing, s.Volume, s.PausedTime}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}
