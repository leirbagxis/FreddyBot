# Plano: implementar-mtproto-executor

## Pedido do usuário
Implementar o MTProtoExecutor para aplicar captions com Message Entities (rich formatting via MTProto) quando o usuário possui conta conectada e configurou legenda com formatação.

## Objetivo
- Criar MTProtoExecutor que implementa TelegramExecutor via gotd/td
- Adicionar suporte a entities no pipeline de channel posts
- Usar conta conectada do usuário para editar mensagens via MTProto
- Converter MessageEntityDTOs para tipos gotd/td

## Contexto atual
- TelegramExecutor interface existe em executor/executor.go
- BotAPIExecutor implementa via telego
- ExecutorFactory recebe nil para MTProto
- SetCaptionHandlerTelego já salva entities como JSON no DefaultCaption
- Fase 1 (handlers + DB) já está completa

## Arquivos analisados
- internal/telegram/executor/executor.go
- internal/telegram/executor/factory.go
- internal/telegram/executor/botapi.go
- internal/telegram/events/channelPost/types.go
- internal/telegram/events/channelPost/pipeline_telego.go
- internal/telegram/events/channelPost/dispatch_telego.go
- internal/telegram/events/channelPost/stage_transform_telego.go
- internal/telegram/events/channelPost/stage_send_telego.go
- internal/telegram/events/channelPost/channelPost.go
- internal/container/appContainer.go
- internal/database/models/mtproto_models.go
- internal/telegram/mtproto/auth/auth.go
- internal/core/services/connected_account.go
- internal/telegram/handlers/callbacks/my_channel/caption.go

## Arquivos que poderão ser modificados
- internal/database/models/mtproto_models.go (AccessHash field)
- internal/telegram/executor/executor.go (EditOptions.Entities)
- internal/telegram/executor/mtproto.go (NOVO - MTProtoExecutor)
- internal/telegram/events/channelPost/pipeline_telego.go (novos campos)
- internal/telegram/events/channelPost/stage_transform_telego.go (UseEntities)
- internal/telegram/events/channelPost/dispatch_telego.go (ExecutorFactory)
- internal/telegram/events/channelPost/channelPost.go (passar factory)
- internal/container/appContainer.go (wire MTProto)

## Estratégia de implementação

### Fluxo de dados com entities:
1. Usuário conecta conta Telegram → sessão salva em connected_accounts
2. Usuário configura legenda com formatação → entities salvos como JSON em DefaultCaption.Entities
3. Novo post no canal → StageTransformTelego detecta UseEntities + HasActiveAccount
4. Se true: FormattedText = raw text, FinalEntities = entities JSON combinados
5. Dispatcher usa ExecutorFactory.ForUser() para obter executor
6. Se FinalEntities setado: MTProtoExecutor.messages.editMessage com entities
7. Senão: BotAPIExecutor com HTML

### MTProtoExecutor design:
- Usa ephemeral gotd client (mesmo pattern de auth.withAuth)
- SessionProvider interface obtém sessão do ConnectedAccountService
- Cria telegram.NewClient com session storage -> client.Run -> API call -> retorna
- Converte MessageEntityDTO JSON -> []tg.MessageEntityClass
- Resolve channel access_hash via ConnectedAccountChannel

## Passos detalhados

1. Adicionar AccessHash int64 a ConnectedAccountChannel
2. Adicionar Entities string a EditOptions
3. Criar funcao dtoToMessageEntity() que converte MessageEntityDTO -> tg.MessageEntityClass
4. Criar MTProtoExecutor com SessionProvider interface
5. Implementar EditMessage, EditCaption, EditReplyMarkup, SendSticker, DeleteMessage, SendMessage
6. Adicionar FinalEntities, ExecutorFactory, OwnerID ao ProcessingContextTelego
7. Modificar StageTransformTelego para detectar UseEntities + HasActiveAccount
8. Modificar dispatch_telego.go para usar ExecutorFactory e entities
9. Modificar channelPost.go para passar ExecutorFactory ao pCtx
10. Modificar appContainer.go para criar e passar MTProtoExecutor

## Riscos
- Sessão MTProto expirada: tratar erro e fallback para BotAPI
- access_hash do canal ausente: resolver via channels.getChannels
- Gotd/td com contextos curtos para ephemeral clients (timeout)
- Alteração em pipeline existente pode quebrar fluxo HTML atual

## Impactos esperados
- Usuários com conta conectada poderão usar formatação rich text nas legendas
- Fallback automático para BotAPI HTML se não houver conta conectada
- Nenhuma mudança na API ou DB schema (apenas campo AccessHash)

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./...
```

### Testes
```bash
go test ./internal/telegram/executor/...
```

### Execução
Rodar o bot normalmente e testar fluxo: conectar conta -> configurar legenda com formatação -> postar no canal

## Rollback
Reverter commits que alteram executor/, events/channelPost/, container/appContainer.go
Manter models/mtproto_models.go (campo AccessHash é aditivo)

## Observações
- MTProtoExecutor usa pattern ephemeral (cria client por chamada) igual auth.go
- BotAPIExecutor continua como fallback padrão
- Separator continua usando Bot API (não afetado)
