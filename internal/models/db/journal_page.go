package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type JournalPage struct {
	ID string

	Name      string
	Type      string
	Src       string
	Sort      int
	Text      PageText
	Title     PageTitle
	Video     PageVideo
	Stats     Stats
	Ownership []OwnershipString
}

func (j *JournalPage) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO journal_page (journal_id, id, name, type, src, sort)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(id) DO UPDATE SET
			journal_id = EXCLUDED.journal_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (j *JournalPage) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: j.ID, fieldName: "journal_page_id"}
	InsertWithCtxParallel(group, tx, j.Text, relId)
	InsertWithCtxParallel(group, tx, j.Title, relId)
	InsertWithCtxParallel(group, tx, j.Video, relId)
	InsertWithCtxParallel(group, tx, j.Stats, relId)

	InsertSliceParallel(group, tx, j.Ownership, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (j *JournalPage) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, j.ID, j.Name, j.Type, j.Src, j.Sort}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return j.InsertObjects(tx)
	}
	return nil
}

func (j *JournalPage) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, j.ID, j.Name, j.Type, j.Src, j.Sort}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return j.InsertObjects(tx)
	}
	return nil
}

type PageText struct {
	ID uint

	Content  string
	Markdown string
	Format   int
}

func (p PageText) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO journal_page_text (%s, content, markdown, format)
		VALUES ($1, $2, $3, $4)`, data.fieldName)
}

func (p PageText) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Content, p.Markdown, p.Format}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p PageText) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Content, p.Markdown, p.Format}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

type PageTitle struct {
	ID uint

	Show  bool
	Level int
}

func (p PageTitle) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO journal_page_title (%s, show, level)
		VALUES ($1, $2, $3)`, data.fieldName)
}

func (p PageTitle) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Show, p.Level}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p PageTitle) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Show, p.Level}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

type PageVideo struct {
	ID uint

	Controls bool
	Volume   float64
}

func (p PageVideo) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO journal_page_video (%s, controls, volume)
		VALUES ($1, $2, $3)`, data.fieldName)
}

func (p PageVideo) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Controls, p.Volume}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p PageVideo) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Controls, p.Volume}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}
