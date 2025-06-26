package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Limpa variáveis de ambiente que podem interferir
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("JWT_SECRET")

	config := LoadConfig()

	if config == nil {
		t.Fatal("LoadConfig não deveria retornar nil")
	}

	// Verifica se as configurações foram carregadas com valores padrão
	if config.Server.Port != 8080 {
		t.Errorf("Porta do servidor deveria ser 8080, mas foi %d", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Host do banco deveria ser localhost, mas foi %s", config.Database.Host)
	}

	if config.JWT.Secret != "your_jwt_secret" {
		t.Errorf("JWT Secret deveria ser 'your_jwt_secret', mas foi %s", config.JWT.Secret)
	}
}

func TestLoadServerConfig(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectedPort  int
		expectedRead  time.Duration
		expectedWrite time.Duration
		expectedIdle  time.Duration
	}{
		{
			name:          "valores padrão",
			envVars:       map[string]string{},
			expectedPort:  8080,
			expectedRead:  5 * time.Second,
			expectedWrite: 10 * time.Second,
			expectedIdle:  120 * time.Second,
		},
		{
			name: "valores customizados",
			envVars: map[string]string{
				"SERVER_PORT":          "3000",
				"SERVER_READ_TIMEOUT":  "10",
				"SERVER_WRITE_TIMEOUT": "20",
				"SERVER_IDLE_TIMEOUT":  "300",
			},
			expectedPort:  3000,
			expectedRead:  10 * time.Second,
			expectedWrite: 20 * time.Second,
			expectedIdle:  300 * time.Second,
		},
		{
			name: "valores inválidos devem usar padrão",
			envVars: map[string]string{
				"SERVER_PORT":          "invalid",
				"SERVER_READ_TIMEOUT":  "invalid",
				"SERVER_WRITE_TIMEOUT": "invalid",
				"SERVER_IDLE_TIMEOUT":  "invalid",
			},
			expectedPort:  8080,
			expectedRead:  5 * time.Second,
			expectedWrite: 10 * time.Second,
			expectedIdle:  120 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpa variáveis de ambiente
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("SERVER_READ_TIMEOUT")
			os.Unsetenv("SERVER_WRITE_TIMEOUT")
			os.Unsetenv("SERVER_IDLE_TIMEOUT")

			// Define variáveis de ambiente para o teste
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := loadServerConfig()

			if config.Port != tt.expectedPort {
				t.Errorf("Porta esperada %d, mas foi %d", tt.expectedPort, config.Port)
			}

			if config.ReadTimeout != tt.expectedRead {
				t.Errorf("ReadTimeout esperado %v, mas foi %v", tt.expectedRead, config.ReadTimeout)
			}

			if config.WriteTimeout != tt.expectedWrite {
				t.Errorf("WriteTimeout esperado %v, mas foi %v", tt.expectedWrite, config.WriteTimeout)
			}

			if config.IdleTimeout != tt.expectedIdle {
				t.Errorf("IdleTimeout esperado %v, mas foi %v", tt.expectedIdle, config.IdleTimeout)
			}
		})
	}
}

func TestLoadDatabaseConfig(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedHost string
		expectedPort int
		expectedUser string
		expectedPass string
		expectedName string
		expectedSSL  string
	}{
		{
			name:         "valores padrão",
			envVars:      map[string]string{},
			expectedHost: "localhost",
			expectedPort: 5432,
			expectedUser: "postgres",
			expectedPass: "postgres",
			expectedName: "auth_system",
			expectedSSL:  "disable",
		},
		{
			name: "valores customizados",
			envVars: map[string]string{
				"DB_HOST":     "myhost.com",
				"DB_PORT":     "5433",
				"DB_USER":     "myuser",
				"DB_PASSWORD": "mypass",
				"DB_NAME":     "mydb",
				"DB_SSLMODE":  "require",
			},
			expectedHost: "myhost.com",
			expectedPort: 5433,
			expectedUser: "myuser",
			expectedPass: "mypass",
			expectedName: "mydb",
			expectedSSL:  "require",
		},
		{
			name: "porta inválida deve usar padrão",
			envVars: map[string]string{
				"DB_PORT": "invalid",
			},
			expectedHost: "localhost",
			expectedPort: 5432,
			expectedUser: "postgres",
			expectedPass: "postgres",
			expectedName: "auth_system",
			expectedSSL:  "disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpa variáveis de ambiente
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_SSLMODE")

			// Define variáveis de ambiente para o teste
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := loadDatabaseConfig()

			if config.Host != tt.expectedHost {
				t.Errorf("Host esperado %s, mas foi %s", tt.expectedHost, config.Host)
			}

			if config.Port != tt.expectedPort {
				t.Errorf("Porta esperada %d, mas foi %d", tt.expectedPort, config.Port)
			}

			if config.User != tt.expectedUser {
				t.Errorf("Usuário esperado %s, mas foi %s", tt.expectedUser, config.User)
			}

			if config.Password != tt.expectedPass {
				t.Errorf("Senha esperada %s, mas foi %s", tt.expectedPass, config.Password)
			}

			if config.Name != tt.expectedName {
				t.Errorf("Nome do banco esperado %s, mas foi %s", tt.expectedName, config.Name)
			}

			if config.SSLMode != tt.expectedSSL {
				t.Errorf("SSL Mode esperado %s, mas foi %s", tt.expectedSSL, config.SSLMode)
			}
		})
	}
}

