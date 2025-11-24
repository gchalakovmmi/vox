package config

import (
	"log/slog"
	"os"
	"time"
	"net/http"
	"strings"
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

	CookieSecure		bool
	CookieSameSite		http.SameSite

	TTSOpenAIURL		string
	TTSOpenAIModel		string
	TTSOpenAIVoice		string
	TTSOpenAIKey		string
}

func (c *Config) Set(key string, dest *string) {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("internal/config/config.go", "Message", "Required env var missing", "var", key)
		os.Exit(1)
	}
	*dest = val
}

func (c *Config) SetDuration(key string, dest *time.Duration) {
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

func (c *Config) SetBool(key string, dest *bool) {
	val := strings.TrimSpace(os.Getenv(key))
	switch strings.ToLower(val) {
	case "true":
		*dest = true
	case "false":
		*dest = false
	default:
		slog.Error("internal/config/config.go", "Message", "env var must be set to either true or false", "var", key)
		os.Exit(1)
	}
}
func (c *Config) SetCookieSameSite() {
	val := strings.TrimSpace(strings.ToLower(os.Getenv("COOKIE_SAME_SITE")))
	switch val {
	case "lax":
		c.CookieSameSite = http.SameSiteLaxMode
	case "strict":
		c.CookieSameSite = http.SameSiteStrictMode
	case "none":
		c.CookieSameSite = http.SameSiteNoneMode
	default:
		slog.Error("internal/config/config.go", "Message", "COOKIE_SAME_SITE must be lax, strict or none")
		os.Exit(1)
	}
}

func (c *Config) GetPort()			string		{ return c.Port }
func (c *Config) GetInstanceName()		string		{ return c.InstanceName }
func (c *Config) GetLogLevel()			slog.Level	{ return c.LogLevel }
func (c *Config) GetAccessTokenSecret()		string		{ return c.AccessTokenSecret }
func (c *Config) GetRefreshTokenSecret()	string		{ return c.RefreshTokenSecret }
func (c *Config) GetEmailLoginEncryptKey()	string		{ return c.EmailLoginEncryptKey }
func (c *Config) GetAccessTokenTTL()		time.Duration	{ return c.AccessTokenTTL }
func (c *Config) GetRefreshTokenTTL()		time.Duration	{ return c.RefreshTokenTTL }
func (c *Config) GetRedisAddress()		string		{ return c.RedisAddress }
func (c *Config) GetRedisPassword()		string		{ return c.RedisPassword }
func (c *Config) GetPostgresUser()		string		{ return c.PostgresUser }
func (c *Config) GetPostgresPassword()		string		{ return c.PostgresPassword }
func (c *Config) GetPostgresDB()		string		{ return c.PostgresDB }
func (c *Config) GetPostgresHost()		string		{ return c.PostgresHost }
func (c *Config) GetTTSOpenAIURL()		string   	{ return c.TTSOpenAIURL }
func (c *Config) GetTTSOpenAIModel()		string		{ return c.TTSOpenAIModel }
func (c *Config) GetTTSOpenAIVoice()		string		{ return c.TTSOpenAIVoice }
func (c *Config) GetTTSOpenAIKey()		string		{ return c.TTSOpenAIKey }

func New() *Config {
	c := &Config{}

	c.Set("PORT", &c.Port)
	c.Set("INSTANCE_NAME", &c.InstanceName)
	c.SetLogLevel()

	c.Set("ACCESS_TOKEN_SECRET", &c.AccessTokenSecret)
	c.Set("REFRESH_TOKEN_SECRET", &c.RefreshTokenSecret)
	c.Set("EMAIL_LOGIN_ENCRYPT_KEY", &c.EmailLoginEncryptKey)
	c.SetDuration("ACCESS_TOKEN_TTL", &c.AccessTokenTTL)
	c.SetDuration("REFRESH_TOKEN_TTL", &c.RefreshTokenTTL)

	c.Set("REDIS_ADDRESS", &c.RedisAddress)
	c.Set("REDIS_PASSWORD", &c.RedisPassword)

	c.Set("POSTGRES_USER", &c.PostgresUser)
	c.Set("POSTGRES_PASSWORD", &c.PostgresPassword)
	c.Set("POSTGRES_DB", &c.PostgresDB)
	c.Set("POSTGRES_HOST", &c.PostgresHost)

	c.SetCookieSameSite()
	c.SetBool("COOKIE_SECURE", &c.CookieSecure)

	c.Set("TTS_OPENAI_URL", &c.TTSOpenAIURL)
	c.Set("TTS_OPENAI_MODEL_NAME", &c.TTSOpenAIModel)
	c.Set("TTS_OPENAI_VOICE", &c.TTSOpenAIVoice)
	c.Set("TTS_OPENAI_API_KEY", &c.TTSOpenAIKey)

	return c
}
