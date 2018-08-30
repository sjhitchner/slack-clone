package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sjhitchner/slack-clone/backend/domain"
	libdb "github.com/sjhitchner/slack-clone/backend/infrastructure/db"
	"github.com/sjhitchner/slack-clone/backend/interfaces/db"
	"github.com/sjhitchner/slack-clone/backend/interfaces/graphql"
	"github.com/sjhitchner/slack-clone/backend/interfaces/resolvers"
	//libsql "github.com/sjhitchner/slack-clone/backend/infrastructure/db/psql"
	libsql "github.com/sjhitchner/slack-clone/backend/infrastructure/db/sqlite"
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

	handler, err := SetupHandler(dbh, string(schema))
	CheckError(err)

	http.Handle("/graphql", handler)
	http.HandleFunc("/cookie", func(w http.ResponseWriter, r *http.Request) {
		const Session = "session"

		cookie, err := r.Cookie(Session)
		if cookie == nil {
			log.Println("No Cookie Found", err)
			cookie = &http.Cookie{
				Name:     Session,
				Value:    "0",
				Path:     "/",
				Domain:   "localhost",
				Expires:  time.Now().Add(time.Minute), //time.Now().Add(7 * 24 * time.Hour),
				MaxAge:   60,                          //24 * 60 * 60 * 7, // 7 days
				Secure:   false,
				HttpOnly: true,
			}
		}
		count, _ := strconv.Atoi(cookie.Value)
		count++
		log.Println("New Cookie Value", count)
		cookie.Value = fmt.Sprintf("%d", count)
		http.SetCookie(w, cookie)

		w.Write([]byte("Ok"))

	})
	http.Handle("/", graphql.NewGraphiQLHandler(string(schema)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupHandler(dbh libdb.DBHandler, schema string) (http.Handler, error) {

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

	return graphql.WrapContext(ctx, handler), nil
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
	, username TEXT NOT NULL
	, email TEXT NOT NULL
	, password TEXT NOT NULL
	, UNIQUE(username)
	, UNIQUE(email)
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
  , (2, "simon", "simon@steve.com", "qwerty")
  , (3, "matt", "matt@steve.com", "qwerty");

INSERT INTO team (id, name, owner_id) VALUES
    (1, "Team Steve", 1)
  , (2, "Team Simon", 2)
  , (3, "Team Matt", 3);

INSERT INTO channel (id, name, team_id, owner_id, is_public) VALUES
    (1, "SteveTalk", 1, 1, false)
  , (2, "SimonTalk", 2, 2, true)
  , (3, "MattTalk", 3, 3, true);

INSERT INTO team_member (team_id, user_id) VALUES
    (1, 1)
  , (1, 2)
  , (2, 2)
  , (2, 3)
  , (3, 1)
  , (3, 3);

INSERT INTO channel_member (channel_id, user_id) VALUES
    (1, 1)
  , (1, 2)
  , (2, 2)
  , (2, 3)
  , (3, 1)
  , (3, 3);
`

	if _, err := dbh.DB().Exec(schema); err != nil {
		return err
	}
	return nil
}
