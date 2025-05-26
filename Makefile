.PHONY: build run test clean docker-build docker-run migrate-up migrate-down prisma-generate prisma-db-push prisma-migrate prisma-studio

# Variáveis
APP_NAME := go-auth-system
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/api
DOCKER_IMAGE := go-auth-system:latest

# Comandos para desenvolvimento
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)/main.go

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean

# Comandos para Docker
docker-build:
	docker build -t $(DOCKER_IMAGE) -f deployments/Dockerfile .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

docker-compose-up:
	docker-compose -f deployments/docker-compose.yml up

docker-compose-down:
	docker-compose -f deployments/docker-compose.yml down

# Comandos para Prisma
prisma-generate:
	cd prisma && go run cmd/run_prisma.go generate

prisma-db-push:
	cd prisma && go run cmd/run_prisma.go db push

prisma-migrate-dev:
	cd prisma && go run cmd/run_prisma.go migrate dev --name $(name)

prisma-migrate-deploy:
	cd prisma && go run cmd/run_prisma.go migrate deploy

prisma-studio:
	cd prisma && go run cmd/run_prisma.go studio

prisma-setup: prisma-generate prisma-db-push

# Outros comandos úteis
lint:
	golangci-lint run ./...

tidy:
	go mod tidy

help:
	@echo "Comandos disponíveis:"
	@echo "  make build              - Compila a aplicação"
	@echo "  make run                - Executa a aplicação"
	@echo "  make test               - Executa os testes"
	@echo "  make clean              - Remove arquivos de build"
	@echo "  make docker-build       - Constrói a imagem Docker"
	@echo "  make docker-run         - Executa a aplicação em um container Docker"
	@echo "  make docker-compose-up  - Inicia todos os serviços com Docker Compose"
	@echo "  make docker-compose-down- Para todos os serviços do Docker Compose"
	@echo "  make prisma-generate    - Gera código do cliente Prisma"
	@echo "  make prisma-db-push     - Atualiza o banco de dados conforme o schema (desenvolvimento)"
	@echo "  make prisma-migrate-dev - Cria uma nova migração (use make prisma-migrate-dev name=nome_da_migracao)"
	@echo "  make prisma-migrate-deploy - Aplica migrações em ambientes de produção"
	@echo "  make prisma-studio      - Abre o Prisma Studio para visualizar o banco de dados"
	@echo "  make prisma-setup       - Configura o banco de dados (gera o cliente e cria tabelas)"
	@echo "  make lint               - Executa linter"
	@echo "  make tidy               - Atualiza dependências" 