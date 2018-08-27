package db

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/sjhitchner/slack-clone/backend/domain"
	"github.com/sjhitchner/slack-clone/backend/infrastructure/db"
)

const SelectTeam = `
SELECT 
    team.id id
  , team.owner_id owner_id
  , team.name name
FROM team`

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

const InsertTeam = `
INSERT INTO team (
    owner_id
  , name
) VALUES ($1, $2)
`

const UpdateTeam = `
UPDATE team SET owner_id = $2, name = $3 WHERE id = $1
`

const InsertTeamMember = `
INSERT INTO team_member (team_id, user_id) VALUES ($1 , $2)
`

const DeleteTeamMember = `
DELETE FROM team_member WHERE team_id = $1 AND user_id = $2
`

type TeamDB struct {
	db db.DBHandler
}

func NewTeamDB(db db.DBHandler) *TeamDB {
	return &TeamDB{db}
}

func (t *TeamDB) GetTeamById(ctx context.Context, id int64) (*domain.Team, error) {
	log.Println(SelectTeamById)

	var obj domain.Team
	err := t.db.GetById(ctx, &obj, SelectTeamById, id)
	return &obj, errors.Wrapf(err, "error getting team '%d'", id)
}

func (t *TeamDB) ListTeamsByOwnerId(ctx context.Context, ownerId int64) ([]*domain.Team, error) {
	log.Println(SelectTeamsByOwnerId)

	var list []*domain.Team
	err := t.db.Select(ctx, &list, SelectTeamsByOwnerId, ownerId)
	return list, errors.Wrapf(err, "error getting teams by owner '%d'", ownerId)
}

func (t *TeamDB) ListTeamsByUserId(ctx context.Context, userId int64) ([]*domain.Team, error) {
	log.Println(SelectTeamsByUserId)

	var list []*domain.Team
	err := t.db.Select(ctx, &list, SelectTeamsByUserId, userId)
	return list, errors.Wrapf(err, "error getting teams by user '%d'", userId)
}

//func (t *TeamDB) ListTeams(ctx context.Context) ([]*Team, error) {
//}

func (t *TeamDB) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	log.Println(InsertTeam)

	id, err := t.db.InsertWithId(
		ctx,
		InsertTeam,
		team.OwnerId,
		team.Name,
	)
	team.Id = id
	return team, errors.Wrapf(err, "unable to insert team")
}

func (t *TeamDB) AddTeamMember(ctx context.Context, member *domain.TeamMember) error {
	log.Println(InsertTeamMember)

	err := t.db.Insert(
		ctx,
		InsertTeamMember,
		member.TeamId,
		member.UserId,
	)
	return errors.Wrap(err, "unabled to insert team member")
}

func (t *TeamDB) DeleteTeamMember(ctx context.Context, member *domain.TeamMember) error {
	log.Println(DeleteTeamMember)

	_, err := t.db.Delete(
		ctx,
		DeleteTeamMember,
		member.TeamId,
		member.UserId,
	)
	return errors.Wrap(err, "unabled to delete team member")
}
