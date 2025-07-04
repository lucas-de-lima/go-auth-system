package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
	"github.com/lucas-de-lima/go-auth-system/prisma/db"
)

// UserRepository implementa a interface domain.UserRepository
type UserRepository struct {
	db *db.PrismaClient
}

// NewUserRepository cria uma nova instância do repositório de usuário
func NewUserRepository(db *db.PrismaClient) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create cria um novo usuário no banco de dados
func (ur *UserRepository) Create(user *domain.User) error {
	ctx := context.Background()

	// Gera um novo UUID se não for fornecido
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Cria o usuário no Prisma
	_, err := ur.db.User.CreateOne(
		db.User.Email.Set(user.Email),
		db.User.Password.Set(user.Password),
		db.User.ID.Set(user.ID),
		db.User.Name.Set(user.Name),
		db.User.CreatedAt.Set(user.CreatedAt),
		db.User.UpdatedAt.Set(user.UpdatedAt),
	).Exec(ctx)

	if err != nil {
		logging.Error("Erro ao criar usuário no banco de dados: %v", err)
		return err
	}

	return nil
}

// GetByID busca um usuário pelo ID
func (ur *UserRepository) GetByID(id string) (*domain.User, error) {
	ctx := context.Background()

	prismaUser, err := ur.db.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, nil
		}
		logging.Error("Erro ao buscar usuário por ID: %v", err)
		return nil, err
	}

	return mapPrismaUserToDomain(prismaUser), nil
}

// GetByEmail busca um usuário pelo email
func (ur *UserRepository) GetByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	prismaUser, err := ur.db.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, nil
		}
		logging.Error("Erro ao buscar usuário por email: %v", err)
		return nil, err
	}

	return mapPrismaUserToDomain(prismaUser), nil
}

// Update atualiza os dados de um usuário
func (ur *UserRepository) Update(user *domain.User) error {
	ctx := context.Background()

	_, err := ur.db.User.FindUnique(
		db.User.ID.Equals(user.ID),
	).Update(
		db.User.Email.Set(user.Email),
		db.User.Password.Set(user.Password),
		db.User.Name.Set(user.Name),
		db.User.UpdatedAt.Set(time.Now()),
	).Exec(ctx)

	if err != nil {
		logging.Error("Erro ao atualizar usuário: %v", err)
		return err
	}

	return nil
}

// Delete remove um usuário pelo ID
func (ur *UserRepository) Delete(id string) error {
	ctx := context.Background()

	_, err := ur.db.User.FindUnique(
		db.User.ID.Equals(id),
	).Delete().Exec(ctx)

	if err != nil {
		logging.Error("Erro ao excluir usuário: %v", err)
		return err
	}

	return nil
}

// List retorna todos os usuários
func (ur *UserRepository) List() ([]*domain.User, error) {
	ctx := context.Background()
	prismaUsers, err := ur.db.User.FindMany().Exec(ctx)
	if err != nil {
		logging.Error("Erro ao listar usuários: %v", err)
		return nil, err
	}
	users := make([]*domain.User, 0, len(prismaUsers))
	for _, pu := range prismaUsers {
		users = append(users, mapPrismaUserToDomain(&pu))
	}
	return users, nil
}

// mapPrismaUserToDomain converte um model Prisma para o modelo de domínio
func mapPrismaUserToDomain(prismaUser *db.UserModel) *domain.User {
	if prismaUser == nil {
		return nil
	}

	name := ""
	if prismaUser.InnerUser.Name != nil {
		name = *prismaUser.InnerUser.Name
	}

	return &domain.User{
		ID:        prismaUser.ID,
		Email:     prismaUser.Email,
		Password:  prismaUser.Password,
		Name:      name,
		Roles:     prismaUser.InnerUser.Roles,
		CreatedAt: prismaUser.CreatedAt,
		UpdatedAt: prismaUser.UpdatedAt,
	}
}
