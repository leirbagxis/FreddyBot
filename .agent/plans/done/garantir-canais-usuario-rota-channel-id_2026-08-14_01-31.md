# Plano: Garantir Lista de Canais do Usuário na Rota /api/channel/:channelId

## Pedido do usuário
"ah e o total de canais do usuario nao esta vindo na rota do /dashboard/:channelID"

## Objetivo
Garantir que a resposta da API na rota `GET /api/channel/:channelId` traga sempre o array completo `user.channels` com todos os canais pertencentes ao proprietário, mesmo quando o objeto do canal for recuperado do cache do Redis/L1.

## Contexto atual
- Na rota `/api/channel/:channelId`, o controller `GetChannelByIDController` chama `ChannelService.GetChannelByID`.
- O `ChannelService` consulta primeiro o cache (L1/Redis). Como os objetos armazenados no cache `channel:v2:<id>` não persistem o relacionamento completo `Owner.Channels`, o campo `user.channels` retornava vazio na resposta JSON quando o canal vinha do cache.
- Como resultado, no frontend (`App.tsx`), `data.user.channels` ficava vazio, fazendo o SideMenu exibir 0 canais.

## Arquivos analisados
- `internal/api/controllers/channelController.go`
- `internal/core/services/channels.go`

## Arquivos que poderão ser modificados
- `internal/api/controllers/channelController.go`

## Estratégia de implementação

1. **Backend (`internal/api/controllers/channelController.go`)**:
   - Em `GetChannelByIDController`, após gerar o `userDTO` via `dto.ToUserDTO(channel.Owner)`:
   - Se `userDTO.Channels` estiver vazio, buscar os canais do usuário utilizando `c.container.ChannelService.GetUserChannels(ctx, channel.OwnerID)` e popular `userDTO.Channels`.
   - Dessa forma, independente de o canal vir do banco de dados ou do cache do Redis, o payload de resposta sempre incluirá a lista completa de canais do usuário.

2. **Validação**:
   - Compilar o backend Go com `go build ./cmd/FreddyBot`.
   - Compilar o frontend com `npm run build`.

## Riscos
- **Nenhum**: Solução resiliente que funciona com e sem cache de canais.

## Impactos esperados
- O campo `user.channels` no JSON retornado pela rota `/api/channel/:channelId` virá preenchido corretamente.
- O contador de canais no SideMenu do dashboard funcionará perfeitamente em `/dashboard/:channelID`.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
go build ./cmd/FreddyBot
cd dashboard && npm run build
```

## Rollback
```bash
git checkout internal/api/controllers/channelController.go
```
