package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/middleware"
)

// UserRoutes define as rotas relacionadas a usuários
type UserRoutes struct {
	userController  *user.UserController
	authMiddleware  *middleware.AuthMiddleware
	adminController *user.AdminController
}

// NewUserRoutes cria uma nova instância de rotas de usuário
func NewUserRoutes(userController *user.UserController, jwtService *auth.JWTService, adminController *user.AdminController) *UserRoutes {
	return &UserRoutes{
		userController:  userController,
		authMiddleware:  middleware.NewAuthMiddleware(jwtService),
		adminController: adminController,
	}
}

// Setup configura as rotas no router fornecido
func (ur *UserRoutes) Setup(router *gin.Engine) {
	// Rotas públicas (não autenticadas)
	publicRoutes := router.Group("/users")
	{
		publicRoutes.POST("/register", ur.userController.Register)
		publicRoutes.POST("/login", ur.userController.Login)
		publicRoutes.POST("/refresh", ur.userController.RefreshToken)
	}

	// Rotas protegidas (requerem autenticação)
	protectedRoutes := router.Group("/users")
	protectedRoutes.Use(ur.authMiddleware.GinAuthenticate())
	{
		protectedRoutes.POST("/logout", ur.userController.Logout)
	}

	// Rotas de admin (protegidas por autenticação e role 'admin')
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(ur.authMiddleware.GinAuthenticate(), ur.authMiddleware.GinRequireRole("admin"))
	{
		adminRoutes.GET("/users", ur.adminController.ListAll)
		adminRoutes.GET("/users/:id", ur.adminController.GetByID)
		adminRoutes.PUT("/users/:id", ur.adminController.Update)
		adminRoutes.DELETE("/users/:id", ur.adminController.Delete)
	}
}
