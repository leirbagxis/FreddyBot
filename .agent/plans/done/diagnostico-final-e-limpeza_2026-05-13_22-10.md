# Plano: diagnostico-final-e-limpeza_2026-05-13_22-10.md

## Pedido do usuário
Resolver o problema de mapeamento inline (cache miss) e garantir que o banco de dados está estável.

## Objetivo
1.  **Diagnosticar updates**: Confirmar se o bot recebe `ChosenInlineResult` do Telegram.
2.  **Estabilizar DB**: Habilitar logs do GORM para garantir visibilidade sobre as tabelas.
3.  **Limpeza**: Organizar os planos pendentes.

## Contexto atual
- O mapeamento inline falha porque o handler `ChosenInlineResultHandler` parece não ser chamado.
- O banco de dados teve problemas recentes com a tabela `users`, e os logs do GORM estão silenciados, dificultando a verificação.

## Arquivos analisados
- `internal/database/database.go`
- `internal/telegram/client.go`
- `internal/telegram/events/loader.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/telegram/client.go`

## Estratégia de implementação
1.  **Logs de DB**: Alterar `logger.Silent` para `logger.Info` em `InitDB`.
2.  **Logs de Update**: Adicionar um log global no `client.go` (middleware de debug) que mostre o tipo de update recebido.

## Passos detalhados

1.  **Modificar `internal/database/database.go`**
    - Alterar `db.Config.Logger = logger.Default.LogMode(logger.Silent)` para `db.Config.Logger = logger.Default.LogMode(logger.Info)`.

2.  **Modificar `internal/telegram/client.go`**
    - Adicionar um log no início do loop de processamento ou um middleware global de depuração para imprimir o tipo de update.

## Riscos
- **Nulo**: Apenas logs adicionais.

## Como testar
1. Reiniciar o bot.
2. Observar os logs de inicialização do banco.
3. Realizar uma postagem via modo inline e verificar se o log "Update recebido: ChosenInlineResult" aparece.

## Rollback
`git checkout internal/database/database.go internal/telegram/client.go`
