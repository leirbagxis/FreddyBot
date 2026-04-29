# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Novo Sistema de Log Centralizado**: Implementação de um logger customizado em `pkg/logger` com suporte a cores e módulos (`[BOT]`, `[API]`, `[DB]`, `[ADMIN]`, `[SYS]`).
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

### Fixed
- Erro de compilação em `admin.go` relacionado à mudança de campos na biblioteca `go-telegram/bot`.
- Vazamentos potenciais de memória em loops de processamento de texto.

### Added (Reconstrução Anterior)
- **Suporte ao SQLite3**: Agora é possível alternar entre SQLite e PostgreSQL usando a variável de ambiente `APP_ENV`.
...
