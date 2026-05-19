# Plano: Correção do Fluxo de Adição de Canal e Permissões

## Pedido do usuário
- A mensagem de sucesso após clicar em "AddYes" não está sendo editada corretamente e os botões (WebApp, Meu Canal, etc.) não funcionam porque faltam variáveis.
- O canal recém adicionado não está recebendo um botão inicial dinâmico (com o nome e link do canal).
- As permissões padrão (`MessagePermission` e `ButtonsPermission`) não estão sendo criadas no banco de dados para os novos canais.
- A configuração das legendas padrão (`packpadrao` e `newpack`) será adiada para uma futura atualização via painel Admin.

## Objetivo
1. Corrigir o handler `AddYesHandlerTelego` para fornecer todas as variáveis necessárias (como `{miniAppUrl}`, `{channelId}`, etc.) ao parser, garantindo que a mensagem de sucesso seja renderizada corretamente.
2. Adicionar o tratamento adequado de erro no `AddYesHandlerTelego` para que edite a mensagem original em caso de falha, em vez de apenas mostrar um alerta.
3. Modificar `CreateChannelWithDefaults` no serviço de canais para gerar as tabelas `MessagePermission` e `ButtonsPermission` junto com o `DefaultCaption`.
4. Modificar `CreateChannelWithDefaults` para criar um `Button` inicial contendo o nome e o link de convite do canal recém-adicionado.

## Contexto atual
- O handler envia apenas o `channelName` para a mensagem `toadd-success-message`. Isso impede a criação dos botões WebApp (que exige `{miniAppUrl}`) e Meu Canal (que exige `{channelId}`).
- O `CreateChannelWithDefaults` apenas inicializa o `DefaultCaption`, ignorando suas relações (`MessagePermission` e `ButtonsPermission`). Se essas tabelas não existirem, os toggles no painel de controle falham ao tentar atualizá-las.
- Não há lógica criando o primeiro botão dinâmico para o canal, o que era um comportamento padrão na versão anterior.

## Arquivos analisados
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/core/services/channels.go`
- `internal/database/models/models.go`
- `config/messages.yml`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/core/services/channels.go`

## Estratégia de implementação
1. Em `channels.go`:
   - Dentro de `CreateChannelWithDefaults`, instanciar `MessagePermission` e `ButtonsPermission` com UUIDs únicos e associá-los ao `DefaultCaption`.
   - Adicionar à slice de `Buttons` do canal um novo botão com `ButtonID` único, `NameButton` igual ao título do canal, e `ButtonURL` igual ao `inviteURL`. O `PositionX` e `PositionY` serão 0.
2. Em `addChannel.go`:
   - Atualizar a passagem de mapa de variáveis no parser `toadd-success-message`, incluindo `botId`, `firstName`, `miniAppUrl` (buscado das configs) e `channelId`.
   - Incluir edição da mensagem também em fluxos de erro durante a adição.

## Passos detalhados

1. **Atualizar `channels.go`:**
   - Adicionar os campos `MessagePermission` e `ButtonsPermission` inicializados com todos os defaults habilitados (true) dentro do `DefaultCaption`.
   - Inicializar a lista de `Buttons` da struct `Channel` com 1 botão, contendo o nome e o URL do canal.
2. **Atualizar `addChannel.go`:**
   - Obter o WebAppURL do `config`.
   - Passar as chaves `{botId}`, `{firstName}`, `{miniAppUrl}`, `{channelId}` e `{channelName}` para a mensagem de sucesso.
   - Refinar a resposta visual caso ocorra um erro ao salvar no banco.

## Riscos
- O banco de dados pode rejeitar a criação do canal se houver problemas de chaves estrangeiras com as permissões. Devemos garantir que o GORM crie as relações corretamente ao salvar o modelo pai (`Channel`).

## Impactos esperados
- Usuários terão a confirmação correta com o botão que abre o WebApp (Dashboard) diretamente.
- O painel de controle carregará corretamente as permissões iniciais.
- O botão dinâmico com o nome do canal aparecerá na primeira vez que uma legenda for gerada.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build -o main ./cmd/FreddyBot/main.go
```

### Execução
1. Inicie o bot.
2. Adicione o bot a um novo canal.
3. Clique em "✅ Sim, vincular".
4. Verifique se a mensagem é substituída pela mensagem de sucesso com os botões "Configure Agora" e "Meu Canal".
5. Verifique no banco de dados se as entradas em `channels`, `default_captions`, `message_permissions`, `buttons_permissions` e `buttons` foram criadas.

## Rollback
Desfazer as alterações usando `git checkout` nos arquivos `channels.go` e `addChannel.go`.

## Observações
- A configuração global de legendas via painel Admin será feita em um plano subsequente.