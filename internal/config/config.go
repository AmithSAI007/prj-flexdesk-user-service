package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	HttpPort string `mapstructure:"HTTP_PORT"`

	DBUser                 string `mapstructure:"DB_USER"`
	DBPassword             string `mapstructure:"DB_PASSWORD"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBName                 string `mapstructure:"DB_NAME"`
	CloudSqlConnectionName string `mapstructure:"CLOUD_SQL_CONNECTION_NAME"`
	DBSSLMode              string `mapstructure:"DB_SSL_MODE"`

	DBMaxConns        int           `mapstructure:"DB_MAX_CONNS"`
	DBMinConns        int           `mapstructure:"DB_MIN_CONNS"`
	DBMaxConnLifetime time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME"`
	DBMaxConnIdleTime time.Duration `mapstructure:"DB_MAX_CONN_IDLE_TIME"`

	PrivateKeyPath string `mapstructure:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `mapstructure:"PUBLIC_KEY_PATH"`

	TokenIssuer string `mapstructure:"TOKEN_ISSUER"`
}

func (c *Config) DBSource() string {
	if c.CloudSqlConnectionName != "" {

		return fmt.Sprintf("user=%s password=%s dbname=%s host=/cloudsql/%s",
			c.DBUser, c.DBPassword, c.DBName, c.CloudSqlConnectionName)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func LoadConfig(path string) (*Config, error) {

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_CONNS", 10)
	viper.SetDefault("DB_MIN_CONNS", 2)
	viper.SetDefault("DB_MAX_CONN_LIFETIME_MINUTES", 30*time.Minute)
	viper.SetDefault("DB_MAX_CONN_IDLE_TIME_MINUTES", 5*time.Minute)
	viper.SetDefault("TOKEN_ISSUER", "rate-my-setup")

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return &config, nil
}
