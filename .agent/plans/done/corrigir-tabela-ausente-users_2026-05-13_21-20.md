# Plano: corrigir-tabela-ausente-users_2026-05-13_21-20.md

## Pedido do usuário
Erro ao acessar "Meu Perfil": `SQL logic error: no such table: users (1)`.

## Objetivo
Garantir que a tabela `users` exista no banco de dados SQLite e que as migrações sejam aplicadas corretamente.

## Contexto atual
- O bot está reportando falta da tabela `users` no SQLite.
- O `database.go` possui `AutoMigrate`, mas o erro persiste.
- Recentemente forçamos o drop de um índice de votos, o que pode ter causado interrupção nas migrações se houveram falhas.

## Arquivos analisados
- `internal/database/database.go`
- `internal/database/models/models.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/database/models/models.go`

## Estratégia de implementação
1.  **Simplificar Nomes**: O GORM já usa o plural como padrão. Vou remover o override explícito de `TableName` no `models.go` para evitar confusões de mapeamento, deixando o GORM usar a convenção padrão (que para `User` é `users`).
2.  **Verificar Migração**: Adicionar logs na função `InitDB` para capturar erros específicos do `AutoMigrate`.
3.  **Forçar Migração**: Garantir que o `AutoMigrate` para `User` seja o primeiro da lista.

## Passos detalhados

1.  **Modificar `internal/database/models/models.go`**
    - Remover a função `func (User) TableName() string { return "users" }`.

2.  **Modificar `internal/database/database.go`**
    - Mover `&models.User{}` para ser o primeiro item no `AutoMigrate`.
    - Capturar o erro do `AutoMigrate` e dar panic com a mensagem completa para debug se falhar.

## Riscos
- **Baixo**: Se a tabela já existir com dados, o GORM apenas manterá. Se não existir, ele criará.

## Impactos esperados
- A tabela `users` será recriada (se ausente) e as consultas de perfil voltarão a funcionar.

## Como testar
1. Reiniciar o bot.
2. Clicar em "Meu Perfil".
3. Verificar se o erro desaparece.

## Rollback
`git checkout internal/database/models/models.go internal/database/database.go`
