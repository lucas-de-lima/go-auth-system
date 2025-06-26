package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
)

// JWTService é o serviço responsável por gerenciar tokens JWT
type JWTService struct {
	secretKey      string
	expirationTime int
	refreshKey     string
	refreshExpTime int
}

// TokenClaims define as claims customizadas para o token JWT
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService cria uma nova instância do serviço JWT
func NewJWTService(secretKey string, expirationHours int, refreshKey string, refreshExpHours int) *JWTService {
	return &JWTService{
		secretKey:      secretKey,
		expirationTime: expirationHours,
		refreshKey:     refreshKey,
		refreshExpTime: refreshExpHours,
	}
}

// GenerateToken gera um novo token JWT para o usuário
func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(s.expirationTime))

	claims := &TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken valida um token JWT e retorna as claims se válido
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

// GenerateRefreshToken gera um token de atualização
func (s *JWTService) GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(s.refreshExpTime))

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.refreshKey))
}

// ValidateRefreshToken valida um refresh token e retorna as claims se válido
func (s *JWTService) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.refreshKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("refresh token inválido")
}
