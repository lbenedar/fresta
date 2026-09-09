package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type GameOptions struct {
	ID uint

	Language      string
	UpdateChannel string
	Port          int
}

func (g *GameOptions) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO game_options (game_id, language, update_channel, port)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
}

func (g *GameOptions) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, g.Language, g.UpdateChannel, g.Port}

	err := tx.QueryRowx(data.query, args...).Scan(&g.ID)
	if err != nil {
		return err
	}

	return nil
}

func (g *GameOptions) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, g.Language, g.UpdateChannel, g.Port}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&g.ID)
	if err != nil {
		return err
	}

	return nil
}

type SetupOptions struct {
	ID uint `db:"id"`

	CSSTheme       string `db:"css_theme"`
	DataPath       string `db:"data_path"`
	Hostname       string `db:"hostname"`
	Language       string `db:"language"`
	LocalHostname  string `db:"local_hostname"`
	UpdateChannel  string `db:"update_channel"`
	Port           int    `db:"port"`
	CompressSocket bool   `db:"compress_socket"`
	CompressStatic bool   `db:"compress_static"`
	Fullscreen     bool   `db:"fullscreen"`
	HotReload      bool   `db:"hot_reload"`
	ProxySSL       bool   `db:"proxy_ssl"`
	Telemetry      bool   `db:"telemetry"`
	Upnp           bool   `db:"upnp"`
	DeleteNEDB     bool   `db:"delete_nedb"`
	NoBackups      bool   `db:"no_backups"`
}

func (s *SetupOptions) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO setup_options (setup_id, css_theme, data_path, hostname, language, local_hostname, 
			update_channel, port, compress_socket, compress_static, fullscreen, hot_reload, proxy_ssl,
			telemetry, upnp, delete_nedb, no_backups)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id`
}

func (s *SetupOptions) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.CSSTheme, s.DataPath, s.Hostname, s.Language, s.LocalHostname,
		s.UpdateChannel, s.Port, s.CompressSocket, s.CompressStatic, s.Fullscreen, s.HotReload,
		s.ProxySSL, s.Telemetry, s.Upnp, s.DeleteNEDB, s.NoBackups}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SetupOptions) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.CSSTheme, s.DataPath, s.Hostname, s.Language, s.LocalHostname,
		s.UpdateChannel, s.Port, s.CompressSocket, s.CompressStatic, s.Fullscreen, s.HotReload,
		s.ProxySSL, s.Telemetry, s.Upnp, s.DeleteNEDB, s.NoBackups}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SetupOptions) GetQuery(data *InsertId[uint]) {
	data.query = `
		SELECT id, css_theme, data_path, hostname, language, local_hostname, update_channel, port, compress_socket,
			compress_static, fullscreen, hot_reload, proxy_ssl, telemetry, upnp, delete_nedb, no_backups
		FROM setup_options
		WHERE setup_id = $1`
}

func (s *SetupOptions) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, s, data.query, data.id)
}
