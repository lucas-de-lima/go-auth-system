package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
)

// UserRoutes define as rotas relacionadas a usuários
type UserRoutes struct {
	userController *user.UserController
}

// NewUserRoutes cria uma nova instância de rotas de usuário
func NewUserRoutes(userController *user.UserController) *UserRoutes {
	return &UserRoutes{
		userController: userController,
	}
}

// Setup configura as rotas no router fornecido
func (r *UserRoutes) Setup(router *gin.Engine) {
	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", r.userController.Register)
		userRouter.POST("/login", r.userController.Login)
		userRouter.POST("/logout", r.userController.Logout)
		userRouter.POST("/refresh", r.userController.RefreshToken)
	}
}
