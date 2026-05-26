# Plano: proteger-link-markdown-parser

## Pedido do usuário
Corrigir definitivamente o link Markdown que continua sem funcionar. O log mostrou que a URL virou `https://t.me/Flor<i>maracuja</i>ofc`, ou seja, a regra de itálico alterou a URL antes da conversão do link.

## Objetivo
Proteger links Markdown `[texto](url)` para que regras de bold/italic/code não modifiquem a URL dentro dos parênteses, e só depois converter o link para HTML `<a href="...">texto</a>`.

## Contexto atual
`DetectParseMode` faz:
1. escape HTML;
2. bold;
3. italic;
4. code;
5. link Markdown.

Como a string ainda contém `[texto](https://t.me/Flor_maracuja_ofc)` durante a etapa de itálico, a regex `_..._` interpreta `_maracuja_` dentro da URL e transforma em `<i>maracuja</i>`. Depois o link é convertido com a URL já corrompida.

## Arquivos analisados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/utils/utils.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/utils/utils.go`

## Estratégia de implementação
1. Em `DetectParseMode`, extrair links Markdown logo após o escape HTML e substituir cada link por placeholder temporário.
2. Rodar bold/italic/code no texto com placeholders, sem URLs expostas às regexes.
3. Converter cada link capturado para `<a href="url">label</a>` com URL normalizada.
4. Restaurar os placeholders pelo HTML final dos links.
5. Aplicar a mesma proteção em `utils.MarkdownToTelegramHTML`, para o broadcast admin não sofrer com underscores em URLs.
6. Manter logs existentes, agora mostrando URL limpa.

## Passos detalhados
1. Criar helper interno em `DetectParseMode` para armazenar links Markdown escapados.
2. Substituir cada link por placeholder como `@@FB_LINK_0@@`.
3. Aplicar bold/italic/code.
4. Restaurar placeholders com HTML de link gerado a partir de URL normalizada.
5. Atualizar `NormalizeMarkdownLinks` em `internal/utils` para também proteger o fluxo usado pelo broadcast.
6. Rodar `gofmt`.
7. Rodar `git diff --check`.
8. Rodar `npm run build`.
9. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando bloqueios locais se persistirem.

## Riscos
- Placeholders precisam ser improváveis de colidir com texto real do usuário.
- Links com parênteses dentro da URL continuam não sendo cobertos pela regex simples atual.
- Há mudanças pendentes no worktree; não devem ser revertidas.

## Impactos esperados
- O log deve mostrar `rawURL="https://t.me/Flor_maracuja_ofc"`, sem `<i>` dentro da URL.
- O dispatch final deve mostrar `<a href="https://t.me/Flor_maracuja_ofc">...`.
- O Telegram deve renderizar o texto como link embutido.

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
Enviar no canal:
```txt
Oi

[𝄗⃝🌻⃝❀⃪֟፝͜͡𝑭𝑳𝑶𝑹 𝑫𝑬 𝑴𝑨𝑹𝑨𝑪𝑼𝑱𝑨𖡼⃟🌻ᬼ⃝⃮࿔꦳꯭ꦿ❀⃪ᰰ᳝ᮀ](https://t.me/Flor_maracuja_ofc)
```

Validar logs:
- URL sem `<i>`;
- `hasHTMLLink=true`;
- preview com `href="https://t.me/Flor_maracuja_ofc"`.

## Rollback
Reverter `internal/telegram/events/channelPost/utils_v2.go` e `internal/utils/utils.go` para o estado anterior.

## Observações
Essa é uma correção em cima do diagnóstico dos logs fornecidos pelo usuário.
