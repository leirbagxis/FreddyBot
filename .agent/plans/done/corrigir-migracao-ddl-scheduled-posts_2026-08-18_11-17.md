# Plano: corrigir-migracao-ddl-scheduled-posts_2026-08-18_11-17

## Pedido do usuário
Resolver o erro de inicialização do banco de dados no PostgreSQL:
`Erro ao inicializar banco de dados: migrate scheduled_posts.pin_message: ERROR: relation "scheduled_posts" does not exist (SQLSTATE 42P01)`

## Objetivo
Garantir que as migrações DDL manuais do PostgreSQL em `internal/database/database.go` verifiquem com segurança a existência da tabela `scheduled_posts` antes de tentar executar comandos `ALTER TABLE`.

## Contexto atual
- As migrações DDL manuais tentam alterar colunas da tabela `scheduled_posts` antes que ela exista no PostgreSQL.
- O bloco `DO $$ BEGIN` atual no PostgreSQL consulta `information_schema.columns` diretamente sem validar se a tabela `scheduled_posts` já foi criada.

## Arquivos analisados
- `internal/database/database.go`

## Arquivos que serão modificados
- `internal/database/database.go`

## Estratégia de implementação

Em `internal/database/database.go`:
Adicionar a checagem `IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='scheduled_posts')` em cada bloco DDL manual no PostgreSQL:

```postgres
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='scheduled_posts') THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='scheduled_posts' AND column_name='id' AND data_type='uuid') THEN
            ALTER TABLE scheduled_posts ALTER COLUMN id TYPE text;
        END IF;

        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='scheduled_posts' AND column_name='pin_message') THEN
            ALTER TABLE scheduled_posts ADD COLUMN pin_message boolean NOT NULL DEFAULT false;
        END IF;
    END IF;
END $$;
```

Dessa forma, se o PostgreSQL estiver zerado ou a tabela ainda não tiver sido instanciada, a migração DDL manual não estourará erro. Em seguida, o `db.AutoMigrate(...)` criará todas as tabelas nativamente com todas as colunas atualizadas!

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/corrigir-migracao-ddl-scheduled-posts_2026-08-18_11-17.md`.
2. Apresentar o resumo ao usuário e solicitar a aprovação explícita.
3. Modificar `internal/database/database.go`.
4. Executar os testes e o comando `make run` / `go build ./cmd/FreddyBot`.

## Riscos
Nenhum. Apenas torna a verificação de migração DDL do PostgreSQL segura e idempotente.

## Compatibilidade
- Linux, PostgreSQL 17, GORM.

## Como testar
```bash
go test ./internal/database/...
go build ./cmd/FreddyBot
```

## Rollback
`git checkout internal/database/database.go`
