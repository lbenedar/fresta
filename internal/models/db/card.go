package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type CardDeck struct {
	ID string

	Name         string
	Type         string
	Description  string
	Img          string
	Folder       string
	Width        int
	Height       int
	Rotation     int
	Sort         int
	DisplayCount bool
	Stats        Stats
	Ownership    []OwnershipString
	Cards        []*Card
}

func (c *CardDeck) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO card_deck (game_id, id, name, type, description, img, folder, width, height, rotation, sort, display_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (c *CardDeck) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: c.ID, fieldName: "card_deck_id"}
	InsertWithCtxParallel(group, tx, c.Stats, relId)
	InsertSliceParallel(group, tx, c.Ownership, relId)
	InsertSliceParallel(group, tx, c.Cards, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (c *CardDeck) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Name, c.Type, c.Description, c.Img, c.Folder, c.Width, c.Height, c.Rotation, c.Sort, c.DisplayCount}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

func (c *CardDeck) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Name, c.Type, c.Description, c.Img, c.Folder, c.Width, c.Height, c.Rotation, c.Sort, c.DisplayCount}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

type Card struct {
	ID string

	Name        string
	Type        string
	Suit        string
	Description string
	Origin      string
	Width       int
	Height      int
	Rotation    int
	Value       int
	Face        int
	Sort        int
	Drawn       bool
	Back        Back
	Stats       Stats
	Faces       []Face
}

func (c *Card) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO card (card_deck_id, id, name, type, suit, description, origin, width, height, rotation, value, face, sort, drawn)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT(id) DO UPDATE SET
			card_deck_id = EXCLUDED.card_deck_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (c *Card) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: c.ID, fieldName: "card_id"}
	InsertWithCtxParallel(group, tx, c.Back, relId)
	InsertWithCtxParallel(group, tx, c.Stats, relId)
	InsertSliceParallel(group, tx, c.Faces, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (c *Card) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Name, c.Type, c.Suit, c.Description, c.Origin, c.Width, c.Height, c.Rotation, c.Value, c.Face, c.Sort, c.Drawn}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

func (c *Card) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Name, c.Type, c.Suit, c.Description, c.Origin, c.Width, c.Height, c.Rotation, c.Value, c.Face, c.Sort, c.Drawn}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

type Face struct {
	ID uint

	Name string
	Img  string
	Text string
}

func (f Face) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO face (card_id, name, img, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
}

func (f Face) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Name, f.Img, f.Text}

	err := tx.QueryRowx(data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f Face) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Name, f.Img, f.Text}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

type Back struct {
	ID uint

	Name string
	Text string
}

func (b Back) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO back (card_id, name, text)
		VALUES ($1, $2, $3)
		RETURNING id`
}

func (b Back) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, b.Name, b.Text}

	err := tx.QueryRowx(data.query, args...).Scan(&b.ID)
	if err != nil {
		return err
	}

	return nil
}

func (b Back) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, b.Name, b.Text}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&b.ID)
	if err != nil {
		return err
	}

	return nil
}
