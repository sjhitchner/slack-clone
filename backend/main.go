package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/graph-gophers/graphql-go"
	"github.com/sjhitchner/slack-clone/backend/domain"
	"github.com/sjhitchner/slack-clone/backend/interfaces/db"
	"github.com/sjhitchner/slack-clone/backend/interfaces/graphql"
	"github.com/sjhitchner/slack-clone/backend/interfaces/resolvers"

	libdb "github.com/sjhitchner/graphql-resolver/lib/db"
	//libsql "github.com/sjhitchner/graphql-resolver/lib/db/psql"
	libsql "github.com/sjhitchner/graphql-resolver/lib/db/sqlite"
)

// https://github.com/graph-gophers/dataloader
var (
	initializeDB bool

	sqlitePath string
	schemaPath string
)

func init() {
	flag.BoolVar(&initializeDB, "initialize", false, "Initialize the DB")
	flag.StringVar(&sqlitePath, "sqlite", ":memory:", "Path to sqlite db")
	flag.StringVar(&schemaPath, "schema", "schema.gql", "Path to graphql schema")
}

func main() {
	flag.Parse()

	//dbh, err := libsql.NewPSQLDBHandler("localhost", "
	dbh, err := libsql.NewSQLiteDBHandler(sqlitePath)
	CheckError(err)

	if initializeDB {
		CheckError(InitializeDBSchema(dbh))
		os.Exit(0)
	}

	schema, err := ioutil.ReadFile(schemaPath)
	CheckError(err)

	aggregator := struct {
		domain.UserRepo
		domain.TeamRepo
		domain.ChannelRepo
		domain.MessageRepo
	}{
		db.NewUserDB(dbh),
		db.NewTeamDB(dbh),
		db.NewChannelDB(dbh),
		db.NewMessageDB(dbh),
	}

	handler := graphql.NewHandler(string(schema), &resolvers.Resolver{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "agg", aggregator)

	http.Handle("/graphql", graphql.WrapContext(ctx, handler))
	http.Handle("/", graphql.NewGraphiQLHandler(string(schema)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func InitializeDBSchema(dbh libdb.DBHandler) error {

	schema := `
CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, username TEXT NOT NULL UNIQUE
	, email TEXT NOT NULL UNIQUE 
	, password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS team (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, owner_id INTEGER NOT NULL
	, name TEXT NOT NULL 
	, FOREIGN KEY (owner_id) REFERENCES user(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
   , UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS channel (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, team_id INTEGER NOT NULL
	, owner_id INTEGER NOT NULL
	, name TEXT NOT NULL
	, is_public BOOLEAN NOT NULL
	, FOREIGN KEY (team_id) REFERENCES team(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
	, FOREIGN KEY (owner_id) REFERENCES user(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
   , UNIQUE(team_id, name)
);

CREATE TABLE IF NOT EXISTS team_member (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, team_id INTEGER NOT NULL
	, user_id INTEGER NOT NULL
	, FOREIGN KEY (team_id) REFERENCES team(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
	, FOREIGN KEY (user_id) REFERENCES user(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
   , UNIQUE(team_id, user_id)
);

CREATE TABLE IF NOT EXISTS channel_member (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, channel_id INTEGER NOT NULL
	, user_id INTEGER NOT NULL
	, FOREIGN KEY (channel_id) REFERENCES channel(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
	, FOREIGN KEY (user_id) REFERENCES user(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
   , UNIQUE(channel_id, user_id)
);

CREATE TABLE IF NOT EXISTS message (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, user_id INTEGER NOT NULL
	, channel_id INTEGER NOT NULL
	, text TEXT NOT NULL
	, timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	, FOREIGN KEY (user_id) REFERENCES user(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
	, FOREIGN KEY (channel_id) REFERENCES channel(id)
	    ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS validate_channel BEFORE INSERT ON channel
BEGIN
  SELECT CASE
  WHEN (
    (
      SELECT COUNT(*) FROM team_member 
      WHERE user_id = NEW.owner_id 
       AND team_id = NEW.team_id
	 ) ISNULL
  )
  THEN RAISE(ABORT, "owner_id not a team member")
  END;
END;

CREATE TRIGGER IF NOT EXISTS validate_channel_member BEFORE INSERT ON channel_member
BEGIN
  SELECT CASE
  WHEN (
    (
      SELECT COUNT(*)
		FROM team_member tm
		JOIN channel c ON c.team_id = tm.team_id
      WHERE tm.user_id = NEW.user_id 
	    AND c.id = NEW.channel_id
	 ) ISNULL
  )
  THEN RAISE(ABORT, "user_id not a team member")
  END;
END;

CREATE TRIGGER IF NOT EXISTS validate_message BEFORE INSERT ON message
BEGIN
  SELECT CASE
  WHEN (
    (
      SELECT COUNT(*)
		FROM team_member tm
		JOIN channel c ON c.team_id = tm.team_id
      WHERE tm.user_id = NEW.user_id 
	    AND c.id = NEW.channel_id
	 ) ISNULL
  )
  THEN RAISE(ABORT, "user_id not a team member")
  END;
END;

INSERT INTO user (id, username, email, password) VALUES
    (1, "steve", "steve@steve.com", "qwerty")
  , (2, "simon", "simon@steve.com", "qwerty");

INSERT INTO team (id, name, owner_id) VALUES
	 (1, "Team Steve", 1)
  , (2, "Team Simon", 2);

INSERT INTO channel (id, name, team_id, owner_id, is_public) VALUES
	 (1, "SteveTalk", 1, 1, false)
  , (2, "SimonTalk", 2, 2, true);

INSERT INTO team_member (user_id, team_id) VALUES
    (1, 1)
  , (1, 2)
  , (2, 1)
  , (2, 2);

INSERT INTO channel_member (user_id, channel_id) VALUES
    (1, 1)
  , (1, 2)
  , (2, 1)
  , (2, 2);
`

	if _, err := dbh.DB().Exec(schema); err != nil {
		return err
	}
	return nil
}
