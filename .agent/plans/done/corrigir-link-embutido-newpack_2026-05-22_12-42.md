# Plano: corrigir link embutido newpack

## Pedido do usuário
Corrigir o New Pack porque o botão usa o link corretamente, mas o texto com link embutido em Markdown, como `[click here to add]($link)`, não vira link clicável. Adicionar logs para depurar o fluxo.

## Objetivo
Garantir que o template de New Pack transforme variáveis e Markdown em HTML válido para Telegram, especialmente links embutidos com `$link`, e registrar logs úteis do processo.

## Contexto atual
O handler `newpack.go` substitui `$link` antes de chamar `DetectParseMode`. O botão usa o `packURL` diretamente e por isso funciona. O texto passa por `DetectParseMode`, que atualmente converte links Markdown para `<a href="...">...</a>` e depois aplica regexes de bold/italic/code sobre o resultado inteiro. Isso pode mexer dentro do `href` gerado, principalmente se o nome técnico do pack tiver underscores ou outros marcadores, quebrando o link embutido mesmo com o botão correto.

O fluxo atual não loga template, caption renderizada nem HTML final, o que dificulta comparar entrada e saída.

## Arquivos analisados
- internal/telegram/events/channelPost/newpack.go
- internal/telegram/events/channelPost/utils_v2.go
- pkg/logger/logger.go

## Arquivos que poderão ser modificados
- internal/telegram/events/channelPost/newpack.go

## Estratégia de implementação
Criar uma conversão específica para New Pack que substitui variáveis e converte Markdown preservando links. Para reduzir risco em outros fluxos, não alterar `DetectParseMode` global agora. No New Pack, usar uma função local que:

1. Recebe o template e variáveis.
2. Substitui `$titulo`, `$title`, `$name`, `$link`, `$count`, `$total`, `$stickers`.
3. Escapa HTML.
4. Protege/converte links Markdown `[texto](url)` sem permitir que regexes posteriores editem o `href`.
5. Aplica formatações simples fora dos links.
6. Retorna o HTML final para `ParseMode: HTML`.

Adicionar logs com `logger.Bot` para mostrar canal, sticker set, pack title, quantidade, flags dos botões, template bruto, caption após variáveis e HTML final. Em erro de edição da mensagem, logar também o HTML final enviado.

## Passos detalhados

1. Adicionar import `html` em `newpack.go` se necessário.
2. Trocar o uso direto de `DetectParseMode(caption)` por um helper local, por exemplo `renderNewPackTemplateHTML`.
3. Manter compatibilidade das variáveis já suportadas.
4. Converter links Markdown depois de escapar HTML e proteger anchors de regexes posteriores.
5. Adicionar logs antes de editar a mensagem.
6. Adicionar log em caso de falha no `EditMessageText` contendo o erro e o HTML final.
7. Rodar `gofmt`.
8. Rodar `git diff --check`.
9. Tentar `go build` e `go test`, documentando bloqueios do toolchain local.

## Riscos
- Baixo a médio: alteração localizada no New Pack.
- Templates muito complexos de Markdown podem continuar limitados ao suporte atual simples do projeto.
- Logs podem expor o texto do template no stdout, mas isso é aceitável para depuração operacional do bot.

## Impactos esperados
- `[click here to add]($link)` passa a virar `<a href="https://t.me/addstickers/...">click here to add</a>` corretamente.
- Botões continuam usando o link direto do pack.
- Logs deixam claro se o problema está no template salvo, na substituição de variáveis ou no HTML enviado ao Telegram.

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
```bash
go run ./cmd/FreddyBot/main.go
```

Teste manual:
1. Salvar na dashboard: `[click here to add]($link) - $count stickers`.
2. Enviar `/newpack` em um canal.
3. Enviar sticker de pack público.
4. Conferir se o texto `click here to add` aparece clicável e se os logs mostram template, caption renderizada e HTML final.

## Rollback
Reverter a alteração em `internal/telegram/events/channelPost/newpack.go` para voltar ao uso direto de `DetectParseMode`.

## Observações
A correção será local ao New Pack para não alterar o comportamento de captions normais, post builder e outros fluxos que dependem de `DetectParseMode`.
