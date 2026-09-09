package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoQuery          = errors.New("Query has not been set")
	ErrorRecordNotFound = errors.New("Record not found")
)

type AllowedIds interface {
	~uint | ~string
}

type InsertId[T AllowedIds] struct {
	id        T
	fieldName string
	tableName string
	query     string
}

type Insertable[T AllowedIds] interface {
	Query(data *InsertId[T])
	Insert(tx *sqlx.Tx, relId *InsertId[T]) error
	InsertCtx(ctx context.Context, tx *sqlx.Tx, relId *InsertId[T]) error
}

func InsertWithCtx[T AllowedIds, I Insertable[T]](tx *sqlx.Tx, data I, relId InsertId[T]) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.Query(&relId)
	return data.InsertCtx(ctx, tx, &relId)
}

func InsertWithCtxParallel[T AllowedIds, I Insertable[T]](g *errgroup.Group, tx *sqlx.Tx, data I, relId InsertId[T]) {
	g.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data.Query(&relId)
		err := data.InsertCtx(ctx, tx, &relId)
		return err
	})
}

func InsertSlice[T AllowedIds, I Insertable[T]](tx *sqlx.Tx, data []I, relId InsertId[T]) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(data) > 0 {
		data[0].Query(&relId)
	}
	for i := range data {
		err := data[i].InsertCtx(ctx, tx, &relId)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertSliceParallelTimeout[T AllowedIds, I Insertable[T]](g *errgroup.Group, tx *sqlx.Tx, data []I, relId InsertId[T], timeout time.Duration) {
	g.Go(func() error {
		wg, _ := errgroup.WithContext(context.Background())

		var err error
		if len(data) > 0 {
			data[0].Query(&relId)
		}
		for i := range data {
			wg.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				err = data[i].InsertCtx(ctx, tx, &relId)
				return err
			})
		}

		return wg.Wait()
	})
}

func InsertSliceParallel[T AllowedIds, I Insertable[T]](g *errgroup.Group, tx *sqlx.Tx, data []I, relId InsertId[T]) {
	g.Go(func() error {
		wg, _ := errgroup.WithContext(context.Background())

		var err error
		if len(data) > 0 {
			data[0].Query(&relId)
		}
		for i := range data {
			wg.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = data[i].InsertCtx(ctx, tx, &relId)
				return err
			})
		}

		return wg.Wait()
	})
}

func InsertSimpleSlice[T AllowedIds, I any](tx *sqlx.Tx, data []I, relId *InsertId[T]) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (%s, value)
		VALUES ($1, $2)`, relId.tableName, relId.fieldName)
	for i := range data {
		_, err := tx.ExecContext(ctx, query, relId.id, data[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertSimpleSliceParallel[T AllowedIds, I any](g *errgroup.Group, tx *sqlx.Tx, data []I, relId *InsertId[T]) {
	g.Go(func() error {
		wg, _ := errgroup.WithContext(context.Background())

		query := fmt.Sprintf(`
			INSERT INTO %s (%s, value)
			VALUES ($1, $2)`, relId.tableName, relId.fieldName)
		for i := range data {
			wg.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err := tx.ExecContext(ctx, query, relId.id, data[i])
				return err
			})
		}

		return wg.Wait()
	})
}

type Gettable[T AllowedIds] interface {
	GetQuery(data *InsertId[T])
	GetCtx(ctx context.Context, db *sqlx.DB, relId *InsertId[T]) error
}

func GetWithCtxParallel[
	T AllowedIds,
	I any,
	PI interface {
		*I
		Gettable[T]
	}](g *errgroup.Group, db *sqlx.DB, data *I, relId InsertId[T]) {
	g.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		PI(data).GetQuery(&relId)

		err := PI(data).GetCtx(ctx, db, &relId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return nil
	})
}

func GetPtrWithCtxParallel[T AllowedIds, I Gettable[T]](g *errgroup.Group, db *sqlx.DB, data *I, relId InsertId[T]) {
	g.Go(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		v := reflect.ValueOf(data).Elem()
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		(*data).GetQuery(&relId)

		err := (*data).GetCtx(ctx, db, &relId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return nil
	})
}

func GetSimpleSliceParallel[T AllowedIds, I any](g *errgroup.Group, db *sqlx.DB, data *[]I, relId *InsertId[T]) {
	g.Go(func() error {
		query := fmt.Sprintf(`
			SELECT value
			FROM %s
			WHERE %s = $1`, relId.tableName, relId.fieldName)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := db.SelectContext(ctx, data, query, relId.id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return nil
	})
}

func IsDoublePointer[T any](v T) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.Ptr {
		return typeOf.Elem().Kind() == reflect.Ptr
	}
	return false
}

func DeleteSetupAll(db *sqlx.DB) error {
	const query = `DELETE FROM setup`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

func DeleteGameAll(db *sqlx.DB) error {
	const query = `DELETE FROM game`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

func DeleteSeqAll(db *sqlx.DB) error {
	const query = `DELETE FROM sqlite_sequence`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

var modelInsertLock sync.Map

func GetMutex(id string) *sync.Mutex {
	val, _ := modelInsertLock.LoadOrStore(id, &sync.Mutex{})

	return val.(*sync.Mutex)
}
