# Plano: postbuilder-fixo-admin-dashboard

## Pedido do usuário
Criar uma postagem PostBuilder fixa a partir da payload fornecida, sem data de expiração, com chave fixa, e permitir ativar/desativar e editar pela Dashboard Admin.

## Objetivo
Adicionar um PostBuilder promocional/permanente administrável, acessível por uma chave estável no inline (`pb <key>`), sem TTL no Redis, com controle de ativo/inativo e edição pelo painel admin.

## Contexto atual
- Sessões normais do PostBuilder são salvas em Redis com chave `pb_session:<id>` por `SavePostBuilderSession`.
- O TTL atual é de 24h.
- O inline handler usa `GetPostBuilderSession(id)` e responde `@bot pb <id>`.
- Configurações globais admin já são persistidas em `ServerConfig` e editadas por `AdminConfigTab`.
- A Dashboard Admin já possui rotas autenticadas para `GET/PUT /api/admin/config`.

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `internal/cache/types.go`
- `internal/cache/cache.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/api/routes/routes.go`
- `internal/api/controllers/adminController/configController.go`
- `internal/database/models/models.go`
- `internal/core/services/server.go`
- `internal/database/repositories/serverConfig.go`
- `dashboard/src/components/AdminConfigTab.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`

## Arquivos que poderão ser modificados
- `internal/cache/cache.go`
- `internal/database/models/models.go`
- `internal/database/database.go`
- `internal/database/repositories/serverConfig.go`
- `internal/core/services/server.go`
- `internal/api/controllers/adminController/configController.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/components/AdminConfigTab.tsx`

## Estratégia de implementação
Persistir a configuração da postagem fixa no banco, dentro de `ServerConfig`, e sincronizar o conteúdo para Redis sem expiração usando uma chave fixa.

Chave proposta:
```txt
legendasbot
```

Chave Redis efetiva:
```txt
pb_session:legendasbot
```

Uso inline:
```txt
@FreddyCaptionBot pb legendasbot
```

Campos novos propostos em `ServerConfig`:
- `FixedPostBuilderEnabled bool`
- `FixedPostBuilderKey string`
- `FixedPostBuilderPayload string`

O payload será JSON compatível com `cache.PostBuilderState`.

## Passos detalhados

1. Adicionar campos novos no model `ServerConfig`.
2. Atualizar `initServerConfig` com valores padrão:
   - enabled: `true`
   - key: `legendasbot`
   - payload: JSON da payload enviada pelo usuário.
3. Adicionar métodos no cache:
   - salvar sessão PostBuilder com ID fixo e sem TTL;
   - remover sessão fixa quando desativada.
4. Atualizar `ServerService.UpdateConfig` para receber e validar os campos da postagem fixa.
5. Atualizar `ConfigController` para aceitar os novos campos no `PUT /api/admin/config`.
6. Ao salvar configuração:
   - se ativo, validar JSON e gravar em Redis como `pb_session:<fixedKey>` sem expiração;
   - se inativo, remover `pb_session:<fixedKey>` do Redis.
7. Atualizar `AdminConfigTab`:
   - toggle Ativo/Inativo da postagem fixa;
   - input da key fixa;
   - textarea/editor JSON para payload;
   - botão salvar junto das configs globais.
8. Atualizar `types.ts` e `api.ts` para os novos campos.
9. Rodar `gofmt`.
10. Rodar build do dashboard.
11. Tentar build/testes Go e registrar limitação se o Go local continuar sem `compile`/`vet`.

## Riscos
- `AutoMigrate` adicionará colunas novas em `server_configs`; é uma alteração de schema não destrutiva.
- Payload JSON inválida deve ser rejeitada pela API para não quebrar o inline handler.
- Redis pode perder a chave em restart; por isso a fonte de verdade será o banco e a API re-sincroniza ao salvar. Opcionalmente, a inicialização pode sincronizar também.
- Chaves fixas conflitantes com sessões aleatórias são improváveis, mas a key `legendasbot` deve ser reservada.

## Impactos esperados
- Admin poderá editar payload e status pela Dashboard.
- `@FreddyCaptionBot pb legendasbot` retornará a postagem fixa enquanto estiver ativa.
- A chave Redis não terá expiração.
- PostBuilder normal continuará usando IDs aleatórios com TTL de 24h.

## Compatibilidade
- Linux: compatível
- macOS: compatível
- Windows: compatível
- Docker: requer migration automática via GORM
- CI/CD: requer build Go e build Vite

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
cd dashboard && npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
make dev
```

Teste manual:
1. Abrir Dashboard Admin > Configurações.
2. Ver a seção da postagem fixa.
3. Confirmar key `legendasbot`.
4. Salvar payload.
5. Usar `@FreddyCaptionBot pb legendasbot`.
6. Desativar no painel.
7. Confirmar que `@FreddyCaptionBot pb legendasbot` retorna postagem não encontrada.
8. Reativar e confirmar que volta a funcionar.

## Rollback
- Reverter alterações nos arquivos listados.
- Se necessário, ignorar as colunas novas adicionadas por `AutoMigrate`; elas não quebram versões antigas.
- Remover manualmente a chave Redis:
```bash
redis-cli DEL pb_session:legendasbot
```

## Observações
- A key fixa proposta é `legendasbot`. Se você quiser outro nome, alterar antes da implementação.
- Como o inline não renderiza custom emoji de forma confiável, os botões usarão fallback Unicode no resultado inline, conforme decisão anterior.
