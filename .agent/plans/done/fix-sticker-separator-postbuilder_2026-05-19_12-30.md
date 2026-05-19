# Plano: fix-sticker-separator-postbuilder

## Pedido do usuário
O usuário relatou três problemas:
1. Ao configurar um sticker separador para um canal, a mensagem exibe a tag `{channelName}` sem substituir pelo nome real do canal.
2. Após enviar o comando para configurar o separador (o Redis registra `awaiting_sticker`), o envio do sticker é interceptado incorretamente pelo `PostBuilder`.
3. O envio de qualquer mídia em um chat privado com o bot inicia instantaneamente uma sessão do `PostBuilder`, e o usuário questiona se esse é o comportamento esperado.

## Objetivo
Corrigir a prioridade dos handlers de sessão interativa sobre os handlers genéricos (como o PostBuilder), refinar o filtro do PostBuilder e garantir a correta renderização de placeholders nas mensagens.

## Contexto atual
- Em `internal/telegram/loader_telego.go`, o handler do PostBuilder (`postbuilder.HandlerTelego`) está sendo registrado **antes** dos handlers de estado ativo (`SetStickerSeparatorHandlerTelego`).
- A função `matchPostBuilderTelego` retorna `true` para qualquer mídia, interceptando o fluxo de separador.
- Na função `SetStickerSeparatorHandlerTelego`, o objeto `Separator` está sendo criado sem um UUID para a chave primária (`ID`), o que pode causar falhas ao salvar.

## Arquivos analisados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/handlers/callbacks/my_channel/channel_actions.go`
- `pkg/parser/parser.go`
- `internal/database/models/models.go`

## Arquivos que poderão ser modificados
- `internal/telegram/loader_telego.go`
- `internal/telegram/handlers/callbacks/my_channel/channel_actions.go`

## Estratégia de implementação
1. **Reordenação de Handlers**: Moveremos o registro dos handlers de sessão ativa (`matchAwaitingStickerSeparatorTelego`, `matchAwaitingTransferAccessTelego`) para o topo, antes da verificação ampla de mídias do PostBuilder.
2. **Refinamento do PostBuilder**: Atualizaremos o `matchPostBuilderTelego` para retornar `false` caso o usuário esteja em algum fluxo de estado pendente no Redis (ex: configurando separador ou transferindo canal).
3. **Correção de UUID no Separador**: Adicionaremos a geração do UUID para o campo `ID` do separador ao salvar no banco.
4. **Resiliência do Parser**: Garantiremos que o nome do canal sempre seja renderizado (fallback para string vazia ou tratamento seguro caso seja nil).

## Passos detalhados

1. **Editar `internal/telegram/loader_telego.go`**:
   - Mover os blocos `bh.Handle(callbackMyChannel.SetStickerSeparatorHandlerTelego(c), ...)` para antes de `bh.Handle(postbuilder.HandlerTelego(c), ...)`.
   - Modificar `matchPostBuilderTelego` para checar `GetAwaitingStickerSeparator` e `GetTransferChannel`, retornando `false` se o usuário estiver nesses estados.

2. **Editar `internal/telegram/handlers/callbacks/my_channel/channel_actions.go`**:
   - Importar `"github.com/google/uuid"`.
   - Na função `SetStickerSeparatorHandlerTelego`, ao instanciar o `Separator`, incluir `ID: uuid.NewString()`.

## Riscos
- Mudar a ordem dos handlers pode afetar outros fluxos. (Risco baixo, as verificações ativas são rigorosas).

## Impactos esperados
- Os usuários conseguirão enviar o sticker e configurar corretamente o separador.
- O PostBuilder continuará funcionando como atalho rápido, mas respeitará sessões ativas de outros comandos.
- A mensagem mostrará o título do canal corretamente.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker

## Como testar

### Build
`go build -o tmp/FreddyBot ./cmd/FreddyBot/`

### Execução
Iniciar o bot, acessar "Meus Canais" -> Escolher um Canal -> Configurar Separador. Verificar substituição da tag e a interceptação correta.

## Rollback
Restaurar os arquivos `loader_telego.go` e `channel_actions.go` via histórico do Git.
