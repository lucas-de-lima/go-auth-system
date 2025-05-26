package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError é o tipo de erro personalizado da aplicação.
type AppError struct {
	// Code é o código de status HTTP
	Code int `json:"-"`
	// Message é a mensagem amigável para o cliente
	Message string `json:"message"`
	// Internal é o erro original para logging/debugging
	Internal error `json:"-"`
}

// Error implementa a interface error
func (e AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Unwrap implementa a interface Unwrapper para compatibilidade com errors.Is/As
func (e AppError) Unwrap() error {
	return e.Internal
}

// StatusCode retorna o código de status HTTP associado a este erro
func (e AppError) StatusCode() int {
	return e.Code
}

// WithError cria uma cópia do erro com um erro interno adicionado
func (e AppError) WithError(err error) AppError {
	return AppError{
		Code:     e.Code,
		Message:  e.Message,
		Internal: err,
	}
}

// WithMessage cria uma cópia do erro com uma mensagem personalizada
func (e AppError) WithMessage(message string) AppError {
	return AppError{
		Code:     e.Code,
		Message:  message,
		Internal: e.Internal,
	}
}

// Is implementa a interface para compatibilidade com errors.Is
func (e AppError) Is(target error) bool {
	t, ok := target.(AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code && e.Message == t.Message
}

// NewAppError cria um novo AppError
func NewAppError(code int, message string, err error) AppError {
	return AppError{
		Code:     code,
		Message:  message,
		Internal: err,
	}
}

// As funções a seguir são wrappers para errors.Is e errors.As
// para facilitar o uso com o pacote errors padrão

// Is relays to errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As relays to errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Wrap envolve um erro com um AppError
func Wrap(err error, code int, message string) AppError {
	return NewAppError(code, message, err)
}

// GetStatusCode obtém o código de status HTTP de um erro
// Se o erro não for um AppError, retorna Internal Server Error
func GetStatusCode(err error) int {
	var appErr AppError
	if As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// GetMessage obtém a mensagem amigável de um erro
func GetMessage(err error) string {
	var appErr AppError
	if As(err, &appErr) {
		return appErr.Message
	}
	return "Erro interno do servidor"
}
