package main

import (
	"database/sql"
	"fmt"
	"os"

	"app/envloader"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	var env envloader.Env
	if err := envloader.Load(&env); err != nil {
		panic(err)
	}

	dbURL := fmt.Sprintf("libsql://%v-%v.aws-us-east-1.turso.io?authToken=%v",
		env.DB_NAME,
		env.DB_TURSO_USER,
		env.DB_TOKEN,
	)

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT)")

}
