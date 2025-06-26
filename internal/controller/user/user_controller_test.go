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

type mockUserService struct {
	CreateFn        func(*domain.User) error
	AuthenticateFn  func(string, string) (string, string, error)
	RefreshTokensFn func(string) (string, string, error)
	GetByIDFn       func(string) (*domain.User, error)
	UpdateFn        func(*domain.User) error
	DeleteFn        func(string) error
	GetByEmailFn    func(string) (*domain.User, error)
	ListFn          func() ([]*domain.User, error)
}

func (m *mockUserService) Create(u *domain.User) error { return m.CreateFn(u) }
func (m *mockUserService) Authenticate(e, p string) (string, string, error) {
	return m.AuthenticateFn(e, p)
}
func (m *mockUserService) RefreshTokens(t string) (string, string, error) {
	return m.RefreshTokensFn(t)
}
func (m *mockUserService) GetByID(id string) (*domain.User, error) { return m.GetByIDFn(id) }
func (m *mockUserService) Update(u *domain.User) error             { return m.UpdateFn(u) }
func (m *mockUserService) Delete(id string) error                  { return m.DeleteFn(id) }
func (m *mockUserService) GetByEmail(email string) (*domain.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(email)
	}
	return nil, nil
}
func (m *mockUserService) List() ([]*domain.User, error) {
	if m.ListFn != nil {
		return m.ListFn()
	}
	return nil, nil
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// Testa o registro de usuário com dados válidos, espera sucesso (201)
func TestUserController_Register_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_Register_Success")

	// Arrange: Configura o mock e dados de entrada
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/register", uc.Register)
	body := map[string]interface{}{"email": "a@b.com", "password": "123", "name": "Lucas"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição
	r.ServeHTTP(w, req)

	// Assert: Verifica o resultado esperado
	assert.Equal(t, http.StatusCreated, w.Code)
	t.Log("[FIM] TestUserController_Register_Success")
}

// Testa o registro de usuário com JSON malformado, espera erro 400
func TestUserController_Register_BadRequest(t *testing.T) {
	t.Log("[INICIO] TestUserController_Register_BadRequest")

	// Arrange: Configura o mock e dados de entrada malformados
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/register", uc.Register)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{malformed}")))
	req.Header.Set("Content-Type", "application/json")

	// Act: Executa a requisição com JSON inválido
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Log("[FIM] TestUserController_Register_BadRequest")
}

// Testa login com credenciais válidas, espera sucesso (200)
func TestUserController_Login_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_Login_Success")

	// Arrange: Configura o mock para retornar tokens válidos
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "access", "refresh", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/login", uc.Login)
	body := map[string]interface{}{"email": "a@b.com", "password": "123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de login
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Login_Success")
}

// Testa login com credenciais inválidas, espera erro 401
func TestUserController_Login_InvalidCredentials(t *testing.T) {
	t.Log("[INICIO] TestUserController_Login_InvalidCredentials")

	// Arrange: Configura o mock para retornar erro de credenciais inválidas
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", pkgerrors.ErrInvalidCredentials },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/login", uc.Login)
	body := map[string]interface{}{"email": "a@b.com", "password": "wrong"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de login com credenciais inválidas
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Log("[FIM] TestUserController_Login_InvalidCredentials")
}

// Testa logout com refresh token válido, espera sucesso (200)
func TestUserController_Logout_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_Logout_Success")

	// Arrange: Configura o mock e dados de entrada
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/logout", uc.Logout)
	body := map[string]interface{}{"refresh_token": "valid-refresh-token"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de logout
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Logout_Success")
}

// Testa logout sem refresh token, espera erro 400
func TestUserController_Logout_NoRefreshToken(t *testing.T) {
	t.Log("[INICIO] TestUserController_Logout_NoRefreshToken")

	// Arrange: Configura o mock e dados de entrada sem refresh token
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/logout", uc.Logout)
	body := map[string]interface{}{} // Sem refresh token
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de logout sem refresh token
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Log("[FIM] TestUserController_Logout_NoRefreshToken")
}

// Testa refresh token com token válido, espera sucesso (200)
func TestUserController_RefreshToken_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_RefreshToken_Success")

	// Arrange: Configura o mock para retornar novos tokens
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "new-access", "new-refresh", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/refresh", uc.RefreshToken)
	body := map[string]interface{}{"refresh_token": "valid-refresh-token"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de refresh token
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_RefreshToken_Success")
}

