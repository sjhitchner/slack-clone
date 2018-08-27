package db

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
	"github.com/sjhitchner/slack-clone/backend/infrastructure/db"
)

const SelectChannel = `
SELECT
    channel.id id
  , channel.team_id team_id
  , channel.owner_id owner_id
  , channel.name name
  , channel.is_public is_public
FROM channel`

const SelectChannelById = SelectChannel + `
WHERE channel.id = $1
`

const SelectChannelsByTeamId = SelectChannel + `
WHERE channel.team_id = $1
`

const InsertChannel = `
INSERT INTO channel (
    team_id
  , owner_id
  , name
  , is_public
) VALUES ($1, $2, $3, $4)
`

const UpdateChannel = `
UPDATE channel SET
	 name = $2
  , owner_id = $3
  , is_public = $4
WHERE id = $1
`

const InsertChannelMember = `
INSERT INTO channel_member (channel_id, user_id) VALUES ($1, $2)
`

const DeleteChannelMember = `
DELETE FROM channel_member WHERE channel_id = $1 AND user_id = $2
`

type ChannelDB struct {
	db db.DBHandler
}

func NewChannelDB(db db.DBHandler) *ChannelDB {
	return &ChannelDB{db}
}

func (t *ChannelDB) GetChannelById(ctx context.Context, id int64) (*domain.Channel, error) {
	log.Println(SelectChannelById)

	var obj domain.Channel
	err := t.db.GetById(ctx, &obj, SelectChannelById, id)
	return &obj, errors.Wrapf(err, "error getting channel '%d'", id)
}

func (t *ChannelDB) ListChannelsByTeamId(ctx context.Context, teamId int64) ([]*domain.Channel, error) {
	log.Println(SelectChannelsByTeamId)

	var list []*domain.Channel
	err := t.db.Select(ctx, &list, SelectChannelsByTeamId, teamId)
	return list, errors.Wrapf(err, "error getting channels by team '%d'", teamId)
}

func (t *ChannelDB) CreateChannel(ctx context.Context, channel *domain.Channel) (*domain.Channel, error) {
	log.Println(InsertChannel)

	id, err := t.db.InsertWithId(
		ctx,
		InsertChannel,
		channel.TeamId,
		channel.OwnerId,
		channel.Name,
		channel.IsPublic,
	)
	channel.Id = id
	return channel, errors.Wrapf(err, "unable to insert channel")
}

func (t *ChannelDB) AddChannelMember(ctx context.Context, member *domain.ChannelMember) error {
	log.Println(InsertChannelMember)

	err := t.db.Insert(
		ctx,
		InsertChannelMember,
		member.ChannelId,
		member.UserId,
	)
	return errors.Wrap(err, "unabled to insert channel member")
}

func (t *ChannelDB) DeleteChannelMember(ctx context.Context, member *domain.ChannelMember) error {
	log.Println(DeleteChannelMember)

	_, err := t.db.Delete(
		ctx,
		DeleteChannelMember,
		member.ChannelId,
		member.UserId,
	)
	return errors.Wrap(err, "unabled to delete channel member")
}
