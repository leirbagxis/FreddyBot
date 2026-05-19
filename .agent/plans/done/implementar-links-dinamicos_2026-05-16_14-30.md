# Plano: implementar-links-dinamicos_2026-05-16_14-30.md

## Pedido do usuário
Implementar a funcionalidade de "links dinâmicos". O usuário deseja um toggle na dashboard que, quando ativado, faça o bot ler a legenda de mensagens em canais e transforme padrões específicos em botões.
Padrões suportados:
1. `!Nome do Botão` (linha 1) e `!url do botão` (linha 2).
2. Links embutidos no texto (hiperlinks HTML ou Markdown).
O bot deve remover esses padrões do texto original e adicioná-los como botões de URL.

## Objetivo técnico
Adicionar suporte à detecção e extração de botões dinâmicos a partir do texto/legenda das postagens em canais, limpando o texto final e preservando a integridade das outras configurações de legenda.

## Contexto atual
O sistema já possui um pipeline de processamento em `internal/telegram/events/channelPost` que lida com legendas padrão, hashtags e custom captions. A lógica de transformação ocorre em `stage_transform.go`.

## Arquivos analisados
- `internal/database/models/models.go`
- `internal/database/repositories/channel.go`
- `internal/core/services/channels.go`
- `internal/telegram/events/channelPost/stage_transform.go`
- `internal/telegram/events/channelPost/formatting.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- **Backend:**
    - `internal/database/models/models.go` (Adição do campo `DynamicLinks`)
    - `internal/database/repositories/channel.go` (Método de atualização)
    - `internal/core/services/channels.go` (Lógica de serviço)
    - `internal/api/controllers/channelController.go` (Endpoint da API)
    - `internal/api/routes/routes.go` (Registro da rota)
    - `internal/telegram/events/channelPost/stage_transform.go` (Lógica de extração)
- **Frontend:**
    - `dashboard/src/types.ts` (Interface do Canal)
    - `dashboard/src/api.ts` (Chamada de API)
    - `dashboard/src/App.tsx` (Estado e renderização do toggle)
    - `dashboard/src/components/DashboardInicioTab.tsx` (Adição do card de configurações)

## Estratégia de implementação
1. **Banco de Dados:** Adicionar `DynamicLinks` ao modelo `Channel`.
2. **API:** Criar endpoint `PUT /api/channel/:channelId/dynamic-links` para alternar a funcionalidade.
3. **Frontend:** Adicionar um toggle no dashboard para o dono do canal.
4. **Pipeline Telegram:** 
    - No `StageTransform`, após processar a legenda base:
    - Se `DynamicLinks` estiver ativo:
        - Usar Regex para encontrar padrões `!Nome\n!URL`.
        - Usar Regex para encontrar tags `<a>` (visto que o texto já estará em HTML).
        - Converter os resultados em objetos `models.Button`.
        - Remover os padrões e links do `FormattedText`.
        - Adicionar os novos botões ao `FinalButtons`.

## Passos detalhados

### 1. Preparação do Banco e API
- Adicionar campo `DynamicLinks` em `internal/database/models/models.go`.
- Implementar `UpdateDynamicLinks` no repositório e serviço de canais.
- Criar controller e registrar rota na API.

### 2. Implementação no Dashboard
- Atualizar `types.ts` e `api.ts`.
- Adicionar o campo `dynamicLinks` no componente de permissões ou em um novo card de "Configurações Extras" na aba de permissões/configurações do canal.
- Integrar a chamada de API no hook de salvamento de configurações.

### 3. Lógica de Extração (Pipeline)
- Criar funções auxiliares em `formatting.go` para extrair links dinâmicos:
    - `ExtractBangLinks(text string) ([]Button, string)`
    - `ExtractEmbeddedLinks(text string) ([]Button, string)`
- Integrar no `StageTransform` em `stage_transform.go`.

## Riscos
- **Falso Positivo:** Usuários podem usar `!` no início de linhas para outros fins. *Mitigação: Exigir que a URL comece com http/https.*
- **Conflito de Formatação:** Remover tags `<a>` pode quebrar a estrutura do HTML se não for feito com cuidado. *Mitigação: Usar parsing de string seguro.*

## Impactos esperados
- Melhora na agilidade de postagem para canais que já usam esses padrões.
- Redução da necessidade de configurar botões manualmente na dashboard para cada postagem.

## Compatibilidade
- Backend: Go 1.24, GORM (SQLite/Postgres).
- Frontend: React + Vite.

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Testes
1. Ativar Links Dinâmicos no Dashboard.
2. Enviar mensagem no canal com:
   ```txt
   Legenda legal
   !Botão Site
   !https://google.com
   ```
3. Verificar se o bot envia o botão e remove as linhas com `!`.

## Rollback
Desativar o toggle no dashboard ou reverter as alterações no `stage_transform.go`.

## Observações
- A prioridade será dada aos links embutidos e depois aos padrões `!`.
- Os botões extraídos serão adicionados ao final da lista de botões existentes.
