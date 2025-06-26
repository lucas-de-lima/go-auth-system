package service

import (
	"errors"
	"testing"

	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}
func (m *mockUserRepo) Create(user *domain.User) error {
	if _, exists := m.users[user.ID]; exists {
		return errors.New("already exists")
	}
	if user.ID == "" {
		user.ID = user.Email // simplificação
	}
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) GetByID(id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}
func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) Update(user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("not found")
	}
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) Delete(id string) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("not found")
	}
	delete(m.users, id)
	return nil
}
func (m *mockUserRepo) List() ([]*domain.User, error) {
	var list []*domain.User
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

type errorRepo struct{}

func (e *errorRepo) Create(user *domain.User) error          { return errors.New("repo error") }
func (e *errorRepo) GetByID(id string) (*domain.User, error) { return nil, errors.New("repo error") }
func (e *errorRepo) GetByEmail(email string) (*domain.User, error) {
	return nil, errors.New("repo error")
}
func (e *errorRepo) Update(user *domain.User) error { return errors.New("repo error") }
func (e *errorRepo) Delete(id string) error         { return errors.New("repo error") }
func (e *errorRepo) List() ([]*domain.User, error)  { return nil, errors.New("repo error") }

func TestUserService_CreateAndGet(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	user := &domain.User{ID: "1", Email: "a@b.com", Password: "senha", Name: "A"}
	err := us.Create(user)
	assert.NoError(t, err)
	// Não permite duplicado
	err = us.Create(user)
	assert.Error(t, err)
	// GetByID
	u, err := us.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "a@b.com", u.Email)
	// GetByEmail
	u, err = us.GetByEmail("a@b.com")
	assert.NoError(t, err)
	assert.Equal(t, "1", u.ID)
}

func TestUserService_UpdateAndDelete(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	user := &domain.User{ID: "2", Email: "b@b.com", Password: "senha", Name: "B"}
	_ = us.Create(user)
	user.Name = "Novo Nome"
	err := us.Update(user)
	assert.NoError(t, err)
	u, _ := us.GetByID("2")
	assert.Equal(t, "Novo Nome", u.Name)
	// Delete
	err = us.Delete("2")
	assert.NoError(t, err)
	u, _ = us.GetByID("2")
	assert.Nil(t, u)
}

func TestUserService_Authenticate(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	user := &domain.User{ID: "3", Email: "c@b.com", Password: "senha123", Name: "C"}
	_ = us.Create(user)
	// Sucesso
	access, refresh, err := us.Authenticate("c@b.com", "senha123")
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	// Senha errada
	_, _, err = us.Authenticate("c@b.com", "errada")
	assert.Error(t, err)
	// Email não existe
	_, _, err = us.Authenticate("nao@existe.com", "senha")
	assert.Error(t, err)
}

func TestUserService_RefreshTokens(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	user := &domain.User{ID: "4", Email: "d@b.com", Password: "senha", Name: "D"}
	_ = us.Create(user)
	_, refresh, _ := us.Authenticate("d@b.com", "senha")
	access2, refresh2, err := us.RefreshTokens(refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, access2)
	assert.NotEmpty(t, refresh2)
	// Token já usado (blacklist)
	_, _, err = us.RefreshTokens(refresh)
	assert.Error(t, err)
}

func TestUserService_BlacklistAndClear(t *testing.T) {
	BlacklistRefreshToken("token1")
	_, _, err := NewUserService(newMockUserRepo(), auth.NewJWTService("s", 1, "r", 1)).RefreshTokens("token1")
	assert.Error(t, err)
	ClearRefreshTokenBlacklist()
	// Agora não está mais na blacklist
	// Não retorna erro de blacklist, mas sim de token inválido
	_, _, err = NewUserService(newMockUserRepo(), auth.NewJWTService("s", 1, "r", 1)).RefreshTokens("token1")
	assert.Error(t, err)
}

