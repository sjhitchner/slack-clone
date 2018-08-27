package domain

import (
	"context"
	"time"
)

type Team struct {
	Id      int64  `db:"id"`
	OwnerId int64  `db:"owner_id"`
	Name    string `db:"name"`
}

type TeamMember struct {
	Id     int64 `db:"id"`
	TeamId int64 `db:"team_id"`
	UserId int64 `db:"user_id"`
}

type Channel struct {
	Id       int64  `db:"id"`
	TeamId   int64  `db:"team_id"`
	OwnerId  int64  `db:"owner_id"`
	Name     string `db:"name"`
	IsPublic bool   `db:"is_public"`
}

type ChannelMember struct {
	Id        int64 `db:"id"`
	UserId    int64 `db:"user_id"`
	ChannelId int64 `db:"channel_id"`
}

type Message struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	ChannelId int64     `db:"channel_id"`
	Text      string    `db:"text"`
	Timestamp time.Time `db:"timestamp"`
}

type User struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type TeamRepo interface {
	GetTeamById(ctx context.Context, id int64) (*Team, error)

	ListTeamsByOwnerId(ctx context.Context, ownerId int64) ([]*Team, error)
	ListTeamsByUserId(ctx context.Context, userId int64) ([]*Team, error)
	// ListTeams(ctx context.Context) ([]*Team, error)

	CreateTeam(ctx context.Context, team *Team) (*Team, error)
	AddTeamMember(ctx context.Context, teamMember *TeamMember) error
	DeleteTeamMember(ctx context.Context, teamMember *TeamMember) error
}

type ChannelRepo interface {
	GetChannelById(ctx context.Context, id int64) (*Channel, error)

	ListChannelsByTeamId(ctx context.Context, teamId int64) ([]*Channel, error)

	CreateChannel(ctx context.Context, channel *Channel) (*Channel, error)
	AddChannelMember(ctx context.Context, channelMember *ChannelMember) error
	DeleteChannelMember(ctx context.Context, channelMember *ChannelMember) error
}

type UserRepo interface {
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	ListUsersByTeamId(ctx context.Context, teamId int64) ([]*User, error)
	ListUsersByChannelId(ctx context.Context, channelId int64) ([]*User, error)

	CreateUser(ctx context.Context, user *User) (*User, error)
}

type MessageRepo interface {
	GetMessageById(ctx context.Context, id int64) (*Message, error)

	ListMessagesByChannelId(ctx context.Context, channelId int64) ([]*Message, error)
	// ListMessagesByUserId(ctx context.Context, userId int64) ([]*Message, error)

	SendMessage(ctx context.Context, message *Message) error
}

type Aggregator interface {
	TeamRepo
	ChannelRepo
	UserRepo
	MessageRepo
}
