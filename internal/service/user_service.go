package service

import (
	"errors"
	"time"

	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/repository"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// UserService implementa a interface domain.UserService
type UserService struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserService cria uma nova instância do serviço de usuário
func NewUserService(userRepo repository.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Create cria um novo usuário
func (s *UserService) Create(user *domain.User) error {
	// Verifica se já existe um usuário com o mesmo email
	existingUser, _ := s.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email já está em uso")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Error("Erro ao gerar hash da senha: %v", err)
		return err
	}

	// Atualiza a senha com o hash
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Salva o usuário no repositório
	return s.userRepo.Create(user)
}

// GetByID busca um usuário pelo ID
func (s *UserService) GetByID(id string) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// GetByEmail busca um usuário pelo email
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(email)
}

// Update atualiza os dados de um usuário
func (s *UserService) Update(user *domain.User) error {
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("usuário não encontrado")
	}

	// Se estiver atualizando a senha, hash da nova senha
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			logging.Error("Erro ao gerar hash da senha: %v", err)
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		// Mantém a senha atual
		user.Password = existingUser.Password
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

// Delete remove um usuário pelo ID
func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}

// Authenticate autentica um usuário e retorna um token JWT
func (s *UserService) Authenticate(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("credenciais inválidas")
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logging.Error("Senha incorreta: %v", err)
		return "", errors.New("credenciais inválidas")
	}

	// Gera token JWT
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		logging.Error("Erro ao gerar token JWT: %v", err)
		return "", err
	}

	return token, nil
}
