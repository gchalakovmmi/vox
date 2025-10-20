package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port		string
	InstanceName	string
	LogLevel	slog.Level
}

func New() *Config {
	c := &Config{}
	c.SetPort()
	c.SetInstanceName()
	c.SetLogLevel()
	return c
}

func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) GetInstanceName() string {
	return c.InstanceName
}

func (c *Config) GetLogLevel() slog.Level {
	return c.LogLevel
}

func (c *Config) SetPort() {
	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		slog.Error("Environment variable 'PORT' is not set! Exiting...")
		os.Exit(1)
	}
}

func (c *Config) SetInstanceName() {
	c.InstanceName = os.Getenv("INSTANCE_NAME")
	if c.InstanceName == "" {
		slog.Error("Environment variable 'InstanceName' is not set! Exiting...")
		os.Exit(1)
	}
}

func (c *Config) SetLogLevel() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		slog.Error("Environment variable 'LogLevel' is not set! Exiting...")
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
