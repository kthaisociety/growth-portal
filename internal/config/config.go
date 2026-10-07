// Package config reads every setting from the environment, once, at boot.
package config

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Config struct {
	// Host is the bind address. Empty binds all interfaces, which the container needs. Set it to
	// 127.0.0.1 for local dev: a dev machine with a public IP exposes anything bound to 0.0.0.0.
	Host string
	Port string

	// DatabaseURL is a Postgres connection string, e.g.
	// postgres://growth:growth@127.0.0.1:5435/growth?sslmode=disable.
	DatabaseURL string

	// LumaServiceURL is luma-service's base URL, no trailing slash, e.g.
	// http://<luma-service app name>:8000 on dokploy-network.
	LumaServiceURL string
	// LumaServiceSecret is sent as X-Service-Secret; luma-service's LUMA_SERVICE_SECRET.
	LumaServiceSecret string
}

func Load() *Config {
	return &Config{
		Host:              getEnv("HOST", ""),
		Port:              getEnv("PORT", "8000"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		LumaServiceURL:    strings.TrimSuffix(getEnv("LUMA_SERVICE_URL", ""), "/"),
		LumaServiceSecret: getEnv("LUMA_SERVICE_SECRET", ""),
	}
}

// Validate fails if a required setting is missing, naming every missing variable at once.
func (c *Config) Validate() error {
	var missing []string
	for name, value := range map[string]string{
		"DATABASE_URL":        c.DatabaseURL,
		"LUMA_SERVICE_URL":    c.LumaServiceURL,
		"LUMA_SERVICE_SECRET": c.LumaServiceSecret,
	} {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
