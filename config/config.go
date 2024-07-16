package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBTimezone    string
	DBMaxIdle     int
	DBMaxConn     int
	DBMaxLifeTime int

	Debug     bool
	PrettyLog bool

	TracerEndpoint string

	RedisURL string

	S3Endpoint string
	S3Bucket   string
	S3Region   string
	S3Key      string
	S3Secret   string
}

func NewConfig() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "fliqt_test"),
		DBTimezone:    getEnv("DB_TIMEZONE", "Etc/UTC"),
		DBMaxIdle:     getEnvInt("DB_MAX_IDLE", 2),
		DBMaxConn:     getEnvInt("DB_MAX_CONN", 10),
		DBMaxLifeTime: getEnvInt("DB_MAX_LIFE", int(time.Minute*60)),

		Debug:     getEnv("DEBUG", "false") == "true",
		PrettyLog: getEnv("PRETTY_LOG", "false") == "true",

		TracerEndpoint: getEnv("TRACER_ENDPOINT", "localhost:4317"),

		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

		S3Endpoint: getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Bucket:   getEnv("S3_BUCKET", "fliqt"),
		S3Region:   getEnv("S3_REGION", "us-east-1"),
		S3Key:      getEnv("S3_KEY", ""),
		S3Secret:   getEnv("S3_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	envStr := getEnv(key, strconv.Itoa(defaultValue))

	env, err := strconv.ParseInt(envStr, 10, 64)
	if err == nil {
		return defaultValue
	}

	return int(env)
}

func (c *Config) GetDBDSN() string {
	// MySQL DSN
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, url.QueryEscape(c.DBTimezone))
}
