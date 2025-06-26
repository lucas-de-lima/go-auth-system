package errors

import (
	"strings"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	details := []ValidationDetail{{Field: "email", Message: "inválido"}}
	err := NewValidationError("Erro de validação", details)

	msg := err.Error()
	if msg == "Erro de validação" {
		t.Error("Mensagem de erro deveria conter detalhes dos campos, mas veio simplificada")
	}
	if !contains(msg, "Erro de validação") {
		t.Errorf("Mensagem deveria conter 'Erro de validação', obteve '%s'", msg)
	}
	if !contains(msg, "email: inválido") {
		t.Errorf("Mensagem deveria conter detalhes do campo, obteve '%s'", msg)
	}
}

func TestFormatValidationResponse(t *testing.T) {
	details := []ValidationDetail{
		{Field: "email", Message: "inválido"},
		{Field: "password", Message: "curta"},
	}
	err := NewValidationError("Erro de validação", details)

	resp := FormatValidationResponse(err)
	if resp["message"] != "Erro de validação" {
		t.Errorf("Esperava mensagem 'Erro de validação', obteve '%v'", resp["message"])
	}

	fields, ok := resp["fields"].(map[string]interface{})
	if !ok {
		t.Fatalf("Esperava fields como map[string]interface{}, obteve %T", resp["fields"])
	}

	expected := map[string]string{"email": "inválido", "password": "curta"}
	for k, v := range expected {
		if fields[k] != v {
			t.Errorf("Esperava fields[%s]=%s, obteve %v", k, v, fields[k])
		}
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
