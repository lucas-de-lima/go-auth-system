package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/routes"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	// outros imports necessários
)

func main() {
	// Inicializar o router do Gin
	router := gin.Default()

	// Inicializar serviços e repositórios
	// Exemplo:
	// userRepository := repository.NewUserRepository(db)
	// jwtService := auth.NewJWTService(secretKey, expTime)
	// userService := service.NewUserService(userRepository, jwtService)

	// Para fins de exemplo, estou criando um serviço mock
	var userService *service.UserService

	// Inicializar o controller
	userController := user.NewUserController(userService)

	// Inicializar e configurar as rotas
	userRoutes := routes.NewUserRoutes(userController)
	userRoutes.Setup(router)

	// Iniciar o servidor
	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
