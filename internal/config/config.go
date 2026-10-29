package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv   string `env:"APP_ENV" env-default:"development"`
	AppName  string `env:"APP_NAME" env-default:"Go Starter Kit"`
	Port     int    `env:"PORT" env-default:"3000"`
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`

	// Database Master Relational
	DBDefaultConnection string `env:"DB_DEFAULT_CONNECTION" env-default:"sqlite1"`

	// SQLite Connections
	DBSqliteEnabled bool   `env:"DB_SQLITE_ENABLED" env-default:"true"`
	DBSqlite1Path   string `env:"DB_SQLITE_1_PATH" env-default:"storages/database/sqlite1.dev.db"`
	DBSqlite2Path   string `env:"DB_SQLITE_2_PATH" env-default:"storages/database/sqlite2.dev.db"`
	DBSqlite3Path   string `env:"DB_SQLITE_3_PATH" env-default:"storages/database/sqlite3.dev.db"`

	// PostgreSQL Connections
	DBPGEnabled bool   `env:"DB_PG_ENABLED" env-default:"false"`
	DBPG1DSN    string `env:"DB_PG_1_DSN" env-default:""`
	DBPG2DSN    string `env:"DB_PG_2_DSN" env-default:""`
	DBPG3DSN    string `env:"DB_PG_3_DSN" env-default:""`

	// MySQL Connections
	DBMySQLEnabled bool   `env:"DB_MYSQL_ENABLED" env-default:"false"`
	DBMySQL1DSN    string `env:"DB_MYSQL_1_DSN" env-default:""`
	DBMySQL2DSN    string `env:"DB_MYSQL_2_DSN" env-default:""`
	DBMySQL3DSN    string `env:"DB_MYSQL_3_DSN" env-default:""`

	// MongoDB Connections
	DBMongoEnabled bool   `env:"DB_MONGO_ENABLED" env-default:"false"`
	DBMongo1URI    string `env:"DB_MONGO_1_URI" env-default:""`
	DBMongo2URI    string `env:"DB_MONGO_2_URI" env-default:""`
	DBMongo3URI    string `env:"DB_MONGO_3_URI" env-default:""`

	// Redis
	RedisURL string `env:"REDIS_URL" env-default:"redis://localhost:6379"`

	// Mail System
	MailDriver    string `env:"MAIL_DRIVER" env-default:"log"` // log, smtp, resend
	MailFromName  string `env:"MAIL_FROM_NAME" env-default:"Go Starter Kit"`
	MailFromEmail string `env:"MAIL_FROM_EMAIL" env-default:"noreply@example.com"`
	MailSMTPHost  string `env:"MAIL_SMTP_HOST" env-default:"localhost"`
	MailSMTPPort  int    `env:"MAIL_SMTP_PORT" env-default:"1025"`
	MailSMTPUser  string `env:"MAIL_SMTP_USER" env-default:""`
	MailSMTPPass  string `env:"MAIL_SMTP_PASS" env-default:""`
	MailResendKey string `env:"MAIL_RESEND_API_KEY" env-default:""`
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsTesting() bool {
	return c.AppEnv == "testing"
}

// Load membaca konfigurasi dari file path (default ".env") dan mengizinkan override dari env var OS.
func Load(envFiles ...string) (*Config, error) {
	var cfg Config
	targetEnvFile := ".env"
	if len(envFiles) > 0 && envFiles[0] != "" {
		targetEnvFile = envFiles[0]
	}

	if _, err := os.Stat(targetEnvFile); err == nil {
		if err := cleanenv.ReadConfig(targetEnvFile, &cfg); err != nil {
			return nil, fmt.Errorf("failed to read env config file %s: %w", targetEnvFile, err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("failed to read env variables: %w", err)
		}
	} else {
		return nil, fmt.Errorf("error accessing env config file %s: %w", targetEnvFile, err)
	}

	return &cfg, nil
}
