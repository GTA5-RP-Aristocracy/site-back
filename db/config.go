package db

type (
	// Config represents the configuration options for the database.
	Config struct {
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASS"`
		Host     string `env:"DB_HOST"`
		Database string `env:"DB_NAME"`
	}
)
