package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   HTTPServer
	Database DatabaseConfig
}

type HTTPServer struct {
	Env             string
	LogLevel        string
	Port            string
	Timeout         time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load(path string) (Config, error) {
	if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env file: %w", err)
	}

	serverCfg, err := loadServerConfig()
	if err != nil {
		return Config{}, err
	}

	databaseCfg, err := loadDatabaseConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Server:   serverCfg,
		Database: databaseCfg,
	}, nil
}

func loadServerConfig() (HTTPServer, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "APP_ENV")
	}

	logLevel := os.Getenv("APP_LOG_LEVEL")
	if logLevel == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "APP_LOG_LEVEL")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "APP_PORT")
	}

	timeoutRaw := os.Getenv("HTTP_TIMEOUT")
	if timeoutRaw == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "HTTP_TIMEOUT")
	}

	timeout, err := time.ParseDuration(timeoutRaw)
	if err != nil {
		return HTTPServer{}, fmt.Errorf("parse %s as duration: %w", "HTTP_TIMEOUT", err)
	}

	idleTimeoutRaw := os.Getenv("HTTP_IDLE_TIMEOUT")
	if idleTimeoutRaw == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "HTTP_IDLE_TIMEOUT")
	}

	idleTimeout, err := time.ParseDuration(idleTimeoutRaw)
	if err != nil {
		return HTTPServer{}, fmt.Errorf("parse %s as duration: %w", "HTTP_IDLE_TIMEOUT", err)
	}

	shutdownTimeoutRaw := os.Getenv("HTTP_SHUTDOWN_TIMEOUT")
	if shutdownTimeoutRaw == "" {
		return HTTPServer{}, fmt.Errorf("env variable %q is not set", "HTTP_SHUTDOWN_TIMEOUT")
	}

	shutdownTimeout, err := time.ParseDuration(shutdownTimeoutRaw)
	if err != nil {
		return HTTPServer{}, fmt.Errorf("parse %s as duration: %w", "HTTP_SHUTDOWN_TIMEOUT", err)
	}

	return HTTPServer{
		Env:             env,
		LogLevel:        logLevel,
		Port:            port,
		Timeout:         timeout,
		IdleTimeout:     idleTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_HOST")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_PORT")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_USER")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_PASSWORD")
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_NAME")
	}

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		return DatabaseConfig{}, fmt.Errorf("env variable %q is not set", "DB_SSLMODE")
	}

	cfg := DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
	}

	return cfg, nil
}

func (c DatabaseConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%s", c.Host, c.Port),
		Path:   c.Name,
	}

	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
