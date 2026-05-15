# Plano: corrigir-erro-duplicidade-voto_2026-05-13_20-30.md

## Pedido do usuário
O bot está registrando erro de `duplicate key value violates unique constraint "idx_vote_user"` ao processar votos.

## Objetivo
Resolver o erro de duplicidade de votos e garantir que a operação de "Toggle" e "Transferência" de voto seja atômica e resiliente a cliques rápidos (race conditions).

## Contexto atual
- O modelo `Vote` possui um índice único `idx_vote_user` que inclui o campo `Emoji`.
- A lógica de `ToggleVote` realiza `Delete` seguido de `Create`, o que não é atômico.
- O erro ocorre no PostgreSQL devido à restrição de unicidade.

## Arquivos analisados
- `internal/database/models/models.go`
- `internal/database/repositories/vote.go`

## Arquivos que poderão ser modificados
- `internal/database/models/models.go`
- `internal/database/repositories/vote.go`

## Estratégia de implementação
1.  **Ajustar Modelo (`models.go`)**: Remover o campo `Emoji` do índice único `idx_vote_user`. Isso garantirá que um usuário só possa ter UM registro por mensagem no banco de dados.
2.  **Refatorar Repositório (`vote.go`)**:
    -   Utilizar uma transação ou lógica de `FirstOrCreate` / `Save` com tratamento de conflito.
    -   Implementar a lógica:
        -   Se o voto para aquele usuário+mensagem+emoji já existe: Deletar (Toggle OFF).
        -   Se não existe: Realizar um Upsert (On Conflict) para atualizar o emoji escolhido ou criar um novo.

## Passos detalhados

1.  **Modificar `internal/database/models/models.go`**
    -   Remover `index:idx_vote_user,unique` da tag do campo `Emoji`.
    -   Manter `index:idx_vote_count` no campo `Emoji` para contagem.

2.  **Modificar `internal/database/repositories/vote.go`**
    -   Reimplementar `ToggleVote` usando uma transação para garantir atomicidade.
    -   Lógica:
        1. Buscar voto existente para o usuário na mensagem.
        2. Se existir:
           - Se for o mesmo emoji: Deletar.
           - Se for emoji diferente: Atualizar campo Emoji.
        3. Se não existir: Criar novo.

## Riscos
- **Migração de Banco:** Alterar índices únicos pode exigir intervenção manual se o GORM não conseguir remover o índice antigo automaticamente no PostgreSQL. (Vou assumir que o GORM AutoMigrate resolverá ou darei a instrução de limpeza).

## Impactos esperados
- Fim dos erros de "duplicate key".
- Votos mais fluidos e sem falhas em cliques rápidos.

## Compatibilidade
- PostgreSQL (Prod)
- SQLite (Dev)

## Como testar
1. Iniciar o bot.
2. Clicar repetidamente e muito rápido em diferentes emojis de reação.
3. Verificar os logs para garantir que não aparecem mais erros de constraint.

## Rollback
`git checkout internal/database/models/models.go internal/database/repositories/vote.go`
