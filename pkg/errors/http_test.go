package errors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithValidationError(t *testing.T) {
	r := httptest.NewRecorder()
	details := map[string]interface{}{"field": "inválido"}
	RespondWithValidationError(r, "Erro de validação", details)

	assertStatusHTTP(t, r.Code, http.StatusBadRequest)

	var resp ErrorResponse
	err := json.Unmarshal(r.Body.Bytes(), &resp)
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

func TestRespondWithJSON(t *testing.T) {
	r := httptest.NewRecorder()
	payload := map[string]string{"ok": "true"}
	RespondWithJSON(r, http.StatusCreated, payload)

	assertStatusHTTP(t, r.Code, http.StatusCreated)

	var resp map[string]string
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta JSON: %v", err)
	}
	if resp["ok"] != "true" {
		t.Errorf("Esperava campo 'ok' igual a 'true', obteve '%s'", resp["ok"])
	}
}

// assertStatusHTTP é um helper para comparar status HTTP
func assertStatusHTTP(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Status esperado %d, mas foi %d", want, got)
	}
}
