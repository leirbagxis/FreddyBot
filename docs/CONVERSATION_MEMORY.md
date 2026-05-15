# Memória de Desenvolvimento: Refatoração FreddyBot (Maio 2026)

Este documento serve como um resumo completo das mudanças realizadas para facilitar a continuidade do desenvolvimento em uma nova sessão ou chat.

## 1. Estado Atual do Projeto
O bot passou por uma grande refatoração de performance e arquitetura, saindo de um modelo monolítico para um modelo baseado em **Pipeline de Processamento (V2)**.

### Principais Tecnologias:
- **Backend:** Go (Golang) com biblioteca `go-telegram/bot`.
- **Banco de Dados:** PostgreSQL (Produção) / SQLite (Desenvolvimento) com GORM.
- **Cache:** Redis.
- **Frontend:** React + Vite + Tailwind CSS.

---

## 2. Mudanças Realizadas

### A. Performance e Banco de Dados
- **Índices GORM:** Adicionados índices em `username` (User), `created_at` (User), `updated_at` (Channel) e um índice composto em `position_x/y` (Button) para acelerar buscas e ordenações.
- **Cache de Duas Camadas (L1/L2):** Implementado cache em memória local (L1 via `go-cache`) integrado ao `ChannelService`. O bot agora consulta `RAM -> Redis -> SQLite`, reduzindo drasticamente o I/O.
- **Paginação Real:** Comandos `/users` e `/channels` no admin agora usam `LIMIT/OFFSET` no SQL, evitando carregar milhares de registros na RAM.
- **Upsert Atômico:** O método `UpsertUser` agora usa `ON CONFLICT` do GORM (1 chamada ao banco em vez de 2).
- **Eager Loading:** Rota de detalhes do canal na API agora carrega o dono (`Owner`) via `Joins`, eliminando consultas N+1.
- **Otimização de Memória:** Substituição massiva de concatenação de strings (`+=`) por `strings.Builder`.
- **Worker Pool:** Aumentado para **20 workers** e fila de **5.000 mensagens** para suportar picos em 380+ canais.

### B. O Novo Motor: ChannelPost V2 (Pipeline)
A lógica de postagem foi totalmente modularizada no pacote `internal/telegram/events/channelPost`.
- **Arquitetura:** Dividida em 6 Estágios (`StagePreflight`, `StageSpecialFlows`, `StageMediaGrouping`, `StageTransform`, `StageDecorate`, `StageSend`).
- **Sincronização de Metadados:** Implementado um sistema proativo (via eventos `MyChatMember`) e reativo (durante o post) para manter Nome e Link do canal sempre atualizados.
- **Debounce de Verificação:** O bot agora só chama `GetChat` no Telegram 1 vez por hora por canal (ou quando detecta mudança visual), economizando API.
- **Prioridade de Links:** Username público (@) agora tem prioridade absoluta sobre links privados.
- **Lógica de Legenda:** 
    - Texto: Adiciona legenda do banco embaixo.
    - Mídia/Álbum: Substitui a legenda original inteiramente pela do banco (paridade com V1).

### C. Correções de Build e Limpeza Arquitetural
- **Correção de Panic (Nil Pointer):** Resolvido erro de desreferenciamento de ponteiro nulo no `ToUserDTO` ao acessar o Dashboard. Adicionadas checagens de segurança no Mapper e corrigida a estratégia de carregamento (`Preload`) do Owner no repositório.
- **Remoção do MessageProcessor:** A estrutura legada `MessageProcessor` foi eliminada. Toda a sua lógica de despacho, formatação e metadados foi convertida em funções integradas nativamente aos estágios do Pipeline V2.
- **Frontend:** Reparados os arquivos `PermissionsCard.tsx` e `ButtonGrid.tsx` que estavam corrompidos.
- **Static Assets:** Corrigida a rota no Go para servir `/assets` corretamente do build do Vite.
- **Testes:** Ajustado o setup de testes unitários para suportar chaves estrangeiras no SQLite.

---

## 3. Estrutura de Arquivos V2
- `pipeline.go`: Orquestrador da linha de montagem.
- `context_v2.go`: Objeto de estado que viaja entre os estágios.
- `stage_*.go`: Lógicas isoladas de cada etapa do processamento.
- `dispatch_v2.go`: Centraliza o envio final e retentativas para o Telegram.
- `utils_v2.go`: Funções utilitárias compartilhadas.
- `TECHNICAL_V2_MOTOR.md`: Documentação técnica detalhada da arquitetura.

---

## 4. Pendências e Próximos Passos
- **Monitoramento de Cache:** Observar o uso de memória do cache L1 se o número de canais ativos crescer exponencialmente (TTL atual: 5 min).
- **Escalabilidade:** O motor V2 está pronto para suportar milhares de canais; o bot já utiliza PostgreSQL em produção, o que garante suporte a alta concorrência. Próximo passo seria otimização de queries lentas caso o volume de dados se torne massivo.

---
**Data:** 15 de Maio de 2026  
**Status:** Build Verde, Arquitetura V2 Limpa e Otimizada.
