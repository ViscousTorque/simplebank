package db

/*

First attempt to follow the Udemy video has an error:

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())

}

cannot use conn (variable of type *sql.DB) as DBTX value in argument to New: *sql.DB does not implement DBTX (wrong type for method Exec)
		have Exec(string, ...any) (sql.Result, error)
		want Exec(context.Context, string, ...interface{})

Wrap *sql.DB with pgxpool.Conn or similar

This error occurs because the New function you're calling expects a type that implements the DBTX interface,
which has methods like Exec that use context.Context as their first argument (e.g., Exec(context.Context, string, ...interface{})),
while *sql.DB does not implement this interface directly.

Here's how to resolve the issue:
Solution 1: Wrap *sql.DB with pgxpool.Conn or similar

If your New function is part of a library (e.g., using pgx with PostgreSQL), you might need to use a connection pool or a wrapper like pgx.Conn. For example:

Explanation

The core issue is that the New function is designed to work with a type that satisfies the DBTX interface,
and *sql.DB does not directly satisfy this interface because its methods (Exec, Query, etc.) do not take
ontext.Context as their first argument. Using a compatible library or wrapping *sql.DB is necessary to resolve this mismatch.
*/

import (
	"context"
	"log"
	"os"
	"testing"

	"main/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
