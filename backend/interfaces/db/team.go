package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/slack-clone/backend/domain"
)

const SelectTeam = `
SELECT 
    team.id id
  , team.owner_id owner_id
  , team.name name
FROM team
`

const SelectTeamById = SelectTeam + `
WHERE team.id = $1
`

const SelectTeamsByOwnerId = SelectTeam + `
WHERE team.owner_id = $1
`

const SelectTeamsByUserId = SelectTeam + `
JOIN team_member ON team.id = team_member.team_id
WHERE team_member.user_id = $1
`

type TeamDB struct {
	db db.DBHandler
}

func NewTeamDB(db db.DBHandler) *TeamDB {
	return &TeamDB{db}
}

func (t *TeamDB) GetTeamById(ctx context.Context, id int64) (*domain.Team, error) {
	var obj domain.Team
	err := t.db.GetById(ctx, &obj, SelectTeamById, id)
	return &obj, errors.Wrapf(err, "error getting team '%d'", id)
}

func (t *TeamDB) ListTeamsByOwnerId(ctx context.Context, ownerId int64) ([]*domain.Team, error) {
	var list []*domain.Team
	err := t.db.Select(ctx, &list, SelectTeamsByOwnerId, ownerId)
	return list, errors.Wrapf(err, "error getting teams by owner '%d'", ownerId)
}

func (t *TeamDB) ListTeamsByUserId(ctx context.Context, userId int64) ([]*domain.Team, error) {
	var list []*domain.Team
	err := t.db.Select(ctx, &list, SelectTeamsByUserId, userId)
	return list, errors.Wrapf(err, "error getting teams by user '%d'", userId)
}

//func (t *TeamDB) ListTeams(ctx context.Context) ([]*Team, error) {
//}
