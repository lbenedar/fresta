package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type World struct {
	ID string `db:"id"`

	Title             string    `db:"title"`
	Description       string    `db:"description"`
	Version           string    `db:"version"`
	System            string    `db:"system"`
	Background        string    `db:"background"`
	JoinTheme         string    `db:"join_theme"`
	CoreVersion       string    `db:"core_version"`
	SystemVersion     string    `db:"system_version"`
	LastPlayed        string    `db:"last_played"`
	Playtime          int       `db:"playtime"`
	Availability      int       `db:"availability"`
	NextSession       time.Time `db:"next_session"`
	Socket            bool      `db:"socket"`
	Protected         bool      `db:"protected"`
	Exclusive         bool      `db:"exclusive_"`
	PersistentStorage bool      `db:"persistent_storage"`
	Locked            bool      `db:"locked"`
	Owned             bool      `db:"owned"`
	HasStorage        bool      `db:"has_storage"`
	Compatibility     Compatibility
	Relationships     Relationships
	Tags              []string
	Scripts           []string
	Esmodules         []string
	Authors           Authors
	Media             MediaSlice
	Styles            StyleSlice
	Languages         Languages
	Packs             Packs
	PackFolders       Folders
}

func (w *World) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO world (%[1]s, id, title, description, version, system, background, join_theme,
			core_version, system_version, last_played, playtime, availability, next_session, socket,
			protected, exclusive_, persistent_storage, locked, owned, has_storage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT(id) DO UPDATE SET
			%[1]s = EXCLUDED.%[1]s,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`,
		data.fieldName)
}

func (w *World) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{id: w.ID, fieldName: "world_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, w.Compatibility, relId)
	InsertWithCtxParallel(group, tx, w.Relationships, relId)

	esModulesRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "es_modules"}
	InsertSimpleSliceParallel(group, tx, w.Esmodules, esModulesRelId)
	scriptRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "scripts"}
	InsertSimpleSliceParallel(group, tx, w.Scripts, scriptRelId)
	tagsRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "tags"}
	InsertSimpleSliceParallel(group, tx, w.Tags, tagsRelId)

	InsertSliceParallel(group, tx, w.Authors, relId)
	InsertSliceParallel(group, tx, w.Media, relId)
	InsertSliceParallel(group, tx, w.Styles, relId)
	InsertSliceParallel(group, tx, w.Languages, relId)
	InsertSliceParallel(group, tx, w.Packs, relId)
	InsertSliceParallel(group, tx, w.PackFolders, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (w *World) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, w.ID, w.Title, w.Description, w.Version, w.System, w.Background, w.JoinTheme,
		w.CoreVersion, w.SystemVersion, w.LastPlayed, w.Playtime, w.Availability, w.NextSession, w.Socket,
		w.Protected, w.Exclusive, w.PersistentStorage, w.Locked, w.Owned, w.HasStorage}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return w.InsertObjects(tx)
	}
	return nil
}

func (w *World) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, w.ID, w.Title, w.Description, w.Version, w.System, w.Background, w.JoinTheme,
		w.CoreVersion, w.SystemVersion, w.LastPlayed, w.Playtime, w.Availability, w.NextSession, w.Socket,
		w.Protected, w.Exclusive, w.PersistentStorage, w.Locked, w.Owned, w.HasStorage}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return w.InsertObjects(tx)
	}
	return nil
}

func (w *World) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, title, description, version, system, background, join_theme,
			core_version, system_version, last_played, playtime, availability, next_session, socket,
			protected, exclusive_, persistent_storage, locked, owned, has_storage
		FROM world
		WHERE %s = $1`, data.fieldName)
}

func (w *World) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error {
		relId := InsertId[string]{id: w.ID, fieldName: "world_id"}
		funcGroup, _ := errgroup.WithContext(context.Background())

		GetWithCtxParallel(funcGroup, db, &w.Relationships, relId)
		GetWithCtxParallel(funcGroup, db, &w.Compatibility, relId)

		scriptRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "scripts"}
		GetSimpleSliceParallel(group, db, &w.Scripts, scriptRelId)
		esModulesRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "es_modules"}
		GetSimpleSliceParallel(group, db, &w.Esmodules, esModulesRelId)
		tagsRelId := &InsertId[string]{id: w.ID, fieldName: "world_id", tableName: "tags"}
		GetSimpleSliceParallel(group, db, &w.Tags, tagsRelId)

		GetWithCtxParallel(group, db, &w.Authors, relId)
		GetWithCtxParallel(group, db, &w.Media, relId)
		GetWithCtxParallel(group, db, &w.Styles, relId)
		GetWithCtxParallel(group, db, &w.Languages, relId)
		GetWithCtxParallel(group, db, &w.Packs, relId)
		GetWithCtxParallel(group, db, &w.PackFolders, relId)

		return funcGroup.Wait()
	})

	return group.Wait()
}

