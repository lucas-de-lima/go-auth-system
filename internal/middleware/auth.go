package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
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
			errors.HandleError(w, errors.ErrMissingToken)
			return
		}

		// Extrai o token do cabeçalho (formato: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errors.HandleError(w, errors.ErrBadRequest.WithMessage("Formato de autorização inválido"))
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			logging.Error("Token inválido: %v", err)
			errors.HandleError(w, errors.ErrInvalidToken.WithError(err))
			return
		}

		// Adiciona informações do usuário ao contexto
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Continua para o próximo handler com o contexto atualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GinAuthenticate é um middleware de autenticação para o Gin
func (m *AuthMiddleware) GinAuthenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		ip := c.ClientIP()
		rota := c.FullPath()
		userAgent := c.Request.UserAgent()

		if authHeader == "" {
			logging.Warning("[%s] [%s] [%s] Tentativa de acesso sem token de autenticação", ip, rota, userAgent)
			errors.GinHandleError(c, errors.ErrMissingToken)
			c.Abort()
			return
		}

		// Extrai o token do cabeçalho (formato: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logging.Warning("[%s] [%s] [%s] Formato de token inválido: '%s'", ip, rota, userAgent, authHeader)
			errors.GinHandleError(c, errors.ErrBadRequest.WithMessage("Formato de autorização inválido"))
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			logging.Warning("[%s] [%s] [%s] Token inválido: %v", ip, rota, userAgent, err)
			errors.GinHandleError(c, errors.ErrInvalidToken.WithError(err))
			c.Abort()
			return
		}

		// Adiciona informações do usuário ao contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("roles", claims.Roles)

		logging.Info("[%s] [%s] [%s] Autenticação bem-sucedida para user_id=%s, email=%s", ip, rota, userAgent, claims.UserID, claims.Email)

		// Continua para o próximo handler
		c.Next()
	}
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

// GinRequireRole verifica se o usuário tem um papel específico (versão Gin)
func (m *AuthMiddleware) GinRequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rota := c.FullPath()
		userAgent := c.Request.UserAgent()
		userID, exists := c.Get("user_id")
		userEmail, _ := c.Get("user_email")

		// Busca as roles do contexto (claims do JWT)
		rolesIface, hasRoles := c.Get("roles")
		var roles []string
		if hasRoles {
			roles, _ = rolesIface.([]string)
		}

		if !exists || !hasRoles || !containsRole(roles, role) {
			logging.Warning("[%s] [%s] [%s] Acesso negado: usuário (id=%v, email=%v) não possui o papel '%s'", ip, rota, userAgent, userID, userEmail, role)
			errors.GinHandleError(c, errors.ErrForbidden.WithMessage("Acesso negado: permissão insuficiente"))
			c.Abort()
			return
		}

		logging.Info("[%s] [%s] [%s] Usuário autorizado (id=%v, email=%v) com papel '%s'", ip, rota, userAgent, userID, userEmail, role)
		c.Next()
	}
}

// containsRole verifica se o slice de roles contém o papel exigido
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
