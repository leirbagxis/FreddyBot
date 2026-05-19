# Plano: aprimorar-regex-links-dinamicos_2026-05-16_16-45.md

## Pedido do usuário
O bot não está lendo o padrão `!Nome` e `!Link`.

## Objetivo técnico
Tornar o Regex de extração de links dinâmicos (bang-style) mais robusto e flexível para variações de espaçamento e quebras de linha.

## Contexto atual
O regex atual é `(?m)^!(.+)\s*\n\s*!(https?://[^\s<>"]+)`. Ele pode estar sendo muito rígido quanto à quebra de linha ou o conteúdo capturado no nome.

## Arquivos analisados
- `internal/telegram/events/channelPost/formatting.go`

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/formatting.go`

## Estratégia de implementação
1. **Refinar Regex:** Alterar o regex para capturar melhor o nome e a URL, permitindo espaços opcionais logo após o `!`.
   Novo Regex sugerido: `(?m)^!\s*([^\n!]+)\s*\n\s*!\s*(https?://[^\s<>"]+)`
2. **Normalização:** Garantir que o texto capturado seja limpo de espaços extras.
3. **Debug:** Adicionar logs temporários (opcional, mas recomendado) para entender o que está chegando no extrator.

## Passos detalhados
1. Atualizar a variável `bangLinkRegex` em `internal/telegram/events/channelPost/formatting.go`.
2. Ajustar a função `ExtractDynamicLinks` se necessário para lidar com as capturas.

## Riscos
- O regex capturar algo que não deveria se for muito permissivo. *Mitigação: Manter a exigência da URL começando com http.*

## Impactos esperados
- O bot passará a reconhecer o padrão mesmo que haja um espaço após o `!` ou variações de quebras de linha (Windows/Unix).

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Teste
1. Enviar no canal:
   ```
   Texto
   ! Botão 
   !https://link.com
   ```
2. Verificar se o botão é criado.

## Rollback
Reverter para o regex anterior.
