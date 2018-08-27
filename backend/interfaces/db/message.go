package db

import (
	"context"
	"log"

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
FROM message`

const SelectMessageById = SelectMessage + `
WHERE id = $1
`

const SelectMessagesByUserId = SelectMessage + `
WHERE user_id = $1
`

const SelectMessagesByChannelId = SelectMessage + `
WHERE channel_id = $1
`

const InsertMessage = `
INSERT INTO message (
	 user_id
  , channel_id
  , text
) VALUES (
    $1
  , $2
  , $3
)
`

type MessageDB struct {
	db db.DBHandler
}

func NewMessageDB(db db.DBHandler) *MessageDB {
	return &MessageDB{db}
}

func (t *MessageDB) GetMessageById(ctx context.Context, id int64) (*domain.Message, error) {
	log.Println(SelectMessageById)

	var obj domain.Message
	err := t.db.GetById(ctx, &obj, SelectMessageById, id)
	return &obj, errors.Wrapf(err, "error getting message '%d'", id)
}

func (t *MessageDB) ListMessagesByUserId(ctx context.Context, userId int64) ([]*domain.Message, error) {
	log.Println(SelectMessagesByUserId)

	var list []*domain.Message
	err := t.db.Select(ctx, &list, SelectMessagesByUserId, userId)
	return list, errors.Wrapf(err, "error getting messages by user '%d'", userId)
}

func (t *MessageDB) ListMessagesByChannelId(ctx context.Context, channelId int64) ([]*domain.Message, error) {
	log.Println(SelectMessagesByChannelId)

	var list []*domain.Message
	err := t.db.Select(ctx, &list, SelectMessagesByChannelId, channelId)
	return list, errors.Wrapf(err, "error getting messages by channel '%d'", channelId)
}

func (t *MessageDB) SendMessage(ctx context.Context, message *domain.Message) error {
	log.Println(InsertMessage)

	err := t.db.Insert(
		ctx,
		InsertMessage,
		message.UserId,
		message.ChannelId,
		message.Text,
	)
	return errors.Wrapf(err, "error inserting message")
}
