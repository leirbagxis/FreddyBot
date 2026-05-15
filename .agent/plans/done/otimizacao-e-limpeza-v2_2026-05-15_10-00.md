# Plano: otimizacao-e-limpeza-v2_2026-05-15_10-00.md

## Pedido do usuário
Realizar a limpeza e otimização final da arquitetura V2 do motor de postagens.

## Objetivo
1. **Limpeza do Legado V1 (MessageProcessor):** Integrar lógicas espalhadas nos métodos de `MessageProcessor` de volta aos Estágios nativos do pipeline V2.
2. **Cache L1 (Local em Memória):** Adicionar cache concorrente local antes das chamadas ao Redis para configurações em massa de canais.
3. **Otimização de Banco de Dados:** Criar índices vitais nos modelos do GORM para escala.

## Contexto atual
- O projeto usa a arquitetura V2 de pipeline para processamento de mensagens.
- `MessageProcessor` atua como ponte com código da V1.
- Falta cache local, causando gargalo no banco (SQLite) ou latência de rede (Redis) se escalado de fato.

## Arquivos analisados
- `internal/telegram/events/channelPost/types.go`
- `internal/telegram/events/channelPost/dispatch_v2.go`
- `internal/telegram/events/channelPost/formatting.go`
- `internal/cache/local.go`
- `internal/database/models/channel.go`
- `internal/database/models/user.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/*.go`
- `internal/cache/local.go`
- `internal/core/services/channels.go`
- `internal/database/models/channel.go`
- `internal/database/models/user.go`

## Estratégia de implementação

1. **Fase 1: Banco de Dados:**
   Adicionar tags `gorm:"index"` para `telegram_id` nas tabelas `channels` e `users`, e para `owner_id` em `channels`.

2. **Fase 2: Cache L1 Local:**
   Construir estrutura na memória (`sync.Map` ou via `RWMutex`) com TTL e instanciá-la em `internal/cache/local.go`. Modificar a obtenção de Canais em `ChannelService` para checar `Cache Local -> Cache Redis -> SQL`.

3. **Fase 3: Refatoração de Pipeline:**
   Remover ou reduzir a mega estrutura `MessageProcessor`. Suas funções internas (de `dispatch_v2.go` e `formatting.go`) passarão para funções privadas dos próprios arquivos de estágio (`stage_decorate.go`, `stage_send.go`, etc), utilizando o `ProcessingContext`.

## Passos detalhados
1. Atualizar structs de `models`.
2. Implementar controle de Map na memória e injetar como dependência de `ChannelService`.
3. Atualizar fluxo de invalidação (`InvalidateChannel`) para limpar cache local além do redis.
4. Mover lógica contida nos processadores do `MessageProcessor` para `Stage`.

## Riscos
- **Gargalo no Bot:** Desacoplar os processos que interagem com o bot pode falhar se alguma dependência for perdida no Pipeline.
- **Sincronia de Cache:** Precisaremos assegurar que a expiração L1 (local) seja de curta duração.

## Compatibilidade
- Linux, Docker

## Como testar
### Build
`go build ./...`
### Execução
Bot testado via Webhook postando conteúdos multimídia em canais configurados.

## Rollback
Caso problemas severos apareçam no despacho de mensagens: `git reset --hard` e reversão dos modelos caso migrações gerem locks.