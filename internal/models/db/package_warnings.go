package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type PackageWarning struct {
	ID uint `db:"id"`

	Key   string `db:"key_"`
	Value *PackageWarningsData
}

func (p *PackageWarning) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO package_warnings (%s, key_)
		VALUES ($1, $2)
		RETURNING id`, data.fieldName)
}

func (p *PackageWarning) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Key}

	err := tx.QueryRowx(data.query, args...).Scan(&p.ID)
	if err != nil {
		return err
	}

	InsertWithCtx(tx, p.Value, InsertId[uint]{id: p.ID})

	return nil
}

func (p *PackageWarning) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.Key}

	mutex := GetMutex("package_warning_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&p.ID)
	if err != nil {
		return err
	}

	return InsertWithCtx(tx, p.Value, InsertId[uint]{id: p.ID})
}

type PackageWarnings []*PackageWarning

func (p *PackageWarnings) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, key_
		FROM package_warnings
		WHERE %s = $1`, data.fieldName)
}

func (p *PackageWarnings) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *p {
		group.Go(func() error {
			pOne := (*p)[i]

			funcGroup, _ := errgroup.WithContext(context.Background())

			GetPtrWithCtxParallel(funcGroup, db, &pOne.Value,
				InsertId[uint]{id: pOne.ID, fieldName: "package_warnings_id"})

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (p *PackageWarnings) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, p, data.query, data.id)
	if err != nil {
		return err
	}

	return p.GetObjects(db)
}

type PackageWarningsData struct {
	ID string `db:"id"`

	Type          string `db:"type"`
	Manifest      string `db:"manifest"`
	Reinstallable bool   `db:"reinstallable"`
	Warning       []string
	Error         []string
}

func (p *PackageWarningsData) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	warningData := &InsertId[string]{id: p.ID, fieldName: "package_warnings_data_id", tableName: "package_warnings_data_warning"}
	InsertSimpleSliceParallel(group, tx, p.Warning, warningData)
	errorData := &InsertId[string]{id: p.ID, fieldName: "package_warnings_data_id", tableName: "package_warnings_data_error"}
	InsertSimpleSliceParallel(group, tx, p.Error, errorData)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (p *PackageWarningsData) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO package_warnings_data (id, type, reinstallable, manifest)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(id) DO NOTHING`
}

func (p *PackageWarningsData) ConnectGameQuery(data *InsertId[uint]) {
	data.query = `
		INSERT INTO package_warnings_to_data (package_warnings_id, package_warnings_data_id)
		VALUES ($1, $2)`
}

func (p *PackageWarningsData) ConnectGame(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID}

	_, err := tx.Exec(data.query, args...)
	return err
}

func (p *PackageWarningsData) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{p.ID, p.Type, p.Reinstallable, p.Manifest}

	res, err := tx.Exec(data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = p.InsertObjects(tx)
	}
	if err != nil {
		return err
	}

	dataCopy := *data
	p.ConnectGameQuery(&dataCopy)

	return p.ConnectGame(tx, &dataCopy)
}

func (p *PackageWarningsData) ConnectGameCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, p.ID}

	_, err := tx.ExecContext(ctx, data.query, args...)
	return err
}

func (p *PackageWarningsData) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{p.ID, p.Type, p.Reinstallable, p.Manifest}

	res, err := tx.ExecContext(ctx, data.query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected != 0 {
		err = p.InsertObjects(tx)
	}
	if err != nil {
		return err
	}

	dataCopy := *data
	p.ConnectGameQuery(&dataCopy)

	return p.ConnectGameCtx(ctx, tx, &dataCopy)
}

func (p *PackageWarningsData) GetQuery(data *InsertId[uint]) {
	data.query = `
		SELECT p.id, p.type, p.reinstallable, p.manifest
		FROM package_warnings_data AS p
		JOIN package_warnings_to_data AS pd ON pd.package_warnings_data_id = p.id
		WHERE pd.package_warnings_id = $1`
}

func (p *PackageWarningsData) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	GetSimpleSliceParallel(group, db, &p.Warning,
		&InsertId[string]{id: p.ID, fieldName: "package_warnings_data_id", tableName: "package_warnings_data_warning"})
	GetSimpleSliceParallel(group, db, &p.Error,
		&InsertId[string]{id: p.ID, fieldName: "package_warnings_data_id", tableName: "package_warnings_data_error"})

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (p *PackageWarningsData) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.GetContext(ctx, p, data.query, data.id)
	if err != nil {
		return err
	}

	return p.GetObjects(db)
}
