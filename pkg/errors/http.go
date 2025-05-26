package errors

import (
	"encoding/json"
	"net/http"

	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

// ErrorResponse é a estrutura da resposta de erro
type ErrorResponse struct {
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HandleError processa o erro e responde adequadamente
func HandleError(w http.ResponseWriter, err error) {
	var appErr AppError
	if !As(err, &appErr) {
		// Se não for um AppError, envolve com ErrInternalServer
		appErr = ErrInternalServer.WithError(err)
	}

	// Loga o erro com o erro interno, se existir
	if appErr.Internal != nil {
		logging.Error("Erro na requisição: %v", appErr)
	}

	// Responde com o erro apropriado
	RespondWithError(w, appErr.Code, appErr.Message)
}

// RespondWithError responde com um erro em formato JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ErrorResponse{
		Message: message,
	})
}

// RespondWithValidationError responde com um erro de validação em formato JSON
func RespondWithValidationError(w http.ResponseWriter, message string, details map[string]interface{}) {
	RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{
		Message: message,
		Details: details,
	})
}

// RespondWithJSON responde com um objeto em formato JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logging.Error("Erro ao serializar resposta JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// WithRecovery é um middleware que recupera de pânicos
func WithRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Converte o pânico em um erro de servidor interno
				err, ok := r.(error)
				if !ok {
					err = ErrInternalServer.WithMessage("Erro interno do servidor")
				}

				logging.Error("Panic recuperado em handler HTTP: %v", r)
				HandleError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
