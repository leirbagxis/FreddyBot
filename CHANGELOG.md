# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Reconstrução Total do Webapp**: Transformação do frontend de um sistema de economia para um Gerenciador de Canais do Telegram robusto.
  - **Dashboard de Canais**: Nova interface que lista automaticamente os canais gerenciados pelo usuário.
  - **Editor de Canal Detalhado**:
    - Gestão de Legendas (Padrão e Novos Packs).
    - CRUD completo de Botões com controle de Layout (Posições X e Y).
    - Aba de Permissões para controle granular de mídias (Preview, Foto, Vídeo, GIF, etc.) com atualizações em tempo real.
  - **Painel Admin**: Refatorado para monitoramento global de todos os usuários e canais do sistema.
  - **Integração Backend**: Autenticação nativa com Telegram WebApp e sincronização total com as rotas do backend Go.
  - **Estética Cyberpunk/Terminal**: Interface unificada com tema Dark, cores neon, animações de scanline e fontes mono.
- Novas tipagens TypeScript baseadas diretamente nos modelos Go (`models.go`).

### Changed
- Refatoração completa da camada de API (`api.ts`) para mapear as novas rotas do servidor.
- Atualização do fluxo de autenticação para utilizar cookies JWT seguros.

### Removed
- Antigo sistema de Economia (Loja, Ranking, Inventário e Itens).
- Código legado e componentes de UI que não pertenciam ao escopo de gerenciamento de canais.

### [0.1.0] - 2024-04-28 (Anterior)
### Added
- **Suporte ao SQLite3**: Agora é possível alternar entre SQLite e PostgreSQL usando a variável de ambiente `APP_ENV`.
...
