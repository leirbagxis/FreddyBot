# Changelog

All notable changes to this project will be documented in this file.

## [1.3.0] - 2026-05-19

### Added
- **Migração para Telego Framework**: Refatoração completa do motor do bot para utilizar o framework `telego`, proporcionando maior estabilidade e tipos nativos.
- **Placeholder Dinâmico `{usernameBot}`**: Implementação de substituição em tempo real do nome de usuário do bot nas legendas globais durante o cadastro de novos canais.
- **Proteção Legal e Docs**: Adição de arquivo de `LICENSE` (Proprietário) e atualização completa do `README.md` com instruções de deploy focadas em PostgreSQL e variáveis de ambiente.
- **Limpeza Profunda de Cache**: Implementação de purga automática de sessões e metadados tanto no Redis quanto na memória RAM local ao excluir ou desconectar um canal.

### Fixed
- **Prioridade de Interceptação**: Reordenação de handlers para garantir que fluxos ativos (ex: Sticker Separador) não sejam roubados pelo PostBuilder.
- **Renderização de Tags**: Adição de fallbacks automáticos para a tag `{channelName}`, evitando textos vazios em canais recém-criados.
- **Persistência de Modelos**: Correção no salvamento do Sticker Separador com geração obrigatória de UUID.
- **Parsing de YAML**: Ajuste no parser de mensagens para lidar com chaves simples e caracteres especiais de forma mais resiliente.

### Changed
- **Arquitetura de Middlewares**: Migração dos middlewares globais (Blacklist, Manutenção, SaveUser) para o padrão Telego.
- **Inicialização de Banco**: Injeção de templates padrão para legendas globais diretamente no `ServerConfig` durante o primeiro boot.

## [Unreleased]

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
