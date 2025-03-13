package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Name     string
		Host     string
		Port     string
		User     string
		Password string
	}

	JWT struct {
		Secret     string
		Expiration time.Duration
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}

	App struct {
		Host  string
		Port  string
		Debug string
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func replaceEnvVars(input string) string {
	replacer := strings.NewReplacer(
		"${DB_NAME}", getEnvOrDefault("DB_NAME", ""),
		"${DB_HOST}", getEnvOrDefault("DB_HOST", "localhost"),
		"${DB_PORT}", getEnvOrDefault("DB_PORT", "5432"),
		"${DB_USER}", getEnvOrDefault("DB_USER", "root"),
		"${DB_PASSWORD}", getEnvOrDefault("DB_PASSWORD", "password"),

		"${JWT_SECRET}", getEnvOrDefault("JWT_SECRET", "jdfseiowjrwe234_please_change_me"),
		"${JWT_EXPIRATION}", getEnvOrDefault("JWT_EXPIRATION", "3600s"),

		"${REDIS_HOST}", getEnvOrDefault("REDIS_HOST", "localhost"),
		"${REDIS_PORT}", getEnvOrDefault("REDIS_PORT", "6379"),
		"${REDIS_PASSWORD}", getEnvOrDefault("REDIS_PASSWORD", "1000"),

		"${APP_HOST}", getEnvOrDefault("APP_HOST", "localhost"),
		"${APP_PORT}", getEnvOrDefault("APP_PORT", "8080"),
		"${DEBUG_APP}", getEnvOrDefault("DEBUG_APP", "false"),
	)

	return replacer.Replace(input)
}

func LoadConfig(path string) Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config

	config.Database.Name = replaceEnvVars(viper.GetString("database.name"))
	config.Database.Host = replaceEnvVars(viper.GetString("database.host"))
	config.Database.Port = replaceEnvVars(viper.GetString("database.port"))
	config.Database.User = replaceEnvVars(viper.GetString("database.user"))
	config.Database.Password = replaceEnvVars(viper.GetString("database.password"))

	config.JWT.Secret = replaceEnvVars(viper.GetString("jwt.secret"))
	config.JWT.Expiration, _ = time.ParseDuration(replaceEnvVars(viper.GetString("jwt.expiration")))

	config.Redis.Host = replaceEnvVars(viper.GetString("redis.host"))
	config.Redis.Port = replaceEnvVars(viper.GetString("redis.port"))
	config.Redis.Password = replaceEnvVars(viper.GetString("redis.password"))

	config.App.Host = replaceEnvVars(viper.GetString("app.host"))
	config.App.Port = replaceEnvVars(viper.GetString("app.port"))
	config.App.Debug = replaceEnvVars(viper.GetString("app.debug"))

	return config
}
