package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Relationships struct {
	ID uint `db:"id"`

	Systems    RelationshipsData
	Requires   RelationshipsData
	Recommends RelationshipsData
	Conflicts  RelationshipsData
}

func (r Relationships) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO relationships (%s)
		VALUES ($1)
		RETURNING id`, data.fieldName)
}

func (r Relationships) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	err := tx.QueryRowx(data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	InsertSliceParallel(group, tx, r.Systems, InsertId[uint]{id: r.ID, tableName: "relationships_systems"})
	InsertSliceParallel(group, tx, r.Requires, InsertId[uint]{id: r.ID, tableName: "relationships_requires"})
	InsertSliceParallel(group, tx, r.Recommends, InsertId[uint]{id: r.ID, tableName: "relationships_recommends"})
	InsertSliceParallel(group, tx, r.Conflicts, InsertId[uint]{id: r.ID, tableName: "relationships_conflicts"})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (r Relationships) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id}

	mutex := GetMutex("relationships_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())

	InsertSliceParallel(group, tx, r.Systems, InsertId[uint]{id: r.ID, tableName: "relationships_systems"})
	InsertSliceParallel(group, tx, r.Requires, InsertId[uint]{id: r.ID, tableName: "relationships_requires"})
	InsertSliceParallel(group, tx, r.Recommends, InsertId[uint]{id: r.ID, tableName: "relationships_recommends"})
	InsertSliceParallel(group, tx, r.Conflicts, InsertId[uint]{id: r.ID, tableName: "relationships_conflicts"})

	err = group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (r *Relationships) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id
		FROM relationships
		WHERE %s = $1`, data.fieldName)
}

func (r *Relationships) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	GetWithCtxParallel(group, db, &r.Systems, InsertId[uint]{id: r.ID, tableName: "relationships_systems"})
	GetWithCtxParallel(group, db, &r.Requires, InsertId[uint]{id: r.ID, tableName: "relationships_requires"})
	GetWithCtxParallel(group, db, &r.Recommends, InsertId[uint]{id: r.ID, tableName: "relationships_recommends"})
	GetWithCtxParallel(group, db, &r.Conflicts, InsertId[uint]{id: r.ID, tableName: "relationships_conflicts"})

	return group.Wait()
}

func (r *Relationships) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.GetContext(ctx, r, data.query, data.id)
	if err != nil {
		return err
	}

	return r.GetObjects(db)
}

type RelationshipData struct {
	ID uint `db:"id"`

	Key           string `db:"key_"`
	Type          string `db:"type"`
	Manifest      string `db:"manifest"`
	Compatibility Compatibility
}

func (r RelationshipData) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO %s (relationships_id, key_, type, manifest)
		VALUES ($1, $2, $3, $4)`, data.tableName)
}

func (r RelationshipData) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Key, r.Type, r.Manifest}

	err := tx.QueryRowx(data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertWithCtxParallel(group, tx, r.Compatibility, InsertId[string]{id: strconv.FormatUint(uint64(r.ID), 10), fieldName: fmt.Sprintf("%s_id", data.tableName)})

	return group.Wait()
}

func (r RelationshipData) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, r.Key, r.Type, r.Manifest}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&r.ID)
	if err != nil {
		return err
	}

	group, _ := errgroup.WithContext(context.Background())
	InsertWithCtxParallel(group, tx, r.Compatibility, InsertId[string]{id: strconv.FormatUint(uint64(r.ID), 10), fieldName: fmt.Sprintf("%s_id", data.tableName)})

	return group.Wait()
}

type RelationshipsData []RelationshipData

func (r *RelationshipsData) GetQuery(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		SELECT id, key_, type, manifest
		FROM %s
		WHERE relationships_id = $1`, data.tableName)
}

func (r *RelationshipsData) GetObjects(db *sqlx.DB, tableName string) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *r {
		rOne := (*r)[i]
		GetWithCtxParallel(group, db, &rOne.Compatibility,
			InsertId[string]{id: strconv.FormatUint(uint64(rOne.ID), 10), fieldName: fmt.Sprintf("%s_id", tableName)})
	}

	return group.Wait()
}

func (r *RelationshipsData) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, r, data.query, data.id)
	if err != nil {
		return err
	}

	return r.GetObjects(db, data.tableName)
}
