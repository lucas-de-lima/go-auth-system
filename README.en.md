🌐 [Versão em Português (BR)](README.md)

# 🔐 Authentication and Authorization System in Go

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Test Coverage](https://img.shields.io/badge/Test%20Coverage-85%25+%20Critical-brightgreen.svg)]()

A Go-based authentication and authorization system with a layered architecture, PostgreSQL database, and Prisma ORM. Ideal for applications that require secure authentication and access control.

## 🚀 Main Features

### 🔑 Authentication
- **User registration** with data validation
- **Secure login** with password hashing (bcrypt)
- **JWT (JSON Web Tokens)** for stateless authentication
- **Refresh tokens** for automatic session renewal
- **Logout** with token invalidation

### 🛡️ Authorization
- **Role-based access control** (admin/user)
- **Authentication middleware** for route protection
- **Authorization middleware** based on roles

### 👥 User Management
- **Full CRUD** for users (admin only)
- **User administration** (admin only)
- **Input data validation**

### 🔧 Technical Highlights
- **Layered architecture** (Controllers, Services, Repositories)
- **Comprehensive tests** (85%+ coverage in critical modules)
- **Logging system** for auditing
- **Standardized error handling**
- **Input validation**
- **Automatic panic recovery**

## 📋 Requirements

- **Go 1.24+**
- **Docker and Docker Compose**
- **Make** (optional, for Makefile commands)
- **PostgreSQL** (via Docker or local)

## 🏗️ Project Structure

```
go-auth-system/
├── cmd/api/                 # Application entry point
├── configs/                 # Configuration files
├── deployments/             # Deployment configs (Docker)
├── docs/                    # Documentation
├── internal/                # Application code
│   ├── api/                 # HTTP handlers
│   ├── auth/                # Authentication logic (JWT)
│   ├── config/              # Internal config
│   ├── controller/          # Controllers (MVC)
│   ├── domain/              # Domain models
│   ├── middleware/          # Custom middlewares
│   ├── repository/          # Persistence layer
│   ├── routes/              # Route definitions
│   └── service/             # Business logic
├── pkg/                     # Public libraries
│   ├── errors/              # Error handling
│   ├── logging/             # Logging system
│   └── validator/           # Data validation
├── prisma/                  # Prisma schema and client
├── scripts/                 # Build/migration scripts
├── test/                    # Integration tests
└── web/                     # Web assets
```

## ⚡ Installation & Usage

### 🐳 Using Docker (Recommended)

```bash
git clone https://github.com/lucas-de-lima/go-auth-system.git
cd go-auth-system
make docker-compose-up
# API available at: http://localhost:8080
```

### 💻 Local Development

```bash
go mod tidy
cp configs/app.env.example configs/app.env
# Edit configs/app.env with your settings
make prisma-setup
make run
make test
```

## 🔧 Useful Commands

<details>
<summary><strong>📋 See all available commands</strong></summary>

### 🏗️ Build & Run
```bash
make build
make run
make clean
```

### 🐳 Docker
```bash
make docker-build
make docker-run
make docker-compose-up
make docker-compose-down
```

### 🗄️ Database (Prisma)
```bash
make prisma-generate
make prisma-db-push
make prisma-studio
make prisma-setup
```

### 🧪 Tests & Quality
```bash
make test
make lint
make tidy
```

</details>

## ⚙️ Configuration

### 🔐 Environment Variables

Create a `configs/app.env` file at the project root:

```env
# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=5
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=120

# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/auth_system?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_system
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_SECRET=your_super_secret_refresh_key_here
JWT_REFRESH_EXPIRATION_HOURS=168

# Default Admin
DEFAULT_ADMIN_EMAIL=admin@admin.com
DEFAULT_ADMIN_PASSWORD=Admin123!@#
```

## 📡 REST API - Complete Documentation

*(See the Portuguese README for full API documentation, or translate as needed.)*

## 🔒 Security

### 🛡️ Security Features

- **Password hashing** with bcrypt (cost 12)
- **JWT with configurable expiration**
- **Refresh tokens** for secure renewal
- **Token blacklist** for logout
- **Input validation**
- **Audit logs** for all operations
- **Robust authentication middleware**
- **Role-based access control**

### 🔐 Authentication & Authorization

```go
// Middleware usage example
router.Use(authMiddleware.GinAuthenticate())           // Requires authentication
router.Use(authMiddleware.GinRequireRole("admin"))     // Requires admin role
```

### 📊 Test Coverage

- **85%+ coverage** in critical modules (auth, config, domain, middleware, service)
- **Comprehensive unit tests**
- **Integration tests**

### Coverage Report by Package
```bash
go test ./internal/auth -cover      # 92.0%
go test ./internal/config -cover    # 100.0%
go test ./internal/domain -cover    # 100.0%
go test ./internal/middleware -cover # 97.3%
go test ./internal/service -cover   # 85.1%
go test ./pkg/errors -cover         # 87.8%
go test ./pkg/validator -cover      # 95.5%
```

## 🧪 Tests

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests
```bash
go test ./test -v
```

## 📝 Logging System

The system uses logging with different levels for auditing:

- **INFO** - Normal operations (login, registration, etc.)
- **WARNING** - Situations that require attention
- **ERROR** - Errors that need investigation

Example logs:
```
INFO: [192.168.1.1] Login successful: user@example.com
WARNING: [192.168.1.1] Login attempt failed for: user@example.com
ERROR: [192.168.1.1] Failed to register user: database error
```

### Log Configuration
```go
// Default config
logging.SetupLogger(logging.DefaultConfig())

// Custom config
config := logging.Config{
    InfoWriter:    os.Stdout,
    WarningWriter: os.Stdout,
    ErrorWriter:   os.Stderr,
    Prefix:        "[AUTH-SYSTEM] ",
    Flag:          log.LstdFlags | log.Lshortfile,
}
logging.SetupLogger(config)
```

## 🚀 Deployment & Infrastructure

### 🐳 Docker (Recommended)
```bash
make docker-build
make docker-run
```

### 🐳 Docker Compose (Development)
```bash
make docker-compose-up
make docker-compose-down
```

### 🐳 Docker Compose (CI/CD)
```bash
docker-compose -f deployments/docker-compose.ci.yml up -d
```

## 🛠️ Technologies Used

- **Go 1.24+** - Main language
- **Gin** - Web framework
- **Prisma** - ORM for PostgreSQL
- **JWT-Go** - JWT implementation
- **bcrypt** - Password hashing
- **PostgreSQL** - Database
- **Docker** - Containerization
- **Testify** - Testing framework

## 🤝 Contributing

1. **Fork** the project
2. **Create** a branch for your feature (`git checkout -b feature/new-feature`)
3. **Commit** your changes (`git commit -m 'Add new feature'`)
4. **Push** to the branch (`git push origin feature/new-feature`)
5. **Open** a Pull Request

### 📋 Contribution Checklist

- [ ] Code follows project standards
- [ ] Tests have been added/updated
- [ ] Documentation has been updated
- [ ] No breaking changes

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Prisma](https://www.prisma.io/) - Modern ORM
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Testify](https://github.com/stretchr/testify) - Testing framework
- **Cursor** - The AI-powered IDE that made this project easier to build

---

**Developed with ❤️ Go and AI** 