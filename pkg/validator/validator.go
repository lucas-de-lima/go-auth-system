package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Init inicializa o validador
func Init() {
	validate = validator.New()
}

// ValidateStruct valida uma estrutura e retorna uma lista de erros de validação
func ValidateStruct(s interface{}) []ValidationError {
	if validate == nil {
		Init()
	}

	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   toSnakeCase(err.Field()),
				Message: getErrorMessage(err),
			})
		}
	}

	return errors
}

// IsEmail valida se uma string é um email válido
func IsEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// toSnakeCase converte uma string de camelCase para snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getErrorMessage retorna uma mensagem de erro baseada na regra de validação
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Este campo é obrigatório"
	case "email":
		return "Email inválido"
	case "min":
		return fmt.Sprintf("Deve ter no mínimo %s caracteres", err.Param())
	case "max":
		return fmt.Sprintf("Deve ter no máximo %s caracteres", err.Param())
	default:
		return fmt.Sprintf("Validação falhou para a regra: %s", err.Tag())
	}
}
