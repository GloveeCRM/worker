package types

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
}

type EmailConfig struct {
	ResendAPIKey string
}

type Config struct {
	Postgres PostgresConfig
	Email    EmailConfig
}
