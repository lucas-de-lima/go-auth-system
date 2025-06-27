package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-de-lima/go-auth-system/internal/auth"
	"github.com/lucas-de-lima/go-auth-system/internal/domain"
	"github.com/stretchr/testify/assert"
)

func getJWT() *auth.JWTService {
	return auth.NewJWTService("secret", 1, "refresh", 1)
}

func TestContainsRole(t *testing.T) {
	assert.True(t, containsRole([]string{"admin", "user"}, "admin"))
	assert.False(t, containsRole([]string{"user"}, "admin"))
}

func TestGinAuthenticate_SuccessAndFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := getJWT()
	user := &domain.User{ID: "1", Email: "a@b.com", Roles: []string{"admin"}}
	token, _ := jwtService.GenerateToken(user)
	mw := NewAuthMiddleware(jwtService)
	r := gin.New()
	r.GET("/protected", mw.GinAuthenticate(), func(c *gin.Context) {
		c.String(200, "ok")
	})
	// Sucesso
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Falha: sem token
	req2 := httptest.NewRequest("GET", "/protected", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 401, w2.Code)
	// Falha: token inválido
	req3 := httptest.NewRequest("GET", "/protected", nil)
	req3.Header.Set("Authorization", "Bearer tokeninvalido")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 401, w3.Code)
}

func TestGinRequireRole_SuccessAndFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := getJWT()
	user := &domain.User{ID: "1", Email: "a@b.com", Roles: []string{"admin"}}
	token, _ := jwtService.GenerateToken(user)
	mw := NewAuthMiddleware(jwtService)
	r := gin.New()
	r.GET("/admin", mw.GinAuthenticate(), mw.GinRequireRole("admin"), func(c *gin.Context) {
		c.String(200, "ok")
	})
	// Sucesso
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Falha: sem role
	user2 := &domain.User{ID: "2", Email: "b@b.com", Roles: []string{"user"}}
	token2, _ := jwtService.GenerateToken(user2)
	req2 := httptest.NewRequest("GET", "/admin", nil)
	req2.Header.Set("Authorization", "Bearer "+token2)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 403, w2.Code)
}

func TestAuthenticate_HTTP_SuccessAndFail(t *testing.T) {
	jwtService := getJWT()
	user := &domain.User{ID: "1", Email: "a@b.com", Roles: []string{"admin"}}
	token, _ := jwtService.GenerateToken(user)
	mw := NewAuthMiddleware(jwtService)
	// Handler que lê o contexto
	handler := mw.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(UserIDKey)
		if id == nil {
			w.WriteHeader(401)
			return
		}
		w.WriteHeader(200)
	}))
	// Sucesso
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Falha: sem token
	req2 := httptest.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, 401, w2.Code)
	// Falha: token inválido
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("Authorization", "Bearer tokeninvalido")
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)
	assert.Equal(t, 401, w3.Code)
}

func TestRequireRole_HTTP_SuccessAndFail(t *testing.T) {
	jwtService := getJWT()
	user := &domain.User{ID: "1", Email: "a@b.com", Roles: []string{"admin"}}
	token, _ := jwtService.GenerateToken(user)
	mw := NewAuthMiddleware(jwtService)
	// Handler que espera usuário autenticado
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	handler := mw.Authenticate(mw.RequireRole("admin", final))
	// Sucesso
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Falha: não autenticado
	req2 := httptest.NewRequest("GET", "/", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, 401, w2.Code)
}

func TestGinAuthenticate_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := getJWT()
	mw := NewAuthMiddleware(jwtService)
	r := gin.New()
	r.GET("/protected", mw.GinAuthenticate(), func(c *gin.Context) {
		c.String(200, "ok")
	})
	// Header sem 'Bearer'
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Token x")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	// Header só 'Bearer'
	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.Header.Set("Authorization", "Bearer")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)
}

func TestAuthenticate_InvalidHeaderFormat(t *testing.T) {
	jwtService := getJWT()
	mw := NewAuthMiddleware(jwtService)
	handler := mw.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	// Header sem 'Bearer'
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Token x")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	// Header só 'Bearer'
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Authorization", "Bearer")
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	assert.Equal(t, 400, w2.Code)
}
