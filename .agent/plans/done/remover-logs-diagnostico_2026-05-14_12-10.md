# Plano: remover-logs-diagnostico_2026-05-14_12-10.md

## Pedido do usuário
Remover todos os logs de diagnóstico que foram adicionados durante a depuração.

## Objetivo
Limpar o código removendo logs do `logger.Bot` e middlewares de debug, mantendo a base de código limpa e funcional.

## Contexto atual
- Logs foram adicionados em múltiplos arquivos para rastrear o fluxo do `userID`.
- O problema foi resolvido e os logs não são mais necessários.

## Arquivos analisados
- `internal/telegram/client.go`
- `internal/api/auth/middleware.go`
- `internal/api/controllers/userController.go`
- `internal/api/controllers/channelController.go`
- `internal/database/repositories/channel.go`

## Arquivos que poderão ser modificados
- `internal/telegram/client.go`
- `internal/api/auth/middleware.go`
- `internal/api/controllers/userController.go`
- `internal/api/controllers/channelController.go`
- `internal/database/repositories/channel.go`

## Estratégia de implementação
Remover cirurgicamente cada linha de log adicionada e reverter as mudanças nos middlewares de debug e imports desnecessários.

## Passos detalhados

1.  **Modificar `internal/telegram/client.go`**
    - Remover `debugUpdateMiddleware`.
    - Remover `tgbotModels` dos imports.
    - Reverter `opts` para não usar o middleware de debug.

2.  **Modificar `internal/api/auth/middleware.go`**
    - Remover log em `AuthMiddlewareJWT`.
    - Remover log em `AuthorizeChannel`.
    - Remover import de `logger`.

3.  **Modificar `internal/api/controllers/userController.go`**
    - Remover log em `GetUserChannelsController`.

4.  **Modificar `internal/api/controllers/channelController.go`**
    - Remover logs em `GetChannelByIDController`.

5.  **Modificar `internal/database/repositories/channel.go`**
    - Remover logs em `GetChannelByID` e `GetAllChannelsByUserID`.

## Riscos
- **Baixo**: Apenas remoção de logs.

## Como testar
1. O bot deve compilar e rodar sem erros.
2. O terminal não deve mais mostrar as mensagens de `🔑 JWT`, `📊 API`, etc.

## Rollback
`git checkout ...` (todos os arquivos acima)
