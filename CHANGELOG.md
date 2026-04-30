# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Gerenciamento de Reações**: Adicionado toggle para ativar/desativar reações globalmente por canal na aba de Permissões.
- **Novo Sistema de Toasts**: Redesign completo das notificações do dashboard com glassmorphism, barra de progresso e animações modernas.
- **Validação de Emojis**: Implementada verificação rigorosa em tempo real para reações (dashboard e bot), impedindo o uso de letras e números.
- **Aba de Configurações Admin**: Nova interface para gerenciar configurações globais do servidor.
- **Modo Force Join**: Implementada obrigatoriedade de inscrição em canal oficial para uso do bot, com verificação em tempo real no comando `/start`.
- **Modo Manutenção**: Adicionado toggle global para colocar o bot em manutenção via dashboard.
- **Novo Sistema de Log Centralizado**: Implementação de um logger customizado em `pkg/logger` com suporte a cores e módulos (`[BOT]`, `[API]`, `[DB]`, `[ADMIN]`, `[SYS]`).
- **Post Builder**: Adicionado botão "🚀 Compartilhar" na mensagem de confirmação de salvamento para facilitar o compartilhamento via modo inline.
- **Otimização de Performance**:
  - **Paralelismo no Bot**: Comando `/channels` agora verifica a administração de múltiplos canais simultaneamente usando um worker pool.
  - **Fila de Mensagens Inteligente**: Aumentada a vazão para 5 workers paralelos com controle de rate limit por chat (cooldown de 500ms), otimizando o envio em massa para múltiplos canais.
  - **Cache de Banco de Dados**: Integração com Redis no `ChannelRepository` (padrão Cache-Aside) e criação de buscas "light" para reduzir o overhead do banco.
  - **Unificação de API**: Novo endpoint `/api/admin/overview` que consolida dados de usuários e canais em uma única requisição para o dashboard.
- **Segurança e Estabilidade**:
  - Adicionados Timeouts (`Read`, `Write`, `Idle`) ao servidor HTTP para prevenir vazamentos de conexões e ataques Slowloris.
  - Implementado Graceful Shutdown para desligamento seguro da API e conexões Redis.

### Changed
- **Otimização de Regex**: Pré-compilação de todas as expressões regulares globais (formatação de texto e detecção de Markdown/HTML) para reduzir uso de CPU.
- **Refatoração do Comando /channels**: Melhoria na UX com resposta direta (reply) e edição da mensagem de status em tempo real.
- **Ajustes de Infraestrutura**: Conversão do banco de dados legado de PostgreSQL para SQLite para facilitação do desenvolvimento local.
- **Processamento de Canais**: Implementada verificação `via_bot` para que o bot ignore postagens enviadas por ele mesmo via modo inline, prevenindo loops de edição ou processamento duplicado.

### Fixed
- **Desconexão de Bot**: Corrigido erro de "channelID inválido" ao mover a rota para o grupo autorizado e padronizar para o padrão RESTful.
- **Sincronização de UI**: Padronização de campos camelCase (ex: `forceJoin`) entre frontend e backend para garantir persistência visual dos toggles.
- **Verificação de Membros**: Refinada a lógica de Force Join para incluir usuários com status `restricted` e melhorar feedback de erros.
- **Middlewares**:
  - `MaintenanceMiddleware`: Corrigida a lógica para permitir updates normalmente quando a manutenção está desativada, evitando o bloqueio de novos usuários ou consultas inline.
  - `AdminMiddleware`: Adicionada verificação de segurança para evitar pânico quando um usuário não existe na base de dados.
- **Modo Inline**:
  - O Post Builder agora fornece um feedback visual ("Postagem não encontrada") caso o ID da sessão seja inválido ou tenha expirado.
  - Corrigido erro onde o caption poderia ficar vazio em resultados do tipo Article.
- **Webhook**: Atualizada a lista de `AllowedUpdates` para incluir `channel_post`, `edited_message` e `edited_channel_post`, garantindo que o bot receba todos os eventos necessários.
- Erro de compilação em `admin.go` relacionado à mudança de campos na biblioteca `go-telegram/bot`.
- Vazamentos potenciais de memória em loops de processamento de texto.

### Added (Reconstrução Anterior)
- **Suporte ao SQLite3**: Agora é possível alternar entre SQLite e PostgreSQL usando a variável de ambiente `APP_ENV`.
...
