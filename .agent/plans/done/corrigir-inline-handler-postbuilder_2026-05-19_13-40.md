# Plano: corrigir-inline-handler-postbuilder

## Pedido do usuário
O usuário relatou que o comando inline `@FreddyCaptionBot pb uo4wkO04` não está exibindo os resultados (a postagem) no Telegram.

## Objetivo
Restaurar o funcionamento do modo inline do PostBuilder, garantindo que o HTML gerado seja válido e aceito pela API do Telegram, além de remover processamentos redundantes que podem corromper a formatação.

## Contexto atual
- O PostBuilder armazena o texto já em HTML seguro (decisão "Conversão JIT").
- O `InlineHandlerTelego` processa o ID recebido e monta o `InlineQueryResult`.
- Atualmente, o `InlineHandlerTelego` chama `channelpost.DetectParseMode` no título, corpo e rodapé antes de enviar para o Telegram, mesmo que o texto já esteja em HTML.
- `DetectParseMode` não escapa caracteres HTML (como `&` e `<`), gerando HTML inválido caso o texto original os contenha e `IsMarkdown` tenha ativado a função original.
- Falhas na chamada `bot.AnswerInlineQuery` não são registradas (erros engolidos por `_ = `), fazendo com que a falha seja invisível nos logs.

## Arquivos analisados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/formatting_telego.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/utils_v2.go`

## Estratégia de implementação
1. **Corrigir Parser:** Alterar `DetectParseMode` em `utils_v2.go` para executar `html.EscapeString(text)` antes de aplicar as regexes de substituição. Isso garante que qualquer caracter como `&` torne-se `&amp;` ANTES das tags HTML (`<b>`, `<i>`, etc.) serem inseridas.
2. **Remover Redundância:** No `postBuilder.go`, nas funções `InlineHandlerTelego` e `sendFinalPostTelego`, remover as chamadas a `channelpost.DetectParseMode`. O estado `state.Title`, `state.Body` e `state.Footer` já foi convertido e validado como HTML na entrada.
3. **Adicionar Rastreabilidade:** Tratar o retorno de erro nas funções `bot.AnswerInlineQuery` no modo inline, enviando-os para o `logger` para facilitar debug futuro.

## Passos detalhados

1. Em `internal/telegram/events/channelPost/utils_v2.go`:
   - Atualizar a função `DetectParseMode` adicionando a linha `res = html.EscapeString(res)` (com importação de `"html"` caso não esteja).

2. Em `internal/telegram/handlers/events/postBuilder/postBuilder.go`:
   - Na função `InlineHandlerTelego`, substituir `sb.WriteString(channelpost.DetectParseMode(state.Title) + "\n\n")` por `sb.WriteString(state.Title + "\n\n")` (aplicar para Body e Footer também).
   - Fazer a mesma substituição na função `sendFinalPostTelego`.
   - Modificar `_ = bot.AnswerInlineQuery(...)` em `InlineHandlerTelego` para capturar o erro:
     ```go
     if err := bot.AnswerInlineQuery(context.Background(), ...); err != nil {
         logger.Error("BOT", "Erro ao responder Inline Query: %v", err)
     }
     ```

## Riscos
- Mudar o `DetectParseMode` afeta todo o sistema que depende dele para converter Markdown puro. A alteração deve ser testada para garantir que não escapa HTML já existente vindo do banco de dados (que deveria estar seguro).

## Impactos esperados
- O modo inline voltará a exibir as postagens do PostBuilder corretamente para todos os caracteres.
- Menor processamento (CPU) ao evitar rodar as regexes repetidamente.
- Melhor visibilidade de erros no backend.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Testes
1. Criar um PostBuilder com a mensagem: `Texto teste com & e < e >`.
2. Salvar e tentar compartilhar via botão (modo inline). O preview deverá abrir normalmente.
3. Testar a conversão básica do parser para Markdown enviando `*negrito* & texto`.

### Execução
```bash
go run cmd/FreddyBot/main.go
```

## Rollback
Desfazer as alterações nos arquivos com `git checkout internal/telegram/events/channelPost/utils_v2.go internal/telegram/handlers/events/postBuilder/postBuilder.go`.

## Observações
O bug é gerado pelo Telegram ser estrito com a sintaxe `ParseMode: HTML`. Qualquer caractere reservado (especialmente o E-comercial `&` não codificado) derruba a renderização da query inteira. Devido às restrições do sistema de planejamento, o plano foi salvo na pasta do sistema do Gemini, mas será copiado para `.agent/plans/pending/` no início da implementação.