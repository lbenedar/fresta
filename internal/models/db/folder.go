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

type Folder struct {
	ID uint `db:"id"`

	Name    string `db:"name"`
	Sorting string `db:"sorting"`
	Color   string `db:"color"`
	Packs   []string
	Folders Folders
}

func (f *Folder) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO folder (%s, name, sorting, color)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (f *Folder) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{
		id:        strconv.FormatUint(uint64(f.ID), 10),
		fieldName: "folder_id",
		tableName: "folder_packs",
	}
	group, _ := errgroup.WithContext(context.Background())

	InsertSimpleSliceParallel(group, tx, f.Packs, &relId)
	InsertSliceParallel(group, tx, f.Folders, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (f *Folder) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Name, f.Sorting, f.Color}

	err := tx.QueryRowx(data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return f.InsertObjects(tx)
}

func (f *Folder) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, f.Name, f.Sorting, f.Color}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&f.ID)
	if err != nil {
		return err
	}

	return f.InsertObjects(tx)
}

type Folders []*Folder

func (m *Folders) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, name, sorting, color
		FROM folder
		WHERE %s = $1`, data.fieldName)
}

func (f *Folders) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	for i := range *f {
		group.Go(func() error {
			fOne := (*f)[i]

			funcGroup, _ := errgroup.WithContext(context.Background())

			GetSimpleSliceParallel(funcGroup, db, &fOne.Packs,
				&InsertId[uint]{id: fOne.ID, fieldName: "folder_id", tableName: "folder_packs"})
			GetWithCtxParallel(funcGroup, db, &fOne.Folders,
				InsertId[string]{id: strconv.FormatUint(uint64(fOne.ID), 10), fieldName: "folder_id"})

			return funcGroup.Wait()
		})
	}

	return group.Wait()
}

func (f *Folders) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	err := db.SelectContext(ctx, f, data.query, data.id)
	if err != nil {
		return err
	}

	return f.GetObjects(db)
}

type WorldFolder struct {
	ID string

	Name        string
	Type        string
	Folder      string
	Sorting     string
	Description string
	Color       string
	Sort        int
	Stats       Stats
}

func (w *WorldFolder) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO world_folder (game_id, id, name, type, folder, sorting, description, color, sort)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (w *WorldFolder) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: w.ID, fieldName: "world_folder_id"}
	InsertWithCtxParallel(group, tx, w.Stats, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (w *WorldFolder) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, w.ID, w.Name, w.Type, w.Folder, w.Sorting, w.Description, w.Color, w.Sort}

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

func (w *WorldFolder) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, w.ID, w.Name, w.Type, w.Folder, w.Sorting, w.Description, w.Color, w.Sort}

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
