# Plano: final-unificacao-e-limpeza-core_2026-05-13_15-00.md

## Pedido do usuário
Unificar as duas tarefas pendentes: refatorar redundâncias de DB/API e unificar 100% o Bot com a camada de Core Services.

## Objetivo
Finalizar a transição para a arquitetura de "Núcleo Único", garantindo que 100% das regras de negócio e acesso a dados passem pela camada de `internal/core/services`, tanto para a API quanto para o Bot de Telegram. Além de padronizar as respostas da API usando Generics.

## Contexto atual
- Já temos a estrutura de `internal/core/services`.
- Controladores da API já usam os serviços, mas ainda têm structs de resposta redundantes.
- O Bot de Telegram (`internal/telegram`) ainda acessa os repositórios diretamente em muitos lugares.
- Existem métodos em repositórios que poderiam estar melhor organizados.

## Arquivos analisados
- `internal/api/controllers/` (Vários)
- `internal/telegram/commands/` e `callbacks/`
- `internal/core/services/`
- `internal/database/repositories/`

## Arquivos que poderão ser modificados
- `internal/api/types/response.go` (Novo)
- `internal/core/services/` (Todos)
- `internal/container/appContainer.go`
- `internal/telegram/` (Quase todos os handlers)
- `internal/api/controllers/` (Refatoração para Generics)

## Estratégia de implementação
1. **Padronização de API:** Criar um tipo genérico `APIResponse[T]` para simplificar os controladores.
2. **Expansão dos Serviços:** Adicionar aos serviços os métodos necessários que o Bot usa (ex: `UpsertUser`, `CountChannels`, `ToggleMaintenance`).
3. **Migração do Bot:** Substituir sistematicamente `c.UserRepo` por `c.UserService` (e similares) em todos os arquivos de Telegram.
4. **Limpeza Final:** Remover qualquer lógica de negócio que tenha sobrado nos repositórios, deixando-os apenas com GORM puro.

## Passos detalhados

1. **Fase 1: Padronização da API (Generics)**
   - Criar `internal/api/types/response.go`.
   - Atualizar `ButtonsController`, `CaptionController`, etc., para usar o novo padrão.

2. **Fase 2: Fortalecimento dos Serviços Core**
   - Mover lógica de manutenção do `ServerRepo` para `ServerService` (Novo).
   - Adicionar métodos de paginação e contagem ao `UserService` e `ChannelService`.
   - Garantir que todos os métodos de escrita invalidem o cache.

3. **Fase 3: Migração do Bot de Telegram**
   - Refatorar `internal/telegram/commands/admin/admin.go`.
   - Refatorar `internal/telegram/callbacks/` (claimChannel, my_channel, etc).
   - Refatorar `internal/telegram/events/` (addChannel).

4. **Fase 4: Limpeza de Repositórios**
   - Remover dependências de `pkg/parser` de dentro dos repositórios.
   - Consolidar métodos duplicados.

## Riscos
- O Bot de Telegram tem muitos arquivos; a migração deve ser cuidadosa para não quebrar ponteiros nil.
- Mudança na estrutura JSON da API pode afetar o Dashboard (manteremos compatibilidade de nomes).

## Impactos esperados
- **Arquitetura 100% Limpa:** API e Bot são apenas "cascas" para os Core Services.
- **Redução de Código:** Menos structs de resposta e menos duplicação de lógica.
- **Facilidade de Teste:** Toda a lógica de negócio estará testável em um único lugar (Core).

## Compatibilidade
- Go 1.24+ (necessário para Generics).

## Como testar
- Build completo: `make build-server`.
- Testar comandos do Bot: `/admin`, `/start`, adicionar canal.
- Testar rotas da API via Dashboard.

## Rollback
- `git checkout .` ou revert para o commit estável anterior.
