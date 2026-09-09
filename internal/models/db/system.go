package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type System struct {
	ID string `db:"id"`

	Title                 string `db:"title"`
	Description           string `db:"description"`
	URL                   string `db:"url"`
	License               string `db:"license"`
	Bugs                  string `db:"bugs"`
	Changelog             string `db:"changelog"`
	Version               string `db:"version"`
	Manifest              string `db:"manifest"`
	Download              string `db:"download"`
	Background            string `db:"background"`
	PrimaryTokenAttribute string `db:"primary_token_attribute"`
	Availability          int    `db:"availability"`
	Socket                bool   `db:"socket"`
	Protected             bool   `db:"protected"`
	Exclusive             bool   `db:"exclusive_"`
	PersistentStorage     bool   `db:"persistent_storage"`
	Locked                bool   `db:"locked"`
	Owned                 bool   `db:"owned"`
	HasStorage            bool   `db:"has_storage"`
	Compatibility         Compatibility
	Relationships         Relationships
	DocumentTypes         DocumentTypes
	Grid                  *Grid
	Esmodules             []string
	Scripts               []string
	Tags                  []string
	Authors               Authors
	Media                 MediaSlice
	Packs                 Packs
	Styles                StyleSlice
	Languages             Languages
	PackFolders           Folders
}

func (s *System) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO system (setup_id, id, title, description, url, license, bugs, changelog, version, manifest, 
			download, background, primary_token_attribute, availability, socket, protected, exclusive_,
			persistent_storage, locked, owned, has_storage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT(id) DO NOTHING`
}

func (s *System) ConnectGameQuery(data *InsertId[string]) {
	data.query = `
		INSERT INTO game_to_systems (game_id, system_id)
		VALUES ($1, $2)`
}

func (s *System) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{id: s.ID, fieldName: "system_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, s.Compatibility, relId)
	InsertWithCtxParallel(group, tx, s.Relationships, relId)
	InsertWithCtxParallel(group, tx, s.DocumentTypes, relId)
	InsertWithCtxParallel(group, tx, s.Grid, relId)

	esModulesRelId := &InsertId[string]{id: s.ID, fieldName: "system_id", tableName: "es_modules"}
	InsertSimpleSliceParallel(group, tx, s.Esmodules, esModulesRelId)
	scriptRelId := &InsertId[string]{id: s.ID, fieldName: "system_id", tableName: "scripts"}
	InsertSimpleSliceParallel(group, tx, s.Scripts, scriptRelId)
	tagsRelId := &InsertId[string]{id: s.ID, fieldName: "system_id", tableName: "tags"}
	InsertSimpleSliceParallel(group, tx, s.Tags, tagsRelId)

	InsertSliceParallel(group, tx, s.Authors, relId)
	InsertSliceParallel(group, tx, s.Media, relId)
	InsertSliceParallel(group, tx, s.Styles, relId)
	InsertSliceParallel(group, tx, s.Languages, relId)
	InsertSliceParallel(group, tx, s.Packs, relId)
	InsertSliceParallel(group, tx, s.PackFolders, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (s *System) ConnectGame(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID}

	_, err := tx.Exec(data.query, args...)
	return err
}

func (s *System) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *string
	if strings.EqualFold(data.fieldName, "setup_id") {
		dataId = &data.id
	}

	args := []any{dataId, s.ID, s.Title, s.Description, s.URL, s.License, s.Bugs,
		s.Changelog, s.Version, s.Manifest, s.Download, s.Background,
		s.PrimaryTokenAttribute, s.Availability, s.Socket, s.Protected, s.Exclusive,
		s.PersistentStorage, s.Locked, s.Owned, s.HasStorage}

	res, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = s.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		s.ConnectGameQuery(&dataCopy)
		err = s.ConnectGame(tx, &dataCopy)
	}

	return err
}

func (s *System) ConnectGameCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, s.ID}

	_, err := tx.ExecContext(ctx, data.query, args...)
	return err
}

func (s *System) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *string
	if strings.EqualFold(data.fieldName, "setup_id") {
		dataId = &data.id
	}

	args := []any{dataId, s.ID, s.Title, s.Description, s.URL, s.License, s.Bugs,
		s.Changelog, s.Version, s.Manifest, s.Download, s.Background,
		s.PrimaryTokenAttribute, s.Availability, s.Socket, s.Protected, s.Exclusive,
		s.PersistentStorage, s.Locked, s.Owned, s.HasStorage}

	res, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = s.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		s.ConnectGameQuery(&dataCopy)
		err = s.ConnectGameCtx(ctx, tx, &dataCopy)
	}

	return err
}

type Systems []*System

func (s *Systems) GetQuery(data *InsertId[string]) {
	if data.fieldName == "game_id" {
		data.query = `
		SELECT s.id, s.title, s.description, s.url, s.license, s.bugs, s.changelog, s.version, s.manifest, 
			s.download, s.background, s.primary_token_attribute, s.availability, s.socket, s.protected, s.exclusive_,
			s.persistent_storage, s.locked, s.owned, s.has_storage
		FROM system AS s
		JOIN game_to_systems AS gs ON gs.system_id = s.id
		WHERE gs.game_id = $1`
	} else {
		data.query = `
		SELECT s.id, s.title, s.description, s.url, s.license, s.bugs, s.changelog, s.version, s.manifest, 
			s.download, s.background, s.primary_token_attribute, s.availability, s.socket, s.protected, s.exclusive_,
			s.persistent_storage, s.locked, s.owned, s.has_storage
		FROM system AS s
		WHERE s.setup_id = $1`
	}
}

func (s *Systems) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *s {
		group.Go(func() error {
			sOne := (*s)[i]
			relId := InsertId[string]{id: sOne.ID, fieldName: "system_id"}
			funcGroup, _ := errgroup.WithContext(context.Background())

			GetWithCtxParallel(funcGroup, db, &sOne.DocumentTypes, relId)
			GetWithCtxParallel(funcGroup, db, &sOne.Relationships, relId)
			GetWithCtxParallel(funcGroup, db, &sOne.Compatibility, relId)
			GetPtrWithCtxParallel(funcGroup, db, &sOne.Grid, relId)

			scriptRelId := &InsertId[string]{id: sOne.ID, fieldName: "system_id", tableName: "scripts"}
			GetSimpleSliceParallel(group, db, &sOne.Scripts, scriptRelId)
			esModulesRelId := &InsertId[string]{id: sOne.ID, fieldName: "system_id", tableName: "es_modules"}
			GetSimpleSliceParallel(group, db, &sOne.Esmodules, esModulesRelId)
			tagsRelId := &InsertId[string]{id: sOne.ID, fieldName: "system_id", tableName: "tags"}
			GetSimpleSliceParallel(group, db, &sOne.Tags, tagsRelId)

			GetWithCtxParallel(group, db, &sOne.Authors, relId)
			GetWithCtxParallel(group, db, &sOne.Media, relId)
			GetWithCtxParallel(group, db, &sOne.Styles, relId)
			GetWithCtxParallel(group, db, &sOne.Languages, relId)
			GetWithCtxParallel(group, db, &sOne.Packs, relId)
			GetWithCtxParallel(group, db, &sOne.PackFolders, relId)

			return funcGroup.Wait()
		})
	}

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (s *Systems) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, s, data.query, data.id)
	if err != nil {
		return err
	}

	return s.GetObjects(db)
}
