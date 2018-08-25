package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/slack-clone/backend/domain"
)

const SelectChannel = `
SELECT
	channel.id id
  , channel.team_id team_id
  , channel.owner_id owner_id
  , channel.name name
  , channel.is_public is_public
FROM channel
`

const SelectChannelById = SelectChannel + `
WHERE channel.id = $1
`

const SelectChannelsByTeamId = SelectChannel + `
WHERE channel.team_id = $1
`

type ChannelDB struct {
	db db.DBHandler
}

func NewChannelDB(db db.DBHandler) *ChannelDB {
	return &ChannelDB{db}
}

func (t *ChannelDB) GetChannelById(ctx context.Context, id int64) (*domain.Channel, error) {
	var obj domain.Channel
	err := t.db.GetById(ctx, &obj, SelectChannelById, id)
	return &obj, errors.Wrapf(err, "error getting channel '%d'", id)
}

func (t *ChannelDB) ListChannelsByTeamId(ctx context.Context, teamId int64) ([]*domain.Channel, error) {
	var list []*domain.Channel
	err := t.db.Select(ctx, &list, SelectChannelsByTeamId, teamId)
	return list, errors.Wrapf(err, "error getting channels by team '%d'", teamId)
}