// Testa refresh token sem token, espera erro 400
func TestUserController_RefreshToken_NoToken(t *testing.T) {
	t.Log("[INICIO] TestUserController_RefreshToken_NoToken")

	// Arrange: Configura o mock e dados de entrada sem refresh token
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/refresh", uc.RefreshToken)
	body := map[string]interface{}{} // Sem refresh token
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de refresh sem token
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Log("[FIM] TestUserController_RefreshToken_NoToken")
}

// Testa refresh token com token inválido, espera erro 401
func TestUserController_RefreshToken_InvalidToken(t *testing.T) {
	t.Log("[INICIO] TestUserController_RefreshToken_InvalidToken")

	// Arrange: Configura o mock para retornar erro de token inválido
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", pkgerrors.ErrUnauthorized },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/refresh", uc.RefreshToken)
	body := map[string]interface{}{"refresh_token": "invalid-token"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de refresh com token inválido
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Log("[FIM] TestUserController_RefreshToken_InvalidToken")
}

// Testa busca de usuário por ID com ID válido, espera sucesso (200)
func TestUserController_GetByID_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_GetByID_Success")

	// Arrange: Configura o mock para retornar usuário válido
	user := &domain.User{ID: "123", Email: "a@b.com", Name: "Lucas"}
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return user, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.GET("/users/:id", uc.GetByID)
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de busca por ID
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_GetByID_Success")
}

// Testa busca de usuário por ID sem ID, espera erro 404 (Gin não faz match da rota)
func TestUserController_GetByID_NoID(t *testing.T) {
	t.Log("[INICIO] TestUserController_GetByID_NoID")

	// Arrange: Configura o mock e rota sem parâmetro ID
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.GET("/users/:id", uc.GetByID)
	req := httptest.NewRequest("GET", "/users/", nil) // Sem ID
	w := httptest.NewRecorder()

	// Act: Executa a requisição sem ID
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404 (Gin não faz match da rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_GetByID_NoID")
}

// Testa busca de usuário por ID inexistente, espera erro 404
func TestUserController_GetByID_NotFound(t *testing.T) {
	t.Log("[INICIO] TestUserController_GetByID_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, pkgerrors.ErrUserNotFound },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.GET("/users/:id", uc.GetByID)
	req := httptest.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_GetByID_NotFound")
}

// Testa atualização de usuário com dados válidos, espera sucesso (200)
func TestUserController_Update_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_Success")

	// Arrange: Configura o mock para retornar usuário existente e permitir atualização
	user := &domain.User{ID: "123", Email: "a@b.com", Name: "Lucas"}
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return user, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"name": "Lucas Updated"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/123", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição de atualização
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Update_Success")
}

// Testa atualização de usuário sem ID, espera erro 404 (Gin não faz match da rota)
func TestUserController_Update_NoID(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_NoID")

	// Arrange: Configura o mock e rota sem parâmetro ID
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"name": "Lucas Updated"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/", bytes.NewBuffer(b)) // Sem ID
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição sem ID
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404 (Gin não faz match da rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Update_NoID")
}

// Testa atualização de usuário inexistente, espera erro 404
func TestUserController_Update_NotFound(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, pkgerrors.ErrUserNotFound },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"name": "Lucas Updated"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/999", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Update_NotFound")
}

// Testa deleção de usuário com ID válido, espera sucesso (200)
func TestUserController_Delete_Success(t *testing.T) {
	t.Log("[INICIO] TestUserController_Delete_Success")

	// Arrange: Configura o mock para permitir deleção
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.DELETE("/users/:id", uc.Delete)
	req := httptest.NewRequest("DELETE", "/users/123", nil)
	w := httptest.NewRecorder()

	// Act: Executa a requisição de deleção
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna sucesso 200
	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Delete_Success")
}

// Testa deleção de usuário sem ID, espera erro 404 (Gin não faz match da rota)
func TestUserController_Delete_NoID(t *testing.T) {
	t.Log("[INICIO] TestUserController_Delete_NoID")

	// Arrange: Configura o mock e rota sem parâmetro ID
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.DELETE("/users/:id", uc.Delete)
	req := httptest.NewRequest("DELETE", "/users/", nil) // Sem ID
	w := httptest.NewRecorder()

	// Act: Executa a requisição sem ID
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404 (Gin não faz match da rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Delete_NoID")
}

// Testa deleção de usuário inexistente, espera erro 404
func TestUserController_Delete_NotFound(t *testing.T) {
	t.Log("[INICIO] TestUserController_Delete_NotFound")

	// Arrange: Configura o mock para retornar usuário não encontrado
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return pkgerrors.ErrUserNotFound },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.DELETE("/users/:id", uc.Delete)
	req := httptest.NewRequest("DELETE", "/users/999", nil) // ID inexistente
	w := httptest.NewRecorder()

	// Act: Executa a requisição com ID inexistente
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Delete_NotFound")
}

