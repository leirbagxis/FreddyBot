# Plano: corrigir-formatacao-postbuilder

## Pedido do usuário
O usuário relatou que o PostBuilder não está aplicando a formatação (ex: links Markdown) no momento do envio via modo Inline, resultando na exibição de código cru (ex: `[texto](url)`) no chat.

## Objetivo
Garantir que todas as postagens criadas ou importadas no PostBuilder sejam corretamente convertidas e armazenadas como HTML seguro, garantindo que o `InlineQueryResult` funcione corretamente com `ParseMode: telego.ModeHTML`.

## Contexto atual
- O PostBuilder deveria armazenar o estado sempre em HTML (Decisão "Conversão JIT").
- Textos importados de canais (`pb-import-apply`) **não passam** pela conversão, mantendo o Markdown cru no estado.
- A função `ProcessTextWithFormattingTelego` pode estar ignorando sintaxes Markdown se houver apenas entidades automáticas (ex: URLs detectadas pelo Telegram).
- O estado armazenado atualmente na sessão do usuário (`IJLLZNIl`) contém Markdown cru, que está falhando ao renderizar via Inline.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/formatting_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/formatting_telego.go`

## Estratégia de implementação
1. **Corrigir Importação (`pb-import-apply`):** Ao importar dados de um canal, passar `channel.DefaultCaption.Caption` por `channelpost.DetectParseMode` (pois já vem sem entidades do DB).
2. **Safeguard no Inline:** Para corrigir sessões antigas que já estão no Redis com Markdown cru, faremos uma verificação final no `InlineHandlerTelego` antes de construir o `InlineQueryResult`. Se o estado tiver marcadores Markdown (e não for HTML explícito), forçaremos a conversão.
3. **Reforçar `ProcessTextWithFormattingTelego`:** Garantir que o fallback aplique Markdown se as entidades automáticas não cobrirem a formatação manual do texto.

## Passos detalhados
1. Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`, na ação `pb-import-apply`:
   ```go
   if channel.DefaultCaption != nil {
       state.Body = channelpost.DetectParseMode(channel.DefaultCaption.Caption)
   }
   ```
2. Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`, na função `InlineHandlerTelego`, adicionar conversão de fallback:
   ```go
   caption := sb.String()
   // Safeguard: se houver Markdown não convertido e não contiver tags HTML seguras
   if channelpost.IsMarkdown(caption) && !strings.Contains(caption, "<a href=") && !strings.Contains(caption, "<b>") {
       caption = channelpost.DetectParseMode(caption)
   }
   logger.Bot("📝 Caption final para Inline: [%s]", caption)
   ```
3. (Opcional/Segurança) Em `internal/telegram/events/channelPost/formatting_telego.go`, ajustar a precedência caso haja entidades automáticas.

## Riscos
- Mínimos. A conversão de Markdown para HTML é segura e essencial para o modo Inline.

## Como testar
1. Iniciar PostBuilder e importar um canal com legenda padrão em Markdown.
2. Compartilhar via Inline e verificar se o link/negrito funciona.
3. Usar uma sessão já existente que falhou e tentar enviar de novo (deve funcionar via Safeguard).