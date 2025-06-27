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
	userService domain.UserService
}

func NewUserController(userService domain.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var user domain.UserRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logging.Error("[%s] Falha ao decodificar corpo da requisição de registro: %v", ctx.ClientIP(), err)
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

		logging.Warning("[%s] Tentativa de registro com campos obrigatórios faltando: %+v", ctx.ClientIP(), details)
		validationErr := errors.NewValidationError("Campos obrigatórios não preenchidos", details)
		errors.GinHandleError(ctx, validationErr)
		return
	}

	newUser := user.FromUserRequest()
	err := uc.userService.Create(newUser)
	if err != nil {
		logging.Error("[%s] Falha ao registrar usuário %s: %v", ctx.ClientIP(), newUser.Email, err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Novo usuário registrado: %s (id: %s)", ctx.ClientIP(), newUser.Email, newUser.ID)
	errors.GinRespondWithJSON(ctx, http.StatusCreated, newUser.ToUserResponse())
}

func (uc *UserController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Error("[%s] Falha ao decodificar corpo da requisição de login: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	accessToken, refreshToken, err := uc.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logging.Warning("[%s] Tentativa de login falhou para: %s (%v)", ctx.ClientIP(), req.Email, err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Login realizado: %s", ctx.ClientIP(), req.Email)
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
		logging.Error("[%s] Falha ao decodificar corpo da requisição de logout: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	if req.RefreshToken == "" {
		logging.Warning("[%s] Tentativa de logout sem refresh token", ctx.ClientIP())
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("Token de atualização não fornecido"))
		return
	}

	service.BlacklistRefreshToken(req.RefreshToken)
	logging.Info("[%s] Logout realizado (rota: %s)", ctx.ClientIP(), ctx.FullPath())
	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logging.Error("[%s] Falha ao decodificar corpo da requisição de refresh: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	if req.RefreshToken == "" {
		logging.Warning("[%s] Tentativa de refresh sem refresh token", ctx.ClientIP())
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("Token de atualização não fornecido"))
		return
	}

	accessToken, newRefreshToken, err := uc.userService.RefreshTokens(req.RefreshToken)
	if err != nil {
		logging.Warning("[%s] Tentativa de refresh token falhou: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Refresh token bem-sucedido (rota: %s)", ctx.ClientIP(), ctx.FullPath())
	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"token":         accessToken,
		"refresh_token": newRefreshToken,
	})
}

// GetByID busca um usuário pelo ID
func (uc *UserController) GetByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		logging.Warning("[%s] Tentativa de busca de usuário sem ID", ctx.ClientIP())
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	user, err := uc.userService.GetByID(userID)
	if err != nil {
		logging.Warning("[%s] Falha ao buscar usuário por ID %s: %v", ctx.ClientIP(), userID, err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Usuário consultado: id=%s", ctx.ClientIP(), userID)
	errors.GinRespondWithJSON(ctx, http.StatusOK, user.ToUserResponse())
}

// Update atualiza os dados de um usuário
func (uc *UserController) Update(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		logging.Warning("[%s] Tentativa de atualização sem ID", ctx.ClientIP())
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	var updateData struct {
		Email string `json:"email,omitempty"`
		Name  string `json:"name,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		logging.Error("[%s] Falha ao decodificar corpo da requisição de update: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}

	currentUser, err := uc.userService.GetByID(userID)
	if err != nil {
		logging.Warning("[%s] Falha ao buscar usuário para atualização: %v", ctx.ClientIP(), err)
		errors.GinHandleError(ctx, err)
		return
	}

	if updateData.Email != "" {
		currentUser.Email = updateData.Email
	}
	if updateData.Name != "" {
		currentUser.Name = updateData.Name
	}

	err = uc.userService.Update(currentUser)
	if err != nil {
		logging.Error("[%s] Falha ao atualizar usuário %s: %v", ctx.ClientIP(), userID, err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Usuário atualizado: id=%s", ctx.ClientIP(), userID)
	errors.GinRespondWithJSON(ctx, http.StatusOK, currentUser.ToUserResponse())
}

// Delete remove um usuário
func (uc *UserController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		logging.Warning("[%s] Tentativa de deleção sem ID", ctx.ClientIP())
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}

	err := uc.userService.Delete(userID)
	if err != nil {
		logging.Error("[%s] Falha ao deletar usuário %s: %v", ctx.ClientIP(), userID, err)
		errors.GinHandleError(ctx, err)
		return
	}

	logging.Info("[%s] Usuário deletado: id=%s", ctx.ClientIP(), userID)
	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{
		"message": "Usuário deletado com sucesso",
	})
}