// Testa registro com campos obrigatórios faltando, espera erro 400
func TestUserController_Register_MissingFields(t *testing.T) {
	t.Log("[INICIO] TestUserController_Register_MissingFields")

	// Arrange: Configura o mock e dados de entrada incompletos
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/register", uc.Register)
	body := map[string]interface{}{"email": "a@b.com"} // Sem senha
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição com campos faltando
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Log("[FIM] TestUserController_Register_MissingFields")
}

// Testa registro com erro do service (email já existe), espera erro 409
func TestUserController_Register_ServiceError(t *testing.T) {
	t.Log("[INICIO] TestUserController_Register_ServiceError")

	// Arrange: Configura o mock para retornar erro
	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return pkgerrors.ErrEmailAlreadyExists },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/register", uc.Register)
	body := map[string]interface{}{"email": "a@b.com", "password": "123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act: Executa a requisição
	r.ServeHTTP(w, req)

	// Assert: Verifica que retorna erro 409
	assert.Equal(t, http.StatusConflict, w.Code)
	t.Log("[FIM] TestUserController_Register_ServiceError")
}

// Testa logout com refresh_token vazio, espera erro 400
func TestUserController_Logout_MissingToken(t *testing.T) {
	uc := NewUserController(&mockUserService{})
	r := setupGin()
	r.POST("/logout", uc.Logout)
	body := map[string]interface{}{"refresh_token": ""}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Testa refresh de token com refresh_token vazio, espera erro 400
func TestUserController_RefreshToken_MissingToken(t *testing.T) {
	ms := &mockUserService{RefreshTokensFn: func(t string) (string, string, error) { return "", "", nil }}
	uc := NewUserController(ms)
	r := setupGin()
	r.POST("/refresh", uc.RefreshToken)
	body := map[string]interface{}{"refresh_token": ""}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Testa GetByID com ID vazio
func TestUserController_GetByID_EmptyID(t *testing.T) {
	t.Log("[INICIO] TestUserController_GetByID_EmptyID")

	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.GET("/users/:id", uc.GetByID)
	req := httptest.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Quando o ID está vazio, o Gin retorna 404 (não encontra a rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_GetByID_EmptyID")
}

// Testa Update com ID vazio
func TestUserController_Update_EmptyID(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_EmptyID")

	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"email": "novo@b.com"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Quando o ID está vazio, o Gin retorna 404 (não encontra a rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Update_EmptyID")
}

// Testa Delete com ID vazio
func TestUserController_Delete_EmptyID(t *testing.T) {
	t.Log("[INICIO] TestUserController_Delete_EmptyID")

	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return nil, nil },
		UpdateFn:        func(*domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.DELETE("/users/:id", uc.Delete)
	req := httptest.NewRequest("DELETE", "/users/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Quando o ID está vazio, o Gin retorna 404 (não encontra a rota)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Log("[FIM] TestUserController_Delete_EmptyID")
}

// Testa Update com apenas email
func TestUserController_Update_OnlyEmail(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_OnlyEmail")

	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return &domain.User{ID: "1", Email: "old@b.com", Name: "Old"}, nil },
		UpdateFn:        func(u *domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"email": "novo@b.com"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Update_OnlyEmail")
}

// Testa Update com apenas nome
func TestUserController_Update_OnlyName(t *testing.T) {
	t.Log("[INICIO] TestUserController_Update_OnlyName")

	ms := &mockUserService{
		CreateFn:        func(u *domain.User) error { return nil },
		AuthenticateFn:  func(string, string) (string, string, error) { return "", "", nil },
		RefreshTokensFn: func(string) (string, string, error) { return "", "", nil },
		GetByIDFn:       func(string) (*domain.User, error) { return &domain.User{ID: "1", Email: "a@b.com", Name: "Old"}, nil },
		UpdateFn:        func(u *domain.User) error { return nil },
		DeleteFn:        func(string) error { return nil },
		GetByEmailFn:    func(string) (*domain.User, error) { return nil, nil },
		ListFn:          func() ([]*domain.User, error) { return nil, nil },
	}
	uc := NewUserController(ms)
	r := setupGin()
	r.PUT("/users/:id", uc.Update)
	body := map[string]interface{}{"name": "Novo Nome"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log("[FIM] TestUserController_Update_OnlyName")
}