func (w *World) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.GetContext(ctx, w, data.query, data.id)
	if err != nil {
		return err
	}

	return w.GetObjects(db)
}

type Worlds []*World

func (w *Worlds) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, title, description, version, system, background, join_theme,
			core_version, system_version, last_played, playtime, availability, next_session, socket,
			protected, exclusive_, persistent_storage, locked, owned, has_storage
		FROM world
		WHERE %s = $1`, data.fieldName)
}

func (w *Worlds) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *w {
		group.Go(func() error {
			wOne := (*w)[i]
			relId := InsertId[string]{id: wOne.ID, fieldName: "world_id"}
			funcGroup, _ := errgroup.WithContext(context.Background())

			GetWithCtxParallel(funcGroup, db, &wOne.Relationships, relId)
			GetWithCtxParallel(funcGroup, db, &wOne.Compatibility, relId)

			scriptRelId := &InsertId[string]{id: wOne.ID, fieldName: "world_id", tableName: "scripts"}
			GetSimpleSliceParallel(group, db, &wOne.Scripts, scriptRelId)
			esModulesRelId := &InsertId[string]{id: wOne.ID, fieldName: "world_id", tableName: "es_modules"}
			GetSimpleSliceParallel(group, db, &wOne.Esmodules, esModulesRelId)
			tagsRelId := &InsertId[string]{id: wOne.ID, fieldName: "world_id", tableName: "tags"}
			GetSimpleSliceParallel(group, db, &wOne.Tags, tagsRelId)

			GetWithCtxParallel(group, db, &wOne.Authors, relId)
			GetWithCtxParallel(group, db, &wOne.Media, relId)
			GetWithCtxParallel(group, db, &wOne.Styles, relId)
			GetWithCtxParallel(group, db, &wOne.Languages, relId)
			GetWithCtxParallel(group, db, &wOne.Packs, relId)
			GetWithCtxParallel(group, db, &wOne.PackFolders, relId)

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (w *Worlds) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, w, data.query, data.id)
	if err != nil {
		return err
	}

	return w.GetObjects(db)
}

func GetNotStartedWorlds(db *sqlx.DB) ([]string, error) {
	const query = `
		SELECT id FROM world WHERE game_id IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var worldNames []string
	err := db.SelectContext(ctx, &worldNames, query)
	if err != nil {
		return nil, err
	}
	return worldNames, nil
}

func IsWorldInserted(db *sqlx.DB, worldName string) (bool, error) {
	const query = `
		SELECT EXISTS(SELECT 1 FROM world WHERE id = $1 AND game_id IS NOT NULL)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exist bool
	err := db.GetContext(ctx, &exist, query, worldName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func GetWorld(db *sqlx.DB, worldName string) (*World, error) {
	const query = `
		SELECT id, title, description, version, system, background, join_theme, core_version, system_version,
			last_played, playtime, availability, next_session, socket, protected, exclusive_, persistent_storage,
			locked, owned, has_storage
		FROM world WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var world World
	err := db.QueryRowxContext(ctx, query, worldName).Scan(
		&world.ID, &world.Title, &world.Description, &world.Version, &world.System, &world.Background, &world.JoinTheme,
		&world.CoreVersion, &world.SystemVersion, &world.LastPlayed, &world.Playtime, &world.Availability, &world.NextSession,
		&world.Socket, &world.Protected, &world.Exclusive, &world.PersistentStorage, &world.Locked, &world.Owned, &world.HasStorage,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return &world, nil
}

func GetWorldsAndUser(db *sqlx.DB) ([]*World, error) {
	const query = `
		SELECT id, title, description, version, system, background, join_theme, core_version, system_version,
			last_played, playtime, availability, next_session, socket, protected, exclusive_, persistent_storage,
			locked, owned, has_storage
		FROM world`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var world []*World
	err := db.SelectContext(ctx, &world, query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return world, nil
}
