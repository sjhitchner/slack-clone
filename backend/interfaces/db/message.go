package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/slack-clone/backend/domain"
)

const SelectMessage = `
SELECT
    id
  , user_id
  , channel_id
  , text
  , timestamp
FROM message
`

const SelectMessageById = `
WHERE id = $1
`

const SelectMessagesByUserId = `
WHERE user_id = $1
`

const SelectMessagesByChannelId = `
WHERE channel_id = $1
`

type MessageDB struct {
	db db.DBHandler
}

func NewMessageDB(db db.DBHandler) *MessageDB {
	return &MessageDB{db}
}

func (t *MessageDB) GetMessageById(ctx context.Context, id int64) (*domain.Message, error) {
	var obj domain.Message
	err := t.db.GetById(ctx, &obj, SelectMessageById, id)
	return &obj, errors.Wrapf(err, "error getting message '%d'", id)
}

func (t *MessageDB) ListMessagesByUserId(ctx context.Context, userId int64) ([]*domain.Message, error) {
	var list []*domain.Message
	err := t.db.Select(ctx, &list, SelectMessagesByUserId, userId)
	return list, errors.Wrapf(err, "error getting messages by user '%d'", userId)
}

func (t *MessageDB) ListMessagesByChannelId(ctx context.Context, channelId int64) ([]*domain.Message, error) {
	var list []*domain.Message
	err := t.db.Select(ctx, &list, SelectMessagesByChannelId, channelId)
	return list, errors.Wrapf(err, "error getting messages by channel '%d'", channelId)
}
