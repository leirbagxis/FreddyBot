# Plano: implementar-logs-canais-admin

## Pedido do usuário
Adicionar uma nova funcionalidade admin de logs por canais, incluindo também eventos do PostBuilder.

## Objetivo
Criar um sistema persistente de eventos administrativos por canal, visível na dashboard admin, para facilitar auditoria e debugging de postagens, erros, skips, permissões e ações do PostBuilder.

## Contexto atual
O projeto possui:
- Backend Go com GORM e `AutoMigrate` em `internal/database/database.go`.
- Modelos centralizados em `internal/database/models/models.go`.
- Repositories em `internal/database/repositories`.
- Services em `internal/core/services`, que devem concentrar regra de negócio.
- API admin protegida por `RequireRole(auth.RoleAdmin, auth.RoleOwner)` em `internal/api/routes/routes.go`.
- Pipeline de postagem em `internal/telegram/events/channelPost`.
- PostBuilder em `internal/telegram/handlers/events/postBuilder/postBuilder.go`.
- Dashboard admin em React com abas `users`, `channels`, `audit`, `notice`, `config`.

Hoje os logs são principalmente `logger.Bot/API/Error` no stdout/stderr. Isso ajuda em tempo real, mas não dá histórico filtrável por canal na dashboard.

## Arquivos analisados
- `.agent/context.md`
- `internal/database/models/models.go`
- `internal/database/database.go`
- `internal/database/repositories/channel.go`
- `internal/database/repositories/vote.go`
- `internal/core/services/channels.go`
- `internal/api/routes/routes.go`
- `internal/api/controllers/adminController/auditController.go`
- `internal/container/appContainer.go`
- `internal/telegram/events/channelPost/pipeline_telego.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `dashboard/src/App.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/components/AdminDashboard.tsx`

## Arquivos que poderão ser modificados
- `internal/database/models/models.go`
- `internal/database/database.go`
- `internal/database/repositories/channel_event.go` (novo)
- `internal/core/services/channel_events.go` (novo)
- `internal/api/controllers/adminController/channelEventsController.go` (novo)
- `internal/api/routes/routes.go`
- `internal/container/appContainer.go`
- `internal/telegram/events/channelPost/pipeline_telego.go`
- `internal/telegram/events/channelPost/stage_preflight_telego.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/stage_decorate_telego.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminLogsTab.tsx` (novo)
- possivelmente `dashboard/src/index.css` se faltar estilo reutilizável
- `.agent/memory/memory.md`
- `.agent/decisions.md`

## Estratégia de implementação
Criar um log de eventos estruturado, não um dump de texto. Cada evento terá campos fixos para filtro e um `metadata` JSON para detalhes variáveis.

Modelo proposto: `ChannelEvent`.

Campos principais:
- `id` UUID/string primary key
- `channel_id` int64 index nullable/zero permitido para eventos sem canal direto
- `channel_title` string snapshot
- `owner_id` int64 index
- `actor_id` int64 index, quando houver usuário/admin que causou a ação
- `source` string index: `channel_post`, `post_builder`, `admin`, `broadcast`, `newpack`
- `event_type` string index: exemplos abaixo
- `status` string index: `success`, `error`, `skipped`, `info`
- `message_type` string: `text`, `photo`, `video`, `album`, `sticker`, etc.
- `telegram_message_id` int nullable
- `session_id` string index para PostBuilder quando houver
- `error_message` text opcional
- `metadata` text/json serializado
- `created_at` index

Eventos iniciais para instrumentar:
- `post_received`
- `post_processed`
- `post_failed`
- `post_skipped`
- `permission_missing`
- `caption_applied`
- `buttons_applied`
- `dynamic_links_extracted`
- `metadata_updated`
- `postbuilder_started`
- `postbuilder_field_updated`
- `postbuilder_button_added`
- `postbuilder_button_deleted`
- `postbuilder_preview_sent`
- `postbuilder_saved`
- `postbuilder_sent_to_channel`
- `postbuilder_failed`

A primeira versão deve focar em eventos úteis e de baixo risco. Não precisa logar todo clique ou todo estado intermediário se isso poluir demais.

## Passos detalhados
1. Criar `ChannelEvent` em `internal/database/models/models.go` com índices em `channel_id`, `owner_id`, `event_type`, `status`, `source`, `created_at`.
2. Adicionar `&models.ChannelEvent{}` no `AutoMigrate` em `internal/database/database.go`.
3. Criar `ChannelEventRepository` com métodos:
   - `Create(ctx, event)`
   - `List(ctx, filters)` com paginação
   - `Count(ctx, filters)` ou retorno total junto
4. Criar `ChannelEventService` com métodos de alto nível:
   - `Record(ctx, input)`
   - `ListAdmin(ctx, filters)`
   - helpers para serializar metadata sem quebrar o fluxo principal
5. Registrar repository/service em `AppContainer`.
6. Criar controller admin `ChannelEventsController` com endpoint:
   - `GET /api/admin/logs`
7. Filtros da API:
   - `channelId`
   - `ownerId`
   - `actorId`
   - `source`
   - `eventType`
   - `status`
   - `dateFrom`
   - `dateTo`
   - `q` para busca simples em `channel_title`/`error_message`
   - `limit` e `offset`
8. Instrumentar o pipeline de canal sem quebrar processamento:
   - criar eventos de recebimento/processamento/erro/skipped
   - registrar ausência de permissão como `permission_missing`/`skipped`
   - registrar extração de dynamic links e aplicação de botões/legendas de forma resumida
9. Instrumentar PostBuilder:
   - início de sessão
   - edição de campos principais
   - adicionar/deletar botão
   - preview
   - save com `session_id`
   - envio para canal
   - erro ao salvar/enviar
10. Garantir que logging persistente seja best-effort: falha ao salvar log não pode derrubar postagem/PostBuilder.
11. Atualizar `dashboard/src/types.ts` com tipos `ChannelEvent`, filtros e resposta paginada.
12. Atualizar `dashboard/src/api.ts` com `fetchAdminLogs(filters)`.
13. Adicionar aba `Logs` em `adminTabs`.
14. Criar `AdminLogsTab` com UI operacional:
   - filtros por canal/owner/tipo/status/source/período
   - lista paginada
   - badges de status/source
   - detalhe expandível mostrando metadata JSON formatado
   - botão para abrir canal na dashboard quando houver `channelId`
15. Integrar `AdminLogsTab` em `AdminDashboard`.
16. Atualizar `.agent/memory/memory.md` com a nova convenção de logs persistentes.
17. Registrar decisão em `.agent/decisions.md` explicando por que logs vão para dashboard + banco com metadata JSON.
18. Rodar validações.

## Riscos
- Volume de logs pode crescer rápido. Mitigação inicial: paginação e filtros. Retenção automática pode ficar para uma segunda etapa se necessário.
- Logar texto completo de mensagens pode expor conteúdo sensível e aumentar muito o banco. Mitigação: metadata deve conter resumo técnico, IDs, contagens e flags; não armazenar caption completa por padrão.
- Falha de banco durante log não pode afetar o bot. Mitigação: service `Record` deve capturar erro e apenas logar em stdout/stderr.
- Instrumentar muitos pontos do PostBuilder pode poluir histórico. Mitigação: começar com eventos principais e metadata enxuta.
- AutoMigrate em produção precisa ser compatível com Postgres. Modelo deve evitar tipos complexos específicos; metadata pode ser `text` com JSON string para compatibilidade simples.
- Worktree já está sujo com mudanças anteriores. Não reverter nada fora do escopo.

## Impactos esperados
- Admin/owner consegue ver histórico por canal na dashboard.
- Debug de falhas de edição, skips por permissão, erros do Telegram e uso do PostBuilder fica mais rápido.
- Eventos do PostBuilder entram no mesmo histórico, filtráveis por `source=post_builder` e `sessionId`.
- O bot continua funcionando mesmo se persistência de log falhar.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD
- SQLite dev
- PostgreSQL produção

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
# Rodar API/bot como de costume
# Abrir /admin/dash e acessar a aba Logs
```

Cenários manuais:
1. Enviar post normal em canal e verificar eventos `post_received` e `post_processed`.
2. Remover permissão de edição/botões e verificar evento `permission_missing` ou `post_skipped`.
3. Usar PostBuilder: iniciar, editar corpo, adicionar botão, preview, salvar e enviar para canal.
4. Filtrar por canal, source `post_builder`, status `error/success` e período.
5. Abrir detalhe de log e confirmar metadata JSON.

## Rollback
- Remover aba e chamadas da dashboard.
- Remover rota/controller/service/repository de logs.
- Remover instrumentação dos eventos.
- Manter ou remover a tabela `channel_events` conforme decisão operacional. Para rollback de código, a tabela órfã não afeta fluxo.

## Observações
Este plano cria uma nova superfície administrativa relevante. Antes de implementar, confirmar se a primeira versão deve ter retenção automática ou se paginação/filtros bastam por enquanto. Minha recomendação é começar sem retenção automática e adicionar limpeza por idade depois, se o volume justificar.

Há mudanças pendentes anteriores no worktree e planos concluídos ainda não commitados. Não devem ser revertidos.
