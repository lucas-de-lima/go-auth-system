package api

import (
	"encoding/json"
	"net/http"

	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
	"github.com/lucas-de-lima/go-auth-system/pkg/errors"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

// Handler contém os manipuladores da API
type Handler struct {
	userService service.UserService
}

// NewHandler cria uma nova instância do Handler
func NewHandler(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// RegisterRoutes registra as rotas da API
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/register", h.RegisterUser)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("GET /api/users/me", h.GetCurrentUser)
}

// RegisterUser manipula o registro de novos usuários
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.HandleError(w, errors.ErrBadRequest.WithError(err))
		return
	}

	// Validação básica
	if req.Email == "" || req.Password == "" || req.Name == "" {
		details := []errors.ValidationDetail{
			{Field: "email", Message: "Email é obrigatório"},
			{Field: "password", Message: "Senha é obrigatória"},
			{Field: "name", Message: "Nome é obrigatório"},
		}
		validationErr := errors.NewValidationError("Campos obrigatórios não preenchidos", details)
		errors.HandleError(w, validationErr)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.userService.Create(user); err != nil {
		logging.Error("Erro ao criar usuário: %v", err)
		errors.HandleError(w, err)
		return
	}

	errors.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Usuário registrado com sucesso",
		"id":      user.ID,
	})
}

// Login manipula a autenticação de usuários
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("Erro ao decodificar corpo da requisição: %v", err)
		errors.HandleError(w, errors.ErrBadRequest.WithError(err))
		return
	}

	token, refreshToken, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logging.Error("Erro na autenticação: %v", err)
		errors.HandleError(w, err)
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, map[string]string{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

// GetCurrentUser retorna os dados do usuário autenticado
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// O middleware de autenticação adiciona o ID do usuário no contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		errors.HandleError(w, errors.ErrUnauthorized)
		return
	}

	user, err := h.userService.GetByID(userID.(string))
	if err != nil {
		logging.Error("Erro ao buscar usuário: %v", err)
		errors.HandleError(w, err)
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, user)
}
