package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type TeamResolver struct {
	obj *domain.Team
}

func (t *TeamResolver) Id() graphql.ID {
	return ToID(t.obj.Id)
}

func (t *TeamResolver) Name() string {
	return t.obj.Name
}

func (t *TeamResolver) Owner(ctx context.Context) (*UserResolver, error) {
	user, err := Interactor(ctx).GetUserById(ctx, t.obj.OwnerId)
	return &UserResolver{user}, errors.Wrapf(err, "failed getting team owner")
}

func (t *TeamResolver) Members(ctx context.Context) ([]*UserResolver, error) {
	list, err := Interactor(ctx).ListUsersByTeamId(ctx, t.obj.Id)
	resolvers := make([]*UserResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &UserResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting team members")
}

func (t *TeamResolver) Channels(ctx context.Context) ([]*ChannelResolver, error) {
	list, err := Interactor(ctx).ListChannelsByTeamId(ctx, t.obj.Id)
	resolvers := make([]*ChannelResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &ChannelResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting team channels")
}
