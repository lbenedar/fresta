package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type DiscordMember struct {
	ID uint

	GuildId  string
	JoinedAt time.Time
	Nick     string
	Deaf     bool
	Mute     bool

	User *DiscordUser
}

func (d *DiscordMember) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO discord_member (guild_id, joined_at, nick, deaf, mute, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(user_id) DO UPDATE SET
			updated_at = datetime('now')
		RETURNING id`
}

func (d *DiscordMember) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	if d.User != nil {
		relId := InsertId[string]{id: strconv.FormatUint(uint64(d.ID), 10), fieldName: "member_id"}
		InsertWithCtxParallel(group, tx, d.User, relId)
	}

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (d *DiscordMember) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := d.InsertObjects(tx)
	if err != nil {
		return err
	}

	var userId *string = nil
	if d.User != nil {
		userId = &d.User.ID
	}

	args := []any{d.GuildId, d.JoinedAt, d.Nick, d.Deaf, d.Mute, userId}

	err = tx.QueryRow(data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordMember) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := d.InsertObjects(tx)
	if err != nil {
		return err
	}

	var userId *string = nil
	if d.User != nil {
		userId = &d.User.ID
	}

	args := []any{d.GuildId, d.JoinedAt, d.Nick, d.Deaf, d.Mute, userId}

	err = tx.QueryRowContext(ctx, data.query, args...).Scan(&d.ID)
	if err != nil {
		return err
	}

	return nil
}
