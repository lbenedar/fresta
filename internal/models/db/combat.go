package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Combat struct {
	ID string

	Type       string
	Scene      string
	Round      int
	Turn       int
	Sort       int
	Active     bool
	Stats      Stats
	Groups     []string
	Combatants []*Combatant
}

func (c *Combat) Query(data *InsertId[uint]) {
	data.query = `
		INSERT INTO combat (game_id, id, type, scene, round, turn, sort, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT(id) DO UPDATE SET
			game_id = EXCLUDED.game_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (c *Combat) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[string]{id: c.ID, fieldName: "combat_id", tableName: "combat_groups"}
	InsertWithCtxParallel(group, tx, c.Stats, relId)
	InsertSimpleSliceParallel(group, tx, c.Groups, &relId)
	InsertSliceParallel(group, tx, c.Combatants, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (c *Combat) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Type, c.Scene, c.Round, c.Turn, c.Sort, c.Active}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

func (c *Combat) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.Type, c.Scene, c.Round, c.Turn, c.Sort, c.Active}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

type Combatant struct {
	ID string

	TokenId    string
	SceneId    string
	ActorId    string
	Type       string
	Img        string
	Group      string
	Initiative int
	Hidden     bool
	Defeated   bool
	Stats      Stats
}

func (c *Combatant) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO combatant (combat_id, id, token_id, scene_id, actor_id, type, img, group_, initiative, hidden, defeated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT(id) DO UPDATE SET
			combat_id = EXCLUDED.combat_id,
			updated_at = datetime('now')
		RETURNING (created_at == updated_at) AS is_inserted`
}

func (c *Combatant) InsertObjects(tx *sqlx.Tx) error {
	relId := InsertId[string]{id: c.ID, fieldName: "combatant_id"}
	group, _ := errgroup.WithContext(context.Background())

	InsertWithCtxParallel(group, tx, c.Stats, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (c *Combatant) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.TokenId, c.SceneId, c.ActorId, c.Type, c.Img, c.Group, c.Initiative, c.Hidden, c.Defeated}

	var isInserted bool
	err := tx.QueryRowx(data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}

func (c *Combatant) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, c.ID, c.TokenId, c.SceneId, c.ActorId, c.Type, c.Img, c.Group, c.Initiative, c.Hidden, c.Defeated}

	var isInserted bool
	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&isInserted)
	if err != nil {
		return err
	}

	if isInserted {
		return c.InsertObjects(tx)
	}
	return nil
}
