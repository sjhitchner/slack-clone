package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type UserResolver struct {
	obj *domain.User
}

func (t *UserResolver) Id() graphql.ID {
	return graphql.ID(t.obj.Id)
}

func (t *UserResolver) Username() string {
	return t.obj.Username
}

func (t *UserResolver) Email() string {
	return t.obj.Email
}

func (t *UserResolver) Password() string {
	return t.obj.Password
}

func (t *UserResolver) Teams(ctx context.Context) ([]*TeamResolver, error) {
	list, err := Aggregator(ctx).ListTeamsByUserId(ctx, t.obj.Id)
	resolvers := make([]*TeamResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &TeamResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting user teams")
}
