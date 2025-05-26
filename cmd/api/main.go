package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/controller/user"
	"github.com/lucas-de-lima/go-auth-system/internal/repository"
	"github.com/lucas-de-lima/go-auth-system/internal/routes"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/lucas-de-lima/go-auth-system/prisma"
	// outros imports necessários
)

func main() {
	// Carregar variáveis de ambiente do arquivo configs/app.env
	if err := godotenv.Load("configs/app.env"); err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo configs/app.env: %v", err)
	}

	// Inicializar o router do Gin
	// Substituindo gin.Default() por uma configuração personalizada
	router := gin.New()

	// Adicionando middleware de log do Gin
	router.Use(gin.Logger())

	// Adicionando nosso middleware de recuperação personalizado
	router.Use(errors.GinMiddlewareRecovery())

	// Inicializar a conexão com o banco de dados
	prisma.Init()
	defer prisma.Disconnect()

	// Inicializar serviços e repositórios
	userRepository := repository.NewUserRepository(prisma.DB)

	// Obter configurações do JWT do arquivo de ambiente
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your_jwt_secret" // Valor padrão do seu app.env
	}

	refreshKey := os.Getenv("JWT_REFRESH_SECRET")
	if refreshKey == "" {
		refreshKey = "your_refresh_secret" // Valor padrão do seu app.env
	}

	jwtService := auth.NewJWTService(
		secretKey,
		24, // Você pode substituir por os.Getenv("JWT_EXPIRATION_HOURS")
		refreshKey,
		168, // Você pode substituir por os.Getenv("JWT_REFRESH_EXPIRATION_HOURS")
	)

	userService := service.NewUserService(userRepository, jwtService)

	// Inicializar o controller
	userController := user.NewUserController(userService)

	// Inicializar e configurar as rotas
	userRoutes := routes.NewUserRoutes(userController, jwtService)
	userRoutes.Setup(router)

	// Iniciar o servidor
	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