func TestLoadJWTConfig(t *testing.T) {
	tests := []struct {
		name               string
		envVars            map[string]string
		expectedSecret     string
		expectedExpHours   int
		expectedRefreshSec string
		expectedRefreshExp int
	}{
		{
			name:               "valores padrão",
			envVars:            map[string]string{},
			expectedSecret:     "your_jwt_secret",
			expectedExpHours:   24,
			expectedRefreshSec: "your_refresh_secret",
			expectedRefreshExp: 168,
		},
		{
			name: "valores customizados",
			envVars: map[string]string{
				"JWT_SECRET":                   "my_secret_key",
				"JWT_EXPIRATION_HOURS":         "48",
				"JWT_REFRESH_SECRET":           "my_refresh_key",
				"JWT_REFRESH_EXPIRATION_HOURS": "336",
			},
			expectedSecret:     "my_secret_key",
			expectedExpHours:   48,
			expectedRefreshSec: "my_refresh_key",
			expectedRefreshExp: 336,
		},
		{
			name: "valores inválidos devem usar padrão",
			envVars: map[string]string{
				"JWT_EXPIRATION_HOURS":         "invalid",
				"JWT_REFRESH_EXPIRATION_HOURS": "invalid",
			},
			expectedSecret:     "your_jwt_secret",
			expectedExpHours:   24,
			expectedRefreshSec: "your_refresh_secret",
			expectedRefreshExp: 168,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpa variáveis de ambiente
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("JWT_EXPIRATION_HOURS")
			os.Unsetenv("JWT_REFRESH_SECRET")
			os.Unsetenv("JWT_REFRESH_EXPIRATION_HOURS")

			// Define variáveis de ambiente para o teste
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := loadJWTConfig()

			if config.Secret != tt.expectedSecret {
				t.Errorf("Secret esperado %s, mas foi %s", tt.expectedSecret, config.Secret)
			}

			if config.ExpirationHours != tt.expectedExpHours {
				t.Errorf("ExpirationHours esperado %d, mas foi %d", tt.expectedExpHours, config.ExpirationHours)
			}

			if config.RefreshSecret != tt.expectedRefreshSec {
				t.Errorf("RefreshSecret esperado %s, mas foi %s", tt.expectedRefreshSec, config.RefreshSecret)
			}

			if config.RefreshExpHours != tt.expectedRefreshExp {
				t.Errorf("RefreshExpHours esperado %d, mas foi %d", tt.expectedRefreshExp, config.RefreshExpHours)
			}
		})
	}
}

func TestGetDatabaseURL(t *testing.T) {
	config := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "mypassword",
		Name:     "auth_system",
		SSLMode:  "disable",
	}

	expectedURL := "postgresql://postgres:mypassword@localhost:5432/auth_system?sslmode=disable"
	actualURL := config.GetDatabaseURL()

	if actualURL != expectedURL {
		t.Errorf("URL esperada %s, mas foi %s", expectedURL, actualURL)
	}

	// Teste com SSL require
	config.SSLMode = "require"
	expectedURLWithSSL := "postgresql://postgres:mypassword@localhost:5432/auth_system?sslmode=require"
	actualURLWithSSL := config.GetDatabaseURL()

	if actualURLWithSSL != expectedURLWithSSL {
		t.Errorf("URL com SSL esperada %s, mas foi %s", expectedURLWithSSL, actualURLWithSSL)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "variável de ambiente definida",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom_value",
			expected:     "custom_value",
		},
		{
			name:         "variável de ambiente não definida",
			key:          "NONEXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "variável de ambiente vazia",
			key:          "EMPTY_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpa a variável de ambiente
			os.Unsetenv(tt.key)

			// Define a variável de ambiente se necessário
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getEnv(tt.key, tt.defaultValue)

			if result != tt.expected {
				t.Errorf("Resultado esperado %s, mas foi %s", tt.expected, result)
			}
		})
	}
}
