# Plano: dinamic-bot-username-captions

## Pedido do usuário
O usuário solicitou que o link `t.me/usernamebot` nas legendas padrões globais seja gerado dinamicamente com o nome de usuário real do bot em execução, em vez de um valor fixo.

## Objetivo
Inserir um placeholder (`{botUser}`) na configuração padrão do banco de dados e substituí-lo em tempo real (runtime) pelo username real do bot sempre que um novo canal for vinculado.

## Contexto atual
- Em `internal/database/database.go`, a string estática `t.me/usernamebot` está sendo salva como valor padrão.
- Em `internal/telegram/handlers/events/addChannel/addChannel.go`, os valores globais são copiados diretamente das configurações do servidor para o novo canal. O bot já realiza uma chamada a `bot.GetMe()` neste fluxo.

## Arquivos analisados
- `internal/database/database.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Estratégia de implementação
1. **Atualizar o Template Padrão**: Modificar `internal/database/database.go` para usar o placeholder `{botUser}` no lugar de `usernamebot`.
2. **Substituição em Tempo Real**: No arquivo `internal/telegram/handlers/events/addChannel/addChannel.go`, após recuperar as configurações `globalDefault` e `globalNewPack`, usar `strings.ReplaceAll` para substituir `{botUser}` pelo username real do bot (obtido via `bot.GetMe()`).

## Passos detalhados
1. Editar `internal/database/database.go`:
   - Alterar `GlobalDefaultCaption` para `"🐈‍⠀៹ [t.me/legendasbot](https://t.me/{botUser})  ‹"`.
2. Editar `internal/telegram/handlers/events/addChannel/addChannel.go` na função `AddYesHandlerTelego`:
   - Chamar `bot.GetMe()` antes da criação do canal no banco de dados.
   - Substituir as ocorrências:
     ```go
     if botInfo != nil {
         globalDefault = strings.ReplaceAll(globalDefault, "{botUser}", botInfo.Username)
         globalNewPack = strings.ReplaceAll(globalNewPack, "{botUser}", botInfo.Username)
     }
     ```
   - Usar os valores modificados no `CreateChannelWithDefaults`.

## Riscos
- Risco mínimo. A substituição só ocorrerá no momento da criação do canal. Canais existentes ou legendas personalizadas posteriormente não serão afetadas, o que é o comportamento correto.

## Impactos esperados
- Qualquer novo canal adicionado ao bot sempre herdará as legendas padrão com o link correto apontando para o próprio bot, independentemente de onde o bot estiver hospedado ou qual seja seu username.

## Compatibilidade
- Linux
- macOS
- Windows

## Como testar

### Build
`go build -o tmp/FreddyBot ./cmd/FreddyBot/`

### Execução
1. Atualizar manualmente o banco de dados via Dashboard para usar `{botUser}` nas configurações ou criar um novo banco.
2. Adicionar um novo canal ao bot.
3. Verificar no banco de dados (tabela `default_captions`) se o valor salvo para aquele canal contém o nome de usuário correto do bot.

## Rollback
Desfazer o `strings.ReplaceAll` no manipulador de adição de canais.