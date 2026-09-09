package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Pack struct {
	ID string `db:"id"`

	Name        string `db:"name"`
	Label       string `db:"label"`
	Banner      string `db:"banner"`
	Path        string `db:"path"`
	Type        string `db:"type"`
	System      string `db:"system"`
	PackageType string `db:"package_type"`
	PackageName string `db:"package_name"`
	Ownership   Ownership
	Index       IndexSlice
	Folders     PackFolders
}

func (p *Pack) Query(data *InsertId[string]) {
	if !strings.EqualFold(data.fieldName, "game_id") {
		data.query = fmt.Sprintf(`
			INSERT INTO pack (%s, id, name, label, banner, path, type, system, package_type, package_name)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT(id) DO NOTHING`, data.fieldName)
	} else {
		data.query = `
			INSERT INTO pack (module_id, id, name, label, banner, path, type, system, package_type, package_name)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT(id) DO NOTHING`
	}
}

func (p *Pack) ConnectGameQuery(data *InsertId[string]) {
	data.query = `
		INSERT INTO game_to_pack (game_id, pack_id)
		VALUES ($1, $2)`
}

func (p *Pack) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, p.Ownership,
		InsertId[string]{id: p.ID, fieldName: "pack_id"})

	relId := InsertId[string]{id: p.ID, fieldName: "pack_id"}
	InsertSliceParallel(group, tx, p.Index, relId)
	InsertSliceParallel(group, tx, p.Folders, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (p *Pack) ConnectGame(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID}

	_, err := tx.Exec(data.query, args...)
	return err
}

func (p *Pack) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *string
	if !strings.EqualFold(data.fieldName, "game_id") {
		dataId = &data.id
	}

	args := []any{dataId, p.ID, p.Name, p.Label, p.Banner, p.Path, p.Type,
		p.System, p.PackageType, p.PackageName}

	mutex := GetMutex("pack_insert")
	mutex.Lock()
	defer mutex.Unlock()

	res, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = p.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		p.ConnectGameQuery(&dataCopy)
		err = p.ConnectGame(tx, &dataCopy)
	}

	return err
}

func (p *Pack) ConnectGameCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID}

	_, err := tx.ExecContext(ctx, data.query, args...)
	return err
}

func (p *Pack) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	var dataId *string
	if !strings.EqualFold(data.fieldName, "game_id") {
		dataId = &data.id
	}

	args := []any{dataId, p.ID, p.Name, p.Label, p.Banner, p.Path, p.Type,
		p.System, p.PackageType, p.PackageName}

	mutex := GetMutex("pack_insert")
	mutex.Lock()
	defer mutex.Unlock()

	res, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = p.InsertObjects(tx)
		if err != nil {
			return err
		}
	}

	if dataId == nil {
		dataCopy := *data
		p.ConnectGameQuery(&dataCopy)
		err = p.ConnectGameCtx(ctx, tx, &dataCopy)
	}

	return err
}

type Packs []*Pack

func (p *Packs) GetQuery(data *InsertId[string]) {
	if data.fieldName == "game_id" {
		data.query = `
		SELECT p.id, p.name, p.label, p.banner, p.path, p.type, p.system, p.package_type, p.package_name
		FROM pack p
		JOIN game_to_pack AS gp ON gp.pack_id = p.id
		WHERE gm.game_id = $1`
	} else {
		data.query = fmt.Sprintf(`
		SELECT p.id, p.name, p.label, p.banner, p.path, p.type, p.system, p.package_type, p.package_name
		FROM pack p
		WHERE p.%s = $1`, data.fieldName)
	}
}

func (p *Packs) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *p {
		group.Go(func() error {
			pOne := (*p)[i]
			relId := InsertId[string]{id: pOne.ID, fieldName: "pack_id"}
			funcGroup, _ := errgroup.WithContext(context.Background())

			GetWithCtxParallel(funcGroup, db, &pOne.Ownership, relId)
			GetWithCtxParallel(funcGroup, db, &pOne.Index, relId)
			GetWithCtxParallel(funcGroup, db, &pOne.Folders, relId)

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (p *Packs) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, p, data.query, data.id)
	if err != nil {
		return err
	}

	return p.GetObjects(db)
}

type PackFolder struct {
	ID string `db:"id"`

	Description string `db:"description"`
	Name        string `db:"name"`
	Sorting     string `db:"sorting"`
	Type        string `db:"type"`
	Sort        int    `db:"sort"`
}

func (p *PackFolder) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO pack_folder (%s, id, description, name, sorting, type, sort)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, data.fieldName)
}

func (p *PackFolder) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID, p.Description, p.Name, p.Sorting, p.Type, p.Sort}

	_, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PackFolder) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID, p.Description, p.Name, p.Sorting, p.Type, p.Sort}

	_, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	return nil
}

type PackFolders []*PackFolder

func (p *PackFolders) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, description, name, sorting, type, sort
		FROM pack_folder 
		WHERE %s = $1`, data.fieldName)
}

func (p *PackFolders) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.SelectContext(ctx, p, data.query, data.id)
}
