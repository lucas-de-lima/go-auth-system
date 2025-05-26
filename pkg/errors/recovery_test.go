package errors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithRecovery(t *testing.T) {
	// Handler que causará um pânico
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("este é um pânico simulado")
	})

	// Aplicando o middleware de recuperação
	recoveredHandler := WithRecovery(panicHandler)

	// Criando um request de teste
	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	// Executando o handler (deve recuperar do pânico)
	recoveredHandler.ServeHTTP(recorder, req)

	// Verificando se o código de status é Internal Server Error
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusInternalServerError, recorder.Code)
	}

	// Verificando o formato da resposta
	var response ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar resposta JSON: %v", err)
	}

	// Verificando se a mensagem de erro está presente
	if response.Message != "Erro interno do servidor" {
		t.Errorf("Esperava mensagem 'Erro interno do servidor', obteve '%s'", response.Message)
	}
}

func TestWithRecovery_NoRecovery(t *testing.T) {
	// Handler normal sem pânico
	normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Aplicando o middleware de recuperação
	recoveredHandler := WithRecovery(normalHandler)

	// Criando um request de teste
	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	// Executando o handler (não deve interferir no comportamento normal)
	recoveredHandler.ServeHTTP(recorder, req)

	// Verificando se o código de status é OK
	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava status code %d, obteve %d", http.StatusOK, recorder.Code)
	}

	// Verificando o corpo da resposta
	if recorder.Body.String() != "OK" {
		t.Errorf("Esperava corpo 'OK', obteve '%s'", recorder.Body.String())
	}
}
