package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_NAME", "test_db_name")
	os.Setenv("DB_HOST", "test_db_host")
	os.Setenv("DB_PORT", "test_db_port")
	os.Setenv("DB_USER", "test_db_user")
	os.Setenv("DB_PASSWORD", "test_db_password")

	os.Setenv("JWT_SECRET", "test_jwt_secret")
	os.Setenv("JWT_EXPIRATION", "5000s")

	os.Setenv("REDIS_HOST", "test_redis_host")
	os.Setenv("REDIS_PORT", "test_redis_port")
	os.Setenv("REDIS_PASSWORD", "test_redis_password")

	os.Setenv("APP_HOST", "test_app_host")
	os.Setenv("APP_PORT", "test_app_port")
	os.Setenv("APP_DEBUG", "false")

	config := LoadConfig("./")

	assert.Equal(t, "test_db_name", config.Database.Name)
	assert.Equal(t, "test_db_host", config.Database.Host)
	assert.Equal(t, "test_db_port", config.Database.Port)
	assert.Equal(t, "test_db_user", config.Database.User)
	assert.Equal(t, "test_db_password", config.Database.Password)

	assert.Equal(t, "test_jwt_secret", config.JWT.Secret)
	assert.Equal(t, 5000*time.Second, config.JWT.Expiration)

	assert.Equal(t, "test_redis_host", config.Redis.Host)
	assert.Equal(t, "test_redis_port", config.Redis.Port)
	assert.Equal(t, "test_redis_password", config.Redis.Password)

	assert.Equal(t, "test_app_host", config.App.Host)
	assert.Equal(t, "test_app_port", config.App.Port)
	assert.Equal(t, "false", config.App.Debug)
}
