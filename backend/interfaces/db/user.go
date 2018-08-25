package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/slack-clone/backend/domain"
)

const SelectUser = `
SELECT
    user.id id
  , user.username username
  , user.email email
  , user.password password
FROM user 
`

const SelectUserById = SelectUser + `
WHERE user.id = $1
`

const SelectUserByEmail = SelectUser + `
WHERE user.email = $1
`

const SelectUserByUsername = SelectUser + `
WHERE user.username = $1
`

const SelectUsersByTeamId = SelectUser + `
JOIN team_member ON team_member.user_id = user.id
WHERE team_member.team_id = $1
`

const SelectUsersByChannelId = SelectUser + `
JOIN channel_member ON channel_member.user_id = user.id
WHERE channel_member.user_id = $1
`

const InsertUser = `
INSERT INTO user (
	 username
  , email
  , password
) VALUES ($1, $2, $3)
`

const UpdateUser = `
UPDATE user SET
    username = $2
  , email = $3
  , password = $4
WHERE id = $1
`

type UserDB struct {
	db db.DBHandler
}

func NewUserDB(db db.DBHandler) *UserDB {
	return &UserDB{db}
}

func (t *UserDB) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	var obj domain.User
	err := t.db.GetById(ctx, &obj, SelectUserById, id)
	return &obj, errors.Wrapf(err, "error getting user '%d'", id)
}

func (t *UserDB) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var obj domain.User
	err := t.db.GetById(ctx, &obj, SelectUserByUsername, username)
	return &obj, errors.Wrapf(err, "error getting user by username '%d'", username)
}

func (t *UserDB) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var obj domain.User
	err := t.db.GetById(ctx, &obj, SelectUserByEmail, email)
	return &obj, errors.Wrapf(err, "error getting user by email '%s'", email)
}

func (t *UserDB) ListUsersByTeamId(ctx context.Context, teamId int64) ([]*domain.User, error) {
	var list []*domain.User
	err := t.db.Select(ctx, &list, SelectUsersByTeamId, teamId)
	return list, errors.Wrapf(err, "error getting users by team '%d'", teamId)
}

func (t *UserDB) ListUsersByChannelId(ctx context.Context, channelId int64) ([]*domain.User, error) {
	var list []*domain.User
	err := t.db.Select(ctx, &list, SelectUsersByChannelId, channelId)
	return list, errors.Wrapf(err, "error getting users by channel '%d'", channelId)
}

func (t *UserDB) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	id, err := t.db.InsertWithId(
		ctx,
		InsertUser,
		user.Username,
		user.Email,
		user.Password,
	)
	user.Id = id
	return user, errors.Wrapf(err, "unable to insert user")
}
