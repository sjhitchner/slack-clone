package resolvers

import (
	"context"

	"github.com/pkg/errors"

	//"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/slack-clone/backend/domain"
)

type Mutation struct {
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUserResolver struct {
	obj *domain.User
}

func (t *CreateUserResolver) User(ctx context.Context) (*UserResolver, error) {
	return &UserResolver{t.obj}, nil
}

func (t *Mutation) CreateUser(ctx context.Context, args struct {
	Input CreateUserInput
}) (*CreateUserResolver, error) {

	user := &domain.User{
		Username: args.Input.Username,
		Email:    args.Input.Email,
		Password: args.Input.Password,
	}

	user, err := Aggregator(ctx).CreateUser(ctx, user)
	return &CreateUserResolver{user}, errors.Wrapf(err, "error creating user")
}

type CreateTeamInput struct {
	Name    string
	OwnerId int32
}

type CreateTeamResolver struct {
	obj *domain.Team
}

func (t *CreateTeamResolver) Team(ctx context.Context) (*TeamResolver, error) {
	return &TeamResolver{t.obj}, nil
}

func (t *Mutation) CreateTeam(ctx context.Context, args struct {
	Input CreateTeamInput
}) (*CreateTeamResolver, error) {

	team := &domain.Team{
		OwnerId: int64(args.Input.OwnerId),
		Name:    args.Input.Name,
	}

	team, err := Aggregator(ctx).CreateTeam(ctx, team)
	return &CreateTeamResolver{team}, errors.Wrapf(err, "error creating team")
}

type CreateChannelInput struct {
	TeamId   int32
	Name     string
	OwnerId  int32
	IsPublic bool
}

type CreateChannelResolver struct {
	obj *domain.Channel
}

func (t *CreateChannelResolver) Channel(ctx context.Context) (*ChannelResolver, error) {
	return &ChannelResolver{t.obj}, nil
}

func (t *Mutation) CreateChannel(ctx context.Context, args struct {
	Input CreateChannelInput
}) (*CreateChannelResolver, error) {

	channel := &domain.Channel{
		TeamId:   int64(args.Input.TeamId),
		OwnerId:  int64(args.Input.OwnerId),
		Name:     args.Input.Name,
		IsPublic: args.Input.IsPublic,
	}
	channel, err := Aggregator(ctx).CreateChannel(ctx, channel)
	return &CreateChannelResolver{channel}, errors.Wrapf(err, "error creating channel")
}

type SendMessageInput struct {
	UserId    int32
	ChannelId int32
	Text      string
}

type SendMessageResolver struct {
	Ok bool
}

func (t *SendMemberResolver) Ok() boolean {
	return t.Ok
}

func (t *Mutation) SendMessage(ctx context.Context, args struct {
	Input SendMessageInput
}) (*SendMessageResolver, error) {
	message := &domain.Message{
		UserId:   int64(args.Input.UserId),
		ChannelId:  int64(args.Input.ChannelId),
		Text:     args.Input.Text,
	}
	err := Aggregator(ctx).SendMessage(ctx, message)
	return &SendMessageResolver{err == nil}, errors.Wrapf(err, "error sending message")
}

type AddTeamMemberInput struct {
	TeamId int32
	UserId int32
}

type AddTeamMemberResolver struct {
	Ok bool
}

func (t *AddTeamMemberResolver) Ok() boolean {
	return t.Ok
}

func (t *Mutation) AddTeamMember(ctx context.Context, args struct {
	Input AddTeamInput
}) (*AddTeamMemberResolver, error) {
}


type AddChannelMemberInput struct {
	ChannelId int32
	UserId    int32
}

type AddChannelMemberResolver struct {
	Ok bool
}

func (t *AddChannelResolver) Ok() boolean {
	return t.Ok
}

func (t *Mutation) AddChannelMember(ctx context.Context, args struct {
	Input AddChannelInput
}) (*AddChannelMessageResolver, error) {
}
