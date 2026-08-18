# Plano: importar-backup-postgres-teste

## Pedido do usuário
Importar o arquivo `pg-dump-postgres-1785777206.dmp` para fins de teste.

## Objetivo
Restaurar o dump PostgreSQL customizado em um banco de teste isolado, validar o conteúdo restaurado e preservar os dados do banco atualmente em uso.

## Contexto atual
- O backup tem aproximadamente 68 MiB e é um dump PostgreSQL customizado versão 1.16.
- Não há `pg_restore` nem `psql` instalados neste ambiente.
- A configuração local não permite identificar com segurança um PostgreSQL de teste; portanto não é seguro executar restore no destino atual.
- O usuário informou que o uso é somente para testes.

## Arquivos analisados
- `pg-dump-postgres-1785777206.dmp`
- `.env` (somente formato e presença de chaves, sem expor valores)
- `pkg/config/config.go`
- `internal/database/database.go`
- `docker-compose.yml`

## Arquivos que poderão ser modificados
- Nenhum arquivo de código ou configuração do projeto.
- `.agent/memory/memory.md` e o histórico deste plano poderão ser atualizados.

## Estratégia de implementação
Usar um banco PostgreSQL de teste dedicado, por exemplo `freddybot_test_import`, e executar `pg_restore` com cliente compatível com o dump. Antes de qualquer gravação, confirmar host, porta, usuário e nome exato do banco. O banco atual não será substituído sem uma autorização explícita adicional do usuário.

## Passos detalhados

1. Confirmar o destino: banco de teste separado recomendado ou banco atual explicitamente autorizado.
2. Obter um cliente PostgreSQL compatível (`pg_restore`/`psql`) sem alterar o código da aplicação.
3. Listar o conteúdo do dump e confirmar versão e objetos antes da restauração.
4. Criar ou recriar somente o banco de teste nomeado e confirmado.
5. Restaurar o dump usando `pg_restore` com limpeza restrita ao banco de teste.
6. Consultar tabelas e contagens básicas para comprovar a importação.
7. Registrar o resultado sem copiar credenciais ou dados sensíveis para o repositório.

## Riscos
- Restaurar no banco errado pode sobrescrever dados reais.
- Um dump customizado requer ferramentas PostgreSQL compatíveis.
- O dump pode conter proprietários, extensões ou permissões que não existem no servidor-alvo.

## Impactos esperados
- Um banco de teste passa a conter os dados do backup.
- Nenhuma alteração em código, compose, `.env` ou banco de produção.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
make build-server
```

### Testes
```bash
pg_restore --list pg-dump-postgres-1785777206.dmp
psql "$TEST_DATABASE_URL" -c '\\dt'
```

### Execução
```bash
DATABASE_URL="$TEST_DATABASE_URL" pg_restore --dbname="$TEST_DATABASE_URL" --clean --if-exists pg-dump-postgres-1785777206.dmp
```

## Rollback
Apagar somente o banco de teste confirmado. Nenhum banco existente será alterado no fluxo recomendado.

## Observações
- Não será executado restore contra produção.
- `docker-compose.yml` continua fora de escopo, conforme orientação anterior do usuário.
