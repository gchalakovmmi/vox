package config

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Port			string
	InstanceName		string
	LogLevel		slog.Level

	AccessTokenSecret	string
	RefreshTokenSecret	string
	EmailLoginEncryptKey	string
	AccessTokenTTL		time.Duration
	RefreshTokenTTL		time.Duration

	RedisAddress		string
	RedisPassword		string

	PostgresUser		string
	PostgresPassword	string
	PostgresDB		string
	PostgresHost		string
}

func (c *Config) MustSet(key string, dest *string) {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("internal/config/config.go", "Message", "Required env var missing", "var", key)
		os.Exit(1)
	}
	*dest = val
}

func (c *Config) MustSetDuration(key string, dest *time.Duration) {
	str := os.Getenv(key)
	if str == "" {
		slog.Error("internal/config/config.go", "Message", "Required env var missing", "var", key)
		os.Exit(1)
	}
	d, err := time.ParseDuration(str)
	if err != nil {
		slog.Error("internal/config/config.go", "Message", "Invalid duration", "var", key, "error", err)
		os.Exit(1)
	}
	*dest = d
}

func New() *Config {
	c := &Config{}

	c.MustSet("PORT", &c.Port)
	c.MustSet("INSTANCE_NAME", &c.InstanceName)
	c.SetLogLevel()

	c.MustSet("ACCESS_TOKEN_SECRET", &c.AccessTokenSecret)
	c.MustSet("REFRESH_TOKEN_SECRET", &c.RefreshTokenSecret)
	c.MustSet("EMAIL_LOGIN_ENCRYPT_KEY", &c.EmailLoginEncryptKey)
	c.MustSetDuration("ACCESS_TOKEN_TTL", &c.AccessTokenTTL)
	c.MustSetDuration("REFRESH_TOKEN_TTL", &c.RefreshTokenTTL)

	c.MustSet("REDIS_ADDRESS", &c.RedisAddress)
	c.MustSet("REDIS_PASSWORD", &c.RedisPassword)

	c.MustSet("POSTGRES_USER", &c.PostgresUser)
	c.MustSet("POSTGRES_PASSWORD", &c.PostgresPassword)
	c.MustSet("POSTGRES_DB", &c.PostgresDB)
	c.MustSet("POSTGRES_HOST", &c.PostgresHost)

	return c
}

func (c *Config) SetLogLevel() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		slog.Error("internal/config/config.go", "Message", "Environment variable 'LogLevel' is not set! Exiting...")
		os.Exit(1)
	}
	switch level {
		case "DEBUG":
			c.LogLevel = slog.LevelDebug
		case "INFO":
			c.LogLevel = slog.LevelInfo
		case "WARN":
			c.LogLevel = slog.LevelWarn
		case "ERROR":
			c.LogLevel = slog.LevelError
		default:
			c.LogLevel = slog.LevelInfo
	}
}

func (c *Config) GetPort() string			{ return c.Port }
func (c *Config) GetInstanceName() string		{ return c.InstanceName }
func (c *Config) GetLogLevel() slog.Level		{ return c.LogLevel }
func (c *Config) GetAccessTokenSecret() string		{ return c.AccessTokenSecret }
func (c *Config) GetRefreshTokenSecret() string		{ return c.RefreshTokenSecret }
func (c *Config) GetEmailLoginEncryptKey() string	{ return c.EmailLoginEncryptKey }
func (c *Config) GetAccessTokenTTL() time.Duration	{ return c.AccessTokenTTL }
func (c *Config) GetRefreshTokenTTL() time.Duration	{ return c.RefreshTokenTTL }
func (c *Config) GetRedisAddress() string		{ return c.RedisAddress }
func (c *Config) GetRedisPassword() string		{ return c.RedisPassword }
func (c *Config) GetPostgresUser() string		{ return c.PostgresUser }
func (c *Config) GetPostgresPassword() string		{ return c.PostgresPassword }
func (c *Config) GetPostgresDB() string			{ return c.PostgresDB }
func (c *Config) GetPostgresHost() string		{ return c.PostgresHost }
