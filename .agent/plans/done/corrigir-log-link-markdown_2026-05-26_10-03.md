# Plano: corrigir-log-link-markdown

## Pedido do usuário
Investigar e corrigir o caso em que um link Markdown do tipo `[texto](https://t.me/Flor_maracuja_ofc)` ou `[texto](t.me/Flor_maracuja_ofc)` não vira link embutido, adicionando logs para entender o que acontece.

## Objetivo
Adicionar logs de diagnóstico nos conversores de Markdown para HTML e corrigir os pontos mais prováveis onde o link embutido deixa de ser convertido ou chega inválido ao Telegram.

## Contexto atual
Existem dois conversores principais:
- `channelpost.DetectParseMode`, usado em captions, postagens e PostBuilder.
- `utils.MarkdownToTelegramHTML`, usado pelo broadcast admin.

`DetectParseMode` já normaliza URLs com `NormalizeTelegramURL`, mas os safeguards do PostBuilder só chamam a conversão quando não há `<a href=`, `<b>` ou `<tg-emoji>`. Se a mensagem tiver custom emoji HTML e também link Markdown cru, a conversão pode ser pulada.

`MarkdownToTelegramHTML` troca `[texto](url)` direto por `<a href="url">texto</a>` sem normalizar nem logar, então pode gerar HTML diferente do fluxo de postagens.

## Arquivos analisados
- `internal/utils/utils.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/formatting_telego.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/api/controllers/adminController/getAllUserAdminController.go`

## Arquivos que poderão ser modificados
- `internal/utils/utils.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`

## Estratégia de implementação
1. Adicionar logs quando um link Markdown for detectado e convertido em `DetectParseMode`.
2. Adicionar logs quando `MarkdownToTelegramHTML` detectar e converter link Markdown no broadcast.
3. Normalizar URLs no conversor `MarkdownToTelegramHTML` para manter comportamento igual ao `DetectParseMode`.
4. Ajustar o safeguard do PostBuilder para converter links Markdown crus mesmo quando a legenda contém `<tg-emoji>` ou outras tags HTML.
5. Não logar o texto completo da mensagem; logar apenas URL original, URL normalizada e tamanho do texto do link, para evitar vazar conteúdo desnecessário.

## Passos detalhados
1. Atualizar `MarkdownToTelegramHTML` para usar `ReplaceAllStringFunc`, `NormalizeTelegramURL` e logs.
2. Atualizar `DetectParseMode` para logar URL original/normalizada quando converter `[texto](url)`.
3. Criar helper simples no PostBuilder para detectar presença de link Markdown cru (`[...](...)`).
4. Trocar as duas condições do PostBuilder que hoje pulam conversão quando há `<tg-emoji>`.
5. Rodar `gofmt`.
6. Rodar `git diff --check`.
7. Rodar `npm run build` para garantir frontend intacto.
8. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando se a toolchain local continuar sem `vet`/`compile`.

## Riscos
- Logs em alto volume podem aparecer bastante se muitos usuários usarem links Markdown. Por isso os logs serão pontuais e sem conteúdo completo da mensagem.
- Regex simples de Markdown não cobre URLs com parênteses dentro, mas cobre o caso reportado.
- Há mudanças pendentes no worktree de tarefas anteriores; elas não devem ser revertidas.

## Impactos esperados
- Logs devem mostrar se o link foi detectado, qual URL entrou e qual URL normalizada saiu.
- `[texto](https://t.me/Flor_maracuja_ofc)` e `[texto](t.me/Flor_maracuja_ofc)` devem chegar ao Telegram como `<a href="https://t.me/Flor_maracuja_ofc">texto</a>`.
- PostBuilder deve converter link Markdown mesmo quando o corpo já contém custom emoji HTML.

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
Testar com:
```txt
⋆.✿❁˚🌻🌼༘ 𑁍ܓ🌻🌼˚❂✿
[𝄗⃝🌻⃝❀⃪֟፝͜͡𝑭𝑳𝑶𝑹 𝑫𝑬 𝑴𝑨𝑹𝑨𝑪𝑼𝑱𝑨𖡼⃟🌻ᬼ⃝⃮࿔꦳꯭ꦿ❀⃪ᰰ᳝ᮀ](https://t.me/Flor_maracuja_ofc)
⋆.✿❁˚🌻🌼༘ 𑁍ܓ🌻🌼˚❂✿
```

Conferir logs procurando conversões de Markdown link e confirmar que o texto sai como link embutido.

## Rollback
Reverter os arquivos modificados neste plano para remover os logs e voltar ao comportamento anterior.

## Observações
Há mudanças pendentes anteriores no worktree, incluindo normalização de URLs de botões, botão miniapp no `/info`, cache da auditoria e broadcast múltiplo. Essas alterações não devem ser revertidas.
