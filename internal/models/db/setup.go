package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Setup struct {
	ID        uint      `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	IsAdmin         bool `db:"is_admin"`
	IsSetup         bool `db:"is_setup"`
	CoreUpdate      CoreUpdate
	FeaturedContent FeaturedContent
	Files           Files
	Options         *SetupOptions
	Release         Release
	Languages       SetupLanguages
	Modules         Modules
	News            NewsSlice
	PackageWarnings PackageWarnings
	Systems         Systems
	Worlds          Worlds
}

func (s *Setup) InsertObjects(tx *sqlx.Tx) error {
	relData := InsertId[uint]{id: s.ID, fieldName: "setup_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, &s.CoreUpdate, relData)
	InsertWithCtxParallel(group, tx, &s.FeaturedContent, relData)
	InsertWithCtxParallel(group, tx, &s.Files, relData)
	InsertWithCtxParallel(group, tx, s.Options, relData)
	InsertWithCtxParallel(group, tx, &s.Release, relData)

	InsertSliceParallel(group, tx, s.Languages, relData)
	InsertSliceParallel(group, tx, s.Modules, relData)
	InsertSliceParallel(group, tx, s.News, relData)
	InsertSliceParallel(group, tx, s.PackageWarnings, relData)

	relDataString := InsertId[string]{
		id:        strconv.FormatUint(uint64(s.ID), 10),
		fieldName: "setup_id",
	}
	InsertSliceParallel(group, tx, s.Systems, relDataString)
	InsertSliceParallel(group, tx, s.Worlds, relDataString)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (s *Setup) Insert(db *sqlx.DB) error {
	tx := db.MustBegin()
	defer tx.Rollback()

	const query = `
		INSERT INTO setup (is_admin, is_setup)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []any{s.IsAdmin, s.IsSetup}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return err
	}

	err = s.InsertObjects(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Setup) GetObjects(db *sqlx.DB) error {
	relData := InsertId[uint]{id: s.ID, fieldName: "setup_id"}
	group, _ := errgroup.WithContext(context.Background())

	GetWithCtxParallel(group, db, &s.CoreUpdate, relData)
	GetWithCtxParallel(group, db, &s.FeaturedContent, relData)
	GetWithCtxParallel(group, db, &s.Files, relData)
	GetPtrWithCtxParallel(group, db, &s.Options, relData)
	GetWithCtxParallel(group, db, &s.Release, relData)

	GetWithCtxParallel(group, db, &s.Languages, relData)
	GetWithCtxParallel(group, db, &s.Modules, relData)
	GetWithCtxParallel(group, db, &s.News, relData)
	GetWithCtxParallel(group, db, &s.PackageWarnings, relData)

	relDataString := InsertId[string]{
		id:        strconv.FormatUint(uint64(s.ID), 10),
		fieldName: "setup_id",
	}
	// GetWithCtxParallel(group, db, &s.Systems, relDataString)
	GetWithCtxParallel(group, db, &s.Worlds, relDataString)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (s *Setup) Get(db *sqlx.DB) error {
	const query = `
		SELECT id, created_at, is_admin, is_setup FROM setup`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.GetContext(ctx, s, query)
	if err != nil {
		return err
	}

	return s.GetObjects(db)
}

type FeaturedContent struct {
	ID uint `db:"id"`

	Title   string `db:"title"`
	Caption string `db:"caption"`
	URL     string `db:"url"`
	Image   string `db:"image"`
}

func (f *FeaturedContent) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO featured_content (setup_id, title, caption, url, image)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
}

func (f *FeaturedContent) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Title, f.Caption, f.URL, f.Image}

	err := tx.QueryRowx(data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *FeaturedContent) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Title, f.Caption, f.URL, f.Image}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *FeaturedContent) GetQuery(data *InsertId[uint]) {
	data.query = `
		SELECT id, title, caption, url, image 
		FROM featured_content 
		WHERE setup_id = $1`
}

func (f *FeaturedContent) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, f, data.query, data.id)
}

type News struct {
	ID uint `db:"id"`

	Title   string `db:"title"`
	Caption string `db:"caption"`
	URL     string `db:"url"`
	Image   string `db:"image"`
}

func (n *News) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO news (setup_id, title, caption, url, image)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
}

func (n *News) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, n.Title, n.Caption, n.URL, n.Image}

	err := tx.QueryRowx(data.query, args...).Scan(&n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (n *News) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, n.Title, n.Caption, n.URL, n.Image}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&n.ID)
	if err != nil {
		return err
	}

	return nil
}

type NewsSlice []*News

func (n *NewsSlice) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, title, caption, url, image
		FROM news 
		WHERE %s = $1`, data.fieldName)
}

func (n *NewsSlice) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, n, data.query, data.id)
}

func HasSetup(db *sqlx.DB) (bool, error) {
	const query = `
		SELECT EXISTS(SELECT 1 FROM setup)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := db.GetContext(ctx, &exists, query)
	if err != nil {
		return false, err
	}
	return exists, nil
}
