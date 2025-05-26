# Documentação CI/CD para Go Auth System

Este documento explica a implementação de CI/CD baseada em Docker para o projeto Go Auth System com Prisma ORM.

## Visão Geral

A solução de CI/CD implementada resolve o problema da geração do cliente Prisma no ambiente de CI, garantindo que o build e os testes sejam executados em um ambiente isolado e consistente.

## Arquivos Principais

### 1. `deployments/Dockerfile.ci`

Este Dockerfile é específico para o ambiente de CI e inclui:

- **Multi-stage build** otimizado para diferentes fases do pipeline
- Estágio específico para geração do cliente Prisma (melhor cache)
- Estágio de build para testes separado do build de produção
- Imagem final mínima para produção

### 2. `deployments/docker-compose.ci.yml`

Este arquivo Docker Compose configura o ambiente de CI com:

- Serviço da aplicação usando o target `test` do Dockerfile.ci
- Banco de dados PostgreSQL efêmero (usando tmpfs)
- Healthcheck para garantir que o PostgreSQL está pronto
- Configuração para gerar relatórios de cobertura de testes
- Timeouts para evitar builds presos

### 3. `scripts/wait-for-db.sh`

Script auxiliar que:

- Aguarda o PostgreSQL estar pronto antes de executar os testes
- Instalado diretamente no container de teste para maior confiabilidade

### 4. `.github/workflows/go.yml`

Workflow do GitHub Actions melhorado com:

- Etapa de linting com golangci-lint
- Verificação de segurança com gosec
- Cache otimizado para camadas Docker
- Extração e upload de relatórios de cobertura de testes
- Timeout global para evitar builds presos
- Limpeza adequada de recursos após os testes

## Como Funciona

1. **Verificação de código**: Linting e análise de segurança
2. **Checkout do código**: O GitHub Actions faz o checkout do código
3. **Preparação do ambiente**: Configura Docker Buildx com cache otimizado
4. **Execução de testes**:
   - O PostgreSQL é iniciado com healthcheck
   - As migrações do Prisma são aplicadas
   - Os testes são executados com timeout de 5 minutos
   - Relatório de cobertura é gerado e armazenado
5. **Build da imagem final**: Se os testes passarem, a imagem Docker final é construída

## Execução Local

Para executar o mesmo ambiente de CI localmente:

```bash
# Executar os testes
docker-compose -f deployments/docker-compose.ci.yml up --build

# Limpar recursos após os testes
docker-compose -f deployments/docker-compose.ci.yml down -v
```

## Otimizações Implementadas

### 1. Cache Eficiente

- Multi-stage build para melhor aproveitamento de cache
- Estágio separado para geração do cliente Prisma
- Cache de camadas Docker entre execuções do CI

### 2. Confiabilidade

- Healthcheck para o PostgreSQL
- Script wait-for-db.sh instalado diretamente no container
- Timeouts configurados para evitar builds presos
- Limpeza adequada de recursos após os testes

### 3. Qualidade de Código

- Linting com golangci-lint
- Verificação de segurança com gosec
- Geração de relatórios de cobertura de testes
- Verificação de dependências

## Solução de Problemas

### Problema: Cliente Prisma não gerado

O erro `no required module provides package github.com/lucas-de-lima/go-auth-system/prisma/db` foi resolvido com:

1. Estágio específico para geração do cliente Prisma
2. Cópia do cliente gerado para os estágios subsequentes

### Problema: Banco de dados não disponível

Resolvido com:
1. Healthcheck no serviço PostgreSQL
2. Script wait-for-db.sh instalado diretamente no container

### Problema: Builds lentos ou presos

Resolvido com:
1. Cache otimizado para camadas Docker
2. Timeouts configurados para testes e linting
3. BuildKit ativado para builds mais rápidos 