package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/middleware"
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

	// Rotas CRUD
	router.GET("/users/:id", userController.GetByID)
	router.PUT("/users/:id", userController.Update)
	router.DELETE("/users/:id", userController.Delete)

	// Criar AuthMiddleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Rota protegida para teste de middleware
	router.GET("/protected", authMiddleware.GinAuthenticate(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(200, gin.H{"message": "Acesso permitido", "user_id": userID})
	})

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

// TestUserCRUD testa as operações CRUD de usuário
func TestUserCRUD(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário para teste
	testUser := &domain.User{
		Email:    "crud@example.com",
		Password: "password123",
		Name:     "CRUD User",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	var userID string

	t.Run("should get user by ID successfully", func(t *testing.T) {
		// Buscar o usuário criado para obter o ID
		user, err := userService.GetByEmail("crud@example.com")
		require.NoError(t, err)
		userID = user.ID

		// Act
		req := httptest.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "crud@example.com", response["email"])
		assert.Equal(t, "CRUD User", response["name"])
		assert.Equal(t, userID, response["id"])
		assert.NotEmpty(t, response["created_at"])
		assert.NotEmpty(t, response["updated_at"])
		// Senha não deve estar na resposta
		assert.Nil(t, response["password"])
	})

	t.Run("should fail to get user with invalid ID", func(t *testing.T) {
		// Act
		req := httptest.NewRequest("GET", "/users/invalid-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should update user successfully", func(t *testing.T) {
		// Arrange
		updateData := map[string]interface{}{
			"name":  "Updated CRUD User",
			"email": "updated.crud@example.com",
		}
		jsonData, _ := json.Marshal(updateData)

		// Act
		req := httptest.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "updated.crud@example.com", response["email"])
		assert.Equal(t, "Updated CRUD User", response["name"])
		assert.Equal(t, userID, response["id"])
	})

	t.Run("should update only name when only name is provided", func(t *testing.T) {
		// Arrange
		updateData := map[string]interface{}{
			"name": "Only Name Updated",
		}
		jsonData, _ := json.Marshal(updateData)

		// Act
		req := httptest.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "updated.crud@example.com", response["email"]) // Email não deve ter mudado
		assert.Equal(t, "Only Name Updated", response["name"])         // Nome deve ter mudado
		assert.Equal(t, userID, response["id"])
	})

	t.Run("should update only email when only email is provided", func(t *testing.T) {
		// Arrange
		updateData := map[string]interface{}{
			"email": "only.email@example.com",
		}
		jsonData, _ := json.Marshal(updateData)

		// Act
		req := httptest.NewRequest("PUT", "/users/"+userID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "only.email@example.com", response["email"]) // Email deve ter mudado
		assert.Equal(t, "Only Name Updated", response["name"])       // Nome não deve ter mudado
		assert.Equal(t, userID, response["id"])
	})

	t.Run("should fail to update user with invalid ID", func(t *testing.T) {
		// Arrange
		updateData := map[string]interface{}{
			"name": "Invalid User",
		}
		jsonData, _ := json.Marshal(updateData)

		// Act
		req := httptest.NewRequest("PUT", "/users/invalid-id", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should delete user successfully", func(t *testing.T) {
		// Act
		req := httptest.NewRequest("DELETE", "/users/"+userID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Usuário deletado com sucesso", response["message"])
	})

	t.Run("should fail to get deleted user", func(t *testing.T) {
		// Act
		req := httptest.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail to delete user with invalid ID", func(t *testing.T) {
		// Act
		req := httptest.NewRequest("DELETE", "/users/invalid-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestUserValidation testa validações de segurança no registro de usuário
func TestUserValidation(t *testing.T) {
	router, _ := setupTestEnvironment()

	t.Run("should fail with short password", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    "shortpass@example.com",
			"password": "12",
			"name":     "Short Pass",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with empty password", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    "emptypass@example.com",
			"password": "",
			"name":     "Empty Pass",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with empty email", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    "",
			"password": "password123",
			"name":     "Empty Email",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    "invalid-email",
			"password": "password123",
			"name":     "Invalid Email",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":    "duplicate@example.com",
			"password": "password123",
			"name":     "First User",
		}
		jsonData, _ := json.Marshal(userData)

		req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Tentar registrar novamente com o mesmo email
		userData2 := map[string]interface{}{
			"email":    "duplicate@example.com",
			"password": "password123",
			"name":     "Second User",
		}
		jsonData2, _ := json.Marshal(userData2)

		req2 := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code)
	})
}

// TestRefreshTokenInvalidCases testa casos de refresh token expirado e malformado
func TestRefreshTokenInvalidCases(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário e obter JWTService
	testUser := &domain.User{
		Email:    "tokeninvalid@example.com",
		Password: "password123",
		Name:     "Token Invalid",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	jwtService := userService.GetJWTService()

	t.Run("should fail with malformatted refresh token", func(t *testing.T) {
		refreshData := map[string]interface{}{
			"refresh_token": "malformed.token.value",
		}
		jsonData, _ := json.Marshal(refreshData)

		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should fail with expired refresh token", func(t *testing.T) {
		// Gerar um refresh token com expiração curta (1 segundo)
		shortLivedToken, err := generateShortLivedRefreshToken(jwtService, testUser.ID, 1)
		require.NoError(t, err)

		// Esperar expirar
		time.Sleep(2 * time.Second)

		refreshData := map[string]interface{}{
			"refresh_token": shortLivedToken,
		}
		jsonData, _ := json.Marshal(refreshData)

		req := httptest.NewRequest("POST", "/users/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// generateShortLivedRefreshToken gera um refresh token com expiração customizada
func generateShortLivedRefreshToken(jwtService *auth.JWTService, userID string, seconds int) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(seconds))
	claims := jwt.MapClaims{
		"exp": expirationTime.Unix(),
		"sub": userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtService.GetRefreshKey()))
}

// TestAuthMiddleware testa o middleware de autenticação
func TestAuthMiddleware(t *testing.T) {
	router, userService := setupTestEnvironment()

	// Criar usuário e obter token
	testUser := &domain.User{
		Email:    "middleware@example.com",
		Password: "password123",
		Name:     "Middleware User",
	}
	err := userService.Create(testUser)
	require.NoError(t, err)

	accessToken, _, err := userService.Authenticate("middleware@example.com", "password123")
	require.NoError(t, err)

	t.Run("should deny access without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should deny access with malformed token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer malformed.token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should deny access with expired token", func(t *testing.T) {
		jwtService := userService.GetJWTService()
		expiredToken, err := generateShortLivedAccessToken(jwtService, testUser, 1) // 1 segundo
		require.NoError(t, err)
		time.Sleep(2 * time.Second)

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should allow access with valid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Acesso permitido", response["message"])
		assert.NotEmpty(t, response["user_id"])
	})
}

// generateShortLivedAccessToken gera um access token com expiração customizada
func generateShortLivedAccessToken(jwtService *auth.JWTService, user *domain.User, seconds int) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(seconds))
	claims := &auth.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtService.GetSecretKey()))
}
