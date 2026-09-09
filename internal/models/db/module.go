package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Module struct {
	ID string `db:"id"`

	Title             string `db:"title"`
	Description       string `db:"description"`
	URL               string `db:"url"`
	License           string `db:"license"`
	Readme            string `db:"readme"`
	Bugs              string `db:"bugs"`
	Changelog         string `db:"changelog"`
	Version           string `db:"version"`
	Download          string `db:"manifest"`
	Manifest          string `db:"download"`
	Socket            bool   `db:"socket"`
	Protected         bool   `db:"protected"`
	Exclusive         bool   `db:"exclusive_"`
	PersistentStorage bool   `db:"persistent_storage"`
	CoreTranslation   bool   `db:"core_translation"`
	Library           bool   `db:"library"`
	Locked            bool   `db:"locked"`
	Owned             bool   `db:"owned"`
	HasStorage        bool   `db:"has_storage"`
	Active            bool   `db:"active"`
	Availability      int    `db:"availability"`
	DocumentTypes     DocumentTypes
	Relationships     Relationships
	Compatibility     Compatibility
	Scripts           []string
	Esmodules         []string
	Tags              []string
	Authors           Authors
	Media             MediaSlice
	Styles            StyleSlice
	Languages         Languages
	Packs             Packs
	PackFolders       Folders
}

func (m *Module) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO module (setup_id, id, title, description, url, license, readme, bugs,
			changelog, version, manifest, download, socket, protected, exclusive_, persistent_storage,
			core_translation, library, locked, owned, has_storage, active, availability)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23)
		ON CONFLICT(id) DO NOTHING`
}

func (m *Module) ConnectGameQuery(data *InsertId[uint]) {
	data.query = `
		INSERT INTO game_to_module (game_id, module_id)
		VALUES ($1, $2)`
}

func (m *Module) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{id: m.ID, fieldName: "module_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, m.DocumentTypes, relId)
	InsertWithCtxParallel(group, tx, m.Relationships, relId)
	InsertWithCtxParallel(group, tx, m.Compatibility, relId)

	scriptRelId := &InsertId[string]{id: m.ID, fieldName: "module_id", tableName: "scripts"}
	InsertSimpleSliceParallel(group, tx, m.Scripts, scriptRelId)
	esModulesRelId := &InsertId[string]{id: m.ID, fieldName: "module_id", tableName: "es_modules"}
	InsertSimpleSliceParallel(group, tx, m.Esmodules, esModulesRelId)
	tagsRelId := &InsertId[string]{id: m.ID, fieldName: "module_id", tableName: "tags"}
	InsertSimpleSliceParallel(group, tx, m.Tags, tagsRelId)

	InsertSliceParallel(group, tx, m.Authors, relId)
	InsertSliceParallel(group, tx, m.Media, relId)
	InsertSliceParallel(group, tx, m.Styles, relId)
	InsertSliceParallel(group, tx, m.Languages, relId)
	InsertSliceParallel(group, tx, m.Packs, relId)
	InsertSliceParallel(group, tx, m.PackFolders, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (m *Module) ConnectGame(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID}

	_, err := tx.Exec(data.query, args...)
	return err
}

func (m *Module) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *uint
	if strings.EqualFold(data.fieldName, "setup_id") {
		dataId = &data.id
	}

	args := []any{dataId, m.ID, m.Title, m.Description, m.URL, m.License, m.Readme, m.Bugs,
		m.Changelog, m.Version, m.Manifest, m.Download, m.Socket, m.Protected, m.Exclusive, m.PersistentStorage,
		m.CoreTranslation, m.Library, m.Locked, m.Owned, m.HasStorage, m.Active, m.Availability}

	res, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = m.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		m.ConnectGameQuery(&dataCopy)
		err = m.ConnectGame(tx, &dataCopy)
	}

	return err
}

func (m *Module) ConnectGameCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, m.ID}

	_, err := tx.ExecContext(ctx, data.query, args...)
	return err
}

func (m *Module) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *uint
	if strings.EqualFold(data.fieldName, "setup_id") {
		dataId = &data.id
	}

	args := []any{dataId, m.ID, m.Title, m.Description, m.URL, m.License, m.Readme, m.Bugs,
		m.Changelog, m.Version, m.Manifest, m.Download, m.Socket, m.Protected, m.Exclusive, m.PersistentStorage,
		m.CoreTranslation, m.Library, m.Locked, m.Owned, m.HasStorage, m.Active, m.Availability}

	res, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = m.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		m.ConnectGameQuery(&dataCopy)
		err = m.ConnectGameCtx(ctx, tx, &dataCopy)
	}

	return err
}

type Modules []*Module

func (m *Modules) GetQuery(data *InsertId[uint]) {
	if data.fieldName == "game_id" {
		data.query = `
		SELECT m.id, m.title, m.description, m.url, m.license, m.readme, m.bugs,
			m.changelog, m.version, m.manifest, m.download, m.socket, m.protected, m.exclusive_, m.persistent_storage,
			m.core_translation, m.library, m.locked, m.owned, m.has_storage, m.active, m.availability
		FROM module AS m
		JOIN game_to_module AS gm ON gm.module_id = m.id
		WHERE gm.game_id = $1`
	} else {
		data.query = `
		SELECT m.id, m.title, m.description, m.url, m.license, m.readme, m.bugs,
			m.changelog, m.version, m.manifest, m.download, m.socket, m.protected, m.exclusive_, m.persistent_storage,
			m.core_translation, m.library, m.locked, m.owned, m.has_storage, m.active, m.availability
		FROM module AS m
		WHERE m.setup_id = $1`
	}
}

func (m *Modules) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *m {
		group.Go(func() error {
			mOne := (*m)[i]
			relId := InsertId[string]{id: mOne.ID, fieldName: "module_id"}
			funcGroup, _ := errgroup.WithContext(context.Background())

			GetWithCtxParallel(funcGroup, db, &mOne.DocumentTypes, relId)
			GetWithCtxParallel(funcGroup, db, &mOne.Relationships, relId)
			GetWithCtxParallel(funcGroup, db, &mOne.Compatibility, relId)

			scriptRelId := &InsertId[string]{id: mOne.ID, fieldName: "module_id", tableName: "scripts"}
			GetSimpleSliceParallel(group, db, &mOne.Scripts, scriptRelId)
			esModulesRelId := &InsertId[string]{id: mOne.ID, fieldName: "module_id", tableName: "es_modules"}
			GetSimpleSliceParallel(group, db, &mOne.Esmodules, esModulesRelId)
			tagsRelId := &InsertId[string]{id: mOne.ID, fieldName: "module_id", tableName: "tags"}
			GetSimpleSliceParallel(group, db, &mOne.Tags, tagsRelId)

			GetWithCtxParallel(group, db, &mOne.Authors, relId)
			GetWithCtxParallel(group, db, &mOne.Media, relId)
			GetWithCtxParallel(group, db, &mOne.Styles, relId)
			GetWithCtxParallel(group, db, &mOne.Languages, relId)
			GetWithCtxParallel(group, db, &mOne.Packs, relId)
			GetWithCtxParallel(group, db, &mOne.PackFolders, relId)

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (m *Modules) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, m, data.query, data.id)
	if err != nil {
		return err
	}

	return m.GetObjects(db)
}
