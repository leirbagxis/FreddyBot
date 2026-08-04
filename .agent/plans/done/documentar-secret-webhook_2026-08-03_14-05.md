# Plano: documentar-secret-webhook

## Pedido do usuário
Explicar o `TELEGRAM_WEBHOOK_SECRET`, preencher um valor de exemplo seguro e adicioná-lo ao `.env-example`.

## Objetivo
Disponibilizar no template de ambiente uma variável de exemplo compatível com a API do Telegram, sem inserir ou expor um segredo real no repositório.

## Contexto atual
- `WEBHOOK_URL` já existe em `.env-example`.
- O backend exige `TELEGRAM_WEBHOOK_SECRET` somente quando webhook está ativo.
- O token do bot é uma credencial diferente e não deve ser reutilizado como secret de webhook.

## Arquivos analisados
- `.env-example`
- `README.md`
- `pkg/config/config.go`
- `internal/telegram/client.go`

## Arquivos que poderão ser modificados
- `.env-example`
- `.agent/memory/memory.md`
- Histórico deste plano em `.agent/plans/`

## Estratégia de implementação
Adicionar a variável imediatamente após `WEBHOOK_URL`, com um valor de exemplo alfanumérico seguro para commit e um comentário curto explicando que cada ambiente deve gerar o seu próprio valor. O `.env` real não será alterado.

## Passos detalhados

1. Inserir `TELEGRAM_WEBHOOK_SECRET=example_webhook_secret_change_me_2026` no `.env-example`.
2. Adicionar comentário com formato e finalidade do valor, diferenciando-o do token do bot.
3. Verificar que não há alteração no `.env` nem em credenciais reais.
4. Atualizar memória e mover o plano para concluído.

## Riscos
- Copiar o valor de exemplo para produção reduziria a proteção do webhook; por isso o comentário exigirá substituição por valor único.

## Impactos esperados
- Quem clonar o projeto terá a variável necessária visível no template.
- Nenhuma mudança de comportamento, migração ou alteração no compose.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
GOCACHE=/tmp/freddybot-go-cache go build ./cmd/FreddyBot
```

### Testes
```bash
rg -n "TELEGRAM_WEBHOOK_SECRET" .env-example pkg/config/config.go internal/telegram/client.go
```

### Execução
```bash
WEBHOOK_URL=https://example.com/webhook TELEGRAM_WEBHOOK_SECRET=um_valor_unico_seguro go run ./cmd/FreddyBot
```

## Rollback
Remover somente a linha e o comentário adicionados ao `.env-example`.

## Observações
- Nenhum segredo real será gerado, salvo ou exibido.
- O `.env` local permanecerá intocado.
