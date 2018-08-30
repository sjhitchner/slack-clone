package interactor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type Interactor struct {
	domain.Aggregator
}

func NewInteractor(a domain.Aggregator) *Interactor {
	return &Interactor{a}
}

func (t *Interactor) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return t.InsertUser(ctx, user)
}

func (t *Interactor) RemoveUser(ctx context.Context, user *domain.User) error {
	return errors.New("Not implemented") // t.DeleteUser(ctx, user)
}

func (t *Interactor) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	return t.InsertTeam(ctx, team)
}

func (t *Interactor) RemoveTeam(ctx context.Context, team *domain.Team) error {
	return errors.New("Not implemented") // t.DeleteTeam(ctx, team)
}

func (t *Interactor) CreateChannel(ctx context.Context, channel *domain.Channel) (*domain.Channel, error) {
	return t.InsertChannel(ctx, channel)
}

func (t *Interactor) RemoveChannel(ctx context.Context, channel *domain.Channel) error {
	return errors.New("Not implemented") // t.DeleteChannel(ctx, channel)
}

func (t *Interactor) AddTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	return t.InsertTeamMember(ctx, teamMember)
}

func (t *Interactor) RemoveTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	return errors.New("Not implemented") // t.DeleteTeamMember(ctx, teamMember)
}

func (t *Interactor) AddChannelMember(ctx context.Context, channelMember *domain.ChannelMember) error {
	return t.InsertChannelMember(ctx, channelMember)
}

func (t *Interactor) RemoveChannelMember(ctx context.Context, channelMember *domain.TeamMember) error {
	return errors.New("Not implemented") // t.RemoveTeamMember(ctx, teamMember)
}

func (t *Interactor) SendMessage(ctx context.Context, message *domain.Message) error {
	return t.InsertMessage(ctx, message)
}
