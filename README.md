# Sistema de Autenticação em Go

Este projeto implementa um sistema de autenticação completo usando Go, com arquitetura limpa, banco de dados PostgreSQL e Prisma como ORM.

## Funcionalidades

- Registro e autenticação de usuários
- Gerenciamento de sessões com JWT
- Proteção de rotas com middleware de autenticação
- Banco de dados PostgreSQL com Prisma ORM
- Arquitetura limpa e modular

## Pré-requisitos

- Go 1.24+
- Docker e Docker Compose
- Make (opcional, para usar comandos do Makefile)

## Estrutura do Projeto

```
project-root/
├── .github/                    # GitHub-specific files (workflows, templates)
├── cmd/                        # Main applications
│   └── api/                    # API server entrypoint
├── configs/                    # Configuration files
├── deployments/                # Deployment configurations (Docker, K8s)
├── docs/                       # Documentation files
│   └── swagger/                # OpenAPI/Swagger documentation
├── internal/                   # Private application code
│   ├── api/                    # API handlers
│   ├── auth/                   # Authentication logic
│   ├── config/                 # Internal configuration
│   ├── domain/                 # Domain models
│   ├── middleware/             # Custom middleware
│   ├── repository/             # Database repositories
│   └── service/                # Business logic services
├── pkg/                        # Public library code
├── prisma/                     # Prisma schema and client
├── scripts/                    # Build/migration scripts
├── test/                       # Additional test applications/test data
├── web/                        # Web assets
```

## Instalação e Execução

### Usando Docker

```bash
# Iniciar todos os serviços com Docker Compose
make docker-compose-up

# Parar todos os serviços
make docker-compose-down
```

### Desenvolvimento Local

```bash
# Instalar dependências
go mod tidy

# Gerar cliente Prisma
make prisma-generate

# Executar aplicação
make run

# Compilar aplicação
make build
```

## Banco de Dados

O projeto utiliza PostgreSQL com Prisma ORM. Para configurar o banco de dados:

1. Configure as variáveis de ambiente para conexão com o banco de dados
2. Gere o cliente Prisma com `make prisma-generate`

## Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```
# Servidor
SERVER_PORT=8080

# Banco de dados
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/auth_system?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_system

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_SECRET=your_refresh_secret
JWT_REFRESH_EXPIRATION_HOURS=168
```

## Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Faça commit das suas alterações (`git commit -m 'Adiciona nova feature'`)
4. Faça push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para mais detalhes. 