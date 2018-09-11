package resolvers

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
	ggg "github.com/sjhitchner/slack-clone/backend/interfaces/context"
)

type MessageResolver struct {
	obj *domain.Message
}

func (t *MessageResolver) Id() graphql.ID {
	return ToID(t.obj.Id)
}

func (t *MessageResolver) Text() string {
	return t.obj.Text
}

func (t *MessageResolver) Timestamp() string {
	return t.obj.Timestamp.Format(time.RFC3339)
}

func (t *MessageResolver) User(ctx context.Context) (*UserResolver, error) {
	user, err := ggg.Interactor(ctx).GetUserById(ctx, t.obj.UserId)
	return &UserResolver{user}, errors.Wrapf(err, "failed getting message user")
}

func (t *MessageResolver) Channel(ctx context.Context) (*ChannelResolver, error) {
	channel, err := ggg.Interactor(ctx).GetChannelById(ctx, t.obj.ChannelId)
	return &ChannelResolver{channel}, errors.Wrapf(err, "failed getting message channel")
}
