package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

type AdminController struct {
	userService *service.UserService
}

func NewAdminController(userService *service.UserService) *AdminController {
	return &AdminController{userService: userService}
}

// ListAll lista todos os usuários
func (ac *AdminController) ListAll(ctx *gin.Context) {
	users, err := ac.userService.ListAll()
	if err != nil {
		logging.Error("Erro ao listar usuários: %v", err)
		errors.GinHandleError(ctx, errors.ErrInternalServer.WithError(err))
		return
	}
	responses := make([]*domain.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToUserResponse())
	}
	errors.GinRespondWithJSON(ctx, http.StatusOK, responses)
}

// GetByID busca um usuário pelo ID
func (ac *AdminController) GetByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}
	user, err := ac.userService.GetByID(userID)
	if err != nil {
		logging.Error("Erro ao buscar usuário por ID: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}
	errors.GinRespondWithJSON(ctx, http.StatusOK, user.ToUserResponse())
}

// Update atualiza os dados de um usuário (incluindo roles)
func (ac *AdminController) Update(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}
	var updateData struct {
		Email string   `json:"email,omitempty"`
		Name  string   `json:"name,omitempty"`
		Roles []string `json:"roles,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithError(err))
		return
	}
	currentUser, err := ac.userService.GetByID(userID)
	if err != nil {
		logging.Error("Erro ao buscar usuário para atualização: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}
	if updateData.Email != "" {
		currentUser.Email = updateData.Email
	}
	if updateData.Name != "" {
		currentUser.Name = updateData.Name
	}
	if updateData.Roles != nil {
		currentUser.Roles = updateData.Roles
	}
	err = ac.userService.Update(currentUser)
	if err != nil {
		logging.Error("Erro ao atualizar usuário: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}
	errors.GinRespondWithJSON(ctx, http.StatusOK, currentUser.ToUserResponse())
}

// Delete remove um usuário
func (ac *AdminController) Delete(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		errors.GinHandleError(ctx, errors.ErrBadRequest.WithMessage("ID do usuário não fornecido"))
		return
	}
	err := ac.userService.Delete(userID)
	if err != nil {
		logging.Error("Erro ao deletar usuário: %v", err)
		errors.GinHandleError(ctx, err)
		return
	}
	errors.GinRespondWithJSON(ctx, http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}
