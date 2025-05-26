package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	// Teste sem erro interno
	err := AppError{
		Code:    http.StatusNotFound,
		Message: "Usuário não encontrado",
	}

	if err.Error() != "Usuário não encontrado" {
		t.Errorf("Esperava 'Usuário não encontrado', obteve '%s'", err.Error())
	}

	// Teste com erro interno
	internalErr := errors.New("erro de banco de dados")
	err = AppError{
		Code:     http.StatusInternalServerError,
		Message:  "Erro interno",
		Internal: internalErr,
	}

	expected := "Erro interno: erro de banco de dados"
	if err.Error() != expected {
		t.Errorf("Esperava '%s', obteve '%s'", expected, err.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	internalErr := errors.New("erro interno")
	err := AppError{
		Code:     http.StatusBadRequest,
		Message:  "Erro de validação",
		Internal: internalErr,
	}

	if err.Unwrap() != internalErr {
		t.Error("Unwrap deveria retornar o erro interno")
	}
}

func TestAppError_WithError(t *testing.T) {
	originalErr := AppError{
		Code:    http.StatusBadRequest,
		Message: "Requisição inválida",
	}

	internalErr := errors.New("dados inválidos")
	newErr := originalErr.WithError(internalErr)

	if newErr.Code != originalErr.Code {
		t.Errorf("Esperava código %d, obteve %d", originalErr.Code, newErr.Code)
	}

	if newErr.Message != originalErr.Message {
		t.Errorf("Esperava mensagem '%s', obteve '%s'", originalErr.Message, newErr.Message)
	}

	if newErr.Internal != internalErr {
		t.Error("Erro interno não foi definido corretamente")
	}
}

func TestAppError_WithMessage(t *testing.T) {
	originalErr := AppError{
		Code:     http.StatusBadRequest,
		Message:  "Requisição inválida",
		Internal: errors.New("erro original"),
	}

	newMessage := "Formato de dados inválido"
	newErr := originalErr.WithMessage(newMessage)

	if newErr.Code != originalErr.Code {
		t.Errorf("Esperava código %d, obteve %d", originalErr.Code, newErr.Code)
	}

	if newErr.Message != newMessage {
		t.Errorf("Esperava mensagem '%s', obteve '%s'", newMessage, newErr.Message)
	}

	if newErr.Internal != originalErr.Internal {
		t.Error("Erro interno não foi mantido corretamente")
	}
}

func TestErrors_Is(t *testing.T) {
	err1 := ErrNotFound
	err2 := ErrNotFound.WithError(errors.New("erro interno"))

	// Mesmo tipo de erro deve retornar true com errors.Is
	if !Is(err2, err1) {
		t.Error("errors.Is deveria retornar true para erros do mesmo tipo")
	}

	// Tipos diferentes devem retornar false
	if Is(err2, ErrBadRequest) {
		t.Error("errors.Is deveria retornar false para erros de tipos diferentes")
	}
}

func TestErrors_As(t *testing.T) {
	originalErr := errors.New("erro original")
	wrappedErr := Wrap(originalErr, http.StatusBadRequest, "Requisição inválida")

	var appErr AppError
	if !As(wrappedErr, &appErr) {
		t.Error("errors.As deveria extrair o AppError corretamente")
	}

	if appErr.Code != http.StatusBadRequest {
		t.Errorf("Código incorreto: esperava %d, obteve %d", http.StatusBadRequest, appErr.Code)
	}

	if appErr.Message != "Requisição inválida" {
		t.Errorf("Mensagem incorreta: esperava 'Requisição inválida', obteve '%s'", appErr.Message)
	}

	if appErr.Internal != originalErr {
		t.Error("Erro interno não foi preservado corretamente")
	}
}

func TestGetStatusCode(t *testing.T) {
	// Teste com AppError
	err := ErrNotFound
	if GetStatusCode(err) != http.StatusNotFound {
		t.Errorf("Esperava status %d, obteve %d", http.StatusNotFound, GetStatusCode(err))
	}

	// Teste com erro padrão
	stdErr := errors.New("erro padrão")
	if GetStatusCode(stdErr) != http.StatusInternalServerError {
		t.Errorf("Esperava status %d para erro padrão, obteve %d",
			http.StatusInternalServerError, GetStatusCode(stdErr))
	}
}

func TestGetMessage(t *testing.T) {
	// Teste com AppError
	err := ErrForbidden
	if GetMessage(err) != "Acesso negado" {
		t.Errorf("Esperava mensagem 'Acesso negado', obteve '%s'", GetMessage(err))
	}

	// Teste com erro padrão
	stdErr := errors.New("erro padrão")
	if GetMessage(stdErr) != "Erro interno do servidor" {
		t.Errorf("Esperava mensagem 'Erro interno do servidor', obteve '%s'",
			GetMessage(stdErr))
	}
}
