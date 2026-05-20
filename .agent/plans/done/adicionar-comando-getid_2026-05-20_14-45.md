# Plano: Adicionar comando /getid para capturar File IDs

## Pedido do usuário
Adicionar um comando ao bot que, ao ser usado como resposta a uma mídia (foto, gif, vídeo, etc.), retorne o File ID dessa mídia. Esse ID será usado no campo de URL da dashboard admin para disparar mensagens com mídia.

## Objetivo técnico
1. Implementar o handler `GetMediaIDHandlerTelego` para processar o comando `/getid`.
2. O comando deve funcionar apenas quando for uma resposta (`reply`) a uma mensagem contendo mídia.
3. Extrair o File ID de diversos tipos de mídia: Foto, Vídeo, Animação (GIF), Áudio, Documento e Sticker.
4. Registrar o comando no grupo de administradores (`adminGroup`).

## Contexto atual
Atualmente, a dashboard permite enviar mensagens com uma "URL de imagem". No backend (`internal/container/appContainer.go`), o campo `ImageUrl` é usado para preencher o campo `Photo` do `SendPhotoParams`. O Telegram aceita tanto uma URL quanto um File ID nesse campo. Adicionar uma forma fácil de obter o File ID via bot facilitará o uso de mídias que já estão nos servidores do Telegram.

## Arquivos analisados
- `internal/telegram/handlers/commands/admin/admin.go`: Local onde os handlers de comandos administrativos residem.
- `internal/telegram/loader_telego.go`: Local onde os comandos são registrados.
- `internal/container/appContainer.go`: Local onde o processamento do broadcast ocorre.

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/commands/admin/admin.go`
- `internal/telegram/loader_telego.go`

## Estratégia de implementação
- Criar a função `GetMediaIDHandlerTelego` em `admin.go`.
- A lógica verificará `update.Message.ReplyToMessage`.
- Percorrerá os campos de mídia (`Photo`, `Video`, etc.) da mensagem respondida.
- Para Fotos, pegará o ID do maior tamanho disponível.
- Retornará o ID formatado em uma mensagem para o admin, preferencialmente dentro de uma tag `<code>` para facilitar a cópia.

## Passos detalhados

1. **Modificar `internal/telegram/handlers/commands/admin/admin.go`**:
    - Adicionar a função `GetMediaIDHandlerTelego`.
    - Lógica de extração:
        - `Photo`: `m.Photo[len(m.Photo)-1].FileID`
        - `Video`: `m.Video.FileID`
        - `Animation`: `m.Animation.FileID`
        - `Audio`: `m.Audio.FileID`
        - `Document`: `m.Document.FileID`
        - `Sticker`: `m.Sticker.FileID`
    - Responder com uma mensagem informativa contendo o ID.

2. **Modificar `internal/telegram/handlers/loader_telego.go`**:
    - Registrar o comando: `adminGroup.Handle(admin.GetMediaIDHandlerTelego(c), telegohandler.CommandEqual("getid"))`.

3. **Ajustar `internal/telegram/handlers/commands/admin/admin.go` (AdminHelp)**:
    - Adicionar `/getid` à lista de comandos mostrada no `/admin`.

## Riscos
- **Nenhum risco técnico significativo.** É uma funcionalidade de leitura e feedback.

## Impactos esperados
- Administradores poderão obter File IDs de qualquer mídia enviada ao bot ou em canais onde o bot está presente (se encaminhado para o bot).
- Facilidade na criação de campanhas de broadcast com mídias específicas sem necessidade de hospedagem externa.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/...
```

### Testes
1. Enviar uma foto para o bot.
2. Responder a essa foto com `/getid`.
3. Verificar se o bot retorna um ID longo (ex: `AgACAgEAAxkBAAID...`).
4. Repetir com um GIF e um Vídeo.
5. Tentar usar `/getid` sem responder a nada e verificar se o bot avisa que precisa de uma resposta.

## Rollback
Remover o handler e o registro do comando.

## Observações
O campo na dashboard diz "ImageUrl", mas o Telegram tratará o File ID corretamente se passado no mesmo parâmetro. No futuro, podemos renomear o campo na UI para "URL ou File ID da Mídia".
