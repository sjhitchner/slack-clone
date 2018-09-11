package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

type Validator interface {
	Validate() error
}

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
	Id       int64    `db:"id"`
	Username Username `db:"username"`
	Email    Email    `db:"email"`
	Password Password `db:"password"`
}

func (t Team) Validate() error {

	if t.OwnerId == 0 {
		return errors.New("Invalid owner id")
	}

	if len(t.Name) < 3 {
		return errors.Errorf("Team name too short '%s'", t.Name)
	}

	return nil
}

// TODO multi errors?
func (t User) Validate() error {
	//err := make([]error, 0, 3)
	//err = append(err, t.Username.Validate()...)
	//err = append(err, t.Email.Validate()...)
	//err = append(err, t.Password.Validate()...)
	if err := t.Username.Validate(); err != nil {
		return err
	}

	if err := t.Email.Validate(); err != nil {
		return err
	}

	if err := t.Password.Validate(); err != nil {
		return err
	}

	return nil
}

// TODO add validation
//go:generate sqltype -type=Username -primative=string
type Username string

func (t Username) String() string { return string(t) }

func (t Username) Validate() error {

	if len(t) < 3 {
		return NewValidationError("username", "too short")
	}

	if !IsAlphanumeric(string(t)) {
		return NewValidationError("username", "not alphanumeric")
	}

	return nil
}

//go:generate sqltype -type=Email -primative=string
type Email string

func (t Email) String() string { return string(t) }

func (t Email) Validate() error {
	if len(t) == 0 {
		return NewValidationError("email", "email is empty")
	}
	return nil
}

//go:generate sqltype -type=Password -primative=string
type Password string

func (t Password) String() string { return string(t) }

func (t Password) Validate() error {
	if len(t) == 0 {
		return NewValidationError("password", "password is empty")
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{field, message}
}

func (t ValidationError) Error() string {
	return fmt.Sprintf(`%s: "%s"`, t.Field, t.Message)
}

type TeamRepo interface {
	GetTeamById(ctx context.Context, id int64) (*Team, error)

	ListTeamsByOwnerId(ctx context.Context, ownerId int64) ([]*Team, error)
	ListTeamsByUserId(ctx context.Context, userId int64) ([]*Team, error)
	// ListTeams(ctx context.Context) ([]*Team, error)

	InsertTeam(ctx context.Context, team *Team) (*Team, error)
	InsertTeamMember(ctx context.Context, teamMember *TeamMember) error
	DeleteTeamMember(ctx context.Context, teamMember *TeamMember) error
}

type ChannelRepo interface {
	GetChannelById(ctx context.Context, id int64) (*Channel, error)

	ListChannelsByTeamId(ctx context.Context, teamId int64) ([]*Channel, error)

	InsertChannel(ctx context.Context, channel *Channel) (*Channel, error)
	InsertChannelMember(ctx context.Context, channelMember *ChannelMember) error
	DeleteChannelMember(ctx context.Context, channelMember *ChannelMember) error
}

type UserRepo interface {
	GetUserById(ctx context.Context, id int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	ListUsers(ctx context.Context) ([]*User, error)
	ListUsersByTeamId(ctx context.Context, teamId int64) ([]*User, error)
	ListUsersByChannelId(ctx context.Context, channelId int64) ([]*User, error)

	InsertUser(ctx context.Context, user *User) (*User, error)
}

type MessageRepo interface {
	GetMessageById(ctx context.Context, id int64) (*Message, error)

	ListMessagesByChannelId(ctx context.Context, channelId int64) ([]*Message, error)
	// ListMessagesByUserId(ctx context.Context, userId int64) ([]*Message, error)

	InsertMessage(ctx context.Context, message *Message) error
}

type Aggregator interface {
	TeamRepo
	ChannelRepo
	UserRepo
	MessageRepo
}

type Interactor interface {
	Aggregator
	CreateUser(ctx context.Context, user *User) (*User, error)
	RemoveUser(ctx context.Context, user *User) error
	CreateTeam(ctx context.Context, team *Team) (*Team, error)
	RemoveTeam(ctx context.Context, team *Team) error
	CreateChannel(ctx context.Context, channel *Channel) (*Channel, error)
	RemoveChannel(ctx context.Context, channel *Channel) error
	AddTeamMember(ctx context.Context, teamMember *TeamMember) error
	RemoveTeamMember(ctx context.Context, teamMember *TeamMember) error
	AddChannelMember(ctx context.Context, channelMember *ChannelMember) error
	RemoveChannelMember(ctx context.Context, channelMember *ChannelMember) error
	SendMessage(ctx context.Context, message *Message) error
}

func (t User) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t Team) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t Channel) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t ChannelMember) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t TeamMember) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (t Message) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
