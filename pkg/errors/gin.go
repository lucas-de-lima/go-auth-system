package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/pkg/logging"
)

// GinHandleError processa o erro e responde adequadamente em contexto Gin
func GinHandleError(c *gin.Context, err error) {
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
	GinRespondWithError(c, appErr.Code, appErr.Message)
}

// GinRespondWithError responde com um erro em formato JSON
func GinRespondWithError(c *gin.Context, code int, message string) {
	GinRespondWithJSON(c, code, ErrorResponse{
		Message: message,
	})
}

// GinRespondWithValidationError responde com um erro de validação em formato JSON
func GinRespondWithValidationError(c *gin.Context, message string, details map[string]interface{}) {
	GinRespondWithJSON(c, http.StatusBadRequest, ErrorResponse{
		Message: message,
		Details: details,
	})
}

// GinRespondWithJSON responde com um objeto em formato JSON
func GinRespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

// GinValidationResponse formata uma resposta de erro de validação para Gin
func GinValidationResponse(err error) gin.H {
	details, ok := GetValidationDetails(err)
	if !ok {
		return gin.H{
			"message": GetMessage(err),
		}
	}

	fields := make(map[string]interface{})
	for _, detail := range details {
		fields[detail.Field] = detail.Message
	}

	return gin.H{
		"message": GetMessage(err),
		"details": gin.H{
			"fields": fields,
		},
	}
}

// GinMiddlewareRecovery é um middleware de recuperação para Gin
func GinMiddlewareRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Converte o pânico em um erro de servidor interno
				err, ok := r.(error)
				if !ok {
					err = ErrInternalServer.WithMessage("Erro interno do servidor")
				}

				logging.Error("Panic recuperado em handler Gin: %v", r)
				GinHandleError(c, err)
				c.Abort()
			}
		}()

		c.Next()
	}
}
