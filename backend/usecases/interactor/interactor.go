package interactor

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
)

type Interactor struct {
	domain.Aggregator
}

func NewInteractor(a domain.Aggregator) *Interactor {
	return &Interactor{a}
}

func (t *Interactor) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	log.Println("GetUserById", id)
	return t.Aggregator.GetUserById(ctx, id)
}

func (t *Interactor) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	return t.Aggregator.InsertUser(ctx, user)
}

func (t *Interactor) RemoveUser(ctx context.Context, user *domain.User) error {
	return errors.New("Not implemented") // t.Aggregator.DeleteUser(ctx, user)
}

func (t *Interactor) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	return t.Aggregator.InsertTeam(ctx, team)
}

func (t *Interactor) RemoveTeam(ctx context.Context, team *domain.Team) error {
	return errors.New("Not implemented") // t.Aggregator.DeleteTeam(ctx, team)
}

func (t *Interactor) CreateChannel(ctx context.Context, channel *domain.Channel) (*domain.Channel, error) {
	return t.Aggregator.InsertChannel(ctx, channel)
}

func (t *Interactor) RemoveChannel(ctx context.Context, channel *domain.Channel) error {
	return errors.New("Not implemented") // t.Aggregator.DeleteChannel(ctx, channel)
}

func (t *Interactor) AddTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	return t.Aggregator.InsertTeamMember(ctx, teamMember)
}

func (t *Interactor) RemoveTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	return errors.New("Not implemented") // t.Aggregator.DeleteTeamMember(ctx, teamMember)
}

func (t *Interactor) AddChannelMember(ctx context.Context, channelMember *domain.ChannelMember) error {
	return t.Aggregator.InsertChannelMember(ctx, channelMember)
}

func (t *Interactor) RemoveChannelMember(ctx context.Context, channelMember *domain.ChannelMember) error {
	return errors.New("Not implemented") // t.Aggregator.RemoveTeamMember(ctx, teamMember)
}

func (t *Interactor) SendMessage(ctx context.Context, message *domain.Message) error {
	return t.Aggregator.InsertMessage(ctx, message)
}
