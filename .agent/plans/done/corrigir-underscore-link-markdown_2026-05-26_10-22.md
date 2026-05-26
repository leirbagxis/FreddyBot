# Plano: corrigir-underscore-link-markdown

## Pedido do usuário
Corrigir o link embutido que é detectado nos logs, mas aparece no Telegram sem link. O caso usa URL com underscores: `https://t.me/Flor_maracuja_ofc`.

## Objetivo
Impedir que o parser de itálico (`_texto_`) altere URLs dentro de `href`, e adicionar log do HTML final enviado para confirmar o payload.

## Contexto atual
O log mostra que `DetectParseMode` detecta e converte o link Markdown:

```txt
rawURL="https://t.me/Flor_maracuja_ofc"
normalizedURL="https://t.me/Flor_maracuja_ofc"
```

Mas em `DetectParseMode`, a conversão de link Markdown acontece antes da conversão de itálico. Depois disso, a regex de itálico `_([^_
]+)_` pode encontrar os underscores dentro do `href` (`Flor_maracuja_ofc`) e inserir `<i>` dentro da URL, deixando o HTML inválido para o Telegram.

## Arquivos analisados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/formatting_telego.go`
- `internal/utils/utils.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/dispatch_telego.go`

## Estratégia de implementação
1. Mudar a ordem em `DetectParseMode`: aplicar bold/italic/code antes de converter links Markdown para HTML.
2. Manter os logs de conversão de link já adicionados.
3. Adicionar um log compacto no dispatch de texto antes de `EditMessageText`, mostrando se o texto final contém `<a href=` e um trecho curto sanitizado do HTML final.
4. Evitar logar mensagem inteira para não poluir logs nem expor conteúdo grande.

## Passos detalhados
1. Reorganizar `DetectParseMode` para converter Markdown inline primeiro e links por último.
2. Garantir que `href` continue usando `html.EscapeString(url)`.
3. Adicionar log em `ProcessTextDispatchTelego` com tamanho do texto, presença de link HTML e preview curto.
4. Rodar `gofmt`.
5. Rodar `git diff --check`.
6. Rodar `npm run build` como validação geral.
7. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando se a toolchain local seguir sem `vet`/`compile`.

## Riscos
- Alterar a ordem do parser pode mudar casos raros em que o texto do link também tem Markdown interno, mas evita quebrar URLs com underscore, que é mais crítico.
- O log do HTML final deve ser curto para não gerar excesso de saída.

## Impactos esperados
- URLs como `https://t.me/Flor_maracuja_ofc` não serão mais modificadas pela regex de itálico.
- O Telegram deve receber `<a href="https://t.me/Flor_maracuja_ofc">texto</a>` válido.
- Os logs devem confirmar que o texto final enviado contém `<a href=`.

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
[𝄗⃝🌻⃝❀⃪֟፝͜͡𝑭𝑳𝑶𝑹 𝑫𝑬 𝑴𝑨𝑹𝑨𝑪𝑼𝑱𝑨𖡼⃟🌻ᬼ⃝⃮࿔꦳꯭ꦿ❀⃪ᰰ᳝ᮀ](https://t.me/Flor_maracuja_ofc)
```

Confirmar nos logs:
- conversão Markdown detectada;
- dispatch final com `hasHTMLLink=true`;
- no Telegram, texto clicável.

## Rollback
Reverter `internal/telegram/events/channelPost/utils_v2.go` e `internal/telegram/events/channelPost/dispatch_telego.go` para o estado anterior.

## Observações
Há mudanças pendentes anteriores no worktree. Não devem ser revertidas.
