package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	pkgerrors "github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockAdminUserService struct {
	ListFn    func() ([]*domain.User, error)
	GetByIDFn func(string) (*domain.User, error)
	UpdateFn  func(*domain.User) error
	DeleteFn  func(string) error
}

func (m *mockAdminUserService) List() ([]*domain.User, error)           { return m.ListFn() }
func (m *mockAdminUserService) GetByID(id string) (*domain.User, error) { return m.GetByIDFn(id) }
func (m *mockAdminUserService) Update(u *domain.User) error             { return m.UpdateFn(u) }
func (m *mockAdminUserService) Delete(id string) error                  { return m.DeleteFn(id) }

// Métodos não usados
func (m *mockAdminUserService) Create(u *domain.User) error                      { return nil }
func (m *mockAdminUserService) GetByEmail(email string) (*domain.User, error)    { return nil, nil }
func (m *mockAdminUserService) Authenticate(e, p string) (string, string, error) { return "", "", nil }
func (m *mockAdminUserService) RefreshTokens(t string) (string, string, error)   { return "", "", nil }

func setupGinAdmin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAdminController_ListAll_Success(t *testing.T) {
	t.Log("[INICIO] TestAdminController_ListAll_Success")

	// Arrange: Configura o mock para retornar lista de usuários
	ms := &mockAdminUserService{ListFn: func() ([]*domain.User, error) {
		return []*domain.User{{ID: "1", Email: "a@b.com"}}, nil
	}}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.GET("/admin/users", ac.ListAll)
	req := httptest.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de listagem
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestAdminController_ListAll_Success")
}

func TestAdminController_ListAll_Error(t *testing.T) {
	t.Log("[INICIO] TestAdminController_ListAll_Error")

	// Arrange: Configura o mock para retornar erro interno
	ms := &mockAdminUserService{ListFn: func() ([]*domain.User, error) {
		return nil, pkgerrors.ErrInternalServer
	}}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.GET("/admin/users", ac.ListAll)
	req := httptest.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de listagem
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	t.Log("[FIM] TestAdminController_ListAll_Error")
}

func TestAdminController_GetByID_Success(t *testing.T) {
	t.Log("[INICIO] TestAdminController_GetByID_Success")

	// Arrange: Configura o mock para retornar usuário válido
	ms := &mockAdminUserService{GetByIDFn: func(id string) (*domain.User, error) {
		return &domain.User{ID: id, Email: "a@b.com"}, nil
	}}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.GET("/admin/users/:id", ac.GetByID)
	req := httptest.NewRequest("GET", "/admin/users/1", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de busca por ID
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestAdminController_GetByID_Success")
}

func TestAdminController_GetByID_NotFound(t *testing.T) {
	t.Log("[INICIO] TestAdminController_GetByID_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockAdminUserService{GetByIDFn: func(id string) (*domain.User, error) {
		return nil, pkgerrors.ErrUserNotFound
	}}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.GET("/admin/users/:id", ac.GetByID)
	req := httptest.NewRequest("GET", "/admin/users/1", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestAdminController_GetByID_NotFound")
}

func TestAdminController_Update_Success(t *testing.T) {
	t.Log("[INICIO] TestAdminController_Update_Success")

	// Arrange: Configura o mock para retornar usuário existente e permitir atualização
	ms := &mockAdminUserService{
		GetByIDFn: func(id string) (*domain.User, error) { return &domain.User{ID: id, Email: "a@b.com"}, nil },
		UpdateFn:  func(u *domain.User) error { return nil },
	}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.PUT("/admin/users/:id", ac.Update)
	body := map[string]interface{}{"email": "novo@b.com", "name": "Novo", "roles": []string{"admin"}}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de atualização
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestAdminController_Update_Success")
}

func TestAdminController_Update_NotFound(t *testing.T) {
	t.Log("[INICIO] TestAdminController_Update_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockAdminUserService{
		GetByIDFn: func(id string) (*domain.User, error) { return nil, pkgerrors.ErrUserNotFound },
	}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.PUT("/admin/users/:id", ac.Update)
	body := map[string]interface{}{"email": "novo@b.com"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/admin/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestAdminController_Update_NotFound")
}

func TestAdminController_Delete_Success(t *testing.T) {
	t.Log("[INICIO] TestAdminController_Delete_Success")

	// Arrange: Configura o mock para permitir deleção
	ms := &mockAdminUserService{DeleteFn: func(id string) error { return nil }}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.DELETE("/admin/users/:id", ac.Delete)
	req := httptest.NewRequest("DELETE", "/admin/users/1", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de deleção
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestAdminController_Delete_Success")
}

func TestAdminController_Delete_NotFound(t *testing.T) {
	t.Log("[INICIO] TestAdminController_Delete_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockAdminUserService{DeleteFn: func(id string) error { return pkgerrors.ErrUserNotFound }}
	ac := NewAdminController(ms)
	r := setupGinAdmin()
	r.DELETE("/admin/users/:id", ac.Delete)
	req := httptest.NewRequest("DELETE", "/admin/users/1", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestAdminController_Delete_NotFound")
}
