# Plano: Corrigir Interceptação Indevida de Mensagens de Texto Comum e Links no PostBuilder

## Pedido do usuário
Quando o usuário envia qualquer texto comum ou link para o bot no chat privado (sem mídia e sem encaminhamento), o bot identifica erroneamente a mensagem como uma chamada para criar um post no `PostBuilder`, enviando o prompt do PostBuilder.

## Objetivo
Garantir que mensagens de texto comum ou links enviados no chat privado **não** acionem o `PostBuilder`, a menos que:
1. A mensagem contenha mídias (foto, vídeo, áudio, sticker, documento, animação).
2. A mensagem seja **encaminhada** de um canal ou chat.
3. O usuário já esteja em uma etapa de edição ativa do PostBuilder (ex: `set_title`, `set_body`, `add_button`).
4. O usuário esteja em um fluxo de agendamento ativo.

## Contexto atual
Atualmente em `internal/telegram/loader_telego.go` no predicado `matchPostBuilderTelego`:
- A condição `update.Message.Text != ""` faz com que **qualquer** mensagem de texto simples seja casada (`true`).
- Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`, `ProcessIncomingContentTelego` define `mediaType = "text"` para qualquer texto simples quando o usuário não possui sessão ativa, enviando a oferta do `🛠️ Post Builder`.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Arquivos que poderão ser modificados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. **Em `loader_telego.go` (`matchPostBuilderTelego`)**:
   - Remover `update.Message.Text != ""` da verificação genérica de disparo.
   - Restringir a captura do PostBuilder a:
     - Sessões ativas de PostBuilder (`state != nil && state.Step != ""`).
     - Fluxos ativos de agendamento (`scheduleState != nil && scheduleState.SessionID != ""`).
     - Mensagens com mídia (`Photo`, `Video`, `Animation`, `Audio`, `Document`, `Sticker`).
     - Mensagens encaminhadas de canais ou usuários (`ForwardOrigin != nil`).

2. **Em `postBuilder.go` (`ProcessIncomingContentTelego`)**:
   - Atualizar a verificação para que textos comuns não encaminhados (sem mídia e sem `ForwardOrigin`) encerrem precocemente sem oferecer o menu do PostBuilder.

## Passos detalhados
1. Atualizar o predicado `matchPostBuilderTelego` em `internal/telegram/loader_telego.go`.
2. Atualizar a guarda de mídia/encaminhamento em `internal/telegram/handlers/events/postBuilder/postBuilder.go`.
3. Testar a compilação do projeto com `go build ./...`.

## Riscos
- Mínimo. Mensagens encaminhadas de canais (incluindo posts texto encaminhados) continuarão disparando o PostBuilder normalmente.

## Impactos esperados
- Textos e links comuns digitados no chat privado com o bot não acionarão mais o PostBuilder indevidamente.

## Compatibilidade
- Linux, Windows, macOS, Docker

## Como testar

### Build
```bash
go build ./...
```

## Rollback
Desfazer as alterações nos arquivos `loader_telego.go` e `postBuilder.go`.
