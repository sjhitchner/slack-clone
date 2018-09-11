package resolvers

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
	ggg "github.com/sjhitchner/slack-clone/backend/interfaces/context"
)

type Mutation struct {
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUserResolver struct {
	obj    *domain.User
	errors []error
}

func NewCreateUserResolver(user *domain.User, err ...error) *CreateUserResolver {
	return &CreateUserResolver{
		user,
		err,
	}
}

func (t *CreateUserResolver) Ok() bool {
	return len(t.errors) == 0
}

func (t *CreateUserResolver) User(ctx context.Context) (*UserResolver, error) {
	return NewUserResolver(t.obj), nil
}

func (t *CreateUserResolver) Errors() *[]*ErrorResolver {
	return Errors(t.errors...)
}

func (t *Mutation) CreateUser(ctx context.Context, args struct {
	Input CreateUserInput
}) (*CreateUserResolver, error) {

	user := &domain.User{
		Username: domain.Username(args.Input.Username),
		Email:    domain.Email(args.Input.Email),
		Password: domain.Password(args.Input.Password),
	}

	// TODO Move this to an ggg.Interactor
	if err := user.Validate(); err != nil {
		return NewCreateUserResolver(nil, err), nil
	}

	user, err := ggg.Interactor(ctx).CreateUser(ctx, user)
	return NewCreateUserResolver(user, err), nil
	//}, errors.Wrapf(err, "error creating user")
}

type CreateTeamInput struct {
	Name string
}

type CreateTeamResolver struct {
	obj    *domain.Team
	errors []error
}

func NewCreateTeamResolver(team *domain.Team, err ...error) *CreateTeamResolver {
	return &CreateTeamResolver{
		team,
		err,
	}
}

func (t *CreateTeamResolver) Ok(ctx context.Context) bool {
	return len(t.errors) == 0
}

func (t *CreateTeamResolver) Team(ctx context.Context) (*TeamResolver, error) {
	return &TeamResolver{t.obj}, nil
}

func (t *CreateTeamResolver) Errors() *[]*ErrorResolver {
	return Errors(t.errors...)
}

func (t *Mutation) CreateTeam(ctx context.Context, args struct {
	Input CreateTeamInput
}) (*CreateTeamResolver, error) {

	team := &domain.Team{
		Name: args.Input.Name,
	}

	team, err := ggg.Interactor(ctx).CreateTeam(ctx, team)
	return NewCreateTeamResolver(team, err), nil
	//, errors.Wrapf(err, "error creating team")
}

type CreateChannelInput struct {
	TeamId   int32
	Name     string
	OwnerId  int32
	IsPublic bool
}

type CreateChannelResolver struct {
	obj    *domain.Channel
	errors []error
}

func NewCreateChannelResolver(obj *domain.Channel, err ...error) *CreateChannelResolver {
	return &CreateChannelResolver{
		obj,
		err,
	}
}

func (t *CreateChannelResolver) Ok(ctx context.Context) bool {
	return len(t.errors) == 0
}

func (t *CreateChannelResolver) Channel(ctx context.Context) (*ChannelResolver, error) {
	return &ChannelResolver{t.obj}, nil
}

func (t *CreateChannelResolver) Errors() *[]*ErrorResolver {
	return Errors(t.errors...)
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
	channel, err := ggg.Interactor(ctx).CreateChannel(ctx, channel)
	return NewCreateChannelResolver(channel, err), nil
	//, errors.Wrapf(err, "error creating channel")
}

type SendMessageInput struct {
	UserId    int32
	ChannelId int32
	Text      string
}

type SendMessageResolver struct {
	ok bool
}

func (t *SendMessageResolver) Ok() bool {
	return t.ok
}

func (t *Mutation) SendMessage(ctx context.Context, args struct {
	Input SendMessageInput
}) (*SendMessageResolver, error) {
	message := &domain.Message{
		UserId:    int64(args.Input.UserId),
		ChannelId: int64(args.Input.ChannelId),
		Text:      args.Input.Text,
	}
	err := ggg.Interactor(ctx).SendMessage(ctx, message)
	return &SendMessageResolver{err == nil}, errors.Wrapf(err, "error sending message")
}

type TeamMemberInput struct {
	TeamId int32
	UserId int32
}

type TeamMemberResolver struct {
	ok bool
}

func (t *TeamMemberResolver) Ok() bool {
	return t.ok
}

func (t *Mutation) AddTeamMember(ctx context.Context, args struct {
	Input TeamMemberInput
}) (*TeamMemberResolver, error) {
	member := &domain.TeamMember{
		TeamId: int64(args.Input.TeamId),
		UserId: int64(args.Input.UserId),
	}
	err := ggg.Interactor(ctx).AddTeamMember(ctx, member)
	return &TeamMemberResolver{err == nil}, errors.Wrapf(err, "error adding team member")
}

func (t *Mutation) DeleteTeamMember(ctx context.Context, args struct {
	Input TeamMemberInput
}) (*TeamMemberResolver, error) {
	member := &domain.TeamMember{
		TeamId: int64(args.Input.TeamId),
		UserId: int64(args.Input.UserId),
	}
	err := ggg.Interactor(ctx).DeleteTeamMember(ctx, member)
	return &TeamMemberResolver{err == nil}, errors.Wrapf(err, "error deleting team member")
}

type ChannelMemberInput struct {
	ChannelId int32
	UserId    int32
}

type ChannelMemberResolver struct {
	ok bool
}

func (t *ChannelMemberResolver) Ok() bool {
	return t.ok
}

func (t *Mutation) AddChannelMember(ctx context.Context, args struct {
	Input ChannelMemberInput
}) (*ChannelMemberResolver, error) {
	member := &domain.ChannelMember{
		ChannelId: int64(args.Input.ChannelId),
		UserId:    int64(args.Input.UserId),
	}
	err := ggg.Interactor(ctx).AddChannelMember(ctx, member)
	return &ChannelMemberResolver{err == nil}, errors.Wrapf(err, "error adding channel member")
}

func (t *Mutation) DeleteChannelMember(ctx context.Context, args struct {
	Input ChannelMemberInput
}) (*ChannelMemberResolver, error) {
	member := &domain.ChannelMember{
		ChannelId: int64(args.Input.ChannelId),
		UserId:    int64(args.Input.UserId),
	}
	err := ggg.Interactor(ctx).DeleteChannelMember(ctx, member)
	return &ChannelMemberResolver{err == nil}, errors.Wrapf(err, "error deleting channel member")
}
