package validator

import (
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Email string `validate:"required,email"`
	Name  string `validate:"required,min=3,max=10"`
}

type fakeFieldError struct {
	tag   string
	param string
}

func (f fakeFieldError) Tag() string                      { return f.tag }
func (f fakeFieldError) Param() string                    { return f.param }
func (f fakeFieldError) Field() string                    { return "Field" }
func (f fakeFieldError) StructField() string              { return "Field" }
func (f fakeFieldError) StructNamespace() string          { return "" }
func (f fakeFieldError) Namespace() string                { return "" }
func (f fakeFieldError) Kind() reflect.Kind               { return reflect.String }
func (f fakeFieldError) Type() reflect.Type               { return reflect.TypeOf("") }
func (f fakeFieldError) Value() interface{}               { return "" }
func (f fakeFieldError) ParamInt() int64                  { return 0 }
func (f fakeFieldError) ActualTag() string                { return f.tag }
func (f fakeFieldError) Error() string                    { return "" }
func (f fakeFieldError) Translate(_ ut.Translator) string { return "" }

func TestIsEmail(t *testing.T) {
	assert.True(t, IsEmail("a@b.com"))
	assert.False(t, IsEmail("a@b"))
	assert.False(t, IsEmail(""))
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, "email_test", toSnakeCase("EmailTest"))
	assert.Equal(t, "nome", toSnakeCase("Nome"))
}

func TestGetErrorMessage(t *testing.T) {
	cases := []struct {
		tag    string
		param  string
		expect string
	}{
		{"required", "", "Este campo é obrigatório"},
		{"email", "", "Email inválido"},
		{"min", "3", "mínimo"},
		{"max", "10", "máximo"},
		{"outra", "", "Validação falhou"},
	}
	for _, c := range cases {
		fake := fakeFieldError{tag: c.tag, param: c.param}
		msg := getErrorMessage(fake)
		assert.Contains(t, msg, c.expect)
	}
}

func TestValidateStruct(t *testing.T) {
	Init()
	// Válido
	v := testStruct{Email: "a@b.com", Name: "Lucas"}
	errs := ValidateStruct(v)
	assert.Len(t, errs, 0)
	// Inválido
	v2 := testStruct{Email: "errado", Name: "Lu"}
	errs = ValidateStruct(v2)
	assert.True(t, len(errs) > 0)
	fields := []string{errs[0].Field, errs[len(errs)-1].Field}
	assert.Contains(t, fields, "email")
	assert.Contains(t, fields, "name")
}
