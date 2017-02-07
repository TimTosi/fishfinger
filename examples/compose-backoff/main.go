package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/timtosi/fishfinger"
)

// -----------------------------------------------------------------------------

// printRes is an helper function that prints the result of a query.
func printRes(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		for i, col := range values {
			fmt.Printf("%s : %s\n", columns[i], string(col))
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

func main() {
	// First of all, let's create a new Compose file handler by using the
	// `fishfinger.NewCompose` function that takes the path to your Compose file
	// as an argument.
	c, err := fishfinger.NewCompose("./docker-compose.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Now, time to start the MySQL service with the `Compose.StartBackoff`
	// function. This function takes two arguments: the Backoff function used
	// and a list of service name that will be started by this function call.
	// The Backoff function used here is the one provided by default by the
	// Fishfinger project but you are expected to provide another one that suits
	// your need.
	if err := c.StartBackoff(fishfinger.SocketBackoff, "datastore"); err != nil {
		log.Fatal(err)
	}

	// It's only keep trying to connect to a specific port exposed by the
	// container. The fact is the function will not find any remote listener
	// until all data is correctly loaded. In this way, you are assured
	// everything is ready to be processed by the rest of your program.
	addr, err := c.Port("datastore", "3306/tcp")
	if err != nil {
		log.Fatal(err)
	}
	user, err := c.Env("datastore", "MYSQL_USER")
	if err != nil {
		log.Fatal(err)
	}
	pass, err := c.Env("datastore", "MYSQL_PASS")
	if err != nil {
		log.Fatal(err)
	}
	proto, err := c.Env("datastore", "MYSQL_PROTO")
	if err != nil {
		log.Fatal(err)
	}
	dbName, err := c.Env("datastore", "MYSQL_DB_NAME")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@%s(%s:%s)/%s",
			user,
			pass,
			proto,
			strings.Split(addr, ":")[0],
			strings.Split(addr, ":")[1],
			dbName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM sinking_ships")
	if err != nil {
		log.Fatal(err)
	}
	if err := printRes(rows); err != nil {
		log.Fatal(err)
	}

	if err := c.Stop(); err != nil {
		log.Fatal(err)
	}
}
