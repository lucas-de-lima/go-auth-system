package service

import (
	"time"

	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// UserService implementa a interface domain.UserService
type UserService struct {
	userRepo   domain.UserRepository
	jwtService *auth.JWTService
}

// NewUserService cria uma nova instância do serviço de usuário
func NewUserService(userRepo domain.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Create cria um novo usuário
func (us *UserService) Create(user *domain.User) error {
	// Verifica se já existe um usuário com o mesmo email
	existingUser, err := us.userRepo.GetByEmail(user.Email)
	if err != nil {
		logging.Error("Erro ao verificar email: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	if existingUser != nil {
		return errors.ErrEmailAlreadyExists
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Error("Erro ao gerar hash da senha: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	// Atualiza a senha com o hash
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Salva o usuário no repositório
	err = us.userRepo.Create(user)
	if err != nil {
		logging.Error("Erro ao criar usuário: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	return nil
}

// GetByID busca um usuário pelo ID
func (us *UserService) GetByID(id string) (*domain.User, error) {
	user, err := us.userRepo.GetByID(id)
	if err != nil {
		logging.Error("Erro ao buscar usuário por ID: %v", err)
		return nil, errors.ErrInternalServer.WithError(err)
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// GetByEmail busca um usuário pelo email
func (us *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := us.userRepo.GetByEmail(email)
	if err != nil {
		logging.Error("Erro ao buscar usuário por email: %v", err)
		return nil, errors.ErrInternalServer.WithError(err)
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// Update atualiza os dados de um usuário
func (us *UserService) Update(user *domain.User) error {
	// Verifica se o usuário existe
	existingUser, err := us.userRepo.GetByID(user.ID)
	if err != nil {
		logging.Error("Erro ao verificar usuário: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	// Atualiza o usuário
	user.UpdatedAt = time.Now()
	err = us.userRepo.Update(user)
	if err != nil {
		logging.Error("Erro ao atualizar usuário: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	return nil
}

// Delete remove um usuário
func (us *UserService) Delete(id string) error {
	// Verifica se o usuário existe
	existingUser, err := us.userRepo.GetByID(id)
	if err != nil {
		logging.Error("Erro ao verificar usuário: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	// Remove o usuário
	err = us.userRepo.Delete(id)
	if err != nil {
		logging.Error("Erro ao excluir usuário: %v", err)
		return errors.ErrInternalServer.WithError(err)
	}

	return nil
}

// Authenticate autentica um usuário e retorna access token e refresh token
func (us *UserService) Authenticate(email, password string) (string, string, error) {
	// Busca o usuário pelo email
	user, err := us.userRepo.GetByEmail(email)
	if err != nil {
		logging.Error("Erro ao buscar usuário para autenticação: %v", err)
		return "", "", errors.ErrInternalServer.WithError(err)
	}

	if user == nil {
		return "", "", errors.ErrInvalidCredentials
	}

	// Verifica a senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logging.Error("Senha inválida para usuário %s: %v", email, err)
		return "", "", errors.ErrInvalidCredentials
	}

	// Gera o token JWT
	accessToken, err := us.jwtService.GenerateToken(user)
	if err != nil {
		logging.Error("Erro ao gerar token JWT: %v", err)
		return "", "", errors.ErrInternalServer.WithError(err)
	}

	refreshToken, err := us.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		logging.Error("Erro ao gerar refresh token: %v", err)
		return "", "", errors.ErrInternalServer.WithError(err)
	}

	return accessToken, refreshToken, nil
}

// refreshTokenBlacklist é um mapa em memória para blacklist de refresh tokens
var refreshTokenBlacklist = make(map[string]struct{})

// RefreshTokens realiza a rotação do refresh token e gera novos tokens
func (us *UserService) RefreshTokens(refreshToken string) (string, string, error) {
	// Verifica se o token está na blacklist
	if _, blacklisted := refreshTokenBlacklist[refreshToken]; blacklisted {
		return "", "", errors.ErrUnauthorized.WithMessage("Refresh token inválido ou já utilizado")
	}

	claims, err := us.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.ErrUnauthorized.WithError(err)
	}

	userID := claims.Subject
	user, err := us.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return "", "", errors.ErrUserNotFound
	}

	// Gera novos tokens
	accessToken, err := us.jwtService.GenerateToken(user)
	if err != nil {
		return "", "", errors.ErrInternalServer.WithError(err)
	}
	newRefreshToken, err := us.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.ErrInternalServer.WithError(err)
	}

	// Adiciona o refresh token antigo à blacklist
	refreshTokenBlacklist[refreshToken] = struct{}{}

	return accessToken, newRefreshToken, nil
}

// BlacklistRefreshToken adiciona um refresh token à blacklist em memória
func BlacklistRefreshToken(token string) {
	refreshTokenBlacklist[token] = struct{}{}
}

// ClearRefreshTokenBlacklist limpa a blacklist de refresh tokens (usado apenas para testes)
func ClearRefreshTokenBlacklist() {
	refreshTokenBlacklist = make(map[string]struct{})
}

// GetJWTService retorna o ponteiro do JWTService (uso exclusivo para testes)
func (us *UserService) GetJWTService() *auth.JWTService {
	return us.jwtService
}

// ListAll retorna todos os usuários (admin)
func (us *UserService) ListAll() ([]*domain.User, error) {
	users, err := us.userRepo.List()
	if err != nil {
		logging.Error("Erro ao listar usuários: %v", err)
		return nil, errors.ErrInternalServer.WithError(err)
	}
	return users, nil
}
