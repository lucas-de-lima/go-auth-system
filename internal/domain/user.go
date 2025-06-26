package domain

import (
	"time"
)

// User representa o modelo de domínio para usuários
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // não expor senha nas respostas JSON
	Name      string    `json:"name,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserService define as operações disponíveis para usuários
type UserService interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	Authenticate(email, password string) (string, error) // retorna token JWT
	List() ([]*User, error)
}

// UserRepository define as operações de persistência para usuários
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	List() ([]*User, error)
}

// UserResponse representa a resposta de um usuário
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRequest representa a requisição de um usuário
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
	Name     string `json:"name,omitempty"`
}

// Mapper functions
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Roles:     u.Roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *UserRequest) FromUserRequest() *User {
	return &User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Roles:    []string{"user"}, // padrão: todo novo usuário é "user"
	}
}
