package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Game struct {
	ID        uint      `db:"id"`
	CreatedAt time.Time `db:"created_at"`

	DemoMode   bool   `db:"demo_mode"`
	IdleLogout bool   `db:"idle_logout"`
	Paused     bool   `db:"paused"`
	UserID     string `db:"user_id"`

	Addresses       Addresses
	Files           Files
	Options         GameOptions
	Release         Release
	World           *World
	System          *System
	CoreUpdate      CoreUpdate
	SystemUpdate    SystemUpdate
	ActiveUsers     []string
	Modules         []*Module
	PackageWarnings PackageWarnings
	Packs           []*Pack
	Messages        []*Message
	Combats         []*Combat
	CardDeck        []*CardDeck
	Users           Users
	Macros          []*Macro
	Folders         []*WorldFolder
	Items           []*Item
	Settings        []*Setting
	Journals        []*Journal
	Tables          []*Table
	Playlists       []*Playlist
	Actors          []*Actor
	// Scenes          []Scene
}

func (g *Game) InsertObjects(db *sqlx.DB) error {
	relData := InsertId[uint]{id: g.ID, fieldName: "game_id"}

	err := func() error {
		tx := db.MustBegin()
		defer tx.Rollback()

		groupFunc, _ := errgroup.WithContext(context.Background())

		InsertWithCtxParallel(groupFunc, tx, &g.Addresses, relData)
		InsertWithCtxParallel(groupFunc, tx, &g.Files, relData)
		InsertWithCtxParallel(groupFunc, tx, &g.Options, relData)
		InsertWithCtxParallel(groupFunc, tx, &g.Release, relData)
		InsertWithCtxParallel(groupFunc, tx, &g.CoreUpdate, relData)
		InsertWithCtxParallel(groupFunc, tx, &g.SystemUpdate, relData)

		relDataString := InsertId[string]{id: strconv.FormatUint(uint64(g.ID), 10), fieldName: "game_id"}
		InsertWithCtxParallel(groupFunc, tx, g.World, relDataString)
		InsertWithCtxParallel(groupFunc, tx, g.System, relDataString)

		InsertSimpleSliceParallel(groupFunc, tx, g.ActiveUsers,
			&InsertId[uint]{id: g.ID, fieldName: "game_id", tableName: "active_users"})

		InsertSliceParallel(groupFunc, tx, g.Modules, relData)
		InsertSliceParallel(groupFunc, tx, g.PackageWarnings, relData)
		InsertSliceParallel(groupFunc, tx, g.Packs, relDataString)
		InsertSliceParallel(groupFunc, tx, g.Messages, relData)
		InsertSliceParallel(groupFunc, tx, g.Combats, relData)
		InsertSliceParallel(groupFunc, tx, g.CardDeck, relData)
		InsertSliceParallel(groupFunc, tx, g.Users, relData)
		InsertSliceParallel(groupFunc, tx, g.Macros, relData)
		InsertSliceParallel(groupFunc, tx, g.Folders, relData)

		errFunc := groupFunc.Wait()
		if errFunc != nil && !errors.Is(errFunc, sql.ErrNoRows) {
			return errFunc
		}
		return tx.Commit()
	}()

	err = func() error {
		tx := db.MustBegin()
		defer tx.Rollback()
		groupFunc, _ := errgroup.WithContext(context.Background())

		InsertSliceParallel(groupFunc, tx, g.Settings, relData)
		InsertSliceParallel(groupFunc, tx, g.Journals, relData)
		InsertSliceParallel(groupFunc, tx, g.Tables, relData)
		InsertSliceParallel(groupFunc, tx, g.Playlists, relData)

		errFunc := groupFunc.Wait()
		if errFunc != nil && !errors.Is(errFunc, sql.ErrNoRows) {
			return errFunc
		}
		return tx.Commit()
	}()

	err = func() error {
		tx := db.MustBegin()
		defer tx.Rollback()
		groupFunc, _ := errgroup.WithContext(context.Background())

		InsertSliceParallelTimeout(groupFunc, tx, g.Items,
			InsertId[string]{id: strconv.FormatUint(uint64(g.ID), 10), fieldName: "game_id"},
			15*time.Second)
		InsertSliceParallelTimeout(groupFunc, tx, g.Actors, relData, 15*time.Second)

		errFunc := groupFunc.Wait()
		if errFunc != nil && !errors.Is(errFunc, sql.ErrNoRows) {
			return errFunc
		}
		return tx.Commit()
	}()

	return err
}

func (g *Game) Insert(db *sqlx.DB) error {
	tx := db.MustBegin()
	defer tx.Rollback()

	const query = `
		INSERT INTO game (demo_mode, idle_logout, paused, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	args := []any{g.DemoMode, g.IdleLogout, g.Paused, g.UserID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&g.ID, &g.CreatedAt)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	err = g.InsertObjects(db)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return err
}

func (g *Game) GetObjects(db *sqlx.DB) error {
	group, _ := errgroup.WithContext(context.Background())

	relDataString := InsertId[uint]{
		id:        g.ID,
		fieldName: "game_id",
	}
	GetPtrWithCtxParallel(group, db, &g.World, relDataString)
	GetWithCtxParallel(group, db, &g.Users, relDataString)

	err := group.Wait()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func (g *Game) Get(db *sqlx.DB, id int) error {
	const query = `
		SELECT id, created_at, demo_mode, idle_logout, paused FROM game WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.GetContext(ctx, g, query, id)
	if err != nil {
		return err
	}

	return g.GetObjects(db)
}

type Games []*Game

func (g *Games) GetObjects(db *sqlx.DB) error {
	// relData := InsertId[uint]{id: v.ID, fieldName: "game_id"}
	group, _ := errgroup.WithContext(context.Background())

	for _, v := range *g {
		group.Go(func() error {
			funcGroup, _ := errgroup.WithContext(context.Background())

			relDataString := InsertId[uint]{
				id:        v.ID,
				fieldName: "game_id",
			}
			GetPtrWithCtxParallel(group, db, &v.World, relDataString)
			GetWithCtxParallel(group, db, &v.Users, relDataString)

			return funcGroup.Wait()
		})

	}

	return group.Wait()
}

func (g *Games) Get(db *sqlx.DB) error {
	const query = `
		SELECT id, created_at, demo_mode, idle_logout, paused FROM game`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.SelectContext(ctx, g, query)
	if err != nil {
		return err
	}

	return g.GetObjects(db)
}
