package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type DiscordInteraction struct {
	ID string

	AppID   string
	GuildId string

	CreatedAt time.Time

	Message *DiscordMessage
	Member  *DiscordMember
	User    *DiscordUser
}

func (d *DiscordInteraction) InsertObjects(db *sqlx.DB) error {
	relData := InsertId[string]{id: d.ID, fieldName: "interaction_id"}

	err := func() error {
		tx := db.MustBegin()
		defer tx.Rollback()

		groupFunc, _ := errgroup.WithContext(context.Background())

		if d.Message != nil {
			InsertWithCtxParallel(groupFunc, tx, d.Message, relData)
		}

		if d.Member != nil {
			InsertWithCtxParallel(groupFunc, tx, d.Member, relData)
		}

		if d.User != nil {
			InsertWithCtxParallel(groupFunc, tx, d.User, relData)
		}

		errFunc := groupFunc.Wait()
		if errFunc != nil && !errors.Is(errFunc, sql.ErrNoRows) {
			return errFunc
		}
		return tx.Commit()
	}()

	return err
}

func (d *DiscordInteraction) Insert(db *sqlx.DB) error {
	err := d.InsertObjects(db)
	if err != nil {
		return err
	}

	tx := db.MustBegin()
	defer tx.Rollback()

	var messageId *string = nil
	var memberId *uint = nil
	var userId *string = nil

	if d.Message != nil {
		messageId = &d.Message.ID
	}

	if d.Member != nil {
		memberId = &d.Member.ID
	}

	if d.User != nil {
		userId = &d.User.ID
	}

	const query = `
		INSERT INTO discord_interaction (id, application_id, guild_id, message_id, member_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`

	args := []any{d.ID, d.AppID, d.GuildId, messageId, memberId, userId}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = tx.QueryRowxContext(ctx, query, args...).Scan(&d.CreatedAt)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}
