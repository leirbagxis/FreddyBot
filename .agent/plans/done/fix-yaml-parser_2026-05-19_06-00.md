# Plano: Correção do Parser YAML e Reply Markup

## Pedido do usuário
Corrigir o erro ao dar `/start` (falha ao parsear YAML e erro 400 Bad Request da API do Telegram).

## Objetivo
1. Corrigir a lógica de unmarshal no arquivo `pkg/parser/parser.go`.
2. Assegurar que `ReplyMarkup` nunca seja enviado como `null` pela biblioteca `telego`.

## Contexto atual
- O arquivo `config/messages.yml` é uma lista de mensagens.
- O erro `cannot unmarshal !!seq into map[string]parser.Message` indica que, em algum momento, o código tenta jogar a lista diretamente em um map.
- O erro `Bad Request: object expected as reply markup` ocorre porque a mensagem de fallback ("Mensagem 'start' não encontrada!") é enviada sem botões, e a serialização JSON envia `null` ao invés de omitir a propriedade.

## Arquivos analisados
- `pkg/parser/parser.go`
- `internal/telegram/handlers/commands/start/start.go`

## Arquivos que poderão ser modificados
- `pkg/parser/parser.go`
- `internal/telegram/handlers/commands/start/start.go`

## Estratégia de implementação
1. **Refatorar Parser YAML**: O código atual de `loadMessages` no disco parece correto (`yaml.Unmarshal(data, &rawList)`), mas como o erro diz `map[string]parser.Message`, o cache (`loadOnce`) ou o binário antigo pode estar causando problemas. Vou reescrever a função de forma limpa e garantir que a conversão para o mapa aconteça sempre sem erros de tipagem.
2. **Correção do Handler Start**: O `ReplyParameters` vazio está causando problemas com o JSON (`chat_id: ""`). Vamos remover o envio de `ReplyParameters` se não houver um `MessageID` explícito para responder, ou simplificar o envio das mensagens.

## Passos detalhados
1. Editar `pkg/parser/parser.go` para assegurar que `yaml.Unmarshal` use apenas o slice, com logs extras.
2. Editar `internal/telegram/handlers/commands/start/start.go` para remover o `ReplyParameters` desnecessário (que tenta responder à mensagem `/start` do usuário, o que causa o erro de `null`/`""`).

## Riscos
- Mínimos, apenas correção de bugs.

## Impactos esperados
- O comando `/start` voltará a funcionar e carregar a mensagem corretamente do arquivo `.yml`.

## Compatibilidade
- Go 1.25+
- telego v1.9.0

## Como testar
### Build
```bash
go build -o main ./cmd/FreddyBot/main.go
```
### Execução
```bash
./main
```

## Rollback
Desfazer as alterações nos arquivos modificados.
