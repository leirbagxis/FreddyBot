# Plano: implementar-autodestruicao-posts_2026-08-15_20-40

## Pedido do usuário
Implementar a funcionalidade de **Auto-Destruição Automática de Posts (Mensagens Temporárias)** no FreddyBot.
Permitir que o usuário defina um tempo de expiração/auto-exclusão para posts (ex: Desativada, 1h, 6h, 12h, 24h, 48h) tanto via bot do Telegram (PostBuilder) quanto via Dashboard Web. Quando a mensagem for enviada para o canal (imediatamente ou por agendamento), o bot deve agendar e executar a exclusão da mensagem no canal Telegram no tempo exato programado.

## Objetivo
Criar a infraestrutura de dados, o serviço em background de auto-exclusão (`AutoDeleteService`), as opções no PostBuilder do Telegram, o suporte no `SchedulerService`, as rotas da API REST e a interface no Dashboard Web para postagens temporárias com exclusão automática.

## Contexto atual
- O FreddyBot publica mensagens em canais do Telegram e agenda postagens com `SchedulerService`.
- No momento, não há nenhum mecanismo para apagar automaticamente mensagens do canal após um tempo determinado.
- O GORM gerencia os modelos do banco de dados em `internal/database/models/`.
- O `AppContainer` gerencia os serviços de background através de `StartBackground(ctx)`.

## Arquivos analisados
- `internal/database/models/scheduled_post.go`
- `internal/database/models/auto_delete_post.go` *(a ser criado)*
- `internal/database/repositories/auto_delete.go` *(a ser criado)*
- `internal/database/repositories/scheduled_post.go`
- `internal/database/database.go`
- `internal/cache/types.go`
- `internal/core/services/auto_delete.go` *(a ser criado)*
- `internal/core/services/scheduler.go`
- `internal/container/appContainer.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `dashboard/src/api.ts`
- `dashboard/src/types/index.ts`

## Arquivos que poderão ser modificados / criados
- `internal/database/models/auto_delete_post.go` *(novo)*
- `internal/database/models/scheduled_post.go`
- `internal/cache/types.go`
- `internal/database/repositories/auto_delete.go` *(novo)*
- `internal/core/services/auto_delete.go` *(novo)*
- `internal/core/services/scheduler.go`
- `internal/database/database.go`
- `internal/container/appContainer.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_schedules/my_schedules.go`
- `dashboard/src/types/index.ts`
- `dashboard/src/api.ts`
- `dashboard/src/components/PostBuilderModal.tsx` / componentes do PostBuilder

## Estratégia de implementação

1. **Modelo de Dados & Migração**:
   - Criar `models.AutoDeletePost` (`id`, `channel_id`, `message_id`, `delete_at`, `status`, `last_error`, `created_at`).
   - Adicionar `AutoDeleteMin int` em `models.ScheduledPost` e em `cache.PostBuilderState`.
   - Adicionar `AutoDeletePost` no `AutoMigrate` em `database.go`.

2. **Repositório & Serviço de Auto-Exclusão (`AutoDeleteService`)**:
   - Criar `AutoDeleteRepository` em `internal/database/repositories/auto_delete.go`:
     - `Create`: insere um novo agendamento de exclusão.
     - `GetDuePosts`: retorna postagens com `delete_at <= now` e `status = 'pending'`.
     - `MarkDeleted`: atualiza status para `'deleted'`.
     - `MarkFailed`: atualiza status para `'failed'` com log do erro.
   - Criar `AutoDeleteService` em `internal/core/services/auto_delete.go`:
     - `Start(ctx context.Context)`: loop de background com ticker (ex: 30 segundos) que busca posts vencidos e chama `bot.DeleteMessage(...)`.
     - `ScheduleAutoDelete(ctx, channelID, messageID, autoDeleteMin)`: calcula `delete_at = time.Now().Add(...)` e registra no banco.

3. **Integração no Container e Workers**:
   - Instanciar `AutoDeleteRepository` e `AutoDeleteService` no `AppContainer`.
   - Iniciar `c.AutoDeleteService.Start(ctx)` dentro de `c.StartBackground(ctx)`.

4. **Integração com o Agendador (`SchedulerService`) e Envio Direto**:
   - Ao disparar um post agendado em `SchedulerService.sendScheduledPost`: se `post.AutoDeleteMin > 0`, agendar a auto-destruição para o `msgID` retornado pelo Telegram.
   - Ao publicar um post direto via Telegram Bot no `postBuilder.go`: se `state.AutoDeleteMin > 0`, agendar a auto-destruição da mensagem publicada.

5. **Interface do Bot do Telegram (PostBuilder)**:
   - Adicionar o botão `⏱️ Auto-Destruição: [Desativada / 1h / 6h / 12h / 24h / 48h]` no teclado do PostBuilder em `postBuilder.go`.
   - Criar o handler para o callback `pb-autodelete` alternando ciclicamente as opções de tempo.
   - Na tela de detalhes dos agendamentos (`my_schedules.go`), exibir a indicação de `⏱️ Auto-destruição em Xh` quando o recurso estiver ativado.

6. **Interface no Dashboard Web (React)**:
   - Atualizar os tipos TypeScript (`dashboard/src/types/index.ts`) para incluir `autoDeleteMin?: number`.
   - Adicionar o campo seletor de "Auto-Destruição da Mensagem" nas telas do PostBuilder e de Agendamentos no Dashboard.

7. **Validação & Testes**:
   - Testes automatizados Go (`go test ./...` e `go vet ./...`).
   - Compilação do binário (`go build ./cmd/FreddyBot`).
   - Build do Dashboard (`npm run build`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/implementar-autodestruicao-posts_2026-08-15_20-40.md` e obter aprovação explícita do usuário.
2. Criar `models.AutoDeletePost` e atualizar `ScheduledPost` e `PostBuilderState`.
3. Criar `AutoDeleteRepository` e `AutoDeleteService`.
4. Atualizar `database.go` com a nova migração e registrar o serviço em `appContainer.go`.
5. Integrar o agendamento de auto-destruição em `SchedulerService` e `postBuilder.go`.
6. Adicionar a opção de Auto-Destruição no menu do PostBuilder do Telegram e em `my_schedules.go`.
7. Atualizar as interfaces e tipos do Dashboard Web.
8. Executar testes e compilações (`go test ./...`, `go build`, `npm run build`).

## Riscos
- Risco de erro do Telegram caso o bot perca permissões de exclusão de mensagens no canal antes do tempo de expiração: o serviço tratará com tratamento de erro gracioso gravando `status = 'failed'` e logando o erro sem travar o worker.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Backend Go
```bash
go test ./...
go vet ./...
go build ./cmd/FreddyBot
```

### Build Frontend React
```bash
cd dashboard && npm run build
```

## Rollback
`git checkout . && git clean -fd`
