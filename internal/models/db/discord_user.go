package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DiscordUser struct {
	ID string

	Email         string
	Username      string
	Avatar        string
	Locale        string
	Discriminator string
	GlobalName    string
}

func (d *DiscordUser) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO discord_user (id, email, username, avatar, locale, discriminator, global_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(id) DO NOTHING`
}

func (d *DiscordUser) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{d.ID, d.Email, d.Username, d.Avatar, d.Locale, d.Discriminator, d.GlobalName}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordUser) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{d.ID, d.Email, d.Username, d.Avatar, d.Locale, d.Discriminator, d.GlobalName}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}
