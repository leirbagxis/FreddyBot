# Plano: fase-4-conclusao-reestruturacao-v3_2026-05-14_13-30.md

## Pedido do usuário
Realizar a fase 4 (validação final e conclusão) da reestruturação V3.

## Objetivo
Finalizar a implementação dos padrões definidos na V3 da arquitetura, garantindo que o tratamento global de erros, o uso de DTOs e a centralização do cache na camada de Serviço estejam consistentes em todo o projeto.

## Contexto atual
- As bases (DTOs, Mappers, pacote `pkg/errors`, e Middleware de Erros) já foram implementadas.
- O repositório de canais (`ChannelRepository`) teve a invalidação de cache transformada em um stub (`InvalidateChannelCache` vazio).
- Alguns Controllers (ex: `GetChannelByIDController`, `DisconectChannel`) ainda possuem respostas de erro embutidas (`ctx.JSON(http.StatusBadRequest, ...)`) misturadas com o novo padrão `ctx.Error(err)`.
- Serviços ainda estão chamando métodos de cache vazios nos repositórios, em vez de dependerem de um `CacheService` injetado via container.

## Arquivos analisados
- `internal/api/controllers/channelController.go`
- `internal/core/services/channels.go`
- `internal/database/repositories/channel.go`
- `pkg/errors/errors.go`

## Arquivos que poderão ser modificados
- `internal/api/controllers/*` (Todos os controllers para limpeza de erros).
- `internal/core/services/*` (Refatoração para usar CacheService diretamente e padronizar retornos).
- `internal/database/repositories/*` (Remoção final de métodos stub de cache).
- `internal/container/appContainer.go` (Caso necessário atualizar as dependências injetadas nos serviços).

## Estratégia de implementação

### 1. Limpeza de Tratamento de Erros nos Controllers
Revisar todos os métodos dos controladores (`channelController.go`, `userController.go`, `ButtonsController.go`, etc.) para garantir que **qualquer erro** seja repassado para o middleware via `ctx.Error(err)` ou que use as respostas de erro de DTO quando for um erro de validação/domínio não mapeado pelo `AppError`.

### 2. Centralização do Cache na Camada de Serviço
Remover as chamadas como `s.channelRepo.InvalidateChannelCache(ctx, channelID)` dos serviços (como em `channels.go`). Em vez disso, o `ChannelService` deve receber o `CacheService` (ou a interface de cache apropriada) em seu construtor e gerenciar o cache diretamente.

### 3. Limpeza de Repositórios
Excluir completamente os métodos stub relacionados a cache que permaneceram nos repositórios (ex: `InvalidateChannelCache`).

### 4. Validação Final (Build e Rotas)
Garantir que os DTOs não estão quebrando os mapeamentos esperados pelo frontend (o frontend espera `id`, `title`, etc, conforme definido no JSON).

## Passos detalhados

1. **Revisar e Refatorar Serviços**:
   - Analisar `internal/core/services/*`.
   - Injetar dependência de Cache nos serviços que precisarem.
   - Remover delegamento de cache para os repositórios.

2. **Revisar e Refatorar Controllers**:
   - Analisar `internal/api/controllers/*`.
   - Substituir os retornos `ctx.JSON` de erros não controlados por retornos customizados usando `pkg/errors` e `ctx.Error()`.

3. **Limpeza e Ajuste em Repositórios**:
   - Analisar `internal/database/repositories/*`.
   - Remover stubs.

4. **Validação**:
   - Executar `go build ./...` e ajustar possíveis falhas de dependências.

## Riscos
- Quebra de contrato da API caso um erro anteriormente retornado com estrutura diferente passe a retornar no padrão do Middleware e o frontend não esteja preparado.
- Erros de compilação por mudanças na injeção de dependência (`appContainer.go`).

## Impactos esperados
- Código mais limpo e dentro do princípio de Responsabilidade Única (SRP).
- Controllers focados apenas em orquestração (receber, chamar serviço, formatar saída).

## Compatibilidade
- Linux, macOS, Windows, Docker.

## Como testar

### Build
```bash
go build ./...
```

### Execução
Iniciar a aplicação e fazer login via frontend para validar a comunicação base.

## Rollback
`git reset --hard` para descartar todas as alterações.

## Observações
O foco desta fase é a limpeza e padronização absoluta da arquitetura estabelecida nos passos anteriores.