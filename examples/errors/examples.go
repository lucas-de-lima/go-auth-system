package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/middleware"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
)

// Esta é uma função que demonstra como usar o sistema de erros em um serviço
func exemploServico() error {
	// Simulando um erro de banco de dados
	dbErr := fmt.Errorf("falha de conexão com o banco de dados")

	// Cenário 1: Tratando um erro do repositório
	err := funcaoRepositorio()
	if err != nil {
		// Envolvendo o erro do repositório com um erro de aplicação
		return errors.ErrInternalServer.WithError(err)
	}

	// Cenário 2: Erro de validação
	if !emailValido("usuario@exemplo") {
		// Criando um erro de validação
		details := []errors.ValidationDetail{
			{Field: "email", Message: "Formato de email inválido"},
		}
		return errors.NewValidationError("Dados inválidos", details)
	}

	// Cenário 3: Erro de regra de negócio
	if usuarioJaExiste("usuario@exemplo.com") {
		return errors.ErrEmailAlreadyExists
	}

	// Cenário 4: Criando um erro personalizado
	if ocorreuErroInesperado() {
		return errors.NewAppError(
			http.StatusInternalServerError,
			"Erro inesperado durante o processamento",
			dbErr,
		)
	}

	return nil
}

// Exemplo de uso do sistema de erros com HTTP padrão
func exemploHandlerHTTP(w http.ResponseWriter, r *http.Request) {
	// Simulando um erro retornado por um serviço
	err := exemploServico()
	if err != nil {
		// Usando o HandleError para tratar o erro adequadamente
		errors.HandleError(w, err)
		return
	}

	// Verificando o tipo específico de erro
	if errors.Is(err, errors.ErrEmailAlreadyExists) {
		// Tratamento específico para este tipo de erro
		errors.RespondWithError(w, http.StatusConflict, "Este email já está cadastrado. Tente recuperar sua senha.")
		return
	}

	// Tratamento de erros de validação
	details, isValidationErr := errors.GetValidationDetails(err)
	if isValidationErr {
		// Criando uma resposta personalizada para erros de validação
		validationResponse := map[string]interface{}{
			"status":  "erro",
			"message": "Erro de validação",
			"fields":  details,
		}
		errors.RespondWithJSON(w, http.StatusBadRequest, validationResponse)
		return
	}

	// Resposta de sucesso
	resposta := map[string]interface{}{
		"status":  "sucesso",
		"message": "Operação realizada com sucesso",
	}
	errors.RespondWithJSON(w, http.StatusOK, resposta)
}

// Exemplo de uso do sistema de erros com Gin
func exemploHandlerGin(c *gin.Context) {
	// Simulando um erro retornado por um serviço
	err := exemploServico()
	if err != nil {
		// Usando o GinHandleError para tratar o erro adequadamente no Gin
		errors.GinHandleError(c, err)
		return
	}

	// Verificando o tipo específico de erro
	if errors.Is(err, errors.ErrEmailAlreadyExists) {
		// Tratamento específico para este tipo de erro
		errors.GinRespondWithError(c, http.StatusConflict, "Este email já está cadastrado. Tente recuperar sua senha.")
		return
	}

	// Tratamento de erros de validação
	validationResponse := errors.GinValidationResponse(err)
	if validationResponse != nil {
		errors.GinRespondWithJSON(c, http.StatusBadRequest, validationResponse)
		return
	}

	// Resposta de sucesso
	errors.GinRespondWithJSON(c, http.StatusOK, gin.H{
		"status":  "sucesso",
		"message": "Operação realizada com sucesso",
	})
}

// Funções auxiliares simuladas apenas para os exemplos
func funcaoRepositorio() error {
	// Simulando algum erro
	return nil
}

func emailValido(email string) bool {
	// Simulando validação de email
	return len(email) > 5 && email[len(email)-4:] == ".com"
}

func usuarioJaExiste(email string) bool {
	// Simulando verificação se usuário existe
	return email == "usuario@exemplo.com"
}

func ocorreuErroInesperado() bool {
	// Simulando erro inesperado
	return false
}

// Exemplo de como configurar o HTTP padrão com o middleware de recuperação
func configurarRotasHTTP() {
	// Criando um serveMux
	mux := http.NewServeMux()

	// Registrando handlers
	mux.HandleFunc("/api/exemplo", exemploHandlerHTTP)

	// Aplicando o middleware de recuperação
	handler := errors.WithRecovery(mux)

	// Iniciar o servidor HTTP com o handler protegido
	http.ListenAndServe(":8080", handler)
}

// Exemplo de como configurar o Gin com o middleware de recuperação e autenticação
func configurarRotasGin() {
	// Criando um novo router Gin
	router := gin.New()

	// Adicionando middleware de log
	router.Use(gin.Logger())

	// Adicionando middleware de recuperação personalizado
	router.Use(errors.GinMiddlewareRecovery())

	// Criando o serviço JWT (exemplo)
	jwtService := auth.NewJWTService(
		"secret-key",  // Chave secreta
		24,            // Expiração em horas
		"refresh-key", // Chave de refresh
		168,           // Expiração do refresh em horas
	)

	// Criando o middleware de autenticação
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Rotas públicas
	publicRoutes := router.Group("/api")
	{
		publicRoutes.POST("/auth/register", exemploHandlerGin)
		publicRoutes.POST("/auth/login", exemploHandlerGin)
	}

	// Rotas protegidas
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(authMiddleware.GinAuthenticate())
	{
		protectedRoutes.GET("/users/me", exemploHandlerGin)
	}

	// Iniciando o servidor
	router.Run(":8080")
}

func main() {
	// Este é apenas um arquivo de exemplo e não precisa ser executado
}
