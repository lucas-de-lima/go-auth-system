package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var user domain.UserRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	// Validação básica
	if user.Email == "" || user.Password == "" {
		details := []errors.ValidationDetail{}

		if user.Email == "" {
			details = append(details, errors.ValidationDetail{Field: "email", Message: "Email é obrigatório"})
		}

		if user.Password == "" {
			details = append(details, errors.ValidationDetail{Field: "password", Message: "Senha é obrigatória"})
		}

		validationErr := errors.NewValidationError("Campos obrigatórios não preenchidos", details)
		errors.GinHandleError(ctx, validationErr)
		return
	}

	newUser := user.FromUserRequest()
	err := uc.userService.Create(newUser)
	if err != nil {
		logging.Error("Erro ao criar usuário: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusCreated, newUser.ToUserResponse())
}

func (uc *UserController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	token, err := uc.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logging.Error("Erro na autenticação: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}

func (uc *UserController) Logout(ctx *gin.Context) {
	// Implementação do método de logout usando o novo sistema de erros
	// Exemplo:
	userID, exists := ctx.Get("user_id")
	if !exists {
		errors.GinHandleError(ctx, errors.ErrUnauthorized)
		return
	}

	// Aqui implementaria a lógica de logout
	logging.Info("Usuário %s realizou logout", userID)

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
	// Implementação do método de refresh token usando o novo sistema de erros
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	if req.RefreshToken == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("Token de atualização não fornecido"))
		return
	}

	// Aqui implementaria a lógica de refresh token
	// Por enquanto apenas um exemplo:
	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Token atualizado com sucesso",
		"token":   "novo-token-jwt",
	})
}
