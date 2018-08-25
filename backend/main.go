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
CREATE TABLE IF NOT EXISTS team (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, owner_id INTEGER NOT NULL
	, name TEXT NOT NULL
	, FOREIGN KEY (owner_id) REFERENCES user(id)
);

CREATE TABLE IF NOT EXISTS channel (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, team_id INTEGER NOT NULL
	, owner_id INTEGER NOT NULL
	, name TEXT NOT NULL
	, is_public BOOLEAN NOT NULL
	, FOREIGN KEY (team_id) REFERENCES team(id)
	, FOREIGN KEY (owner_id) REFERENCES user(id)
);

CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, username TEXT NOT NULL
	, email TEXT NOT NULL
	, password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS message (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, user_id INTEGER NOT NULL
	, channel_id INTEGER NOT NULL
	, text TEXT NOT NULL
	, timestamp DATETIME NOT NULL
	, FOREIGN KEY (user_id) REFERENCES user(id)
	, FOREIGN KEY (channel_id) REFERENCES channel(id)
);

CREATE TABLE IF NOT EXISTS team_member (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, team_id INTEGER NOT NULL
	, user_id INTEGER NOT NULL
	, FOREIGN KEY (team_id) REFERENCES team(id)
	, FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE IF NOT EXISTS channel_member (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	, channel_id INTEGER NOT NULL
	, user_id INTEGER NOT NULL
	, FOREIGN KEY (channel_id) REFERENCES channel(id)
	, FOREIGN KEY (user_id) REFERENCES user(id)
);
`

	if _, err := dbh.DB().Exec(schema); err != nil {
		return err
	}
	return nil
}
