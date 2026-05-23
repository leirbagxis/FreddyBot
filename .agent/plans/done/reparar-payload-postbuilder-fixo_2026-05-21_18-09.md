# Plano: reparar-payload-postbuilder-fixo

## Pedido do usuário
Corrigir erro na inicialização: `Payload do PostBuilder fixo inválida: unexpected end of JSON input`.

## Objetivo
Garantir que configurações existentes com payload vazio/inválido sejam reparadas automaticamente com o payload padrão e sincronizadas no Redis sem TTL.

## Contexto atual
- `ServerConfig` ganhou `FixedPostBuilderPayload`.
- Em instalações existentes, a coluna pode estar vazia mesmo com `FixedPostBuilderEnabled=true`.
- A sincronização no container tenta `json.Unmarshal` do payload e falha com `unexpected end of JSON input`.
- O bootstrap atual só repara payload exatamente vazio, mas pode não cobrir string com espaços ou payload inválido já salvo.

## Arquivos analisados
- `internal/database/database.go`
- `internal/container/appContainer.go`
- `internal/api/controllers/adminController/configController.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/container/appContainer.go`

## Estratégia de implementação
Adicionar reparo defensivo no bootstrap do banco e na sincronização do container. Se o payload estiver vazio, apenas espaços, ou JSON inválido, usar o payload padrão e salvar de volta em `ServerConfig`.

## Passos detalhados

1. Exportar uma função em `internal/database` que retorne o payload padrão.
2. No `initServerConfig`, validar o payload com `json.Unmarshal`; se inválido, preencher com padrão.
3. No `syncFixedPostBuilderSession`, se o payload estiver vazio/inválido, substituir pelo padrão, persistir a configuração e sincronizar Redis.
4. Rodar `gofmt`.
5. Validar com `git diff --check`.
6. Tentar `go build`, registrando a limitação do ambiente se continuar sem `compile`.

## Riscos
- Se o admin salvar uma payload inválida manualmente fora da API, ela será substituída pelo padrão na próxima inicialização.
- Payload inválida pela Dashboard já deve ser bloqueada pela API.

## Impactos esperados
- A inicialização deixa de logar erro de payload vazia.
- `pb_session:legendasbot` volta a ser sincronizada sem TTL.
- Configuração existente é auto-reparada.

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
make dev
```

## Rollback
Reverter alterações em `internal/database/database.go` e `internal/container/appContainer.go`.

## Observações
- Sem alteração de schema.
