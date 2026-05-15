# Plano: reestruturacao-arquitetural-v3_2026-05-14_13-00.md

## Pedido do usuário
Realizar uma limpeza profunda e melhoria arquitetural na API e nos serviços.

## Objetivo
Transformar a API em uma estrutura profissional e escalável, aplicando padrões de **DTO**, **Validação Centralizada**, **Tratamento Global de Erros** e **Desacoplamento de Cache**.

## Contexto atual
- **Vazamento de Modelos**: Structs do GORM são devolvidas diretamente via JSON.
- **Validação Manual**: `if len(name) > 64` espalhados nos serviços.
- **Cache no Repositório**: Repositórios invalidando cache, o que dificulta testes e quebra o SRP (Single Responsibility Principle).
- **Tratamento de Erro Repetitivo**: Controllers cheios de `if err != nil` manuais.
- **Inconsistência de Tipos**: Campos como `UserId` vs `user_id` e `userID` geram confusão.

## Estratégia de implementação

### 1. Camada de DTO (Data Transfer Objects)
- Criar `internal/api/dto` para definir structs de Request e Response.
- **Regra**: O Controller recebe um DTO, chama o Serviço, e o Serviço devolve um Modelo ou DTO. O Controller mapeia para o DTO de saída.
- Isso protege os campos internos do banco (`TokenVersion`, `CreatedAt`, etc).

### 2. Validação via Struct Tags
- Remover validações manuais dos serviços.
- Usar tags do Gin/Validator (`binding:"required,min=3,max=64"`) nas structs de Request.
- Adicionar validações customizadas (ex: `validateEmoji`) se necessário.

### 3. Centralização do Cache (Service Layer)
- Mover as chamadas de `InvalidateCache` dos repositórios para os Serviços.
- Repositórios devem ser "burros": apenas SQL.

### 4. Tratamento Global de Erros
- Criar `pkg/errors` com erros padrão: `ErrNotFound`, `ErrUnauthorized`, `ErrInternal`.
- Criar um Middleware no Gin que captura esses erros e gera o JSON de erro padronizado.

### 5. Unificação de Controladores e Serviços
- Mesclar `ButtonsCustomCaptionController` em `CustomCaptionController`.
- Mesclar lógicas de botões no `ButtonService` (evitar duplicação entre botões normais e customizados).

## Arquivos que serão modificados (Principais)
- `internal/api/routes/routes.go` (registro de novos controllers e middlewares)
- `internal/api/controllers/*` (refatoração para DTOs e Erros)
- `internal/core/services/*` (remoção de validações e inclusão de cache)
- `internal/database/repositories/*` (remoção de cache)
- `internal/api/dto/` (pasta nova)
- `pkg/errors/` (pasta nova)

## Passos detalhados (Fase 1: Estrutura Base)

1.  **Criar `pkg/errors/errors.go`**
    - Definir tipos de erro e funções auxiliares.
2.  **Criar DTOs para Canais e Usuários**
    - Criar `UserDTO`, `ChannelDTO`, `ButtonDTO` para evitar vazamento do GORM.
3.  **Refatorar `ButtonService` e `CustomCaptionService`**
    - Unificar lógica de cálculo de posição.
4.  **Refatorar Controllers**
    - Aplicar os DTOs e simplificar o tratamento de erro.

## Riscos
- **Compatibilidade com Frontend**: Os nomes dos campos JSON nos DTOs devem ser idênticos aos modelos atuais para não quebrar o React.

## Como testar
1. `go build ./...`
2. Testar fluxos de criação e listagem no Dashboard.
3. Forçar erros (ex: canal inexistente) e ver se o JSON de erro segue o novo padrão.

## Rollback
`git checkout .` (em blocos por arquivo)
