package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Ownership struct {
	ID uint `db:"id"`

	Player    string `db:"player"`
	Trusted   string `db:"trusted"`
	Assistant string `db:"assistant"`
}

func (o Ownership) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO ownership (%s, player, trusted, assistant)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, data.fieldName)
}

func (o Ownership) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, o.Player, o.Trusted, o.Assistant}

	err := tx.QueryRowx(data.query, args...).Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}

func (o Ownership) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, o.Player, o.Trusted, o.Assistant}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}

func (o *Ownership) GetQuery(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		SELECT id, player, trusted, assistant
		FROM ownership
		WHERE %s = $1`, data.fieldName)

}

func (o *Ownership) GetCtx(ctx context.Context, db *sqlx.DB, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	return db.GetContext(ctx, o, data.query, data.id)
}

type OwnershipString struct {
	ID uint

	Key   string
	Value int
}

func (o OwnershipString) Query(data *InsertId[string]) {
	data.query = fmt.Sprintf(`
		INSERT INTO ownership_string (%s, key_, value)
		VALUES ($1, $2, $3)
		RETURNING id`, data.fieldName)
}

func (o OwnershipString) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, o.Key, o.Value}

	err := tx.QueryRowx(data.query, args...).Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}

func (o OwnershipString) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, o.Key, o.Value}

	mutex := GetMutex("ownership_string_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}
