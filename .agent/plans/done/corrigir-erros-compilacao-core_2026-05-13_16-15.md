# Plano: corrigir-erros-compilacao-core_2026-05-13_16-15.md

## Pedido do usuário
O usuário relatou 4 erros de compilação após a refatoração do Core.

## Objetivo
Corrigir os erros de compilação nos arquivos `custom_caption.go`, `permissions.go` e `vote.go` na camada de serviços.

## Contexto atual
- `custom_caption.go`: Chamadas para `UpdateCustomCaptionButton` e `DeleteCustomCaptionButton` não estão passando o `context.Context`.
- `permissions.go`: A variável `rowsAffected` está declarada mas não utilizada em `UpdateReactionsActive`.
- `vote.go`: Incompatibilidade de tipo no retorno de `GetVoteCounts` (`map[string]int64` vs `map[string]int`).

## Arquivos analisados
- `internal/core/services/custom_caption.go`
- `internal/core/services/permissions.go`
- `internal/core/services/vote.go`
- `internal/database/repositories/custom_caption.go`
- `internal/database/repositories/vote.go`

## Arquivos que poderão ser modificados
- `internal/core/services/custom_caption.go`
- `internal/core/services/permissions.go`
- `internal/core/services/vote.go`

## Estratégia de implementação
1.  **Custom Caption:** Adicionar `ctx` como primeiro argumento nas chamadas do repositório.
2.  **Permissions:** Adicionar `rows_affected` ao mapa de retorno no `PermissionsService`.
3.  **Vote:** Atualizar a assinatura do método `GetVoteCounts` no serviço para retornar `map[string]int64`, alinhando com o repositório.

## Passos detalhados

1.  **Corrigir `internal/core/services/custom_caption.go`**
    - Localizar chamada de `s.customCaptionRepo.UpdateCustomCaptionButton` e adicionar `ctx`.
    - Localizar chamada de `s.customCaptionRepo.DeleteCustomCaptionButton` e adicionar `ctx`.

2.  **Corrigir `internal/core/services/permissions.go`**
    - No método `UpdateReactionsActive`, incluir `rows_affected: rowsAffected` no objeto `Data` da resposta.

3.  **Corrigir `internal/core/services/vote.go`**
    - Alterar o tipo de retorno de `GetVoteCounts` de `(map[string]int, error)` para `(map[string]int64, error)`.

## Riscos
- **Baixo:** São correções de sintaxe e tipos que impedem o build.

## Impactos esperados
- Build do projeto voltará a funcionar (`go run` ou `make build`).

## Compatibilidade
- Go 1.24+

## Como testar

### Build
```bash
go run cmd/FreddyBot/main.go
```

## Rollback
`git checkout internal/core/services/`
