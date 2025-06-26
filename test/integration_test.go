package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment configura o ambiente de teste com repositório em memória
func setupTestEnvironment() (*gin.Engine, *service.UserService) {
	// Configurar Gin para modo de teste
	gin.SetMode(gin.TestMode)

	// Limpar blacklist antes de cada teste
	clearRefreshTokenBlacklist()

	// Criar repositório em memória
	memRepo := NewInMemoryUserRepository()

	// Criar JWT service com chaves de teste
	jwtService := auth.NewJWTService(
		"test-secret-key",
		24, // 24 horas
		"test-refresh-key",
		168, // 7 dias
	)

	// Criar user service
	userService := service.NewUserService(memRepo, jwtService)

	// Criar user controller
	userController := user.NewUserController(userService)

	// Configurar rotas
	router := gin.New()
	router.POST("/users/register", userController.Register)
	router.POST("/users/login", userController.Login)
	router.POST("/users/logout", userController.Logout)
	router.POST("/users/refresh", userController.RefreshToken)

	return router, userService
}

// clearRefreshTokenBlacklist limpa a blacklist de refresh tokens para isolamento dos testes
func clearRefreshTokenBlacklist() {
	// Acessar a blacklist através de uma função pública no service
	// Como a blacklist é privada, vamos usar uma abordagem diferente
	// Vou criar uma função no service para limpar a blacklist
	service.ClearRefreshTokenBlacklist()
}

// InMemoryUserRepository implementa um repositório em memória para testes
type InMemoryUserRepository struct {
	users map[string]*domain.User
}

func NewInMemoryUserRepository() domain.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) Update(user *domain.User) error {
	if _, exists := r.users[user.ID]; exists {
		r.users[user.ID] = user
		return nil
	}
	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	delete(r.users, id)
	return nil
}

// TestUserRegistration testa o fluxo de registro de usuário
func TestUserRegistration(t *testing.T) {
	router, _ := setupTestEnvironment()

	t.Run("should register user successfully", func(t *testing.T) {
		// Arrange
		userData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
			"name":     "Test User",
		}
		jsonData, _ := json.Marshal(userData)

		// Act
		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "test@example.com", response["email"])
		assert.Equal(t, "Test User", response["name"])
		assert.NotEmpty(t, response["id"])
		assert.NotEmpty(t, response["created_at"])
		assert.NotEmpty(t, response["updated_at"])
		// Senha não deve estar na resposta
		assert.Nil(t, response["password"])
	})

	t.Run("should fail with empty email", func(t *testing.T) {
		// Arrange
		userData := map[string]interface{}{
			"email":    "",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		// Act
		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with empty password", func(t *testing.T) {
		// Arrange
		userData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "",
		}
		jsonData, _ := json.Marshal(userData)

		// Act
		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestUserLogin testa o fluxo de login
func TestUserLogin(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário para teste
	testUser := &domain.User{
		Email:    "login@example.com",
		Password: "password123",
		Name:     "Login User",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	t.Run("should login successfully", func(t *testing.T) {
		// Arrange
		loginData := map[string]interface{}{
			"email":    "login@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		// Act
		req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["token"])
		assert.NotEmpty(t, response["refresh_token"])
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		// Arrange
		loginData := map[string]interface{}{
			"email":    "login@example.com",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		// Act
		req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should fail with non-existent email", func(t *testing.T) {
		// Arrange
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		// Act
		req := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestRefreshToken testa o fluxo de refresh token
func TestRefreshToken(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário e obter tokens
	testUser := &domain.User{
		Email:    "refresh@example.com",
		Password: "password123",
		Name:     "Refresh User",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	// Fazer login para obter tokens
	accessToken, refreshToken, err := userService.Authenticate("refresh@example.com", "password123")
	require.NoError(t, err)

	t.Run("should refresh tokens successfully", func(t *testing.T) {
		// Arrange
		time.Sleep(2 * time.Second) // Garante que o novo token terá timestamp diferente
		refreshData := map[string]interface{}{
			"refresh_token": refreshToken,
		}
		jsonData, _ := json.Marshal(refreshData)

		// Act
		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotEmpty(t, response["token"])
		assert.NotEmpty(t, response["refresh_token"])

		// Os novos tokens devem ser diferentes dos originais
		assert.NotEqual(t, accessToken, response["token"])
		assert.NotEqual(t, refreshToken, response["refresh_token"])
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		// Arrange
		refreshData := map[string]interface{}{
			"refresh_token": "invalid-token",
		}
		jsonData, _ := json.Marshal(refreshData)

		// Act
		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should fail with empty refresh token", func(t *testing.T) {
		// Arrange
		refreshData := map[string]interface{}{
			"refresh_token": "",
		}
		jsonData, _ := json.Marshal(refreshData)

		// Act
		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail when using same refresh token twice", func(t *testing.T) {
		// Obter um token fresco para este teste específico
		_, freshRefreshToken, err := userService.Authenticate("refresh@example.com", "password123")
		require.NoError(t, err)

		// Primeiro uso do refresh token
		refreshData := map[string]interface{}{
			"refresh_token": freshRefreshToken,
		}
		jsonData, _ := json.Marshal(refreshData)

		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Segundo uso do mesmo refresh token deve falhar
		req2 := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})
}

// TestLogout testa o fluxo de logout
func TestLogout(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário e obter refresh token
	testUser := &domain.User{
		Email:    "logout@example.com",
		Password: "password123",
		Name:     "Logout User",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	_, refreshToken, err := userService.Authenticate("logout@example.com", "password123")
	require.NoError(t, err)

	t.Run("should logout successfully", func(t *testing.T) {
		// Arrange
		logoutData := map[string]interface{}{
			"refresh_token": refreshToken,
		}
		jsonData, _ := json.Marshal(logoutData)

		// Act
		req := httptest.NewRequest("POST", "/users/logout", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Logout realizado com sucesso", response["message"])
	})

	t.Run("should fail with empty refresh token", func(t *testing.T) {
		// Arrange
		logoutData := map[string]interface{}{
			"refresh_token": "",
		}
		jsonData, _ := json.Marshal(logoutData)

		// Act
		req := httptest.NewRequest("POST", "/users/logout", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should invalidate refresh token after logout", func(t *testing.T) {
		// Fazer logout
		logoutData := map[string]interface{}{
			"refresh_token": refreshToken,
		}
		jsonData, _ := json.Marshal(logoutData)

		req := httptest.NewRequest("POST", "/users/logout", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Tentar usar o refresh token após logout deve falhar
		refreshData := map[string]interface{}{
			"refresh_token": refreshToken,
		}
		refreshJson, _ := json.Marshal(refreshData)

		req2 := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(refreshJson))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})
}
