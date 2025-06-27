package domain

import (
	"testing"
	"time"
)

func TestToUserResponse(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Name:      "Test User",
		Roles:     []string{"user", "admin"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := user.ToUserResponse()

	if response.ID != user.ID {
		t.Errorf("ID esperado %s, mas foi %s", user.ID, response.ID)
	}

	if response.Email != user.Email {
		t.Errorf("Email esperado %s, mas foi %s", user.Email, response.Email)
	}

	if response.Name != user.Name {
		t.Errorf("Name esperado %s, mas foi %s", user.Name, response.Name)
	}

	if len(response.Roles) != len(user.Roles) {
		t.Errorf("Número de roles esperado %d, mas foi %d", len(user.Roles), len(response.Roles))
	}

	for i, role := range user.Roles {
		if response.Roles[i] != role {
			t.Errorf("Role[%d] esperado %s, mas foi %s", i, role, response.Roles[i])
		}
	}

	if response.CreatedAt != user.CreatedAt {
		t.Errorf("CreatedAt esperado %v, mas foi %v", user.CreatedAt, response.CreatedAt)
	}

	if response.UpdatedAt != user.UpdatedAt {
		t.Errorf("UpdatedAt esperado %v, mas foi %v", user.UpdatedAt, response.UpdatedAt)
	}
}

func TestToUserResponseWithEmptyFields(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Name:      "",         // campo vazio
		Roles:     []string{}, // slice vazio
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := user.ToUserResponse()

	if response.ID != user.ID {
		t.Errorf("ID esperado %s, mas foi %s", user.ID, response.ID)
	}

	if response.Email != user.Email {
		t.Errorf("Email esperado %s, mas foi %s", user.Email, response.Email)
	}

	if response.Name != "" {
		t.Errorf("Name deveria ser vazio, mas foi %s", response.Name)
	}

	if len(response.Roles) != 0 {
		t.Errorf("Roles deveria ser vazio, mas tem %d elementos", len(response.Roles))
	}
}

func TestFromUserRequest(t *testing.T) {
	request := &UserRequest{
		Email:    "newuser@example.com",
		Password: "mypassword123",
		Name:     "New User",
	}

	user := request.FromUserRequest()

	if user.Email != request.Email {
		t.Errorf("Email esperado %s, mas foi %s", request.Email, user.Email)
	}

	if user.Password != request.Password {
		t.Errorf("Password esperado %s, mas foi %s", request.Password, user.Password)
	}

	if user.Name != request.Name {
		t.Errorf("Name esperado %s, mas foi %s", request.Name, user.Name)
	}

	// Verifica se o role padrão foi definido
	if len(user.Roles) != 1 {
		t.Errorf("Número de roles esperado 1, mas foi %d", len(user.Roles))
	}

	if user.Roles[0] != "user" {
		t.Errorf("Role padrão esperado 'user', mas foi %s", user.Roles[0])
	}

	// Verifica se os campos de tempo estão zerados (não definidos)
	if !user.CreatedAt.IsZero() {
		t.Errorf("CreatedAt deveria estar zerado, mas foi %v", user.CreatedAt)
	}

	if !user.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt deveria estar zerado, mas foi %v", user.UpdatedAt)
	}

	// Verifica se o ID está vazio (não definido)
	if user.ID != "" {
		t.Errorf("ID deveria estar vazio, mas foi %s", user.ID)
	}
}

func TestFromUserRequestWithEmptyName(t *testing.T) {
	request := &UserRequest{
		Email:    "user@example.com",
		Password: "password123",
		Name:     "", // nome vazio
	}

	user := request.FromUserRequest()

	if user.Email != request.Email {
		t.Errorf("Email esperado %s, mas foi %s", request.Email, user.Email)
	}

	if user.Password != request.Password {
		t.Errorf("Password esperado %s, mas foi %s", request.Password, user.Password)
	}

	if user.Name != "" {
		t.Errorf("Name deveria ser vazio, mas foi %s", user.Name)
	}

	// Verifica se o role padrão foi definido
	if len(user.Roles) != 1 || user.Roles[0] != "user" {
		t.Errorf("Role padrão deveria ser 'user', mas foi %v", user.Roles)
	}
}

func TestUserStructFields(t *testing.T) {
	user := &User{
		ID:        "test-id",
		Email:     "test@example.com",
		Password:  "password",
		Name:      "Test User",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Testa se os campos estão acessíveis
	if user.ID == "" {
		t.Error("ID não deveria estar vazio")
	}

	if user.Email == "" {
		t.Error("Email não deveria estar vazio")
	}

	if user.Password == "" {
		t.Error("Password não deveria estar vazio")
	}

	if user.Name == "" {
		t.Error("Name não deveria estar vazio")
	}

	if len(user.Roles) == 0 {
		t.Error("Roles não deveria estar vazio")
	}

	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt não deveria estar zerado")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt não deveria estar zerado")
	}
}

func TestUserResponseStructFields(t *testing.T) {
	response := &UserResponse{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Testa se os campos estão acessíveis
	if response.ID == "" {
		t.Error("ID não deveria estar vazio")
	}

	if response.Email == "" {
		t.Error("Email não deveria estar vazio")
	}

	if response.Name == "" {
		t.Error("Name não deveria estar vazio")
	}

	if len(response.Roles) == 0 {
		t.Error("Roles não deveria estar vazio")
	}

	if response.CreatedAt.IsZero() {
		t.Error("CreatedAt não deveria estar zerado")
	}

	if response.UpdatedAt.IsZero() {
		t.Error("UpdatedAt não deveria estar zerado")
	}
}

func TestUserRequestStructFields(t *testing.T) {
	request := &UserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	// Testa se os campos estão acessíveis
	if request.Email == "" {
		t.Error("Email não deveria estar vazio")
	}

	if request.Password == "" {
		t.Error("Password não deveria estar vazio")
	}

	if request.Name == "" {
		t.Error("Name não deveria estar vazio")
	}
}
