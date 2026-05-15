# Plano: fix-admin-blacklist-functional-impact_2026-05-14_15-00.md

## Pedido do usuário
Mesmo colocando o usuário como admin ou colocando ele na blacklist, não tem efeito algum.

## Objetivo
Garantir que as alterações de status de Admin e Blacklist tenham impacto imediato nas permissões do usuário. Atualmente, o `Role` (cargo) do usuário é codificado no JWT durante o login, o que significa que se um usuário for promovido a Admin, ele precisará fazer logout e login novamente para que o token reflita o novo cargo. Além disso, precisamos verificar se o middleware de blacklist está funcionando corretamente.

## Contexto atual
- O `Role` é armazenado no `Claims` do JWT (`internal/api/auth/jwt.go`).
- O `AuthController.Login` busca o status de admin/blacklist no banco apenas durante o login.
- O `AuthMiddlewareJWT` verifica a blacklist no banco em cada requisição, o que é bom, mas o `RequireRole` e `AuthorizeChannel` confiam apenas no `Role` que está dentro do token.
- Se o status de um usuário mudar, o token dele permanece válido com o cargo antigo até expirar (12h).

## Arquivos analisados
- `internal/api/auth/middleware.go`
- `internal/api/auth/jwt.go`
- `internal/api/controllers/authController.go`
- `internal/database/repositories/user.go`

## Arquivos que poderão ser modificados
- `internal/api/auth/middleware.go`

## Estratégia de implementação
Ajustar o `AuthMiddlewareJWT` para re-validar o `Role` do usuário no banco de dados em cada requisição (ou injetar o cargo atualizado no contexto do Gin). Isso garante que mudanças de permissão sejam instantâneas sem exigir novo login.

## Passos detalhados

1.  **Modificar `internal/api/auth/middleware.go`**
    - No `AuthMiddlewareJWT`, após validar o token e buscar o usuário no banco para checar a blacklist, atualizar o `role` injetado no contexto com base no status real do banco.
    - Se o usuário for o Owner (definido no config), manter `RoleOwner`.
    - Se o usuário tiver `IsAdmin = true` no banco, setar `RoleAdmin` no contexto, mesmo que o token diga `RoleUser`.
    - Se o usuário tiver `IsAdmin = false` no banco, setar `RoleUser` no contexto, mesmo que o token diga `RoleAdmin`.

## Riscos
- **Performance**: Adiciona uma consulta ao banco por requisição autenticada. Como usamos SQLite e as tabelas são pequenas, o impacto é negligenciável para este projeto.
- **Consistência**: Garante que a verdade venha sempre do banco de dados.

## Como testar
1. Logar como um usuário comum.
2. Tentar acessar `/api/admin/overview` (deve dar erro 403).
3. Pela Dashboard Admin (com outra conta Admin/Owner), promover o usuário a Admin.
4. Tentar acessar `/api/admin/overview` com o usuário promovido (deve funcionar agora, sem precisar de re-login).
5. Mesma lógica para a Blacklist: ao ser banido, a próxima requisição deve retornar 403 Forbidden.

## Rollback
`git checkout internal/api/auth/middleware.go`
