package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/middleware"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAdminTestEnvironment configura o ambiente de teste com um admin
func setupAdminTestEnvironment() (*gin.Engine, *service.UserService, *auth.JWTService, string) {
	gin.SetMode(gin.TestMode)
	service.ClearRefreshTokenBlacklist()
	memRepo := NewInMemoryUserRepository()
	jwtService := auth.NewJWTService(
		"test-secret-key",
		24,
		"test-refresh-key",
		168,
	)
	userService := service.NewUserService(memRepo, jwtService)
	userController := user.NewUserController(userService)
	adminController := user.NewAdminController(userService)
	router := gin.New()
	// Rotas públicas
	router.POST("/users/register", userController.Register)
	router.POST("/users/login", userController.Login)
	// Rotas admin
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(authMiddleware.GinAuthenticate(), authMiddleware.GinRequireRole("admin"))
	{
		adminRoutes.GET("/users", adminController.ListAll)
		adminRoutes.GET("/users/:id", adminController.GetByID)
		adminRoutes.PUT("/users/:id", adminController.Update)
		adminRoutes.DELETE("/users/:id", adminController.Delete)
	}
	// Criar usuário admin
	adminUser := &domain.User{
		Email:    "admin@example.com",
		Password: "adminpass",
		Name:     "Admin User",
		Roles:    []string{"admin"},
	}
	err := userService.Create(adminUser)
	require.NoError(nil, err)
	// Obter token admin
	accessToken, _, err := userService.Authenticate("admin@example.com", "adminpass")
	require.NoError(nil, err)
	return router, userService, jwtService, accessToken
}

func TestAdminListAllUsers(t *testing.T) {
	router, userService, _, adminToken := setupAdminTestEnvironment()
	// Criar usuário comum
	user := &domain.User{
		Email:    "user1@example.com",
		Password: "userpass",
		Name:     "User 1",
	}
	err := userService.Create(user)
	require.NoError(t, err)
	// Requisição
	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, len(response) >= 1)
}

func TestAdminGetUserByID(t *testing.T) {
	router, userService, _, adminToken := setupAdminTestEnvironment()
	user := &domain.User{
		Email:    "user2@example.com",
		Password: "userpass",
		Name:     "User 2",
	}
	err := userService.Create(user)
	require.NoError(t, err)
	// Buscar pelo ID
	req := httptest.NewRequest("GET", "/admin/users/"+user.ID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, user.Email, response["email"])
	// Caso de usuário não encontrado
	req2 := httptest.NewRequest("GET", "/admin/users/inexistente", nil)
	req2.Header.Set("Authorization", "Bearer "+adminToken)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestAdminUpdateUser(t *testing.T) {
	router, userService, _, adminToken := setupAdminTestEnvironment()
	user := &domain.User{
		Email:    "user3@example.com",
		Password: "userpass",
		Name:     "User 3",
	}
	err := userService.Create(user)
	require.NoError(t, err)
	// Atualizar nome e roles
	updateData := map[string]interface{}{
		"name":  "Novo Nome",
		"roles": []string{"admin", "user"},
	}
	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/admin/users/"+user.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Novo Nome", response["name"])
	assert.Contains(t, response["roles"].([]interface{}), "admin")
	// Atualizar usuário inexistente
	req2 := httptest.NewRequest("PUT", "/admin/users/inexistente", bytes.NewBuffer(jsonData))
	req2.Header.Set("Authorization", "Bearer "+adminToken)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestAdminDeleteUser(t *testing.T) {
	router, userService, _, adminToken := setupAdminTestEnvironment()
	user := &domain.User{
		Email:    "user4@example.com",
		Password: "userpass",
		Name:     "User 4",
	}
	err := userService.Create(user)
	require.NoError(t, err)
	req := httptest.NewRequest("DELETE", "/admin/users/"+user.ID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Usuário deletado com sucesso", response["message"])
	// Deletar usuário inexistente
	req2 := httptest.NewRequest("DELETE", "/admin/users/inexistente", nil)
	req2.Header.Set("Authorization", "Bearer "+adminToken)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestAdminAccessDenied(t *testing.T) {
	router, userService, _, _ := setupAdminTestEnvironment()
	// Criar usuário comum
	user := &domain.User{
		Email:    "user5@example.com",
		Password: "userpass",
		Name:     "User 5",
	}
	err := userService.Create(user)
	require.NoError(t, err)
	// Obter token de usuário comum
	userToken, _, err := userService.Authenticate("user5@example.com", "userpass")
	require.NoError(t, err)
	// Tentar acessar rota admin sem token
	req := httptest.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	// Tentar acessar rota admin com token de usuário comum
	req2 := httptest.NewRequest("GET", "/admin/users", nil)
	req2.Header.Set("Authorization", "Bearer "+userToken)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)
	// Token inválido
	req3 := httptest.NewRequest("GET", "/admin/users", nil)
	req3.Header.Set("Authorization", "Bearer tokeninvalido")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusUnauthorized, w3.Code)
}
