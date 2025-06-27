# 🔐 Sistema de Autenticação e Autorização em Go

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Test%20Coverage-85%25+%20Critical-brightgreen.svg)]()

Um sistema de autenticação e autorização desenvolvido em Go, com arquitetura estruturada, banco de dados PostgreSQL e Prisma ORM. Ideal para aplicações que necessitam de autenticação segura e controle de acesso.

## 🚀 Funcionalidades Principais

### 🔑 Autenticação
- **Registro de usuários** com validação de dados
- **Login seguro** com hash de senhas (bcrypt)
- **JWT (JSON Web Tokens)** para autenticação stateless
- **Refresh tokens** para renovação automática de sessões
- **Logout** com invalidação de tokens

### 🛡️ Autorização
- **Sistema de roles** (papéis) para controle de acesso
- **Middleware de autenticação** para proteção de rotas
- **Middleware de autorização** baseado em roles
- **Controle de acesso** por role (admin/user)

### 👥 Gerenciamento de Usuários
- **CRUD completo** de usuários (via admin)
- **Administração** de usuários (apenas admins)
- **Validação de dados** de entrada

### 🔧 Recursos Técnicos
- **Arquitetura estruturada** (Controllers, Services, Repositories)
- **Testes abrangentes** (85%+ de cobertura nos componentes críticos)
- **Sistema de logs** para auditoria
- **Tratamento de erros** padronizado
- **Validação de dados** de entrada
- **Recovery de panics** automático

## 📋 Pré-requisitos

- **Go 1.24+**
- **Docker e Docker Compose**
- **Make** (opcional, para usar comandos do Makefile)
- **PostgreSQL** (via Docker ou local)

## 🏗️ Arquitetura do Projeto

```
go-auth-system/
├── 📁 cmd/api/                 # Ponto de entrada da aplicação
├── 📁 configs/                 # Arquivos de configuração
├── 📁 deployments/             # Configurações de deploy (Docker)
├── 📁 docs/                    # Documentação
├── 📁 internal/                # Código interno da aplicação
│   ├── 📁 api/                 # Handlers HTTP
│   ├── 📁 auth/                # Lógica de autenticação (JWT)
│   ├── 📁 config/              # Configuração interna
│   ├── 📁 controller/          # Controladores (MVC)
│   ├── 📁 domain/              # Modelos de domínio
│   ├── 📁 middleware/          # Middlewares customizados
│   ├── 📁 repository/          # Camada de persistência
│   ├── 📁 routes/              # Definição de rotas
│   └── 📁 service/             # Lógica de negócio
├── 📁 pkg/                     # Bibliotecas públicas
│   ├── 📁 errors/              # Tratamento de erros
│   ├── 📁 logging/             # Sistema de logs
│   └── 📁 validator/           # Validação de dados
├── 📁 prisma/                  # Schema e cliente Prisma
├── 📁 scripts/                 # Scripts de build/migração
├── 📁 test/                    # Testes de integração
└── 📁 web/                     # Assets web
```

## ⚡ Instalação e Execução

### 🐳 Usando Docker (Recomendado)

```bash
# Clonar o repositório
git clone https://github.com/lucas-de-lima/go-auth-system.git
cd go-auth-system

# Iniciar todos os serviços (API + PostgreSQL)
make docker-compose-up

# A API estará disponível em: http://localhost:8080
```

### 💻 Desenvolvimento Local

```bash
# Instalar dependências
go mod tidy

# Configurar variáveis de ambiente
cp configs/app.env.example configs/app.env
# Editar configs/app.env com suas configurações

# Gerar cliente Prisma e configurar banco
make prisma-setup

# Executar aplicação
make run

# Executar testes
make test
```

## 🔧 Comandos Úteis

<details>
<summary><strong>📋 Ver todos os comandos disponíveis</strong></summary>

### 🏗️ Build e Execução
```bash
make build              # Compila a aplicação
make run                # Executa a aplicação
make clean              # Remove arquivos de build
```

### 🐳 Docker
```bash
make docker-build       # Constrói a imagem Docker
make docker-run         # Executa em container Docker
make docker-compose-up  # Inicia todos os serviços
make docker-compose-down # Para todos os serviços
```

