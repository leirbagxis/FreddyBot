# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Post Builder (Novas Funcionalidades)**:
  - **Suporte a Stickers**: Agora é possível criar postagens a partir de stickers, com suporte total a botões e reações (sem legenda, conforme limitação do Telegram).
  - **Enviar para Canais**: Novo fluxo pós-salvamento que permite enviar a postagem diretamente para qualquer canal configurado pelo usuário no bot, sem depender exclusivamente do modo inline.
  - **Importação de Canal**: Recurso para copiar legendas, reações e botões de canais já cadastrados diretamente para o PostBuilder.
  - **Gerenciador de Botões**: Interface dedicada para visualização e exclusão individual de botões durante a montagem do post.
- **Conversão JIT HTML**: Implementada conversão automática de Markdown e Entidades do Telegram para HTML no PostBuilder, garantindo maior fidelidade visual e estabilidade no envio.

### Fixed
- **Estabilidade do PostBuilder**: Corrigidos erros de "Bad Request" ao lidar com caracteres especiais em MarkdownV2 através da migração para ParseMode HTML.
- **Persistência de Sessão**: Resolvido bug de "Sessão expirada" ao realizar ações após o salvamento da postagem.
- **Layout de Teclado**: Otimização do posicionamento de botões de ação para melhor usabilidade em dispositivos móveis.
- **Interceptação Global de Stickers**: Corrigido comportamento do bot ao receber stickers fora do contexto de comandos.

### Changed
- **Arquitetura PostBuilder**: Refatoração interna para maior modularidade e suporte a múltiplos tipos de mídia de forma consistente.
- **Cache de Sessão**: Otimização do tempo de vida e estrutura de dados das sessões temporárias do PostBuilder no Redis.


### Changed
- **Limpeza do Legado V1 (Motor V2 Puro)**: A estrutura legada `MessageProcessor` foi eliminada. Toda a sua lógica de despacho, formatação e metadados foi convertida em funções independentes e integradas nativamente aos estágios do Pipeline V2, tornando a arquitetura 100% modular.
- **Segurança (.gitignore)**: Atualização massiva das regras de ignorância para proteger arquivos sensíveis (`.env`), bancos de dados temporários e diretórios de ferramentas de automação.
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
