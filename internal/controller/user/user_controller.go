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

	accessToken, refreshToken, err := uc.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logging.Error("Erro na autenticação: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) Logout(ctx *gin.Context) {
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

	// Adiciona o refresh token à blacklist
	service.BlacklistRefreshToken(req.RefreshToken)

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
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

	accessToken, newRefreshToken, err := uc.userService.RefreshTokens(req.RefreshToken)
	if err != nil {
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": newRefreshToken,
	})
}

// GetByID busca um usuário pelo ID
func (uc *UserController) GetByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	user, err := uc.userService.GetByID(userID)
	if err != nil {
		logging.Error("Erro ao buscar usuário por ID: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, user.ToUserResponse())
}

// Update atualiza os dados de um usuário
func (uc *UserController) Update(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	var updateData struct {
		Email string `json:"email,omitempty"`
		Name  string `json:"name,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	// Busca o usuário atual
	currentUser, err := uc.userService.GetByID(userID)
	if err != nil {
		logging.Error("Erro ao buscar usuário para atualização: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	// Atualiza apenas os campos fornecidos
	if updateData.Email != "" {
		currentUser.Email = updateData.Email
	}
	if updateData.Name != "" {
		currentUser.Name = updateData.Name
	}

	err = uc.userService.Update(currentUser)
	if err != nil {
		logging.Error("Erro ao atualizar usuário: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, currentUser.ToUserResponse())
}

// Delete remove um usuário
func (uc *UserController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	err := uc.userService.Delete(userID)
	if err != nil {
		logging.Error("Erro ao deletar usuário: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}

	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Usuário deletado com sucesso",
	})
}
