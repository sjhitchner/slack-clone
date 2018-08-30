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

func NewUserResolver(obj *domain.User) *UserResolver {
	if obj == nil {
		return nil
	}
	return &UserResolver{obj}
}

func (t *UserResolver) Id() graphql.ID {
	return ToID(t.obj.Id)
}

func (t *UserResolver) Username() string {
	return t.obj.Username.String()
}

func (t *UserResolver) Email() string {
	return t.obj.Email.String()
}

func (t *UserResolver) Password() string {
	return t.obj.Password.String()
}

func (t *UserResolver) Teams(ctx context.Context) ([]*TeamResolver, error) {
	list, err := Interactor(ctx).ListTeamsByUserId(ctx, t.obj.Id)
	resolvers := make([]*TeamResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &TeamResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "failed getting user teams")
}
