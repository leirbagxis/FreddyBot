# Plano: Corrigir Botões do PostBuilder (Enviar para Canais e Agendar Envio)

## Pedido do usuário
Quando o usuário clica nos botões `pb-send-to-channels:<sessionID>` ("📢 Enviar para Canais") ou `pb-schedule:<sessionID>` ("📅 Agendar Envio") no PostBuilder, nada acontece.

## Objetivo
Garantir que os botões "Enviar para Canais" e "Agendar Envio" funcionem perfeitamente, exibindo o alerta nativo e/ou a interface de seleção de canais quando o usuário possui canais, ou uma mensagem explicativa com botão de retorno quando o usuário não possui canais cadastrados ou a sessão do post expirou.

## Contexto atual
- No arquivo `internal/telegram/handlers/events/postBuilder/postBuilder.go`, a função `CallbackHandlerTelego` executa uma chamada incondicional a `bot.AnswerCallbackQuery` logo no início do fluxo (linha 714).
- De acordo com a API de Bots do Telegram, um `callback_query_id` só pode ser respondido uma única vez.
- Quando as funções `handleSendToChannelsTelego` ou `handleSchedulePost` tentam chamar `AnswerCallbackQuery` com `ShowAlert: true` (por exemplo, quando o usuário tem 0 canais cadastrados), o Telegram rejeita a chamada secundária com erro `QUERY_ID_INVALID`.
- Como a função retorna imediatamente após essa falha sem alterar o texto ou o teclado da mensagem, o indicador de carregamento do botão é encerrado sem emitir alerta nem alterar a tela, fazendo parecer que "nada aconteceu".
- Além disso, faltava a validação da sessão no Redis antes de renderizar a lista de canais em `handleSendToChannelsTelego` e `handleSchedulePost`.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/loader_telego.go`
- `internal/core/services/channels.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação

1. **Remover a resposta prematura global em `CallbackHandlerTelego`**:
   - Remover a chamada genérica antecedente `_ = bot.AnswerCallbackQuery(...)` na linha 714 da função principal, permitindo que sub-handlers emitam seus alertas customizados (`ShowAlert: true`) sem sofrer rejeição `QUERY_ID_INVALID`.
   - Adicionar chamadas de `AnswerCallbackQuery` pontuais para as ações simples que não possuem handler dedicado.

2. **Validar existência da sessão no Redis em `handleSendToChannelsTelego` e `handleSchedulePost`**:
   - Verificar se `c.CacheService.GetPostBuilderSession(ctx, sessionID)` retorna um estado válido antes de tentar buscar canais.
   - Se a sessão expirou ou não for encontrada, exibir um alerta e atualizar a mensagem na tela informando a expiração com botão para voltar.

3. **Tratar o cenário de 0 canais cadastrados com alerta E atualização de interface**:
   - Quando `len(channels) == 0` ou ocorrer erro ao buscar canais, responder o `AnswerCallbackQuery` com `ShowAlert: true` informando que o usuário não possui canais.
   - Atualizar a mensagem editando o texto para um aviso claro (`"⚠️ Você não possui nenhum canal cadastrado..."`) com o botão `🔙 Voltar ao Post`, garantindo feedback visual completo no chat.

## Passos detalhados

1. **Editar `internal/telegram/handlers/events/postBuilder/postBuilder.go`**:
   - Remover a chamada global `_ = bot.AnswerCallbackQuery(...)` da linha 714.
   - Em `handleSendToChannelsTelego`:
     - Adicionar validação da sessão em cache `GetPostBuilderSession`.
     - Se `len(channels) == 0`, emitir `AnswerCallbackQuery` com `ShowAlert: true` e editar a mensagem na tela com mensagem informativa e botão `🔙 Voltar ao Post`.
   - Em `handleSchedulePost`:
     - Adicionar validação da sessão em cache `GetPostBuilderSession`.
     - Se `len(channels) == 0`, emitir `AnswerCallbackQuery` com `ShowAlert: true` e editar a mensagem na tela com mensagem informativa e botão `🔙 Voltar ao Post`.
   - Adicionar respostas pontuais de `AnswerCallbackQuery` nas ramificações que não possuem handler auxiliar.

2. **Compilar e verificar a aplicação**:
   - Executar `go build ./cmd/FreddyBot` para garantir que o código compila sem erros.

## Riscos
- **Baixo**: A alteração é restrita à camada de apresentação dos callbacks do PostBuilder e preserva a regra de negócios do Telegram API.

## Impactos esperados
- Os botões "📢 Enviar para Canais" e "📅 Agendar Envio" responderão imediatamente.
- Se o usuário tiver canais cadastrados, a lista de seleção de canais será exibida na tela.
- Se o usuário não tiver canais cadastrados, um alerta nativo + tela informativa amigável com botão de retorno serão apresentados.
- Se a sessão do post tiver expirado, o usuário receberá aviso claro em vez de um clique silencioso sem efeito.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Testes
```bash
go test ./...
```

### Execução
```bash
go run cmd/FreddyBot/main.go
```

## Rollback
Em caso de problemas, reverter a alteração no arquivo via git:
```bash
git checkout internal/telegram/handlers/events/postBuilder/postBuilder.go
```

## Observações
Nenhuma observação extra.
