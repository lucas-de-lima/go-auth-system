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
}

// UserRepository define as operações de persistência para usuários
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

// UserResponse representa a resposta de um usuário
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRequest representa a requisição de um usuário
type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

// Mapper functions
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *UserRequest) FromUserRequest() *User {
	return &User{
		Email:    u.Email,
		Password: u.Password,
	}
}

