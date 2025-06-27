package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig armazena configurações do servidor HTTP
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig armazena configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig armazena configurações para autenticação JWT
type JWTConfig struct {
	Secret          string
	ExpirationHours int
	RefreshSecret   string
	RefreshExpHours int
}

// LoadConfig carrega as configurações a partir de variáveis de ambiente
func LoadConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
	}
}

func loadServerConfig() ServerConfig {
	port := mustAtoi(getEnv("SERVER_PORT", "8080"), 8080)

	readTimeout := mustAtoi(getEnv("SERVER_READ_TIMEOUT", "5"), 5)
	writeTimeout := mustAtoi(getEnv("SERVER_WRITE_TIMEOUT", "10"), 10)
	idleTimeout := mustAtoi(getEnv("SERVER_IDLE_TIMEOUT", "120"), 120)

	return ServerConfig{
		Port:         port,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}
}

func loadDatabaseConfig() DatabaseConfig {
	dbPort := mustAtoi(getEnv("DB_PORT", "5432"), 5432)

	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "auth_system"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func loadJWTConfig() JWTConfig {
	expHours := mustAtoi(getEnv("JWT_EXPIRATION_HOURS", "24"), 24)
	refreshExpHours := mustAtoi(getEnv("JWT_REFRESH_EXPIRATION_HOURS", "168"), 168)

	return JWTConfig{
		Secret:          getEnv("JWT_SECRET", "your_jwt_secret"),
		ExpirationHours: expHours,
		RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your_refresh_secret"),
		RefreshExpHours: refreshExpHours,
	}
}

// mustAtoi tenta converter uma string para int, retornando o valor padrão em caso de erro
func mustAtoi(s string, defaultValue int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultValue
}

// GetDatabaseURL retorna a URL de conexão com o banco de dados
func (c *DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// getEnv recupera uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
