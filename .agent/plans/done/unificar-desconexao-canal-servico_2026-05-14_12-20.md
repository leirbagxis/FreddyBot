# Plano: unificar-desconexao-canal-servico_2026-05-14_12-20.md

## Pedido do usuário
Garantir que a API e o Bot usem os mesmos serviços e remover redundâncias (unificação do núcleo).

## Objetivo
Mover o método `DisconnectChannel` do `AppContainer` para o `ChannelService`. O container deve ser apenas para injeção de dependências, enquanto a lógica de negócio (enviar mensagem de adeus, sair do chat e deletar do banco) deve morar no serviço.

## Contexto atual
- `AppContainer` possui o método `DisconnectChannel` que contém lógica de negócio.
- `ChannelController` chama `c.container.DisconnectChannel`.
- O Bot (em `addChannel/callback.go` ou similar) também deve usar essa mesma lógica unificada.

## Arquivos analisados
- `internal/container/appContainer.go`
- `internal/core/services/channels.go`
- `internal/api/controllers/channelController.go`

## Arquivos que poderão ser modificados
- `internal/container/appContainer.go`
- `internal/core/services/channels.go`
- `internal/api/controllers/channelController.go`

## Estratégia de implementação
1.  **Ajustar `ChannelService`**: Adicionar o campo `Bot *bot.Bot` e `Cache *cache.Service` à struct (se necessário) e mover a lógica de desconexão para lá.
2.  **Ajustar `AppContainer`**: Remover o método `DisconnectChannel` e garantir que o `ChannelService` receba a instância do Bot após ele ser inicializado.
3.  **Ajustar `ChannelController`**: Mudar a chamada para `c.container.ChannelService.DisconnectChannel`.

## Passos detalhados

1.  **Modificar `internal/core/services/channels.go`**
    - Adicionar `Bot *bot.Bot` e `cache *cache.Service` na struct `ChannelService`.
    - Implementar `func (s *ChannelService) DisconnectChannel(ctx context.Context, userID int64, channelID int64) error`.
    - Integrar a mensagem de "adeus" e o comando de sair do canal nessa função.

2.  **Modificar `internal/container/appContainer.go`**
    - Remover a função `DisconnectChannel`.
    - Na criação do `ChannelService`, passar o `cacheService`.
    - No `StartBot` (onde o bot é injetado), atualizar a referência do Bot no `ChannelService`.

3.  **Modificar `internal/api/controllers/channelController.go`**
    - Mudar `c.container.DisconnectChannel` para `c.container.ChannelService.DisconnectChannel`.

## Riscos
- **Injeção Circular**: Precisamos garantir que o Bot seja injetado no serviço apenas após a inicialização.

## Como testar
1. Desconectar um canal pelo Dashboard.
2. Verificar se o bot envia a mensagem de adeus e sai do canal.
3. Verificar se o canal é removido do banco de dados.

## Rollback
`git checkout ...`
