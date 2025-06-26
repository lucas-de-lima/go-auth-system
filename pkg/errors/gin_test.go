package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GinMiddlewareRecovery())
	return router
}

func TestGinHandleError(t *testing.T) {
	router := setupGinTest()

	// Handler que usa GinHandleError com um AppError
	router.GET("/test/app-error", func(c *gin.Context) {
		err := ErrNotFound
		GinHandleError(c, err)
	})

	// Handler que usa GinHandleError com um erro padrão
	router.GET("/test/std-error", func(c *gin.Context) {
		err := errors.New("erro padrão")
		GinHandleError(c, err)
	})

	// Teste com AppError
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/app-error", nil)
	router.ServeHTTP(w, req)

	// Verificando o código de status
	if w.Code != http.StatusNotFound {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusNotFound, w.Code)
	}

	// Verificando o corpo da resposta
	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}
	if response.Message != "Recurso não encontrado" {
		t.Errorf("Esperava mensagem 'Recurso não encontrado', obteve '%s'", response.Message)
	}

	// Teste com erro padrão
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test/std-error", nil)
	router.ServeHTTP(w, req)

	// Verificando o código de status
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusInternalServerError, w.Code)
	}

	// Verificando o corpo da resposta
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}
	if response.Message != "Erro interno do servidor" {
		t.Errorf("Esperava mensagem 'Erro interno do servidor', obteve '%s'", response.Message)
	}
}

func TestGinMiddlewareRecovery(t *testing.T) {
	router := setupGinTest()

	// Handler que causa um pânico
	router.GET("/test/panic", func(c *gin.Context) {
		panic("este é um pânico simulado")
	})

	// Handler normal
	router.GET("/test/normal", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Teste com pânico
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/panic", nil)
	router.ServeHTTP(w, req)

	// Verificando o código de status
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusInternalServerError, w.Code)
	}

	// Teste sem pânico
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test/normal", nil)
	router.ServeHTTP(w, req)

	// Verificando o código de status
	if w.Code != http.StatusOK {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "OK" {
		t.Errorf("Esperava corpo 'OK', obteve '%s'", w.Body.String())
	}
}

func TestGinValidationResponse(t *testing.T) {
	router := setupGinTest()

	// Handler que retorna erro de validação
	router.GET("/test/validation", func(c *gin.Context) {
		details := []ValidationDetail{
			{Field: "email", Message: "Email inválido"},
			{Field: "password", Message: "Senha muito curta"},
		}
		err := NewValidationError("Erro de validação", details)
		response := GinValidationResponse(err)
		GinRespondWithJSON(c, http.StatusBadRequest, response)
	})

	// Teste com erro de validação
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/validation", nil)
	router.ServeHTTP(w, req)

	// Verificando o código de status
	if w.Code != http.StatusBadRequest {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusBadRequest, w.Code)
	}

	// Verificando o corpo da resposta
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Verificando se a mensagem está presente
	message, ok := response["message"]
	if !ok || message != "Erro de validação" {
		t.Errorf("Esperava mensagem 'Erro de validação', obteve '%v'", message)
	}

	// Verificando se os detalhes estão presentes
	details, ok := response["details"]
	if !ok {
		t.Error("Detalhes de validação não encontrados na resposta")
	}

	detailsMap, ok := details.(map[string]interface{})
	if !ok {
		t.Error("Detalhes de validação não estão no formato esperado")
	}

	fields, ok := detailsMap["fields"]
	if !ok {
		t.Error("Campo 'fields' não encontrado nos detalhes")
	}

	fieldsMap, ok := fields.(map[string]interface{})
	if !ok {
		t.Error("Campo 'fields' não está no formato esperado")
	}

	// Verificando os campos específicos
	email, ok := fieldsMap["email"]
	if !ok || email != "Email inválido" {
		t.Errorf("Campo 'email' não encontrado ou com valor incorreto: %v", email)
	}

	password, ok := fieldsMap["password"]
	if !ok || password != "Senha muito curta" {
		t.Errorf("Campo 'password' não encontrado ou com valor incorreto: %v", password)
	}
}

func TestGinRespondWithValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/validation", func(c *gin.Context) {
		details := map[string]interface{}{"email": "inválido"}
		GinRespondWithValidationError(c, "Erro de validação", details)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/validation", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusBadRequest)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}
	if resp.Message != "Erro de validação" {
		t.Errorf("Esperava mensagem 'Erro de validação', obteve '%s'", resp.Message)
	}
	if resp.Details == nil {
		t.Error("Esperava detalhes de validação, mas veio nil")
	}
}

// assertStatus é um helper para comparar status HTTP
func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Status esperado %d, mas foi %d", want, got)
	}
}
