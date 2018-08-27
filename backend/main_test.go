package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "gopkg.in/check.v1"

	db "github.com/sjhitchner/slack-clone/backend/infrastructure/db"
	libsql "github.com/sjhitchner/slack-clone/backend/infrastructure/db/sqlite"
)

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&MainSuite{})

type MainSuite struct {
	dbh    db.DBHandler
	ts     *httptest.Server
	client *http.Client
}

const DBPath = "/tmp/test.sqlite"

const UserTeamListQuery = `
{
  userTeamList(userId: 1) {
    name
    owner {
      username
    }
    members {
      username
    }
    channels {
      name
      owner {
        username
      }
      members {
        username
      }
      messages {
        text
      }
    }
  }
}
`

type Query struct {
	Query string `json:"query"`
}

func (t Query) Reader() io.Reader {
	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(t)
	return buf
}

func (s *MainSuite) SetUpSuite(c *C) {
	var err error

	s.dbh, err = libsql.NewSQLiteDBHandler(DBPath)
	c.Assert(err, IsNil)
	c.Assert(InitializeDBSchema(s.dbh), IsNil)

	schema, err := ioutil.ReadFile(schemaPath)
	c.Assert(err, IsNil)

	handler, err := SetupHandler(s.dbh, string(schema))
	c.Assert(err, IsNil)

	s.ts = httptest.NewServer(handler)
	s.client = &http.Client{}
}

func (s *MainSuite) TearDownSuite(c *C) {
	s.ts.Close()
	s.dbh.Close()
	c.Assert(os.Remove(DBPath), IsNil)
}

func (s *MainSuite) Test_UserTeamListQuery(c *C) {
	var v struct {
		Data struct {
			UserTeamList []struct {
				Name  string `json:"name"`
				Owner struct {
					Username string `json:"username"`
				} `json:"owner"`
				Channels []struct {
					Name string `json:"name"`
				} `json:"channels"`
			} `json:"userTeamList"`
		} `json:"data"`
	}

	s.RunQuery(c, UserTeamListQuery, &v)

	c.Assert(v.Data.UserTeamList, HasLen, 2)

	list1 := v.Data.UserTeamList[0]
	c.Assert(list1.Name, Equals, "Team Steve")
	c.Assert(list1.Owner.Username, Equals, "steve")
	c.Assert(list1.Channels, HasLen, 1)
	c.Assert(list1.Channels[0].Name, Equals, "SteveTalk")

	list2 := v.Data.UserTeamList[1]
	c.Assert(list2.Name, Equals, "Team Simon")
	c.Assert(list2.Owner.Username, Equals, "simon")
	c.Assert(list2.Channels, HasLen, 1)
	c.Assert(list2.Channels[0].Name, Equals, "SimonTalk")
}

func (s *MainSuite) RunQuery(c *C, queryStr string, v interface{}) {
	query := Query{queryStr}
	request, err := http.NewRequest("POST", s.ts.URL, query.Reader())
	c.Assert(err, IsNil)

	response, err := s.client.Do(request)
	c.Assert(err, IsNil)
	defer response.Body.Close()

	c.Assert(json.NewDecoder(response.Body).Decode(v), IsNil)
}
