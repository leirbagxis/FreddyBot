# Plano: injetar default captions config existente

## Pedido do usuário
Corrigir o bootstrap da configuração global porque, em um banco Postgres já existente com muitos dados, `default caption` e `new pack caption` não são injetadas, enquanto o PostBuilder fixo é.

## Objetivo
Fazer com que `globalDefaultCaption` e `globalNewPackCaption` sejam preenchidos automaticamente quando a linha `ServerConfig` já existir, mas esses campos estiverem vazios.

## Contexto atual
`initServerConfig` monta um `models.ServerConfig` com valores padrão e chama `FirstOrCreate`. Quando a linha `ID=1` não existe, tudo é criado com defaults. Quando a linha já existe, o GORM carrega o registro existente e não aplica os valores do struct aos campos vazios.

Depois do `FirstOrCreate`, existe uma etapa de reparo, mas ela só cobre `FixedPostBuilderKey` e `FixedPostBuilderPayload`. Por isso o PostBuilder fixo é preenchido em bancos antigos, mas `GlobalDefaultCaption` e `GlobalNewPackCaption` podem continuar vazios.

## Arquivos analisados
- .agent/context.md
- internal/database/database.go
- internal/database/models/models.go
- internal/telegram/handlers/events/addChannel/addChannel.go
- internal/core/services/channels.go
- internal/api/controllers/adminController/configController.go

## Arquivos que poderão ser modificados
- internal/database/database.go

## Estratégia de implementação
Extrair os textos padrão de legenda global para constantes/funções locais no pacote `database` e reutilizá-los tanto na criação inicial quanto no reparo de configurações existentes. Após `FirstOrCreate`, verificar `strings.TrimSpace(config.GlobalDefaultCaption)` e `strings.TrimSpace(config.GlobalNewPackCaption)`. Se algum estiver vazio, preencher com o default correspondente e salvar.

A correção não deve sobrescrever valores já configurados pelo admin, apenas preencher campos vazios.

## Passos detalhados

1. Criar constantes para `defaultGlobalDefaultCaption` e `defaultGlobalNewPackCaption` em `internal/database/database.go`.
2. Usar essas constantes no struct inicial de `initServerConfig`.
3. No bloco de reparo pós-`FirstOrCreate`, verificar se `GlobalDefaultCaption` está vazio.
4. Se vazio, preencher com `defaultGlobalDefaultCaption` e marcar `changed = true`.
5. Repetir a lógica para `GlobalNewPackCaption`.
6. Manter a lógica existente do PostBuilder fixo.
7. Rodar `gofmt`.
8. Rodar `git diff --check`.
9. Tentar rodar `go build ./cmd/FreddyBot/main.go` e `go test ./...`, documentando qualquer limitação do toolchain local.

## Riscos
- Baixo risco: alteração restrita ao bootstrap de configuração global.
- Não deve sobrescrever customizações existentes porque só age quando o campo está vazio após `strings.TrimSpace`.
- Em bancos antigos com campo contendo apenas espaços, o valor será considerado vazio e substituído pelo default.

## Impactos esperados
- Bancos Postgres existentes passam a receber `globalDefaultCaption` e `globalNewPackCaption` automaticamente no próximo startup.
- A dashboard admin deve passar a mostrar esses valores após carregar `/api/admin/config`.
- Novos canais vinculados depois disso receberão as legendas globais padrão.
- Canais antigos não serão alterados retroativamente.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
```

### Execução
```bash
go run ./cmd/FreddyBot/main.go
```

## Rollback
Reverter as alterações em `internal/database/database.go`. Em banco já reparado, remover manualmente os valores de `globalDefaultCaption` e `globalNewPackCaption` se for realmente necessário.

## Observações
Essa correção é para a configuração global do servidor. Ela não propaga automaticamente as legendas para canais que já foram criados antes, porque cada canal possui sua própria `DefaultCaption` persistida.