### 🗄️ Banco de Dados (Prisma)
```bash
make prisma-generate    # Gera código do cliente Prisma
make prisma-db-push     # Atualiza schema do banco
make prisma-studio      # Abre interface visual do banco
make prisma-setup       # Setup completo do banco
```

### 🧪 Testes e Qualidade
```bash
make test               # Executa todos os testes
make lint               # Executa linter
make tidy               # Atualiza dependências
```

</details>

## ⚙️ Configuração

### 🔐 Variáveis de Ambiente

Crie um arquivo `configs/app.env` na raiz do projeto:

```env
# 🖥️ Servidor
SERVER_PORT=8080
SERVER_READ_TIMEOUT=5
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=120

# 🗄️ Banco de dados
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/auth_system?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_system
DB_SSLMODE=disable

# 🔑 JWT
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_SECRET=your_super_secret_refresh_key_here
JWT_REFRESH_EXPIRATION_HOURS=168

# 👨‍💼 Admin Padrão
DEFAULT_ADMIN_EMAIL=admin@admin.com
DEFAULT_ADMIN_PASSWORD=Admin123!@#
```

## 📡 API REST - Documentação Completa

### 🔗 Base URL
```
http://localhost:8080
```

<details>
<summary><strong>🔐 Autenticação - Rotas Públicas</strong></summary>

### 📝 Registro de Usuário
**POST** `/users/register`

Registra um novo usuário no sistema.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123",
  "name": "Nome do Usuário"
}
```

**Validações:**
- Email: obrigatório e formato válido
- Senha: obrigatória, mínimo 3 caracteres
- Nome: opcional

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "usuario@exemplo.com",
  "name": "Nome do Usuário",
  "roles": ["user"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Erros possíveis:**
- `400` - Dados inválidos (email já existe, campos obrigatórios faltando)
- `500` - Erro interno do servidor

---

### 🔑 Login
**POST** `/users/login`

Autentica um usuário e retorna tokens de acesso.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Erros possíveis:**
- `401` - Credenciais inválidas
- `500` - Erro interno do servidor

---

### 🔄 Refresh Token
**POST** `/users/refresh`

Renova os tokens de acesso usando um refresh token válido.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Erros possíveis:**
- `401` - Refresh token inválido ou expirado
- `500` - Erro interno do servidor

</details>

<details>
<summary><strong>👤 Usuários - Rotas Protegidas</strong></summary>

### 🚪 Logout
**POST** `/users/logout`

**Headers necessários:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
```json
{
  "message": "Logout realizado com sucesso"
}
```

**Erros possíveis:**
- `401` - Token de acesso inválido
- `400` - Refresh token não fornecido

</details>

<details>
<summary><strong>👨‍💼 Administração - Rotas de Admin</strong></summary>

> ⚠️ **Atenção:** Estas rotas requerem autenticação e role `admin`.

### 📋 Listar Todos os Usuários
**GET** `/admin/users`

**Headers necessários:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "usuario@exemplo.com",
    "name": "Nome do Usuário",
    "roles": ["user"],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

**Erros possíveis:**
- `401` - Token de acesso inválido
- `403` - Acesso negado (role admin necessário)

---

### 👤 Obter Usuário por ID (Admin)
**GET** `/admin/users/:id`

**Headers necessários:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "usuario@exemplo.com",
  "name": "Nome do Usuário",
  "roles": ["user"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Erros possíveis:**
- `401` - Token de acesso inválido
- `403` - Acesso negado (role admin necessário)
- `404` - Usuário não encontrado

---

### ✏️ Atualizar Usuário (Admin)
**PUT** `/admin/users/:id`

**Headers necessários:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "email": "novo@exemplo.com",
  "name": "Novo Nome"
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "novo@exemplo.com",
  "name": "Novo Nome",
  "roles": ["user"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

**Erros possíveis:**
- `401` - Token de acesso inválido
- `403` - Acesso negado (role admin necessário)
- `404` - Usuário não encontrado
- `400` - Dados inválidos

---

### 🗑️ Deletar Usuário (Admin)
**DELETE** `/admin/users/:id`

**Headers necessários:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "message": "Usuário deletado com sucesso"
}
```

**Erros possíveis:**
- `401` - Token de acesso inválido
- `403` - Acesso negado (role admin necessário)
- `404` - Usuário não encontrado

</details>

