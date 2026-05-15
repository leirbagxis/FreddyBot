# Plano: diagnosticar-fluxo-de-dados-dashboard_2026-05-14_11-40.md

## Pedido do usuário
A dashboard não lista os canais e não entra nos detalhes do canal, mesmo com dados presentes no `freddybot.db`.

## Objetivo
Rastrear o `userID` desde a autenticação até a query no banco de dados para entender por que os canais (que existem no banco) não estão chegando ao Dashboard.

## Contexto atual
- O banco `freddybot.db` possui o usuário `7595607953` e o canal `-1002676384505`.
- O Dashboard reporta que não encontra canais.
- Redis é usado para cache e pode estar servindo dados antigos ou vazios.

## Arquivos analisados
- `internal/api/auth/middleware.go`
- `internal/api/controllers/userController.go`
- `internal/database/repositories/channel.go`

## Arquivos que poderão ser modificados
- `internal/api/auth/middleware.go`
- `internal/api/controllers/userController.go`
- `internal/database/repositories/channel.go`

## Estratégia de implementação
1.  **Log de Autenticação**: Confirmar qual ID está sendo extraído do token JWT.
2.  **Log de Controller**: Confirmar qual ID está sendo usado para buscar canais na API.
3.  **Log de Repositório**: Confirmar o resultado da query SQL e se há interferência de cache.

## Passos detalhados

1.  **Modificar `internal/api/auth/middleware.go`**
    - Adicionar log no `AuthMiddlewareJWT`: `logger.Bot("🔑 JWT: Usuário autenticado ID=%d", claims.UserID)`.

2.  **Modificar `internal/api/controllers/userController.go`**
    - Adicionar log no `GetUserChannelsController`: `logger.Bot("📊 API: Buscando canais para UserID=%d", userID)`.

3.  **Modificar `internal/database/repositories/channel.go`**
    - Adicionar log no `GetAllChannelsByUserID`: `logger.Bot("🗄️ DB: Query GetAllChannelsByUserID(%d) retornou %d canais", userID, len(channel))`.

## Riscos
- **Nulo**: Apenas logs.

## Como testar
1. Reiniciar o bot.
2. Acessar o Dashboard.
3. Verificar no terminal os logs na ordem:
    - `🔑 JWT: Usuário autenticado ID=...`
    - `📊 API: Buscando canais para UserID=...`
    - `🗄️ DB: Query GetAllChannelsByUserID(...) retornou ... canais`

## Rollback
`git checkout internal/api/auth/middleware.go internal/api/controllers/userController.go internal/database/repositories/channel.go`
