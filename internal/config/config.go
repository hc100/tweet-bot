package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Credentials struct {
	APIKey            string
	APISecret         string
	AccessToken       string
	AccessTokenSecret string
}

type Config struct {
	Credentials Credentials
	Location    *time.Location
}

func Load(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		if err := loadDotEnv(dotenvPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
	}

	cfg := Config{
		Credentials: Credentials{
			APIKey:            os.Getenv("X_API_KEY"),
			APISecret:         os.Getenv("X_API_SECRET"),
			AccessToken:       os.Getenv("X_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("X_ACCESS_TOKEN_SECRET"),
		},
		Location: loadLocation(),
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func loadLocation() *time.Location {
	name := os.Getenv("TZ")
	if name == "" {
		name = "Asia/Tokyo"
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.FixedZone("JST", 9*60*60)
	}

	return loc
}

func (c Config) validate() error {
	missing := make([]string, 0, 4)

	if c.Credentials.APIKey == "" {
		missing = append(missing, "X_API_KEY")
	}
	if c.Credentials.APISecret == "" {
		missing = append(missing, "X_API_SECRET")
	}
	if c.Credentials.AccessToken == "" {
		missing = append(missing, "X_ACCESS_TOKEN")
	}
	if c.Credentials.AccessTokenSecret == "" {
		missing = append(missing, "X_ACCESS_TOKEN_SECRET")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid .env line: %q", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			return fmt.Errorf("invalid .env line: %q", line)
		}

		current, exists := os.LookupEnv(key)
		if !exists || current == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set env %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read .env: %w", err)
	}

	return nil
}
