package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type ChannelResolver struct {
	obj *domain.Channel
}

func (t *ChannelResolver) Id() graphql.ID {
	return ToID(t.obj.Id)
}

func (t *ChannelResolver) Name() string {
	return t.obj.Name
}

func (t *ChannelResolver) IsPublic() bool {
	return t.obj.IsPublic
}

func (t *ChannelResolver) Owner(ctx context.Context) (*UserResolver, error) {
	user, err := Aggregator(ctx).GetUserById(ctx, t.obj.OwnerId)
	return &UserResolver{user}, errors.Wrapf(err, "failed getting team owner")
}

func (t *ChannelResolver) Team(ctx context.Context) (*TeamResolver, error) {
	user, err := Aggregator(ctx).GetTeamById(ctx, t.obj.TeamId)
	return &TeamResolver{user}, errors.Wrapf(err, "failed getting team owner")
}

func (t *ChannelResolver) Messages(ctx context.Context) ([]*MessageResolver, error) {
	list, err := Aggregator(ctx).ListMessagesByChannelId(ctx, t.obj.Id)
	resolvers := make([]*MessageResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &MessageResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting channel members")
}

func (t *ChannelResolver) Members(ctx context.Context) ([]*UserResolver, error) {
	list, err := Aggregator(ctx).ListUsersByChannelId(ctx, t.obj.Id)
	resolvers := make([]*UserResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &UserResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting channel members")
}
