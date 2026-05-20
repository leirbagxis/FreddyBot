# Plano: Varredura Completa e Melhorias Gerais (Audit & Refactor)

## Pedido do usuário
Realizar uma varredura completa no bot (funções, banco de dados, mensagens) para identificar e implementar melhorias.

## Objetivo técnico
1. **Refatoração:** Reduzir o tamanho de arquivos críticos e melhorar a legibilidade.
2. **Correções Técnicas:** Eliminar valores hardcoded e resolver débitos técnicos identificados.
3. **Otimização de Banco:** Revisar índices e garantir integridade referencial.
4. **Padronização:** Unificar o estilo das mensagens e emojis no bot.

## Contexto atual
O projeto cresceu significativamente. Arquivos como `admin.go` e `postBuilder.go` ultrapassaram 1000 linhas, dificultando a manutenção. Existem referências a bots antigos (XavolaBot) e alguns comentários de "TODO" espalhados pelo código.

## Arquivos analisados
- `internal/telegram/handlers/commands/admin/admin.go`
- `internal/telegram/handlers/events/postBuilder/postBuilder.go`
- `internal/database/models/models.go`
- `config/messages.yml`
- `internal/core/services/channels.go`

## Arquivos que poderão ser modificados
- `internal/telegram/handlers/commands/admin/admin.go`
- `internal/telegram/handlers/commands/admin/admin_users.go` (Novo)
- `internal/telegram/handlers/commands/admin/admin_channels.go` (Novo)
- `internal/telegram/handlers/commands/admin/admin_utils.go` (Novo)
- `internal/core/services/channels.go`
- `config/messages.yml`
- `internal/database/repositories/channel.go`

## Estratégia de implementação

### 1. Refatoração de Handlers (Modularização)
- Dividir `admin.go` em múltiplos arquivos dentro da pasta `admin/` para separar lógica de Usuários, Canais e Utilidades.
- Mover funções auxiliares de `admin.go` para um arquivo de utilidades.

### 2. Correção de Débitos Técnicos
- **Comando `/checkbot`**: Manter a lógica atual de verificação do @XavolaBot conforme solicitado pelo usuário (comando imutável).
- **Mensagem de Ajuda**: Atualizar referências de nomes em outros comandos, mas manter a descrição do `/checkbot` apontando para o XavolaBot.
- **Service Channels**: Implementar o método `Update` no repositório mencionado no comentário "TODO" em `channels.go`.

### 3. Padronização de Mensagens
- Revisar `messages.yml` para garantir que o uso de `tg-emoji` e HTML esteja consistente em todas as seções.
- Melhorar a clareza das mensagens de erro.

### 4. Melhorias no Banco de Dados
- Garantir que todos os `Deletes` de canais realmente limpem `Separator`, `Buttons` e `CustomCaptions` via CASCADE no SQLite (verificar se as FKs estão ativas).

## Passos detalhados

1. **Modularizar `admin.go`**:
    - Criar `admin_users.go`, `admin_channels.go` e `admin_utils.go`.
    - Distribuir os handlers existentes.
2. **Corrigir `/checkbot`**:
    - Atualizar `CheckBotAdminHandlerTelego` em `admin_utils.go` para ser dinâmico.
3. **Limpeza de TODOs**:
    - Resolver a pendência de `UpdateChannel` em `internal/core/services/channels.go`.
4. **Auditoria de Mensagens**:
    - Revisar e unificar `config/messages.yml`.

## Riscos
- **Regressões no Bot**: A modularização de arquivos grandes pode quebrar referências se não for feita com cuidado.
- **Banco de Dados**: Mudanças em FKs ou lógica de delete precisam ser testadas para não apagar dados órfãos ou falhar por restrição.

## Impactos esperados
- Código mais limpo e fácil de expandir.
- Bot totalmente dinâmico, sem referências a IDs ou nomes antigos.
- Melhor performance em consultas de canais devido à limpeza de código redundante.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/...
```

### Testes
1. Testar comandos de admin: `/users`, `/channels`, `/checkbot`, `/getid`.
2. Verificar se o bot responde corretamente após a modularização.
3. Simular a exclusão de um canal e verificar se os dados relacionados sumiram do banco.

## Rollback
O rollback envolveria restaurar a versão única do arquivo `admin.go` e reverter as mudanças nos serviços.
