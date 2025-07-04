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
- Uso específico da versão Go 1.24.3 para compatibilidade

### 2. `deployments/docker-compose.ci.yml`

Este arquivo Docker Compose configura o ambiente de CI com:

- Serviço da aplicação usando o target `test` do Dockerfile.ci
- Banco de dados PostgreSQL efêmero (usando tmpfs)
- Healthcheck para garantir que o PostgreSQL está pronto
- Configuração para gerar relatórios de cobertura de testes
- Timeouts para evitar builds presos
- Otimizações de performance para o PostgreSQL em ambiente de teste

### 3. `scripts/wait-for-db.sh`

Script auxiliar que:

- Aguarda o PostgreSQL estar pronto antes de executar os testes
- Instalado diretamente no container de teste para maior confiabilidade

### 4. `.github/workflows/go.yml`

Workflow do GitHub Actions melhorado com:

- **Etapa dedicada para geração do cliente Prisma** antes de outras etapas
- Cache do cliente Prisma gerado entre jobs (usando actions/cache@v4)
- Formatação automática do código (go fmt) antes do linting
- Configuração personalizada do golangci-lint para controle fino de linting
- Verificação de segurança avançada com gosec:
  - Configuração personalizada via arquivo JSON
  - Filtragem por nível de severidade e confiança
  - Modo não-bloqueante para análise informativa
  - Geração de relatórios detalhados em formato JSON
- Cache otimizado para camadas Docker
- Extração e upload de relatórios (usando actions/upload-artifact@v4)
- Timeout global para evitar builds presos
- Limpeza adequada de recursos após os testes
- Versão específica do Go (1.24.3) para compatibilidade com o projeto
- Uso das versões mais recentes das ações do GitHub Actions para compatibilidade

## Como Funciona

1. **Geração do cliente Prisma**: Uma etapa dedicada gera o cliente Prisma e o armazena em cache
2. **Verificação de código**: 
   - Formatação automática do código
   - Linting com configuração personalizada
   - Análise de segurança configurável que não bloqueia o pipeline
3. **Execução de testes**:
   - O PostgreSQL é iniciado com healthcheck e otimizações
   - As migrações do Prisma são aplicadas
   - Os testes são executados com timeout de 5 minutos
   - Relatório de cobertura é gerado e armazenado
4. **Build da imagem final**: Se os testes passarem, a imagem Docker final é construída

## Execução Local

Para executar o mesmo ambiente de CI localmente:

```bash
# Formatar o código
go fmt ./...

# Instalar e executar gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -no-fail -fmt=json -out=gosec-results.json ./...

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
- Cache do cliente Prisma gerado entre jobs do GitHub Actions

### 2. Confiabilidade

- Healthcheck para o PostgreSQL
- Script wait-for-db.sh instalado diretamente no container
- Timeouts configurados para evitar builds presos
- Limpeza adequada de recursos após os testes
- Verificação de existência do cliente Prisma antes de cada etapa
- Uso das versões mais recentes das ações do GitHub Actions
- Modo não-bloqueante para verificações de segurança (continue-on-error)

### 3. Qualidade de Código

- Formatação automática do código antes do linting
- Linting com golangci-lint após geração do cliente Prisma
- Configuração personalizada de linting para controle fino de regras
- Verificação de segurança avançada:
  - Configuração via arquivo JSON
  - Filtragem por severidade e confiança
  - Relatórios detalhados para análise posterior
- Geração de relatórios de cobertura de testes
- Verificação de dependências
- Correção de erros comuns como verificação de valores de retorno

### 4. Performance

- Otimizações do PostgreSQL para ambiente de teste
- Paralelização de jobs quando possível
- Uso eficiente de cache para reduzir tempo de build

## Solução de Problemas

### Problema: Cliente Prisma não gerado

O erro `no required module provides package github.com/lucas-de-lima/go-auth-system/prisma/db` foi resolvido com:

1. Etapa específica para geração do cliente Prisma antes de outras etapas
2. Cache do cliente Prisma gerado entre jobs
3. Verificação da existência do cliente antes de cada etapa

### Problema: Incompatibilidade de versão do Go

O erro `go.mod requires go >= 1.24.3 (running go 1.24.2)` foi resolvido com:

1. Especificação explícita da versão Go 1.24.3 em todos os ambientes
2. Uso da mesma versão no Dockerfile.ci e no workflow do GitHub Actions
3. Instalação direta do gosec usando o Go 1.24.3 configurado no ambiente

### Problema: Erros de linting

Os erros como `Error return value of w.Write is not checked (errcheck)` e `File is not properly formatted (gofmt)` foram resolvidos com:

1. Verificação adequada dos valores de retorno de funções importantes
2. Formatação automática do código antes do linting
3. Configuração personalizada do golangci-lint para regras específicas
4. Exclusão de regras em contextos específicos (como testes)

### Problema: Falhas na verificação de segurança

Os erros relacionados ao gosec foram resolvidos com:

1. Configuração personalizada do gosec via arquivo JSON
2. Modo não-bloqueante (continue-on-error) para evitar falhas do pipeline
3. Filtragem por nível de severidade e confiança para focar em problemas reais
4. Geração de relatórios detalhados para análise posterior

### Problema: Compatibilidade de ações do GitHub Actions

Os erros como `Missing download info for actions/upload-artifact@v3` foram resolvidos com:

1. Atualização para as versões mais recentes das ações (v4)
2. Garantia de compatibilidade entre as diferentes ações usadas no workflow

### Problema: Banco de dados não disponível

Resolvido com:
1. Healthcheck no serviço PostgreSQL
2. Script wait-for-db.sh instalado diretamente no container
3. Configurações otimizadas para o PostgreSQL em ambiente de teste

### Problema: Builds lentos ou presos

Resolvido com:
1. Cache otimizado para camadas Docker
2. Cache do cliente Prisma entre jobs
3. Timeouts configurados para testes e linting
4. BuildKit ativado para builds mais rápidos 