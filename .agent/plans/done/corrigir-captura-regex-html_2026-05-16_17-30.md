# Plano: corrigir-captura-regex-html_2026-05-16_17-30.md

## Pedido do usuário
O bot não consegue capturar o nome do botão porque a URL está sendo processada primeiro pelo conversor HTML, quebrando a detecção.

## Objetivo técnico
Atualizar a função `ExtractDynamicLinks` e seu `bangLinkRegex` para suportar URLs que já foram encapsuladas em tags `<a>`.

## Contexto atual
O log mostra que o texto antes da extração de links dinâmicos se parece com isto:
```html
!💸 Compre Agora
!<a href="https://google.com/">https://Google.com</a>
```
O `bangLinkRegex` esperava `!https://...`, então ignorava a tag `<a...` e falhava na captura.

## Arquivos analisados
- `internal/telegram/events/channelPost/formatting.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/formatting.go`

## Estratégia de implementação
Em vez de depender puramente do Regex para capturar tudo de uma vez com o HTML bagunçado, podemos ajustar a função para:
1. **Refatorar o Regex:**
   Modificar `bangLinkRegex` para tolerar a presença de tags HTML.
   Novo padrão sugerido: `(?m)^!\s*(.+?)\s*[\r\n]+\s*!\s*(?:<a[^>]*>)?\s*(https?://[^\s<>"]+)(?:</a>)?`
   Isso permite que o link esteja ou não dentro de um `<a href="...">`.
2. **Sanitização do Nome:** Adicionar um regex para limpar qualquer tag HTML que tenha sido inserida acidentalmente na linha do nome do botão (ex: se o usuário mandou o nome em negrito `!**Nome**` que virou `!<b>Nome</b>`).
3. **Manutenção dos logs de depuração:** Manter os logs para confirmar que a alteração resolveu o problema de vez.

## Passos detalhados
1. Atualizar a variável `bangLinkRegex`.
2. No `ExtractDynamicLinks`, capturar o nome e limpar quaisquer tags HTML que estejam nele usando um regex simples de remoção de tags.
3. Testar a montagem.

## Riscos
- O regex ficar muito complexo e lento. *Mitigação: Manter a estrutura simples, apenas tornando as tags HTML opcionais `(?:...)?`.*

## Impactos esperados
- O bot conseguirá unir corretamente o nome do botão com a URL, mesmo que o Telegram tenha formatado a URL como um hyperlink.

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Teste
1. Enviar mensagem com o padrão:
   ```
   !💸 Compre Agora
   !https://google.com
   ```
2. O botão será criado com o nome e URL certos, sem ser interceptado apenas como um "embedded link".

## Rollback
Desfazer as alterações no `formatting.go`.