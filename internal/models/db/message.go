package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Message struct {
	ID string

	Blind     bool
	Emote     bool
	Style     int
	Timestamp int64
	Content   string
	Author    string
	Type      string
	Flavor    string
	Sound     string
	Stats     Stats
	Speaker   Speaker
	Whisper   []string
	Rolls     []string
}

func (m *Message) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO message (game_id, id, blind, emote, style, timestamp, content, author, type, flavor, sound)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (m *Message) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{id: m.ID, fieldName: "message_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, m.Stats, relId)
	InsertWithCtxParallel(group, tx, m.Speaker, relId)

	InsertSimpleSliceParallel(group, tx, m.Whisper, &InsertId[string]{id: m.ID, fieldName: "message_id", tableName: "message_whisper"})
	InsertSimpleSliceParallel(group, tx, m.Rolls, &InsertId[string]{id: m.ID, fieldName: "message_id", tableName: "message_rolls"})

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (m *Message) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID, m.Blind, m.Emote, m.Style, m.Timestamp, m.Content, m.Author, m.Type, m.Flavor, m.Sound}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return m.InsertObjects(tx)
	}
	return nil
}

func (m *Message) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID, m.Blind, m.Emote, m.Style, m.Timestamp, m.Content, m.Author, m.Type, m.Flavor, m.Sound}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return m.InsertObjects(tx)
	}
	return nil
}

type Speaker struct {
	ID uint

	Scene string
	Actor string
	Token string
	Alias string
}

func (s Speaker) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO speaker (%s, scene, actor, token, alias)
		VALUES ($1, $2, $3, $4, $5)`, data.fieldName)
}

func (s Speaker) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Scene, s.Actor, s.Token, s.Alias}

	err := tx.QueryRowx(data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Speaker) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.Scene, s.Actor, s.Token, s.Alias}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&s.ID)
	if err != nil {
		return err
	}

	return nil
}
