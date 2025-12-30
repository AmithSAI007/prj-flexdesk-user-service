package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv   string
	HttpPort string

	DBUser                 string
	DBPassword             string
	DBHost                 string
	DBPort                 string
	DBName                 string
	CloudSqlConnectionName string
	DBSSLMode              string

	DBMaxConns        int
	DBMinConns        int
	DBMaxConnLifetime time.Duration
	DBMaxConnIdleTime time.Duration

	PrivateKeyPath string
	PublicKeyPath  string

	TokenIssuer string
}

func (c *Config) DBSource() string {
	if c.CloudSqlConnectionName != "" {

		return fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s",
			c.DBUser, c.DBPassword, c.DBName, c.CloudSqlConnectionName)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func LoadConfig() *Config {

	dbMaxConns, err := strconv.Atoi(os.Getenv("DB_MAX_CONNS"))
	if err != nil {
		dbMaxConns = 10 // default value if not set or invalid
	}
	dbMinConns, err := strconv.Atoi(os.Getenv("DB_MIN_CONNS"))
	if err != nil {
		dbMinConns = 2 // default value if not set or invalid
	}

	dbMaxConnLifetimeMinutes, err := strconv.Atoi(os.Getenv("DB_MAX_CONN_LIFETIME_MINUTES"))
	if err != nil {
		dbMaxConnLifetimeMinutes = 30 // default value if not set or invalid
	}

	dbMaxConnIdleTimeMinutes, err := strconv.Atoi(os.Getenv("DB_MAX_CONN_IDLE_TIME_MINUTES"))
	if err != nil {
		dbMaxConnIdleTimeMinutes = 5 // default value if not set or invalid
	}

	return &Config{
		AppEnv:   os.Getenv("APP_ENV"),
		HttpPort: os.Getenv("HTTP_PORT"),

		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		DBName:                 os.Getenv("DB_NAME"),
		CloudSqlConnectionName: os.Getenv("CLOUD_SQL_CONNECTION_NAME"),
		DBSSLMode:              os.Getenv("DB_SSL_MODE"),

		DBMaxConns:        dbMaxConns,
		DBMinConns:        dbMinConns,
		DBMaxConnLifetime: time.Duration(dbMaxConnLifetimeMinutes) * time.Minute,
		DBMaxConnIdleTime: time.Duration(dbMaxConnIdleTimeMinutes) * time.Minute,

		PrivateKeyPath: os.Getenv("PRIVATE_KEY_PATH"),
		PublicKeyPath:  os.Getenv("PUBLIC_KEY_PATH"),

		TokenIssuer: os.Getenv("TOKEN_ISSUER"),
	}
}
