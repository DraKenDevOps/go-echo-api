package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Cwd        string
	EnvMode    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string

	JWTPrivateKey string
	JWTPublicKey  string

	EncryptionKey      string
	UploadLimitSize    int64
	ImageCompressLevel int
	RedisURI           string
	LogLevel           string
	Environment        string
	ServiceName        string
	Host               string
	BasePath           string
	TimeZone           string
	LimitMaxBalance    bool
	AppVersion         string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s: %+v", err.Error(), err)
	}
	cfg := &Config{
		Cwd:                cwd,
		EnvMode:            getEnv("ENV_MODE", "development"),
		ServerPort:         getEnv("PORT", "8080"),
		JWTPrivateKey:      getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKey:       getEnv("JWT_PUBLIC_KEY", ""),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", ""),
		RedisURI:           getEnv("REDIS_URI", ""),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		Environment:        getEnv("ENV_MODE", "development"),
		ServiceName:        getEnv("SERVICE_NAME", "go-echo-api"),
		Host:               getEnv("HOST", "0.0.0.0"),
		BasePath:           getEnv("BASE_PATH", "api"),
		TimeZone:           getEnv("TZ", "Asia/Bangkok"),
		UploadLimitSize:    int64(getEnvInt("UPLOAD_LIMIT_SIZE", 10)),
		ImageCompressLevel: getEnvInt("IMAGE_COMPRESS_LEVEL", 70),
		LimitMaxBalance:    getEnvBool("LIMIT_MAX_BALANCE", false),
		AppVersion:         getEnv("APP_VERSION", "1.0.0"),
	}

	dbURI := getEnv("DB_URI", "")
	if dbURI != "" {
		if parsed, err := url.Parse(dbURI); err == nil {
			if user := parsed.User; user != nil {
				cfg.DBUser = user.Username()
				cfg.DBPassword, _ = user.Password()
			}
			cfg.DBHost = parsed.Hostname()
			cfg.DBPort = parsed.Port()
			if cfg.DBPort == "" {
				cfg.DBPort = "3306"
			}
			cfg.DBName = strings.TrimPrefix(parsed.Path, "/")
		}
	}

	if cfg.DBHost == "" {
		cfg.DBHost = getEnv("DB_HOST", "localhost")
		cfg.DBPort = getEnv("DB_PORT", "3306")
		cfg.DBUser = getEnv("DB_USER", "root")
		cfg.DBPassword = getEnv("DB_PASSWORD", "")
		cfg.DBName = getEnv("DB_NAME", "echo_api")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		switch strings.ToLower(v) {
		case "true", "1", "yes":
			return true
		case "false", "0", "no":
			return false
		}
	}
	return fallback
}