## 🔒 Segurança

### 🛡️ Recursos de Segurança Implementados

- **Hash de senhas** com bcrypt (custo 12)
- **JWT com expiração** configurável
- **Refresh tokens** para renovação segura
- **Blacklist de tokens** para logout
- **Validação de entrada** de dados
- **Logs de auditoria** para todas as operações
- **Middleware de autenticação** robusto
- **Controle de acesso baseado em roles**

### 🔐 Autenticação e Autorização

```go
// Exemplo de uso do middleware
router.Use(authMiddleware.GinAuthenticate())           // Requer autenticação
router.Use(authMiddleware.GinRequireRole("admin"))     // Requer role específico
```

### 📊 Cobertura de Testes

- **85%+ de cobertura** nos componentes críticos (auth, config, domain, middleware, service)
- **Testes unitários** abrangentes para lógica de negócio
- **Testes de integração** completos

### Relatório de Cobertura por Pacote
```bash
# Cobertura dos pacotes críticos
go test ./internal/auth -cover      # 92.0%
go test ./internal/config -cover    # 100.0%
go test ./internal/domain -cover    # 100.0%
go test ./internal/middleware -cover # 97.3%
go test ./internal/service -cover   # 85.1%
go test ./pkg/errors -cover         # 87.8%
go test ./pkg/validator -cover      # 95.5%
```

## 🧪 Testes

### Executar Todos os Testes
```bash
make test
```

### Executar Testes com Cobertura
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Testes de Integração
```bash
go test ./test -v
```

## 📝 Sistema de Logs

O sistema utiliza logging com diferentes níveis para auditoria:

- **INFO** - Operações normais (login, registro, etc.)
- **WARNING** - Situações que merecem atenção
- **ERROR** - Erros que precisam de investigação

Exemplo de logs:
```
INFO: [192.168.1.1] Login realizado: usuario@exemplo.com
WARNING: [192.168.1.1] Tentativa de login falhou para: usuario@exemplo.com
ERROR: [192.168.1.1] Falha ao registrar usuário: erro de banco de dados
```

### Configuração de Logs
```go
// Configuração padrão
logging.SetupLogger(logging.DefaultConfig())

// Configuração customizada
config := logging.Config{
    InfoWriter:    os.Stdout,
    WarningWriter: os.Stdout,
    ErrorWriter:   os.Stderr,
    Prefix:        "[AUTH-SYSTEM] ",
    Flag:          log.LstdFlags | log.Lshortfile,
}
logging.SetupLogger(config)
```

## 🚀 Deploy e Infraestrutura

### 🐳 Docker (Recomendado)
```bash
# Build da imagem
make docker-build

# Executar container
make docker-run
```

### 🐳 Docker Compose (Desenvolvimento)
```bash
# Iniciar todos os serviços
make docker-compose-up

# Parar todos os serviços
make docker-compose-down
```

### 🐳 Docker Compose (CI/CD)
```bash
# Usar configuração de CI
docker-compose -f deployments/docker-compose.ci.yml up -d
```

## 🛠️ Tecnologias Utilizadas

- **Go 1.24+** - Linguagem principal
- **Gin** - Framework web
- **Prisma** - ORM para PostgreSQL
- **JWT-Go** - Implementação JWT
- **bcrypt** - Hash de senhas
- **PostgreSQL** - Banco de dados
- **Docker** - Containerização
- **Testify** - Framework de testes

## 🤝 Contribuição

1. **Fork** o projeto
2. **Crie** uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. **Commit** suas alterações (`git commit -m 'Adiciona nova feature'`)
4. **Push** para a branch (`git push origin feature/nova-feature`)
5. **Abra** um Pull Request

### 📋 Checklist para Contribuições

- [ ] Código segue os padrões do projeto
- [ ] Testes foram adicionados/atualizados
- [ ] Documentação foi atualizada
- [ ] Não há quebras de compatibilidade

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🙏 Agradecimentos

- [Gin](https://github.com/gin-gonic/gin) - Framework web
- [Prisma](https://www.prisma.io/) - ORM moderno
- [JWT-Go](https://github.com/golang-jwt/jwt) - Implementação JWT
- [Testify](https://github.com/stretchr/testify) - Framework de testes

---

**Desenvolvido com ❤️ em Go** 