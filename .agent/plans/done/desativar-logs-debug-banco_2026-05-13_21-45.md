# Plano: desativar-logs-debug-banco_2026-05-13_21-45.md

## Pedido do usuário
O problema foi resolvido e os logs do banco de dados podem ser desativados.

## Objetivo
Retornar o log do banco de dados (GORM) para o modo silencioso para evitar poluição no terminal.

## Contexto atual
- O log está em modo `logger.Info` para depurar a criação das tabelas.
- O usuário confirmou que as tabelas foram criadas com sucesso.

## Arquivos analisados
- `internal/database/database.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`

## Estratégia de implementação
1.  **Alterar Nível de Log**: Mudar `logger.Info` para `logger.Silent`.

## Passos detalhados

1.  **Modificar `internal/database/database.go`**
    - Localizar a linha `db.Config.Logger = logger.Default.LogMode(logger.Info)`.
    - Alterar para `db.Config.Logger = logger.Default.LogMode(logger.Silent)`.

## Riscos
- **Nulo**: Apenas reduz a verbosidade dos logs.

## Como testar
1. Reiniciar o bot.
2. Verificar se o terminal não mostra mais as queries SQL de cada operação.

## Rollback
`git checkout internal/database/database.go`
