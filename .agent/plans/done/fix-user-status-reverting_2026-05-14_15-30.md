# Plano: fix-user-status-reverting_2026-05-14_15-30.md

## Pedido do usuário
Mesmo colocando o usuário como admin ou na blacklist, não tem efeito algum.

## Objetivo
Corrigir o bug onde o status de Admin e Blacklist de um usuário é resetado para `false` toda vez que ele interage com o bot. Isso ocorre porque o middleware `SaveUserMiddleware` realiza um "Upsert" que sobrescreve todas as colunas do banco, incluindo as de permissão que não estão presentes no objeto temporário.

## Contexto atual
- `internal/middleware/saveUserMiddleware.go`: Cria um objeto `models.User` apenas com `UserId`, `FirstName` e `Username`.
- `internal/database/repositories/user.go`: O método `UpsertUser` utiliza `clause.OnConflict{UpdateAll: true}`.
- Como o objeto do middleware tem `IsAdmin: false` e `IsBlacklisted: false` (valores padrão de Go), o banco de dados é atualizado para `false` em cada interação do usuário, anulando qualquer alteração feita pela Dashboard Admin.

## Arquivos analisados
- `internal/middleware/saveUserMiddleware.go`
- `internal/database/repositories/user.go`

## Arquivos que poderão ser modificados
- `internal/database/repositories/user.go`

## Estratégia de implementação
Alterar a cláusula `OnConflict` no repositório para atualizar apenas as colunas informativas (`first_name`, `username`, `updated_at`) e preservar as colunas de estado/permissão (`is_admin`, `is_blacklisted`, `is_contribute`).

## Passos detalhados

1.  **Modificar `internal/database/repositories/user.go`**
    - Atualizar `UpsertUser` para usar `DoUpdates` especificando apenas as colunas `first_name`, `username` e `updated_at`.

## Riscos
- **Baixo**: Garante a persistência de dados críticos.

## Como testar
1. Promover um usuário a Admin pela Dashboard.
2. O usuário deve interagir com o bot (ex: enviar uma mensagem).
3. Verificar na Dashboard se ele continua como Admin (não deve voltar para `false`).
4. Verificar se ele consegue executar comandos de Admin no bot.

## Rollback
`git checkout internal/database/repositories/user.go`
