package api

import (
	"encoding/json"
	"net/http"

	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/lucas-de-lima/go-auth-system/internal/service"
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
		http.Error(w, "Erro ao processar solicitação", http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.userService.Create(user); err != nil {
		logging.Error("Erro ao criar usuário: %v", err)
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
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
		http.Error(w, "Erro ao processar solicitação", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logging.Error("Erro na autenticação: %v", err)
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// GetCurrentUser retorna os dados do usuário autenticado
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// O middleware de autenticação adiciona o ID do usuário no contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetByID(userID.(string))
	if err != nil {
		logging.Error("Erro ao buscar usuário: %v", err)
		http.Error(w, "Erro ao buscar dados do usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
