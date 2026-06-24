package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	config   Config
	configMu sync.RWMutex
)

type Config struct {
	Environment string
	ServiceName string
	Port        int

	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	GRPCTLSCertPath string
	GRPCTLSKeyPath  string

	StorageProvider        string
	StorageBucket          string
	StorageRegion          string
	StorageEndpoint        string
	StoragePublicEndpoint  string
	StorageAccessKeyID     string
	StorageSecretAccessKey string
	StorageForcePathStyle  bool
	StorageUploadURLTTL    time.Duration
	StorageDownloadURLTTL  time.Duration
}

func LoadConfig(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		_ = godotenv.Load(dotenvPath)
	}

	cfg := Config{
		Environment:            getEnv("ENVIRONMENT", "development"),
		ServiceName:            getEnv("SERVICE_NAME", "neuraclinic-file-management"),
		Port:                   getEnvInt("PORT", 8000),
		DBHost:                 getEnv("DB_HOST", ""),
		DBPort:                 getEnvInt("DB_PORT", 5432),
		DBUser:                 getEnv("DB_USER", ""),
		DBPass:                 getEnv("DB_PASS", ""),
		DBName:                 getEnv("DB_NAME", ""),
		DBSSLMode:              getEnv("DB_SSLMODE", "disable"),
		GRPCTLSCertPath:        getEnv("GRPC_TLS_CERT_PATH", ""),
		GRPCTLSKeyPath:         getEnv("GRPC_TLS_KEY_PATH", ""),
		StorageProvider:        getEnv("STORAGE_PROVIDER", "s3"),
		StorageBucket:          getEnv("STORAGE_BUCKET", ""),
		StorageRegion:          getEnv("STORAGE_REGION", "us-east-1"),
		StorageEndpoint:        getEnv("STORAGE_ENDPOINT", ""),
		StoragePublicEndpoint:  getEnv("STORAGE_PUBLIC_ENDPOINT", ""),
		StorageAccessKeyID:     getEnv("STORAGE_ACCESS_KEY_ID", ""),
		StorageSecretAccessKey: getEnv("STORAGE_SECRET_ACCESS_KEY", ""),
		StorageForcePathStyle:  getEnvBool("STORAGE_FORCE_PATH_STYLE", true),
		StorageUploadURLTTL:    getEnvDuration("STORAGE_UPLOAD_URL_TTL", 15*time.Minute),
		StorageDownloadURLTTL:  getEnvDuration("STORAGE_DOWNLOAD_URL_TTL", 15*time.Minute),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	setConfig(cfg)
	return cfg, nil
}

func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

func setConfig(cfg Config) {
	configMu.Lock()
	config = cfg
	configMu.Unlock()
}

func (c Config) Validate() error {
	required := map[string]string{
		"DB_HOST":            c.DBHost,
		"DB_USER":            c.DBUser,
		"DB_PASS":            c.DBPass,
		"DB_NAME":            c.DBName,
		"GRPC_TLS_CERT_PATH": c.GRPCTLSCertPath,
		"GRPC_TLS_KEY_PATH":  c.GRPCTLSKeyPath,
		"STORAGE_BUCKET":     c.StorageBucket,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required config key: %s", key)
		}
	}

	if c.Port <= 0 {
		return fmt.Errorf("PORT must be greater than zero")
	}
	if c.DBPort <= 0 {
		return fmt.Errorf("DB_PORT must be greater than zero")
	}
	if c.StorageProvider != "s3" {
		return fmt.Errorf("STORAGE_PROVIDER must be s3")
	}
	if c.StorageRegion == "" {
		return fmt.Errorf("STORAGE_REGION is required")
	}
	if c.StorageUploadURLTTL <= 0 || c.StorageDownloadURLTTL <= 0 {
		return fmt.Errorf("storage URL TTLs must be greater than zero")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
