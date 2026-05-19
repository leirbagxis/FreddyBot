# Plano: Correções Finais da Migração (Callbacks, Emojis, Formatação e AddChannel)

## Pedido do usuário
- Botões como Help, Sticker Separador e Transferir Acesso não estão funcionando.
- As legendas configuradas no bot não estão sendo formatadas (negrito, links, etc.) na hora de editar a postagem.
- O ID de emoji customizado (`custom_emoji`) nos botões não está funcionando.
- Não é possível adicionar um novo canal (Erro `UNIQUE constraint failed: default_captions.owner_channel_id` e "Mensagem 'toadd-message' não encontrada!").

## Objetivo técnico
1. Corrigir o mapeamento de `callback_data` no arquivo `loader_telego.go` para alinhar com os valores definidos em `config/messages.yml`.
2. Atualizar o parser YAML (`pkg/parser/parser.go`) para reconhecer o campo `custom_emoji` e injetá-lo na criação do teclado (`IconCustomEmojiID`).
3. Melhorar o conversor de Markdown para HTML em `DetectParseMode` (`utils_v2.go`) para que suporte links (`[Texto](URL)`) e não aborte prematuramente a conversão se encontrar os caracteres `<` ou `>`.
4. Corrigir o serviço de canais para gerar UUIDs primários (`CaptionID`, `MessagePermissionID`, `ButtonsPermissionID`) ao criar o `DefaultCaption` na função `CreateChannelWithDefaults`.
5. Corrigir os nomes das mensagens YAML chamadas em `addChannel.go` (`toadd-message` -> `toadd-require-message`, `toadd-sucess-message` -> `toadd-success-message`).

## Contexto atual
- Na migração, o nome dos `callback_data` foi presumido a partir dos handlers (ex: `ask-separator`), mas no YAML original o sistema usa chaves curtas (`sptc`, `paccess-info`, `del`, `help`). Isso faz com que os cliques sejam ignorados.
- O campo na struct `Button` está mapeado como `yaml:"icon_custom_emoji_id"`, mas as configurações do usuário usam `yaml:"custom_emoji"`.
- Como o `telego` exige um `ParseMode` definido (neste caso `ModeHTML`), qualquer texto do banco de dados que vier em Markdown (ex: `[Meu Site](http...)`) será renderizado como texto puro. O `DetectParseMode` precisa converter links Markdown para tags `<a>`.
- Ao criar um canal, o GORM tenta inserir o `DefaultCaption` com Primary Keys em branco, violando a integridade do SQLite ao adicionar o segundo canal.

## Arquivos analisados e a serem modificados
- `internal/telegram/loader_telego.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `pkg/parser/parser.go`
- `internal/core/services/channels.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Estratégia de implementação
1. **Callbacks:** Ajustar `loader_telego.go` para escutar `sptc`, `sptc-config`, `spex`, `paccess-info`, `transfer`, `del`, e `help`.
2. **Emojis:** Mudar a tag yaml na struct `Button` em `pkg/parser/parser.go` para `yaml:"custom_emoji,omitempty"`.
3. **Conversão HTML:** Adicionar uma Regex para converter `[texto](url)` em `<a href="url">texto</a>` em `DetectParseMode`.
4. **Constraints:** Importar `github.com/google/uuid` em `channels.go` e gerar UUIDs no bloco de inicialização de `DefaultCaption`.
5. **Nomes de Mensagens:** Alterar os identificadores em `addChannel.go`.

## Passos detalhados
1. Atualizar `pkg/parser/parser.go`.
2. Atualizar `internal/telegram/events/channelPost/utils_v2.go`.
3. Atualizar `internal/telegram/loader_telego.go`.
4. Atualizar `internal/core/services/channels.go`.
5. Atualizar `internal/telegram/handlers/events/addChannel/addChannel.go`.

## Impactos esperados
- Todos os botões da interface voltarão a funcionar.
- Links e formatações escritas em Markdown se tornarão HTML válido.
- Emojis personalizados aparecerão nos botões novamente.
- Usuários conseguirão adicionar canais perfeitamente.

## Como testar
### Build
```bash
go build -o main ./cmd/FreddyBot/main.go
```
### Execução
Tentar vincular um novo canal, clicar no botão de configuração, testar a transferência e verificar os stickers e ajudas.

## Rollback
Desfazer as alterações nos arquivos correspondentes.