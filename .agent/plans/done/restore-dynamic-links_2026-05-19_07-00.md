# Plano: Restaurar Links Dinâmicos Duplos e Formatação HTML

## Pedido do usuário
Restaurar o suporte duplo para links dinâmicos (tanto `[Nome](Link)` quanto `!Nome \n !Link`) e garantir que a formatação (negrito, itálico) configurada nas legendas do banco ou via painel não se perca.

## Objetivo
1.  Atualizar a função `ExtractDynamicLinks` em `utils_v2.go` para capturar os dois formatos de links dinâmicos simultaneamente.
2.  Corrigir a limpeza de tags no regex de links dinâmicos (remover tags HTML caso o botão chegue formatado com `<a href="...">` ou tags de negrito).
3.  Revisar `DetectParseMode` e a lógica de processamento de entidades para garantir que `ModeHTML` seja ativado corretamente sempre que a mensagem tiver qualquer tag HTML ou formatação gerada a partir das entidades do Telegram.

## Contexto atual
- `ExtractDynamicLinks` em `utils_v2.go` só está capturando `\[(.*?)\]\((https?://[^\s)]+)\)`. O código original da v2 tinha suporte para `(?m)^!(.+)\s*\n\s*!(https?://[^\s<>"]+)` que foi perdido/ignorado na migração.
- A função `IsMarkdown` é muito genérica e `DetectParseMode` apenas retorna o texto. Como `telego` espera `ParseMode: telego.ModeHTML` fixo nos envios do pipeline (`ProcessTextDispatchTelego` etc), as legendas perdem formatação se retornarem Markdown ou se o parser não fechar as tags corretamente.
- O logger de debug do telego já deve ser removido para limpar o console em dev.

## Arquivos analisados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/client.go` (Para remover o debug logger)

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/client.go`

## Estratégia de implementação
1.  **Client:** Remover `telego.WithDefaultDebugLogger()` de `internal/telegram/client.go` para parar de floodar o terminal.
2.  **Regex de Links:** Em `utils_v2.go`, introduzir `bangLinkRegex` para suportar `!Botão \n !Link` lado a lado com `linkRegex` (`[Botão](Link)`). A função irá buscar pelos dois, adicionar os botões à lista e remover os padrões correspondentes do texto.
3.  **Remoção de Tags:** Se o link dinâmico vier como `!<a href="...">Link</a>`, usaremos um utilitário interno (`utils.RemoveHTMLTags`) para garantir que o "NameButton" seja um texto limpo, senão a API do Telegram retorna erro 400 ao tentar criar o botão.

## Passos detalhados
1.  Editar `internal/telegram/client.go` (linha ~20) para remover a opção de debug.
2.  Editar `internal/telegram/events/channelPost/utils_v2.go`:
    *   Adicionar regex para estilo "bang" (`!Nome \n !Link`).
    *   Iterar o `ExtractDynamicLinks` para usar ambos.

## Riscos
- Nomes de botões com tags HTML misturadas podem causar crash se não forem limpos. (Será mitigado com `RemoveHTMLTags`).

## Impactos esperados
- O bot não floodará mais o log de dev.
- Links com `!Nome` voltarão a ser extraídos corretamente e não aparecerão na legenda.
- A legenda preservará HTML adequadamente.

## Como testar
### Build
```bash
go build -o main ./cmd/FreddyBot/main.go
```
### Teste
1. Enviar um post com `!Botão \n !https://google.com`
2. Enviar um post com `[Outro Botão](https://google.com)`
3. Ambos devem gerar botões e não exibir o texto na legenda.

## Rollback
Desfazer as alterações nos regexes.
