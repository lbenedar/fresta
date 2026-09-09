package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Token struct {
	ID uint

	Name             string
	ActorLink        bool
	AppendNumber     bool
	PrependAdjective bool
	LockRotation     bool
	RandomImg        bool
	DisplayName      int
	DisplayBars      int
	Disposition      int
	Rotation         int
	Alpha            int
	Width            float64
	Height           float64
	Ring             *Ring
	Sight            *TokenSight
	Texture          *TokenTexture
	Bar1             TokenBar
	Bar2             TokenBar
	Light            *Light
	Occludable       TokenOccludable
	TurnMarker       TokenTurnMarker
}

func (t *Token) Query(data *InsertId[string]) {
	data.query = `
		INSERT INTO token (actor_id, name, actor_link, append_number, prepend_adjective, lock_rotation, random_img, display_name,
			display_bars, disposition, rotation, alpha, width, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`
}

func (t *Token) InsertObjects(tx *sqlx.Tx) error {
	group, _ := errgroup.WithContext(context.Background())

	relId := InsertId[uint]{id: t.ID, fieldName: "token_id"}
	InsertWithCtxParallel(group, tx, t.Ring, relId)
	InsertWithCtxParallel(group, tx, t.Sight, relId)
	InsertWithCtxParallel(group, tx, t.Texture, relId)
	InsertWithCtxParallel(group, tx, t.Bar1, InsertId[uint]{id: t.ID, fieldName: "token_id", tableName: "token_bar_1"})
	InsertWithCtxParallel(group, tx, t.Bar1, InsertId[uint]{id: t.ID, fieldName: "token_id", tableName: "token_bar_2"})
	InsertWithCtxParallel(group, tx, t.Light, relId)
	InsertWithCtxParallel(group, tx, t.Occludable, relId)
	InsertWithCtxParallel(group, tx, t.TurnMarker, relId)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (t *Token) Insert(tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Name, t.ActorLink, t.AppendNumber, t.PrependAdjective, t.LockRotation, t.RandomImg,
		t.DisplayName, t.DisplayBars, t.Disposition, t.Rotation, t.Alpha, t.Width, t.Height}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return t.InsertObjects(tx)
}

func (t *Token) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[string]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Name, t.ActorLink, t.AppendNumber, t.PrependAdjective, t.LockRotation, t.RandomImg,
		t.DisplayName, t.DisplayBars, t.Disposition, t.Rotation, t.Alpha, t.Width, t.Height}

	mutex := GetMutex("token_insert")
	mutex.Lock()
	defer mutex.Unlock()

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return t.InsertObjects(tx)
}

type TokenTexture struct {
	ID uint

	Src            string
	Fit            string
	Tint           string
	ScaleX         float64
	ScaleY         float64
	OffsetX        float64
	OffsetY        float64
	Rotation       float64
	AnchorX        float64
	AnchorY        float64
	AlphaThreshold float64
}

func (t *TokenTexture) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_texture (%s, src, fit, tint, scale_x, scale_y, offset_x, offset_y, rotation, anchor_x, anchor_y, alpha_threshold)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`, data.fieldName)
}

func (t *TokenTexture) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Src, t.Fit, t.Tint, t.ScaleX, t.ScaleY, t.OffsetX, t.OffsetY, t.Rotation, t.AnchorX, t.AnchorY, t.AlphaThreshold}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenTexture) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Src, t.Fit, t.Tint, t.ScaleX, t.ScaleY, t.OffsetX, t.OffsetY, t.Rotation, t.AnchorX, t.AnchorY, t.AlphaThreshold}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

type TokenSight struct {
	ID uint

	Color       string
	VisionMode  string
	Range       int
	Angle       int
	Attenuation float64
	Brightness  float64
	Enabled     bool
}

func (t *TokenSight) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_sight (%s, color, vision_mode, range_, angle, attenuation, brightness, enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`, data.fieldName)
}

func (t *TokenSight) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Color, t.VisionMode, t.Range, t.Angle, t.Attenuation, t.Brightness, t.Enabled}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenSight) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Color, t.VisionMode, t.Range, t.Angle, t.Attenuation, t.Brightness, t.Enabled}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

type TokenBar struct {
	ID uint

	Attribute string
}

func (t TokenBar) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO %s (%s, attribute)
		VALUES ($1, $2)
		RETURNING id`, data.tableName, data.fieldName)
}

func (t TokenBar) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Attribute}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t TokenBar) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Attribute}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

type TokenOccludable struct {
	ID uint

	Radius int
}

func (t TokenOccludable) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_occludable (%s, radius)
		VALUES ($1, $2)
		RETURNING id`, data.fieldName)
}

func (t TokenOccludable) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Radius}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t TokenOccludable) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Radius}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

type TokenTurnMarker struct {
	ID uint

	Mode        int
	Animation   string
	Src         string
	Disposition bool
}

func (t TokenTurnMarker) Query(data *InsertId[uint]) {
	data.query = fmt.Sprintf(`
		INSERT INTO token_turn_maker (%s, mode, animation, src, disposition)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, data.fieldName)
}

func (t TokenTurnMarker) Insert(tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Mode, t.Animation, t.Src, t.Disposition}

	err := tx.QueryRowx(data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t TokenTurnMarker) InsertCtx(ctx context.Context, tx *sqlx.Tx, data *InsertId[uint]) error {
	if data.query == "" {
		return ErrNoQuery
	}

	args := []any{data.id, t.Mode, t.Animation, t.Src, t.Disposition}

	err := tx.QueryRowxContext(ctx, data.query, args...).Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}
