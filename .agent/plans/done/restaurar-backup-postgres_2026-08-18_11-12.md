# Plano: restaurar-backup-postgres_2026-08-18_11-12

## Pedido do usuário
Restaurar o arquivo de backup de banco de dados `pg-dump-postgres-1787051120.dmp` no novo container PostgreSQL (`postgres_freddybot`) recém-criado.

## Objetivo
Restaurar completamente todas as tabelas, registros, relacionamentos e sequências do backup no banco de dados PostgreSQL `freddybot` rodando no container `postgres_freddybot`.

## Contexto atual
- Arquivo de backup localizado em `/home/malbs/Opencode/FreddyBot/pg-dump-postgres-1787051120.dmp`.
- Container PostgreSQL ativo: `postgres_freddybot` (PostgreSQL 17).
- Credenciais configuradas no `.env`: `postgres://postgres:12345@localhost:5432/freddybot`.

## Arquivos analisados
- `/home/malbs/Opencode/FreddyBot/pg-dump-postgres-1787051120.dmp`
- `.env`
- Estado dos containers Docker (`postgres_freddybot`)

## Estratégia de implementação

1. **Verificação do Formato do Arquivo `.dmp`**:
   - Inspecionar se o arquivo é um dump customizado/binário de `pg_restore` ou script SQL direto (`psql`).

2. **Criação do Banco de Dados Target (se necessário)**:
   - Garantir que o banco de dados `freddybot` exista dentro do PostgreSQL (`CREATE DATABASE freddybot;`).

3. **Restauração dos Dados**:
   - Se for formato `pg_restore` (custom/directory/tar):
     ```bash
     cat pg-dump-postgres-1787051120.dmp | docker exec -i postgres_freddybot pg_restore -U postgres -d freddybot --clean --if-exists --no-owner
     ```
   - Se for formato texto/SQL plano:
     ```bash
     cat pg-dump-postgres-1787051120.dmp | docker exec -i postgres_freddybot psql -U postgres -d freddybot
     ```

4. **Validação da Restauração**:
   - Verificar as tabelas criadas e a contagem de registros (`SELECT count(*) FROM ...`).
   - Testar inicialização do backend Go com o banco restaurado.

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/restaurar-backup-postgres_2026-08-18_11-12.md`.
2. Apresentar o resumo ao usuário e solicitar a aprovação explícita.
3. Executar o comando de restauração do backup após confirmação.
4. Validar os dados restaurados com consultas de contagem de registros.

## Riscos
- O banco de dados de destino `freddybot` (recém-criado) terá seus esquemas preenchidos/sobrescritos pelos dados do backup.

## Compatibilidade
- Linux, PostgreSQL 17, Docker.

## Como testar

### Comando de Verificação pós-restauração
```bash
docker exec -it postgres_freddybot psql -U postgres -d freddybot -c "\dt"
```

## Rollback
Em caso de falha, recrear o banco vazio via `docker exec -it postgres_freddybot psql -U postgres -c "DROP DATABASE freddybot; CREATE DATABASE freddybot;"`.
