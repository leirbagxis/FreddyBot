# Plano: corrigir-inicializacao-banco-e-usuario_2026-05-13_21-30.md

## Pedido do usuário
O bot não está criando o banco de dados/tabelas e falha ao salvar o usuário no `/start`.

## Objetivo
Garantir que a tabela `users` seja criada no início e que os dados do usuário sejam salvos de forma síncrona para evitar erros de "tabela não encontrada" ou "usuário não encontrado".

## Contexto atual
- `SaveUserMiddleware` usa goroutine, causando race condition no `/start`.
- `InitDB` silencia erros de migração, dificultando o diagnóstico de por que a tabela `users` sumiu.
- Erro `no such table: users` indica que o `AutoMigrate` falhou ou não foi executado para esse modelo.

## Arquivos analisados
- `internal/database/database.go`
- `internal/middleware/saveUserMiddleware.go`
- `internal/database/models/models.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`
- `internal/middleware/saveUserMiddleware.go`
- `internal/database/models/models.go`

## Estratégia de implementação
1.  **Tornar o salvamento de usuário síncrono**: Remover `go func()` do `SaveUserMiddleware`. Salvar o usuário é rápido e essencial para o funcionamento dos handlers seguintes.
2.  **Habilitar Logs do GORM**: Alterar `logger.Silent` para `logger.Info` em `InitDB` para vermos as queries de criação de tabela no terminal.
3.  **Priorizar Tabela de Usuários**: No `AutoMigrate`, colocar `&models.User{}` como o primeiro item.
4.  **Remover Override de Nome**: Remover `TableName()` da struct `User` para usar o padrão do GORM.

## Passos detalhados

1.  **Modificar `internal/database/models/models.go`**
    - Remover `func (User) TableName()`.

2.  **Modificar `internal/database/database.go`**
    - Mudar `logger.Silent` para `logger.Info`.
    - Garantir que `&models.User{}` seja o primeiro no `AutoMigrate`.

3.  **Modificar `internal/middleware/saveUserMiddleware.go`**
    - Remover o `go func(...) { ... }(...)` e executar o `UpsertUser` diretamente na thread principal do middleware.

## Riscos
- **Performance**: O salvamento síncrono adiciona alguns milissegundos por mensagem, mas é insignificante perto do ganho de consistência.

## Como testar
1. Deletar o arquivo `.db` atual (se existir).
2. Reiniciar o bot.
3. Verificar se o log mostra a criação da tabela `users`.
4. Dar `/start` e verificar se o perfil é carregado corretamente.

## Rollback
`git checkout internal/database/database.go internal/middleware/saveUserMiddleware.go internal/database/models/models.go`
