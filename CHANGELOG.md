# Changelog

All notable changes to this project will be documented in this file.

## [1.4.0] - 2026-05-20

### Added
- **Auditoria Ativa e Gerenciamento em Massa**:
  - Nova aba **Auditoria** no Dashboard Admin para varredura em tempo real de permissões de bots legados (@XavolaBot).
  - Implementação de **Deleção em Massa (Bulk Delete)** de canais diretamente pela interface de auditoria, com confirmação de segurança.
  - Otimização da varredura com **20 workers paralelos** no backend, garantindo respostas rápidas mesmo para grandes volumes de dados.
- **Suporte e Comunicação Direta**:
  - Nova funcionalidade de **Mensagem Individual** na aba Broadcast, permitindo contatar usuários específicos via ID.
  - Injeção automática de cabeçalho oficial (`# 📨 MENSAGEM DO SUPORTE`) em disparos individuais.
  - Comando bot **`/getid`**: Facilidade para administradores capturarem o File ID de qualquer mídia (foto, vídeo, gif) para uso imediato em campanhas de broadcast.
- **Preview de Mídia Real (Dashboard)**:
  - Implementação de um **Media Proxy** seguro no backend (`/api/admin/media-proxy`) que permite ao Dashboard exibir imagens reais do Telegram sem expor o Token do bot ao navegador.
  - Ajuste inteligente no preview para detectar links externos vs. File IDs e exibir a imagem completa sem recortes (`object-contain`).

### Fixed
- **Compatibilidade PostgreSQL**:
  - Implementação de **Limpeza Manual em Cascata** no repositório de canais através de transações SQL. Isso resolve erros de violação de chave estrangeira (`23503`) no Postgres, garantindo que botões e legendas sejam limpos antes da remoção do canal.
  - Normalização de tipos booleanos e IDs para suporte transparente entre SQLite (dev) e Postgres (prod).
- **Fim do Hardcode**: Removidas todas as referências fixas a @LegendasBrBot do código-fonte e do banco de dados inicial, substituindo-as por placeholders dinâmicos como `{botUser}` e `{usernameBot}`.

### Changed
- **Modularização do Core Admin**: O arquivo massivo `admin.go` foi fragmentado em múltiplos arquivos especializados (`admin_users.go`, `admin_channels.go`, `admin_broadcast.go`, `admin_utils.go`), facilitando a manutenção e futuras expansões.
- **Padronização de Mensagens**: Unificação do estilo visual e uso de emojis em todo o arquivo `messages.yml`.

## [1.3.1] - 2026-05-19

### Added
- **Design Dashboard (Sentri Soft + Liquid Glass)**:
  - Implementação do tema **Sentri Soft** (Violeta/Indigo) focado em conforto visual e modernidade.
  - Adição do efeito **Liquid Glass** de alta performance, utilizando transparência seletiva e brilho de borda (specular highlight).
  - Nova animação fluida para troca de tema (claro/escuro) com rotação de ícones e transição de fundo.
- **Identidade Visual**:
  - Unificação do card de identidade do usuário, integrando saudação, ID e o botão de desconexão em um único local limpo e funcional.

### Fixed
- **Performance Extrema (Dashboard)**:
  - Remoção das bibliotecas pesadas **GSAP**, **Three.js** e **Lenis** (smooth scroll) para garantir 60fps constantes no navegador do Telegram.
  - Refatoração profunda do motor React: Implementação massiva de `memo` e `useCallback` para eliminar re-renders desnecessários.
  - Otimização de CSS: Troca de `transition: all` por transições cirúrgicas e redução drástica do uso de `backdrop-filter`.
- **Estabilidade da UI**:
  - Correção de erro fatal (tela preta) na aba de Permissões devido a ícones não importados e tratamento de dados nulos.
  - Implementação de proteção contra falhas (`optional chaining`) em todos os cards de configuração.
- **Navegação**:
  - Ajuste na lógica do botão de voltar nativo do Telegram para navegação consistente entre Meus Canais e Painel Admin.

### Changed
- **Arquitetura de Componentes**: Mover estados locais (como inputs de transferência) para dentro dos componentes filhos, reduzindo a carga no componente principal `App.tsx`.
- **Simplificação de UX**: Remoção do sistema de "Revelar Link" e HUD de Telemetria complexo em favor de uma interface mais direta e rápida.

## [1.3.0] - 2026-05-19

## [1.2.0] - 2026-05-15

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
