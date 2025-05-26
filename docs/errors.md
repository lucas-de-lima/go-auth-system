# Sistema de Tratamento de Erros

Este documento descreve o sistema de tratamento de erros implementado para nossa API Go, que segue as melhores práticas da linguagem e fornece uma experiência consistente para desenvolvedores e usuários da API.

## Visão Geral

O sistema de erros é baseado em um tipo `AppError` personalizado que encapsula:
- Código de status HTTP
- Mensagem amigável para o cliente
- Erro original para logging/debugging

Este sistema foi projetado seguindo estes princípios:
- **Simplicidade**: Evitando abstrações desnecessárias
- **Clareza**: Interfaces autoexplicativas
- **Idiomático**: Seguindo convenções do Go
- **Prático**: Focado em facilitar o desenvolvimento diário

## Estrutura do Sistema

### 1. Pacote de Erros (`pkg/errors`)

O pacote de erros fornece:

- Tipo `AppError` principal
- Catálogo de erros comuns
- Funções auxiliares para tratamento de erros
- Integração com o pacote `errors` do Go 1.13+
- Adaptadores para HTTP padrão e framework Gin

#### Tipo AppError

```go
type AppError struct {
    Code     int    // Código de status HTTP
    Message  string // Mensagem amigável para o cliente
    Internal error  // Erro original para logging
}
```

Este tipo implementa:
- Interface `error` via método `Error()`
- Interface `Unwrapper` via método `Unwrap()`
- Método `Is()` para compatibilidade com `errors.Is()`

#### Catálogo de Erros

O sistema define um conjunto de erros comuns prontos para uso:

```go
var (
    ErrInternalServer = AppError{...}
    ErrBadRequest = AppError{...}
    ErrUnauthorized = AppError{...}
    // E outros...
)
```

### 2. Utilitários para API HTTP

#### Para HTTP Padrão

O sistema fornece funções para responder com HTTP padrão:

- `HandleError(w, err)`: Processa erros e envia resposta apropriada
- `RespondWithError(w, code, message)`: Responde com erro simples
- `RespondWithJSON(w, code, payload)`: Responde com qualquer payload JSON
- `WithRecovery`: Middleware para recuperação de pânicos

#### Para Gin Framework

O sistema também fornece adaptadores para o framework Gin:

- `GinHandleError(c, err)`: Processa erros no contexto do Gin
- `GinRespondWithError(c, code, message)`: Responde com erro no formato Gin
- `GinRespondWithJSON(c, code, payload)`: Responde com qualquer payload JSON
- `GinMiddlewareRecovery()`: Middleware de recuperação para Gin

### 3. Tratamento de Erros de Validação

O sistema possui recursos específicos para erros de validação:

- Tipo `ValidationDetail` para descrever erros em campos específicos
- Função `NewValidationError` para criar erros de validação estruturados
- Função `GetValidationDetails` para extrair detalhes de validação
- Função `GinValidationResponse` para formatar resposta para Gin

## Como Usar

### Em Serviços

O uso em serviços é o mesmo independente do framework HTTP:

```go
// Exemplo de uso em um serviço
func (s *UserService) GetByID(id string) (*User, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        return nil, errors.ErrInternalServer.WithError(err)
    }
    
    if user == nil {
        return nil, errors.ErrUserNotFound
    }
    
    return user, nil
}
```

### Em Handlers HTTP Padrão

```go
// Exemplo de uso em um handler HTTP padrão
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := h.service.GetByID(id)
    if err != nil {
        errors.HandleError(w, err)
        return
    }
    
    errors.RespondWithJSON(w, http.StatusOK, user)
}
```

### Em Controllers Gin

```go
// Exemplo de uso em um controller Gin
func (uc *UserController) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := uc.userService.GetByID(id)
    if err != nil {
        errors.GinHandleError(c, err)
        return
    }
    
    errors.GinRespondWithJSON(c, http.StatusOK, user)
}
```

### Middleware de Autenticação para Gin

```go
// Uso do middleware de autenticação com Gin
func SetupRoutes(router *gin.Engine, jwtService *auth.JWTService) {
    // Criar o middleware de autenticação
    authMiddleware := middleware.NewAuthMiddleware(jwtService)
    
    // Rotas públicas
    publicRoutes := router.Group("/api")
    {
        publicRoutes.POST("/auth/login", authController.Login)
    }
    
    // Rotas protegidas
    protectedRoutes := router.Group("/api")
    protectedRoutes.Use(authMiddleware.GinAuthenticate())
    {
        protectedRoutes.GET("/users/me", userController.GetCurrentUser)
    }
}
```

### Criando Erros Personalizados

```go
// Criando um erro personalizado
err := errors.NewAppError(http.StatusBadRequest, "Formato inválido", originalErr)

// Personalizando um erro existente
err := errors.ErrNotFound.WithMessage("Produto não encontrado")
err := errors.ErrBadRequest.WithError(validationErr)
```

### Verificando Tipos de Erro

```go
// Verificando o tipo de erro com errors.Is
if errors.Is(err, errors.ErrUserNotFound) {
    // Trata erro de usuário não encontrado
}

// Extraindo um AppError com errors.As
var appErr errors.AppError
if errors.As(err, &appErr) {
    // Usa informações do appErr
}
```

## Formato da Resposta JSON

As respostas de erro seguem um formato consistente:

```json
{
    "message": "Mensagem descritiva do erro",
    "details": {  // Campo opcional para erros de validação
        "fields": {
            "email": "Email inválido",
            "password": "Senha muito curta"
        }
    }
}
```

## Configuração do Gin

Para usar corretamente o sistema de erros com Gin, configure o router da seguinte maneira:

```go
// Inicializar o router do Gin com nosso middleware personalizado
router := gin.New()
router.Use(gin.Logger())
router.Use(errors.GinMiddlewareRecovery())
```

## Vantagens do Sistema

1. **Consistência**: Todas as respostas de erro seguem o mesmo formato
2. **Rastreabilidade**: Erros internos são registrados, mas não expostos aos clientes
3. **Compatibilidade**: Funciona com HTTP padrão e Gin Framework
4. **Extensibilidade**: Fácil de estender para novos tipos de erro
5. **Idiomático**: Usa padrões do Go, como interfaces e wrappers

## Considerações para o Futuro

- Adição de suporte para internacionalização de mensagens de erro
- Integração com sistemas de monitoramento para rastreamento de erros
- Expansão dos tipos de erro específicos do domínio
- Adaptadores para outros frameworks web populares 