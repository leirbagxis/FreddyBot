# Plano: normalizar links botoes telegram

## Pedido do usuario
Corrigir links de botoes quando forem informados como `@nomedocanal` ou `t.me/nomedocanal`, transformando em URL valida do Telegram, como `https://t.me/nomedocanal`. O usuario tambem observou que isso provavelmente nao e tratado quando um novo canal e salvo.

## Objetivo
Normalizar formatos curtos de link Telegram em todos os caminhos que criam ou enviam botoes:
- botoes criados na dashboard/API;
- botoes de legenda customizada;
- botoes do PostBuilder;
- botoes dinamicos extraidos de captions;
- link do canal e primeiro botao criado ao cadastrar novo canal;
- sincronizacao de URL do canal quando username/invite muda.

## Contexto atual
- A dashboard (`ButtonGrid`) ja normaliza `@canal` e `t.me/canal`, mas isso so protege o frontend.
- O PostBuilder rejeita URL que nao comeca com `http`.
- O backend `ButtonService` valida com `net/url` sem normalizar antes.
- `CreateChannelWithDefaults` salva `InviteURL` e primeiro botao com o valor recebido, que pode ser `@username` ou `t.me/username`.
- `UpdateChannelBasicInfoTelego` tambem gera `@username` para canal publico.
- `CreateInlineKeyboardTelego` envia `ButtonURL` direto ao Telegram; se estiver `@canal`, o Telegram rejeita.
- `ExtractDynamicLinks` so reconhece `https?://`, entao nao extrai `@canal` ou `t.me/canal`.

## Arquivos analisados
- `dashboard/src/components/ButtonGrid.tsx`
- `internal/utils/utils.go`
- `internal/core/services/buttons.go`
- `internal/core/services/channels.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/keyboard.go`
- `internal/telegram/events/channelPost/metadata.go`
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/telegram/handlers/commands/admin/admin_channels.go`

## Arquivos que poderao ser modificados
- `internal/utils/utils.go`
- `internal/core/services/buttons.go`
- `internal/core/services/channels.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/keyboard.go`
- `internal/telegram/events/channelPost/metadata.go`

## Estrategia de implementacao
Criar uma normalizacao compartilhada em `internal/utils`, para evitar ciclo de imports:
- `NormalizeTelegramURL(raw string) string`
- `IsValidButtonURL(raw string) bool`

Regras:
- trim de espacos;
- `@canal` vira `https://t.me/canal`;
- `t.me/canal` vira `https://t.me/canal`;
- `telegram.me/canal` vira `https://t.me/canal`;
- `http://`, `https://` e `tg://` permanecem validos;
- strings vazias continuam vazias;
- validacao final exige URL parseavel com scheme/host para `http(s)` ou scheme `tg`.

Aplicar a normalizacao antes de salvar e antes de enviar:
- `ButtonService` normaliza payload antes de validar e persistir.
- `ChannelService.CreateChannelWithDefaults` normaliza `inviteURL` antes de salvar canal e primeiro botao.
- `UpdateChannelBasicInfoTelego` gera URL normalizada para username publico.
- `syncChannelButtons` salva URL normalizada ao sincronizar primeiro botao.
- `CreateInlineKeyboardTelego` normaliza em ultimo caso antes de montar teclado.
- `PostBuilder` normaliza a segunda linha do botao antes da validacao.
- `ExtractDynamicLinks` passa a aceitar formatos curtos e salva URL normalizada.

## Passos detalhados

1. Adicionar funcoes de normalizacao/validacao em `internal/utils/utils.go`.
2. Atualizar `ButtonService` para normalizar `ButtonURL` em create/update de botoes default.
3. Atualizar `ButtonService` para normalizar `ButtonURL` em create/update de botoes de legenda customizada.
4. Atualizar `ChannelService.CreateChannelWithDefaults` para normalizar `inviteURL`.
5. Atualizar `UpdateChannelBasicInfoTelego` e `syncChannelButtons` para trabalhar com URL normalizada.
6. Atualizar `CreateInlineKeyboardTelego` para normalizar antes de enviar, como ultima defesa.
7. Atualizar `PostBuilder` para aceitar `@canal`, `t.me/canal`, `telegram.me/canal`, `http(s)` e `tg://`.
8. Atualizar mensagem de erro do PostBuilder com exemplos aceitos.
9. Atualizar `ExtractDynamicLinks` para capturar e normalizar links curtos.
10. Rodar `gofmt` nos arquivos alterados.
11. Rodar `git diff --check`.
12. Tentar `go test ./...` e `go build ./cmd/FreddyBot/main.go`, registrando falha se o toolchain local continuar sem `vet`/`compile`.

## Riscos
- Regex de links dinamicos precisa evitar capturar caracteres de fechamento como `)` em Markdown.
- `tg://` deve continuar aceito.
- Links completos existentes nao devem ser alterados.
- `Release` esta nao rastreado no workspace e nao deve entrar em commit.

## Impactos esperados
- Usuario consegue usar `@nomedocanal` como link de botao.
- Usuario consegue usar `t.me/nomedocanal` como link de botao.
- Novo canal salvo com username publico gera primeiro botao com URL Telegram valida.
- Botoes enviados ao Telegram deixam de falhar por URL curta.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
go build ./cmd/FreddyBot/main.go
```

### Testes
```bash
go test ./...
git diff --check
```

### Execucao
```bash
make run
```

## Rollback
Reverter alteracoes nos arquivos listados em "Arquivos que poderao ser modificados".

## Observacoes
- Nao sera necessario alterar o frontend da dashboard inicialmente, porque ele ja normaliza os principais formatos.