func TestUserService_ListAll(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	_ = us.Create(&domain.User{ID: "5", Email: "e@b.com", Password: "senha", Name: "E"})
	_ = us.Create(&domain.User{ID: "6", Email: "f@b.com", Password: "senha", Name: "F"})
	users, err := us.ListAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserService_Create_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	user := &domain.User{ID: "x", Email: "x@x.com", Password: "senha"}
	err := us.Create(user)
	assert.Error(t, err)
}

func TestUserService_GetByID_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	_, err := us.GetByID("x")
	assert.Error(t, err)
}

func TestUserService_GetByEmail_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	_, err := us.GetByEmail("x@x.com")
	assert.Error(t, err)
}

func TestUserService_Update_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	user := &domain.User{ID: "x", Email: "x@x.com", Password: "senha"}
	err := us.Update(user)
	assert.Error(t, err)
}

func TestUserService_Delete_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	err := us.Delete("x")
	assert.Error(t, err)
}

func TestUserService_ListAll_RepoError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	_, err := us.ListAll()
	assert.Error(t, err)
}

// Simular erro de hash
func TestUserService_Create_HashError(t *testing.T) {
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	repo := newMockUserRepo()
	us := NewUserService(repo, jwtService)
	// Forçar senha muito longa para estourar o bcrypt
	user := &domain.User{ID: "y", Email: "y@y.com", Password: string(make([]byte, 10000))}
	err := us.Create(user)
	assert.Error(t, err)
}

func TestUserService_RefreshTokens_InvalidToken(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	_, _, err := us.RefreshTokens("tokeninvalido")
	assert.Error(t, err)
}

func TestUserService_RefreshTokens_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	// Gera refresh token válido para um ID que não existe no repo
	token, _ := jwtService.GenerateRefreshToken("naoexiste")
	_, _, err := us.RefreshTokens(token)
	assert.Error(t, err)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	user := &domain.User{ID: "naoexiste", Email: "x@x.com", Password: "senha"}
	err := us.Update(user)
	assert.Error(t, err)
}

func TestUserService_Delete_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	err := us.Delete("naoexiste")
	assert.Error(t, err)
}

func TestUserService_Authenticate_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	_, _, err := us.Authenticate("naoexiste@x.com", "senha")
	assert.Error(t, err)
}

func TestUserService_GetJWTService(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	assert.Equal(t, jwtService, us.GetJWTService())
}

func TestUserService_ClearAndBlacklistRefreshToken(t *testing.T) {
	BlacklistRefreshToken("tokentest")
	// Está na blacklist
	_, _, err := NewUserService(newMockUserRepo(), auth.NewJWTService("s", 1, "r", 1)).RefreshTokens("tokentest")
	assert.Error(t, err)
	ClearRefreshTokenBlacklist()
	// Não está mais na blacklist, mas token inválido
	_, _, err = NewUserService(newMockUserRepo(), auth.NewJWTService("s", 1, "r", 1)).RefreshTokens("tokentest")
	assert.Error(t, err)
}

func TestUserService_ListAll_ErrorAndEmpty(t *testing.T) {
	// Erro
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(&errorRepo{}, jwtService)
	_, err := us.ListAll()
	assert.Error(t, err)
	// Lista vazia
	repo := newMockUserRepo()
	us2 := NewUserService(repo, jwtService)
	users, err := us2.ListAll()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestUserService_List(t *testing.T) {
	repo := newMockUserRepo()
	jwtService := auth.NewJWTService("secret", 1, "refresh", 1)
	us := NewUserService(repo, jwtService)
	_ = us.Create(&domain.User{ID: "7", Email: "g@b.com", Password: "senha", Name: "G"})
	_ = us.Create(&domain.User{ID: "8", Email: "h@b.com", Password: "senha", Name: "H"})
	users, err := us.List()
	assert.NoError(t, err, "Erro inesperado ao listar usuários")
	assert.Len(t, users, 2, "Deveria retornar 2 usuários")
}
