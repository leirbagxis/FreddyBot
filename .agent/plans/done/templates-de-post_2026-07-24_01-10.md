# Plano: Post Templates

## Pedido do usuário
Criar sistema de templates de postagem — usuário pode salvar o estado atual do PostBuilder como um template nomeado, carregar templates salvos e gerenciá-los.

## Objetivo
Adicionar modelo DB `UserPostTemplate`, repositório, serviço, e callbacks no bot para salvar/carregar/deletar templates via PostBuilder.

## Arquivos que serão modificados
- `internal/database/models/user_post_template.go` (novo)
- `internal/database/database.go` (add AutoMigrate)
- `internal/database/repositories/user_post_template.go` (novo)
- `internal/core/services/user_post_template.go` (novo)
- `internal/container/appContainer.go` (wire service)
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` (callbacks)

## Estratégia
Templates por usuário (não por canal). `PostBuilderState` salvo como JSON em `template_data`. Limite de 50 templates por usuário.

## Passos
1. Model `UserPostTemplate`
2. AutoMigrate
3. Repository
4. Service
5. Wire in container
6. Callbacks: save, load (with list), delete
7. Build verification
