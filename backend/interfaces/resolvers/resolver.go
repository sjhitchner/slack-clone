package resolvers

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	//"github.com/sjhitchner/graphql-resolver/lib/db"
	//"github.com/sjhitchner/slack-clone/backend/domain"
	ggg "github.com/sjhitchner/slack-clone/backend/interfaces/context"
)

func ToID(i int64) graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", i))
}

type Resolver struct {
	*Mutation
}

func (t *Resolver) Ping(ctx context.Context) string {
	return "Pong"
}

func (t *Resolver) UserTeamList(ctx context.Context, args struct {
	UserId int32
}) ([]*TeamResolver, error) {
	list, err := ggg.Interactor(ctx).ListTeamsByUserId(ctx, int64(args.UserId))

	resolvers := make([]*TeamResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &TeamResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "error getting team %d", args.UserId)
}

func (t *Resolver) ChannelMessageList(ctx context.Context, args struct {
	ChannelId int32
}) ([]*MessageResolver, error) {
	list, err := ggg.Interactor(ctx).ListMessagesByChannelId(ctx, int64(args.ChannelId))

	resolvers := make([]*MessageResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &MessageResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "error getting messasge by channel %d", args.ChannelId)
}

func (t *Resolver) TeamChannelList(ctx context.Context, args struct {
	TeamId int32
}) ([]*ChannelResolver, error) {
	list, err := ggg.Interactor(ctx).ListChannelsByTeamId(ctx, int64(args.TeamId))

	resolvers := make([]*ChannelResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &ChannelResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "error getting channel by team %d", args.TeamId)
}

func (t *Resolver) TeamList(ctx context.Context, args struct {
	OwnerId int32
}) ([]*TeamResolver, error) {
	list, err := ggg.Interactor(ctx).ListTeamsByOwnerId(ctx, int64(args.OwnerId))

	resolvers := make([]*TeamResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &TeamResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "error getting team by owner_id %d", args.OwnerId)
}

func (t *Resolver) Team(ctx context.Context, args struct {
	Id int32
}) (*TeamResolver, error) {
	obj, err := ggg.Interactor(ctx).GetTeamById(ctx, int64(args.Id))
	return &TeamResolver{obj}, errors.Wrapf(err, "error getting team %d", args.Id)
}

func (t *Resolver) Channel(ctx context.Context, args struct {
	Id int32
}) (*ChannelResolver, error) {
	obj, err := ggg.Interactor(ctx).GetChannelById(ctx, int64(args.Id))
	return &ChannelResolver{obj}, errors.Wrapf(err, "error getting channel %d", args.Id)
}

func (t *Resolver) User(ctx context.Context, args struct {
	Id int32
}) (*UserResolver, error) {
	obj, err := ggg.Interactor(ctx).GetUserById(ctx, int64(args.Id))
	return &UserResolver{obj}, errors.Wrapf(err, "error getting user %d", args.Id)
}

func (t *Resolver) UserList(ctx context.Context) ([]*UserResolver, error) {
	list, err := ggg.Interactor(ctx).ListUsers(ctx)

	resolvers := make([]*UserResolver, len(list))
	for i := range resolvers {
		resolvers[i] = &UserResolver{list[i]}
	}
	return resolvers, errors.Wrapf(err, "error getting users")
}
