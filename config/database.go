package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var ConnDB *pgx.Conn

func ConnectDB() {
	
	DBurl := "postgres://postgres:132435@localhost:5432/DBproject"

	var err error

	ConnDB, err = pgx.Connect(context.Background(), DBurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v", err)
		os.Exit(1)
	}
	fmt.Println("Success connect to database")
}

// postgres://{user}:{password}@{host}:{port}/{database}
// user = user postgres
// password = password postgres
// host = host postgres
// port = port postgres
// database = database postgres
