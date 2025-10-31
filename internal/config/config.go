package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	// Database
	DatabaseURL string
	
	// Redis
	RedisURL string
	
	// S3/MinIO
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
	
	// AI APIs
	OpenAIKey    string
	AnthropicKey string
	
	// Server
	APIPort          int
	WorkerConcurrency int
	
	// Security
	JWTSecret      string
	AllowedOrigins []string
	
	// Processing
	MaxFileSizeMB int
	TempDir       string
}

// Load carrega as configurações das variáveis de ambiente
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgresql://typecraft:dev_password@localhost:5433/typecraft_db"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/0"),
		S3Endpoint:        getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3AccessKey:       getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey:       getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:          getEnv("S3_BUCKET", "typecraft-files"),
		S3Region:          getEnv("S3_REGION", "us-east-1"),
		OpenAIKey:         getEnv("OPENAI_API_KEY", ""),
		AnthropicKey:      getEnv("ANTHROPIC_API_KEY", ""),
		APIPort:           getEnvInt("API_PORT", 8000),
		WorkerConcurrency: getEnvInt("WORKER_CONCURRENCY", 5),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-in-production"),
		AllowedOrigins:    []string{
			getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		},
		MaxFileSizeMB:     getEnvInt("MAX_FILE_SIZE_MB", 100),
		TempDir:           getEnv("TEMP_DIR", "/tmp/typecraft"),
	}
	
	// Validar configurações críticas
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL é obrigatório")
	}
	
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL é obrigatório")
	}
	
	return cfg, nil
}

// getEnv retorna o valor da variável de ambiente ou o padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retorna o valor inteiro da variável de ambiente ou o padrão
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
