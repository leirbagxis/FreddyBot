# Plano: Redesenhar Página de Logs da Dashboard Admin

## Pedido do usuário
A página de logs está confusa, quase impossível entender o que é cada evento, de quem é e o que se trata. O usuário quer uma apresentação mais clara e legível.

## Objetivo
Redesenhar o componente `AdminLogsTab.tsx` para tornar cada entrada de log imediatamente compreensível, mostrando claramente: o que aconteceu, em qual canal, quem fez, quando, e qual o resultado.

## Contexto atual
- O componente `AdminLogsTab.tsx` exibe eventos em uma lista onde cada item mostra o `eventType` bruto (ex: `post_received`, `postbuilder_started`, `post_skipped`), uma badge de status, e uma badge de fonte.
- A função `eventLabel` apenas substitui underscores por espaços: `"post_received" → "post received"` — sem tradução, sem ícone, sem contexto.
- O `metadata` expandido mostra JSON cru em um `<pre>`, difícil de interpretar.
- Não há separação visual clara entre canal, ator (quem fez) e o que aconteceu.
- Os labels "Owner", "Actor", "Sessão", "Mensagem ID" estão em inglês misturado e sem ícones.

## Arquivos analisados
- `dashboard/src/components/AdminLogsTab.tsx` — renderização dos logs
- `dashboard/src/types.ts` — interface `ChannelEvent`
- `internal/core/services/channel_events.go` — constantes de source/status
- `internal/telegram/events/channelPost/event_logs.go` — eventos de canal
- `internal/telegram/handlers/events/postBuilder/event_logs.go` — eventos de PostBuilder

## Tipos de eventos registrados no backend

### Source: `channel_post`
| eventType | Status | Significado |
|---|---|---|
| `post_received` | info | Post recebido no canal |
| `post_skipped` | skipped | Post ignorado (via_bot, bot_own_message, maintenance, channel_not_found, owner_blacklisted, unsupported_message_type, bot_separator_message) |
| `permission_missing` | skipped | Bot sem permissões necessárias |
| `metadata_updated` | info/error | Título ou username do canal mudou |
| `dynamic_links_extracted` | info | Links dinâmicos extraídos |
| `caption_applied` | info | Legenda aplicada |
| `buttons_applied` | info | Botões aplicados |
| `post_processed` | success | Post processado com sucesso |
| `post_failed` | error | Falha ao processar post |

### Source: `post_builder`
| eventType | Status | Significado |
|---|---|---|
| `postbuilder_started` | info | PostBuilder iniciado |
| `postbuilder_field_updated` | info | Campo editado (título, corpo, etc.) |
| `postbuilder_button_added` | info | Botão adicionado |
| `postbuilder_button_deleted` | info | Botão removido |
| `postbuilder_preview_sent` | info/error | Prévia enviada |
| `postbuilder_saved` | success | Post salvo |
| `postbuilder_sent_to_channel` | success | Post enviado ao canal |
| `postbuilder_failed` | error | Falha |
| `template_deleted` | info | Template deletado |

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminLogsTab.tsx`

## Estratégia de implementação

### 1. Mapa de tradução e ícones por eventType
Criar um dicionário `eventConfig` que mapeia cada `eventType` para:
- **Título em português** (ex: `post_received` → "Post Recebido")
- **Descrição curta** (ex: "Novo post detectado no canal")
- **Ícone Lucide** adequado (ex: `Inbox`, `Send`, `AlertTriangle`, `Check`, `X`, `Eye`, `Wrench`)
- **Cor de acento** para o ícone

### 2. Card de evento redesenhado
Cada evento vira um card compacto que mostra:
- **Linha 1**: Ícone colorido + Título traduzido + Badge de status + Timestamp alinhado à direita
- **Linha 2**: Nome do canal (com link clicável) + ID do canal + "por UserID XXXXX"
- **Linha 3 (condicional)**: Mensagem de erro em vermelho se houver `errorMessage`
- **Expansão**: Metadados formatados em uma tabela legível de chave/valor em vez de JSON cru

### 3. Formatação inteligente do metadata
Em vez de mostrar JSON cru, criar uma função `formatMetadataReadable` que:
- Converte chaves snake_case para labels legíveis (ex: `media_type` → "Tipo de mídia")
- Exibe valores booleanos como "Sim"/"Não"
- Formata números e strings de forma amigável
- Destaca a `reason` de skips com label traduzido

### 4. Filtros melhorados
- Adicionar filtro por `eventType` (dropdown com os tipos traduzidos)
- Manter os filtros existentes (busca, canal, source, status, datas)

### 5. Agrupamento visual
- Usar cores distintas para cada source:
  - `channel_post` → Azul
  - `post_builder` → Roxo
- Ícone de status mais proeminente (bolinha maior ou badge mais visível)

## Passos detalhados
1. Criar mapa `eventConfig` com traduções, ícones e cores para cada `eventType`.
2. Criar mapa `skipReasonLabels` para traduzir razões de skip.
3. Criar função `formatMetadataReadable` para exibir metadata de forma legível.
4. Redesenhar o card de cada evento com hierarquia visual clara.
5. Adicionar filtro de eventType ao painel de filtros.
6. Compilar e verificar com `cd dashboard && npm run build`.

## Riscos
- Nenhum risco ao backend. Alterações são exclusivamente no frontend.
- Novos eventTypes futuros terão fallback para exibição genérica.

## Impactos esperados
- Logs imediatamente compreensíveis sem precisar expandir cada item.
- Identificação clara de canal, ator, evento e resultado.
- Metadata legível sem necessidade de ler JSON cru.

## Compatibilidade
- Linux, macOS, Windows, Docker

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Reverter as alterações em `AdminLogsTab.tsx`.

## Observações
- O backend não precisa ser alterado. Toda a melhoria é no frontend.
- O fallback para eventTypes desconhecidos garante compatibilidade futura.
