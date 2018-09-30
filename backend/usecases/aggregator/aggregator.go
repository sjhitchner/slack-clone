package aggregator

import (
	"github.com/sjhitchner/slack-clone/backend/domain"
)

type Aggregator struct {
	domain.UserRepo
	domain.TeamRepo
	domain.ChannelRepo
	domain.MessageRepo
}

func NewAggregator(
	u domain.UserRepo,
	t domain.TeamRepo,
	c domain.ChannelRepo,
	m domain.MessageRepo) *Aggregator {
	return &Aggregator{
		u,
		t,
		c,
		m,
	}
}
