# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Suporte ao SQLite3**: Agora é possível alternar entre SQLite e PostgreSQL usando a variável de ambiente `APP_ENV`.
  - `APP_ENV=dev`: Utiliza SQLite (ideal para desenvolvimento local).
  - `APP_ENV=prod`: Utiliza PostgreSQL (ideal para produção).
- **Funcionalidade Post Builder**: Um sistema interativo para criar postagens personalizadas com mídias.
  - Detecção automática de mídias (Fotos, Vídeos, GIFs, Áudios e Documentos).
  - Edição de Título, Corpo e Rodapé com suporte a toda a formatação original do Telegram.
  - Adição rápida de botões (Nome e Link em uma única mensagem).
  - Função de **Preview** para visualizar a postagem antes de salvar, sem fechar o editor.
  - Opção de **Salvar** para finalizar a postagem e encerrar a sessão.
- Novas variáveis no `.env-example`: `APP_ENV`.

### Fixed
- Erro ao tentar enviar mídias sem botões configurados (Bad Request no Telegram).
- Problema de interface nula no `ReplyMarkup` que impedia o envio de mídias em certos casos.
- Preservação de entidades de formatação (negrito, itálico, etc.) ao capturar textos do usuário para o Post Builder.
