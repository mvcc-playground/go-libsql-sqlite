package envloader

type Env struct {
	DATABASE_URL  string `env:"DB_URL" validate:"omitempty,url"`
	DB_TOKEN      string
	DB_NAME       string
	DB_TURSO_USER string
}
