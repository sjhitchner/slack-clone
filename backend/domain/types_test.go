package domain

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&DomainSuite{})

type DomainSuite struct {
}

func (s *DomainSuite) Test_User(c *C) {
	{
		user := &User{
			Username: "",
			Email:    "",
			Password: "",
		}

		c.Assert(user.Validate(), NotNil)
	}
	{
		user := &User{
			Username: "username",
			Email:    "",
			Password: "",
		}

		c.Assert(user.Validate(), NotNil)
	}
	{
		user := &User{
			Username: "username",
			Email:    "email",
			Password: "",
		}

		c.Assert(user.Validate(), NotNil)
	}
	{
		user := &User{
			Username: "username",
			Email:    "email",
			Password: "password",
		}

		c.Assert(user.Validate(), NotNil)
	}
}
