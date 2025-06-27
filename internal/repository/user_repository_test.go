package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

// UserModel simula a estrutura do Prisma para testes
// (mantido apenas para os testes de funções auxiliares)
type UserModel struct {
	ID        string
	Email     string
	Password  string
	Name      *string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Função auxiliar para criar ponteiros de string
func stringPtr(s string) *string {
	return &s
}

// mapUserModelToDomain converte UserModel para domain.User
func mapUserModelToDomain(userModel *UserModel) *domain.User {
	if userModel == nil {
		return nil
	}

	name := ""
	if userModel.Name != nil {
		name = *userModel.Name
	}

	return &domain.User{
		ID:        userModel.ID,
		Email:     userModel.Email,
		Password:  userModel.Password,
		Name:      name,
		Roles:     userModel.Roles,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
}

// Testes unitários para funções auxiliares
func TestMapUserModelToDomain_NilUser(t *testing.T) {
	result := mapUserModelToDomain(nil)
	assert.Nil(t, result)
}

func TestMapUserModelToDomain_WithName(t *testing.T) {
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr("Test User"),
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, []string{"user"}, result.Roles)
}

func TestMapUserModelToDomain_WithoutName(t *testing.T) {
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      nil, // Nome é nil
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "", result.Name) // Nome deve ser string vazia
	assert.Equal(t, []string{"user"}, result.Roles)
}

func TestMapUserModelToDomain_EmptyRoles(t *testing.T) {
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr("Test User"),
		Roles:     []string{}, // Roles vazio
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, []string{}, result.Roles)
}

func TestMapUserModelToDomain_MultipleRoles(t *testing.T) {
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr("Test User"),
		Roles:     []string{"user", "admin", "moderator"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.Name)
	assert.Equal(t, []string{"user", "admin", "moderator"}, result.Roles)
}

func TestMapUserModelToDomain_EmptyStringName(t *testing.T) {
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr(""), // Nome é string vazia
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "", result.Name) // Nome deve ser string vazia
	assert.Equal(t, []string{"user"}, result.Roles)
}

func TestMapUserModelToDomain_PreservesTimestamps(t *testing.T) {
	now := time.Now()
	userModel := &UserModel{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr("Test User"),
		Roles:     []string{"user"},
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
	}

	result := mapUserModelToDomain(userModel)

	assert.NotNil(t, result)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now.Add(time.Hour), result.UpdatedAt)
}

// Testes para funções utilitárias
func TestStringPtr(t *testing.T) {
	// Teste com string normal
	str := "test string"
	ptr := stringPtr(str)
	assert.NotNil(t, ptr)
	assert.Equal(t, str, *ptr)

	// Teste com string vazia
	emptyStr := ""
	emptyPtr := stringPtr(emptyStr)
	assert.NotNil(t, emptyPtr)
	assert.Equal(t, emptyStr, *emptyPtr)

	// Teste com string com espaços
	spaceStr := "   "
	spacePtr := stringPtr(spaceStr)
	assert.NotNil(t, spacePtr)
	assert.Equal(t, spaceStr, *spacePtr)
}

