# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Suporte para Arquivos (Arquivos)**: O bot agora suporta a edição e inclusão de botões em postagens do tipo "Arquivo". Suporte adicionado em todas as camadas (Backend, Bot e Dashboard).
- **Controle Granular de Arquivos**: Adicionadas permissões específicas para "Arquivos" na aba de permissões da Dashboard, permitindo ativar/desativar a edição de legenda e adição de botões para este tipo de mídia.
- **Despedida Dramática**: O bot agora envia uma mensagem melodramática e sarcástica antes de sair automaticamente de um canal ao ser desconectado.
- **Alertas de Callback**: Adicionada verificação de existência de canal em todos os handlers de callback do Telegram, exibindo um alerta (ShowAlert) caso o canal não esteja vinculado ou o usuário não tenha permissão.
- **Testes Unitários**: Implementação inicial de suite de testes para o middleware de autorização e repositórios de banco de dados.
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
- **Workflows do Makefile**:
  - `make build`: Agora realiza o build completo (UI + Servidor) e **executa o binário** automaticamente.
  - `make dev`: Agora realiza o build da UI e **executa o bot via `go run`**, facilitando o desenvolvimento do backend.
  - Removido o alvo `run` que se tornou redundante e atualizada a ajuda (`make help`).
- **Dockerfile Multi-stage Build**: Implementada construção em múltiplos estágios (Node.js para frontend e Go para backend), automatizando a geração dos arquivos estáticos da dashboard e reduzindo o tamanho da imagem final.
- **Otimização do Makefile**:
  - O comando `make` agora realiza apenas o build por padrão (em vez de build e run), evitando bloqueios na CLI.
  - Adicionada detecção inteligente de dependências do frontend para evitar `npm install` desnecessários em cada build.
- **Otimização de Regex**: Pré-compilação de todas as expressões regulares globais (formatação de texto e detecção de Markdown/HTML) para reduzir uso de CPU.
- **Refatoração do Comando /channels**: Melhoria na UX com resposta direta (reply) e edição da mensagem de status em tempo real.
- **Ajustes de Infraestrutura**: Conversão do banco de dados legado de PostgreSQL para SQLite para facilitação do desenvolvimento local.
- **Processamento de Canais**: Implementada verificação `via_bot` para que o bot ignore postagens enviadas por ele mesmo via modo inline, prevenindo loops de edição ou processamento duplicado.

### Fixed
- **Comando "Sobre"**: Corrigido erro de "falta de ação" ao clicar no botão "Sobre" devido a uma tag `<blockquote>` não fechada no arquivo de mensagens.
- **Lógica de Versão**: Melhorada a detecção da versão do bot para incluir um fallback automático para o comando `git rev-parse --short HEAD` caso o binário não tenha sido compilado com as flags de versão ou não contenha informações de VCS.
- **Log de Erro de Callback**: Adicionado logging de erros ao handler do comando "Sobre" para facilitar o diagnóstico de falhas na API do Telegram.
- **Log de Porta da API**: Corrigida a exibição da URL no log que apresentava dois pontos extras (ex: `http://localhost::7000`).
- **Camada de UI (Z-Index)**: Corrigido problema onde o toast de confirmação e modais ficavam atrás do menu inferior (TabBar) e de outros elementos, garantindo visibilidade total em todas as resoluções.
- **Payload de Permissões**: Corrigida falha no dashboard que impedia a ativação de novas permissões devido à falta do campo `document` no payload enviado para a API.
- **Contexto de Canal**: Corrigida falha no middleware `AuthorizeChannel` que impedia a injeção do `channelID` no contexto para usuários com cargo Admin/Owner, resolvendo o erro "channelID invalido no contexto" na desconexão de canais.
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
