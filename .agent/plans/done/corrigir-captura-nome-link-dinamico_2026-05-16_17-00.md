# Plano: corrigir-captura-nome-link-dinamico_2026-05-16_17-00.md

## Pedido do usuário
O bot está usando o link como nome do botão, ao invés de pegar o nome que está na linha de cima.

## Objetivo técnico
Corrigir a extração dos grupos de captura do Regex para garantir que o Nome e a URL sejam identificados corretamente em suas respectivas linhas.

## Contexto atual
O regex atual é `(?m)^!\s*([^\n!]+)\s*\n\s*!\s*(https?://[^\s<>"]+)`. Embora pareça correto, o comportamento relatado sugere que o grupo de captura 1 (Nome) e 2 (URL) podem estar sendo mal interpretados ou o regex está falhando em casar a quebra de linha corretamente em alguns ambientes.

## Arquivos analisados
- `internal/telegram/events/channelPost/formatting.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/formatting.go`

## Estratégia de implementação
1. **Refinar o Regex:** Usar uma abordagem mais direta para as linhas.
   Novo Regex: `(?m)^!\s*(.+?)\s*[\r\n]+\s*!\s*(https?://\S+)`
2. **Ajustar Captura:** Garantir que o `ReplaceAllStringFunc` utilize os grupos corretos e adicione logs de depuração para visibilidade no terminal do desenvolvedor.
3. **Limpeza de Nome:** Garantir que o nome capturado não contenha espaços em branco desnecessários ou caracteres de controle.

## Passos detalhados
1. Modificar `bangLinkRegex` em `formatting.go`.
2. Adicionar `logger.Bot` dentro de `ExtractDynamicLinks` para mostrar o que foi capturado como "Name" e o que foi capturado como "URL".

## Riscos
- Se o regex for muito genérico, pode capturar linhas seguidas que começam com `!` mas não são links. *Mitigação: A segunda linha obrigatoriamente deve começar com http/https.*

## Impactos esperados
- O botão será criado com o texto correto (ex: "Visitar Site") e o link correto (ex: "https://google.com").

## Como testar
1. Build do bot.
2. Enviar:
   ```
   ! Clique Aqui
   ! https://google.com
   ```
3. Verificar o log do bot e o resultado no Telegram.

## Rollback
Reverter para a versão anterior do arquivo.
