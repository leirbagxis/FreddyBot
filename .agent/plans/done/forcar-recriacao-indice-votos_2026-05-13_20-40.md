# Plano: forcar-recriacao-indice-votos_2026-05-13_20-40.md

## Pedido do usuário
O erro de duplicidade de votos (`idx_vote_user`) persiste mesmo após a alteração no código.

## Objetivo
Forçar a atualização do índice único no PostgreSQL, removendo a constraint antiga que incluía o campo `emoji`.

## Contexto atual
- O GORM `AutoMigrate` não remove colunas de índices únicos no PostgreSQL automaticamente.
- O índice `idx_vote_user` no banco de dados está "sujo" com a definição antiga.

## Arquivos analisados
- `internal/database/database.go`
- `internal/database/models/models.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`

## Estratégia de implementação
1.  **Executar SQL Directo**: Adicionar uma instrução `tx.Exec("DROP INDEX IF EXISTS idx_vote_user")` no início da inicialização do banco de dados.
2.  **Remigrar**: Deixar o `AutoMigrate` recriar o índice com a nova definição (sem o campo emoji).
3.  **Remover SQL temporário**: Após a primeira execução bem-sucedida, o código pode ser limpo (ou mantido como uma migration de segurança).

## Passos detalhados

1.  **Modificar `internal/database/database.go`**
    - Localizar a função de inicialização/migração.
    - Adicionar: `db.Exec("DROP INDEX IF EXISTS idx_vote_user")` antes do `AutoMigrate`.

## Riscos
- **Baixo**: O índice será recriado imediatamente pelo GORM. Durante os milissegundos entre o DROP e o CREATE, uma duplicata rara poderia entrar, mas é improvável e o benefício de limpar o índice "preso" é maior.

## Impactos esperados
- A restrição de unicidade passará a ser apenas `(user_id, chat_id, message_id, inline_message_id)`.
- O erro de duplicidade desaparecerá definitivamente.

## Compatibilidade
- PostgreSQL
- SQLite

## Como testar
1. Rodar o bot.
2. Verificar se o log não apresenta erros ao iniciar.
3. Testar o voto rápido novamente.

## Rollback
O índice será recriado automaticamente pelo GORM, então o rollback seria apenas remover a linha do `DROP INDEX`.
