package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type SetupLanguage struct {
	ID string `db:"id"`

	Label   string `db:"label"`
	Modules SetupLanguageModules
}

func (l *SetupLanguage) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO setup_language (setup_id, label)
		VALUES ($1, $2)
		RETURNING id`
}

func (l *SetupLanguage) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Label}

	err := tx.QueryRowx(data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertSliceParallel(group, tx, l.Modules, InsertId[string]{id: l.ID})

	return group.Wait()
}

func (l *SetupLanguage) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Label}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertSliceParallel(group, tx, l.Modules, InsertId[string]{id: l.ID})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

type SetupLanguages []*SetupLanguage

func (l *SetupLanguages) GetQuery(data *InsertId[uint]) {
	data.query = `
		SELECT id, label 
		FROM setup_language
		WHERE setup_id = $1`
}

func (l *SetupLanguages) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, l, data.query, data.id)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	for i := range *l {
		GetWithCtxParallel(group, db, &(*l)[i].Modules, InsertId[string]{id: (*l)[i].ID})
	}

	return group.Wait()
}

type SetupLanguageModule struct {
	ID string `db:"id"`

	Label string `db:"label"`
	Path  string `db:"path"`
}

func (l SetupLanguageModule) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO setup_language_module (setup_language_id, id, label, path)
		VALUES ($1, $2, $3, $4)`
}

func (l SetupLanguageModule) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.ID, l.Label, l.Path}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (l SetupLanguageModule) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.ID, l.Label, l.Path}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

type SetupLanguageModules []SetupLanguageModule

func (l *SetupLanguageModules) GetQuery(data *InsertId[string]) {
	data.query = `
		SELECT id, label, path 
		FROM setup_language_module
		WHERE setup_language_id = $1`
}

func (l *SetupLanguageModules) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, l, data.query, data.id)
}

type Language struct {
	ID uint `db:"id"`

	Lang string `db:"lang"`
	Name string `db:"name"`
	Path string `db:"path"`
}

func (l *Language) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO language (%s, lang, name, path)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (l *Language) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Lang, l.Name, l.Path}

	err := tx.QueryRowx(data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l *Language) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, l.Lang, l.Name, l.Path}

	mutex := GetMutex("language_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&l.ID)
	if err != nil {
		return err
	}

	return nil
}

type Languages []*Language

func (l *Languages) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, lang, name, path
		FROM language
		WHERE %s = $1`, data.fieldName)
}

func (l *Languages) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, l, data.query, data.id)
}
