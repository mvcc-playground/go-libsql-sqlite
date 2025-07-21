package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql" // Mude para este (puro Go)
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID    int
	Title string
}

func main() {
	dbName := "file:./local-3.db"

	conn, err := sql.Open("libsql", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        dbName,
		Conn:       conn,
	}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gorm filed %s", err)
		os.Exit(1)
	}
	db.AutoMigrate(&Task{})
	taskId := 1000 + (int)(os.Getpid()) + (int)(os.Getppid()) + (int)(os.Getuid())
	task := Task{ID: taskId, Title: "My first task"}
	if err := db.Create(&task).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to create task %s", err)
		os.Exit(1)
	}
}
