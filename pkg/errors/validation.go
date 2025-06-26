package errors

import (
	"net/http"
	"strings"
)

// ValidationDetail representa um erro de validação em um campo específico
type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError cria um erro de validação com detalhes
func NewValidationError(message string, details []ValidationDetail) AppError {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Internal: &ValidationErrors{
			Details: details,
		},
	}
}

// ValidationErrors encapsula vários erros de validação
type ValidationErrors struct {
	Details []ValidationDetail
}

// Error implementa a interface error
func (e *ValidationErrors) Error() string {
	if len(e.Details) == 0 {
		return "erro de validação"
	}

	var messages []string
	for _, detail := range e.Details {
		messages = append(messages, detail.Field+": "+detail.Message)
	}

	return "erros de validação: " + strings.Join(messages, "; ")
}

// GetValidationDetails extrai os detalhes de validação de um erro
func GetValidationDetails(err error) ([]ValidationDetail, bool) {
	var appErr AppError
	if !As(err, &appErr) {
		return nil, false
	}

	var validationErrors *ValidationErrors
	if !As(appErr.Internal, &validationErrors) {
		return nil, false
	}

	return validationErrors.Details, true
}

// FormatValidationResponse formata uma resposta de erro de validação
func FormatValidationResponse(err error) map[string]interface{} {
	resp := make(map[string]interface{})
	resp["message"] = GetMessage(err)

	if details, ok := GetValidationDetails(err); ok {
		fields := make(map[string]interface{})
		for _, detail := range details {
			fields[detail.Field] = detail.Message
		}
		resp["fields"] = fields
	}
	return resp
}
