package errors

import "net/http"

var (
	// ErrInternalServer representa um erro interno do servidor
	ErrInternalServer = AppError{
		Code:    http.StatusInternalServerError,
		Message: "Erro interno do servidor",
	}

	// ErrBadRequest representa um erro de requisição inválida
	ErrBadRequest = AppError{
		Code:    http.StatusBadRequest,
		Message: "Requisição inválida",
	}

	// ErrUnauthorized representa um erro de autenticação
	ErrUnauthorized = AppError{
		Code:    http.StatusUnauthorized,
		Message: "Não autorizado",
	}

	// ErrForbidden representa um erro de permissão
	ErrForbidden = AppError{
		Code:    http.StatusForbidden,
		Message: "Acesso negado",
	}

	// ErrNotFound representa um erro de recurso não encontrado
	ErrNotFound = AppError{
		Code:    http.StatusNotFound,
		Message: "Recurso não encontrado",
	}

	// ErrConflict representa um erro de conflito
	ErrConflict = AppError{
		Code:    http.StatusConflict,
		Message: "Conflito de recursos",
	}

	// ErrValidation representa um erro de validação
	ErrValidation = AppError{
		Code:    http.StatusBadRequest,
		Message: "Erro de validação",
	}

	// Erros específicos de usuário
	ErrUserNotFound = AppError{
		Code:    http.StatusNotFound,
		Message: "Usuário não encontrado",
	}

	ErrEmailAlreadyExists = AppError{
		Code:    http.StatusConflict,
		Message: "Email já está em uso",
	}

	ErrInvalidCredentials = AppError{
		Code:    http.StatusUnauthorized,
		Message: "Credenciais inválidas",
	}

	ErrInvalidToken = AppError{
		Code:    http.StatusUnauthorized,
		Message: "Token inválido ou expirado",
	}

	ErrMissingToken = AppError{
		Code:    http.StatusUnauthorized,
		Message: "Token de autenticação não fornecido",
	}

	ErrPasswordTooWeak = AppError{
		Code:    http.StatusBadRequest,
		Message: "A senha não atende aos requisitos mínimos de segurança",
	}

	// Outros erros específicos da aplicação podem ser adicionados aqui
)
