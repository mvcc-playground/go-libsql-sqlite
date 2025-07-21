package main

import (
	"fmt"

	"app/envloader"
)

type Env struct {
	DATABASE_URL  string `env:"DB_URL" validate:"omitempty,url"`
	DB_TOKEN      string
	DB_NAME       string
	DB_TURSO_USER string
}

func main() {
	var env Env
	if err := envloader.Load(&env); err != nil {
		panic(err)
	}

	dbURL := fmt.Sprintf("libsql://%v-%v.aws-us-east-1.turso.io?authToken=%v",
		env.DB_NAME,
		env.DB_TURSO_USER,
		env.DB_TOKEN,
	)

	fmt.Println("Database URL:", dbURL == env.DATABASE_URL)
	// fmt.Println("Database URL:", dbURL)
}
