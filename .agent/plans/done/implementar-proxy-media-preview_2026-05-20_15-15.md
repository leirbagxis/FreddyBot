# Plano: Implementar Proxy de Mídia para Preview de File ID

## Pedido do usuário
O usuário deseja que o dashboard mostre a imagem real quando um File ID do Telegram é inserido no campo de URL da Dashboard Admin, em vez de apenas dar erro de rota não encontrada.

## Objetivo técnico
1. Criar um endpoint de proxy no backend que recebe um File ID, busca o caminho do arquivo via API do Telegram e redireciona (ou serve) os bytes da imagem.
2. Atualizar o frontend do dashboard para detectar quando o valor inserido no campo de URL da mídia é um File ID e usar o novo endpoint de proxy para carregar a imagem.

## Contexto atual
Atualmente, se um usuário cola um File ID no campo de imagem, o navegador tenta carregá-lo como um caminho relativo (`/admin/AgAC...`), gerando erros de log. O Telegram não fornece URLs diretas públicas para File IDs sem o token do bot.

## Arquivos analisados
- `internal/api/routes/routes.go`: Para registrar a nova rota de proxy.
- `internal/api/controllers/adminController/`: Local para o novo controlador.
- `dashboard/src/components/AdminNoticeTab.tsx`: Componente que exibe o preview da imagem.
- `internal/container/appContainer.go`: O container tem o cliente do Telego necessário para chamar `GetFile`.

## Arquivos que poderão ser modificados
- `internal/api/routes/routes.go`
- `internal/api/controllers/adminController/mediaController.go` (novo arquivo)
- `dashboard/src/components/AdminNoticeTab.tsx`

## Estratégia de implementação
1. **Backend:**
    - Criar `MediaController` com o método `GetMediaPreview`.
    - O método usará `bot.GetFile(fileID)` para obter o `FilePath`.
    - Construirá a URL `https://api.telegram.org/file/bot<token>/<file_path>`.
    - Fará um `Redirect` (302) ou um `Proxy` (mais seguro para esconder o token se o redirecionamento expuser o bot_token na URL final para o cliente - o Telegram permite baixar via bot_token na URL, mas é melhor o servidor baixar e enviar os bytes).
    - Decisão: O servidor fará o download e servirá os bytes com o `Content-Type` correto para evitar expor o token no browser.

2. **Frontend:**
    - No componente `AdminNoticeTab`, adicionar uma lógica que verifica se a `noticeImageUrl` parece um File ID (não começa com `http` e tem um tamanho razoável).
    - Se for um File ID, transformar a URL do `<img>` em `/api/admin/media-proxy/:fileID`.

## Passos detalhados

1. **Criar `internal/api/controllers/adminController/mediaController.go`**:
    - Implementar `NewMediaController`.
    - Implementar `GetMediaPreview(ctx *gin.Context)`:
        - Pegar `fileID` do parâmetro da URL.
        - Chamar `c.container.TelegoBot.GetFile`.
        - Fazer um GET na URL do arquivo do Telegram.
        - Retornar o stream de bytes para o cliente.

2. **Modificar `internal/api/routes/routes.go`**:
    - Adicionar `adminRoute.GET("/media-proxy/:fileId", mediaController.GetMediaPreview)`.

3. **Modificar `dashboard/src/components/AdminNoticeTab.tsx`**:
    - Ajustar a lógica da tag `<img>` para usar o proxy se o valor for um File ID.

## Riscos
- **Consumo de Banda:** O servidor atuará como proxy para imagens do Telegram. Como é apenas para o admin ver o preview, o impacto é mínimo.
- **Cache:** Seria bom adicionar cache para não chamar a API do Telegram toda vez que o admin abrir o preview da mesma imagem.

## Impactos esperados
- O administrador poderá ver o que está enviando mesmo usando File IDs obtidos pelo comando `/getid`.
- Fim dos erros de `Rota não encontrada` nos logs ao usar File IDs.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/...
```

### Testes
1. Ir na dashboard admin -> Aba Broadcast.
2. Colar um File ID (ex: `AgAC...`) no campo de URL da imagem.
3. Verificar se a imagem aparece no preview do lado direito.
4. Verificar nos logs se o endpoint `/api/admin/media-proxy/` foi chamado com sucesso.

## Rollback
Remover a rota e o controlador. Reverter as mudanças no componente React.