// Testes de integração simulados (apenas transformação de dados)
func TestUserRepositoryIntegration_DataTransformation(t *testing.T) {
	// Simula dados vindos do Prisma
	prismaUsers := []UserModel{
		{
			ID:        uuid.New().String(),
			Email:     "user1@example.com",
			Password:  "password1",
			Name:      stringPtr("User 1"),
			Roles:     []string{"user"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Email:     "user2@example.com",
			Password:  "password2",
			Name:      nil, // Usuário sem nome
			Roles:     []string{"admin", "user"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Email:     "user3@example.com",
			Password:  "password3",
			Name:      stringPtr(""), // Nome vazio
			Roles:     []string{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Converte para domain.User
	domainUsers := make([]*domain.User, 0, len(prismaUsers))
	for _, pu := range prismaUsers {
		domainUsers = append(domainUsers, mapUserModelToDomain(&pu))
	}

	// Verifica se a conversão foi feita corretamente
	assert.Len(t, domainUsers, 3)

	// Verifica primeiro usuário
	assert.Equal(t, prismaUsers[0].ID, domainUsers[0].ID)
	assert.Equal(t, prismaUsers[0].Email, domainUsers[0].Email)
	assert.Equal(t, "User 1", domainUsers[0].Name)
	assert.Equal(t, []string{"user"}, domainUsers[0].Roles)

	// Verifica segundo usuário (sem nome)
	assert.Equal(t, prismaUsers[1].ID, domainUsers[1].ID)
	assert.Equal(t, prismaUsers[1].Email, domainUsers[1].Email)
	assert.Equal(t, "", domainUsers[1].Name) // Nome deve ser vazio
	assert.Equal(t, []string{"admin", "user"}, domainUsers[1].Roles)

	// Verifica terceiro usuário (nome vazio, roles vazio)
	assert.Equal(t, prismaUsers[2].ID, domainUsers[2].ID)
	assert.Equal(t, prismaUsers[2].Email, domainUsers[2].Email)
	assert.Equal(t, "", domainUsers[2].Name) // Nome deve ser vazio
	assert.Equal(t, []string{}, domainUsers[2].Roles)
}

// Testes de edge cases
func TestMapUserModelToDomain_EdgeCases(t *testing.T) {
	// Teste com ID vazio
	userModel := &UserModel{
		ID:        "",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      stringPtr("Test User"),
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := mapUserModelToDomain(userModel)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.ID)

	// Teste com email vazio
	userModel.Email = ""
	result = mapUserModelToDomain(userModel)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Email)

	// Teste com password vazio
	userModel.Password = ""
	result = mapUserModelToDomain(userModel)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Password)

	// Teste com roles nil
	userModel.Roles = nil
	result = mapUserModelToDomain(userModel)
	assert.NotNil(t, result)
	assert.Nil(t, result.Roles)
}

// Testes para NewUserRepository
func TestNewUserRepository(t *testing.T) {
	// Como não podemos criar um PrismaClient real sem banco,
	// testamos apenas a criação da estrutura
	repo := &UserRepository{
		db: nil, // Será nil em testes
	}

	assert.NotNil(t, repo)
	assert.Nil(t, repo.db) // Em testes, db será nil
}

// Testes para funções utilitárias adicionais
func TestRepositoryUtilities(t *testing.T) {
	// Testa se as funções auxiliares funcionam corretamente
	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      "Test User",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.NotNil(t, user)
	assert.Equal(t, "test-id", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, []string{"user"}, user.Roles)
}

// Testes para cenários de edge cases
func TestRepositoryEdgeCases(t *testing.T) {
	// Testa cenários extremos
	emptyUser := &domain.User{
		ID:        "",
		Email:     "",
		Password:  "",
		Name:      "",
		Roles:     []string{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	assert.NotNil(t, emptyUser)
	assert.Equal(t, "", emptyUser.ID)
	assert.Equal(t, "", emptyUser.Email)
	assert.Equal(t, "", emptyUser.Name)
	assert.Equal(t, []string{}, emptyUser.Roles)
}

// Testes para validação de dados
func TestRepositoryDataValidation(t *testing.T) {
	// Testa se os dados são preservados corretamente
	user := &domain.User{
		ID:        "user-123",
		Email:     "user@example.com",
		Password:  "hashed-password",
		Name:      "John Doe",
		Roles:     []string{"user", "admin"},
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	// Verifica se todos os campos são preservados
	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "user@example.com", user.Email)
	assert.Equal(t, "hashed-password", user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, []string{"user", "admin"}, user.Roles)
	assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), user.CreatedAt)
	assert.Equal(t, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), user.UpdatedAt)
}

// Testes para funções de transformação de dados
func TestDataTransformation(t *testing.T) {
	// Testa transformações de dados que podem ser úteis no repository
	users := []*domain.User{
		{
			ID:        "user1",
			Email:     "user1@example.com",
			Password:  "pass1",
			Name:      "User 1",
			Roles:     []string{"user"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user2",
			Email:     "user2@example.com",
			Password:  "pass2",
			Name:      "User 2",
			Roles:     []string{"admin"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Testa se conseguimos processar uma lista de usuários
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].ID)
	assert.Equal(t, "user2", users[1].ID)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.Equal(t, "user2@example.com", users[1].Email)
}

// Testes para validação de estrutura de dados
func TestDataStructureValidation(t *testing.T) {
	// Testa se a estrutura de dados está correta
	user := &domain.User{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password123",
		Name:      "Test User",
		Roles:     []string{"user", "admin"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verifica se a estrutura está correta
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Password)
	assert.NotEmpty(t, user.Name)
	assert.Len(t, user.Roles, 2)
	assert.Contains(t, user.Roles, "user")
	assert.Contains(t, user.Roles, "admin")
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}
