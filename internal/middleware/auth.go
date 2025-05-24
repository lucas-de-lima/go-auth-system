package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

type contextKey string

const (
	// UserIDKey é a chave para o ID do usuário no contexto
	UserIDKey contextKey = "user_id"
	// UserEmailKey é a chave para o email do usuário no contexto
	UserEmailKey contextKey = "user_email"
)

// AuthMiddleware é um middleware que verifica a autenticação JWT
type AuthMiddleware struct {
	jwtService *auth.JWTService
}

// NewAuthMiddleware cria uma nova instância do middleware de autenticação
func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate verifica se o token JWT é válido e adiciona as claims no contexto
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Não autorizado: token não fornecido", http.StatusUnauthorized)
			return
		}

		// Extrai o token do cabeçalho (formato: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Formato de autorização inválido", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			logging.Error("Token inválido: %v", err)
			http.Error(w, "Não autorizado: token inválido", http.StatusUnauthorized)
			return
		}

		// Adiciona informações do usuário ao contexto
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Continua para o próximo handler com o contexto atualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole verifica se o usuário tem um papel específico
// Esta é uma função de exemplo que pode ser expandida conforme necessário
func (m *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Aqui você pode implementar a verificação de papéis/permissões
		// Por exemplo, buscar o usuário no banco de dados e verificar seus papéis

		// Por enquanto, apenas verificamos se o usuário está autenticado
		if r.Context().Value(UserIDKey) == nil {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		// Continua para o próximo handler
		next.ServeHTTP(w, r)
	})
}
