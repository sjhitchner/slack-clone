package resolvers

import (
	"github.com/sjhitchner/slack-clone/backend/domain"
)

type ErrorResolver struct {
	name    string
	field   string
	message string
}

func (t *ErrorResolver) Type() string {
	return t.name
}

func (t *ErrorResolver) Field() string {
	return t.field
}

func (t *ErrorResolver) Message() string {
	return t.message
}

func Errors(errs ...error) *[]*ErrorResolver {
	resolvers := make([]*ErrorResolver, 0, len(errs))
	for _, ierr := range errs {
		switch err := ierr.(type) {
		case *domain.ValidationError:
			resolvers = append(resolvers, &ErrorResolver{
				name:    "validation",
				field:   err.Field,
				message: err.Message,
			})
		default:
		}
	}
	return &resolvers
}
