package db

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type DiscordMessage struct {
	ID string

	GuildId string
	Content string

	Timestamp       time.Time
	EditedTimestamp *time.Time

	MsgReference discordgo.MessageReference

	Command string

	History []DiscordCommandHistory
}

func (d *DiscordMessage) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO discord_message (id, guild_id, content, timestamp, edited_timestamp)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(id) DO NOTHING`
}

func (d *DiscordMessage) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{d.ID, d.GuildId, d.Content, d.Timestamp, d.EditedTimestamp}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordMessage) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{d.ID, d.GuildId, d.Content, d.Timestamp, d.EditedTimestamp}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	firstCommand, err := GetFirstCommandTx(tx, d.ID)
	if err != nil {
		return err
	}

	prevCommand, err := GetPrevCommandTx(tx, d.ID)
	if err != nil {
		return err
	}

	switch d.Command {
	case firstCommand.Command:
		err = DeleteAllCommandsTx(tx, d.ID)
		if err != nil {
			return err
		}
	case prevCommand.Command:
		err = DeleteLastCommandTx(tx, d.ID)
		if err != nil {
			return err
		}
	default:
		history := DiscordCommandHistory{Command: d.Command}
		data := &InsertId[string]{id: d.ID}

		history.Query(data)
		err = history.InsertCtx(ctx, tx, data)
		if err != nil {
			return err
		}
	}

	return nil
}
