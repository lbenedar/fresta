package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type DiscordCommandHistory struct {
	ID uint `db:"id"`

	Number  uint   `db:"num"`
	Command string `db:"command"`
}

func (d *DiscordCommandHistory) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO discord_command_history (msg_id, num, command)
		VALUES ($1, (SELECT coalesce(max(num), 0) + 1 FROM discord_command_history WHERE msg_id = $2), $3)
		ON CONFLICT(id) DO NOTHING`
}

func (d *DiscordCommandHistory) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, data.id, d.Command}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (d *DiscordCommandHistory) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, data.id, d.Command}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func GetPrevCommandTx(tx *sqlx.Tx, msgId string) (*DiscordCommandHistory, error) {
	const query = `
		SELECT id, num, command 
		FROM discord_command_history 
		WHERE msg_id = $1 
			AND num = (SELECT coalesce(max(num), 0) - 1 FROM discord_command_history WHERE msg_id = $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId, msgId}

	var history DiscordCommandHistory
	err := tx.GetContext(ctx, &history, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &history, nil
}

func GetFirstCommandTx(tx *sqlx.Tx, msgId string) (*DiscordCommandHistory, error) {
	const query = `
		SELECT id, num, command 
		FROM discord_command_history 
		WHERE msg_id = $1 AND num = 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId}

	var history DiscordCommandHistory
	err := tx.GetContext(ctx, &history, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &history, nil
}

func GetLastCommand(db *sqlx.DB, msgId string) (*DiscordCommandHistory, error) {
	const query = `
		SELECT id, num, command 
		FROM discord_command_history 
		WHERE msg_id = $1 
			AND num = (SELECT coalesce(max(num), 0) FROM discord_command_history WHERE msg_id = $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId, msgId}

	var history DiscordCommandHistory
	err := db.GetContext(ctx, &history, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &history, nil
}

func GetPrevCommand(db *sqlx.DB, msgId string) (*DiscordCommandHistory, error) {
	const query = `
		SELECT id, num, command 
		FROM discord_command_history 
		WHERE msg_id = $1 
			AND num = (SELECT coalesce(max(num), 0) - 1 FROM discord_command_history WHERE msg_id = $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId, msgId}

	var history DiscordCommandHistory
	err := db.GetContext(ctx, &history, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &history, nil
}

func GetFirstCommand(db *sqlx.DB, msgId string) (*DiscordCommandHistory, error) {
	const query = `
		SELECT id, num, command 
		FROM discord_command_history 
		WHERE msg_id = $1 AND num = 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId}

	var history DiscordCommandHistory
	err := db.GetContext(ctx, &history, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &history, nil
}

func DeleteLastCommandTx(tx *sqlx.Tx, msgId string) error {
	const query = `
		DELETE FROM discord_command_history 
		WHERE msg_id = $1 AND num = (SELECT coalesce(max(num), 0) FROM discord_command_history WHERE msg_id = $2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId, msgId}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func DeleteAllCommandsTx(tx *sqlx.Tx, msgId string) error {
	const query = `
		DELETE FROM discord_command_history 
		WHERE msg_id = $1 AND num != 1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{msgId}

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
