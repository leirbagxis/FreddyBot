# Plano: redesenhar-broadcast-alertas-admin

## Pedido do usuário

Adicionar mais dados e notificações úteis para administrador/owner e redesenhar a aba Broadcast com uma interface minimalista e mais clara.

## Objetivo

Transformar a visão geral administrativa em um painel operacional com métricas e alertas acionáveis, e redesenhar o Broadcast como um fluxo de composição, público, revisão e pré-visualização coerente com o Minimalism UI já usado no admin.

## Contexto atual

- O admin carrega `users` e `channels` reais por `GET /api/admin/overview`.
- A visão geral atual exibe quatro métricas e uma fila derivada de admins, blacklist e usuários sem canais.
- A aba Broadcast envia a requisição para `POST /api/admin/notice` e o servidor enfileira o envio em segundo plano. A API não expõe histórico, progresso, fila, entrega ou uma central persistente de alertas.
- Há um componente legado de alertas, mas ele não está conectado à visão geral atualmente usada (`OperationsOverview`).
- O design do admin possui tokens e convenções minimalistas escopados a `.admin-layout-v2`: papel claro, linhas discretas, tipografia neutra, botões escuros e raios de 8–10px.

## Arquivos analisados

- `.agent/context.md`
- `.agent/memory/memory.md`
- `dashboard/src/index.css`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/components/admin/OperationsOverview.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/crmSelectors.ts`
- `dashboard/src/components/admin/AlertsPanel.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `internal/api/controllers/adminController/getAllUserAdminController.go`
- `internal/api/routes/routes.go`
- `internal/container/appContainer.go`

## Arquivos que poderão ser modificados

- `dashboard/src/components/admin/OperationsOverview.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/crmSelectors.ts`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `dashboard/src/index.css`
- `.agent/memory/memory.md`
- `.agent/decisions.md` (somente se a derivação de alertas for registrada como decisão relevante)

## Estratégia de implementação

Usar exclusivamente os dados já fornecidos pela visão geral para evitar métricas fictícias. A central de notificações será in-app e derivada desses dados: usuários bloqueados, usuários sem canais, novos usuários recentes e base sem ativação. Cada alerta terá contagem, gravidade e navegação para a área apropriada.

O Broadcast será reestruturado em duas colunas no desktop e uma coluna no celular: um fluxo principal de composição e segmentação à esquerda, e um resumo de alcance com pré-visualização Telegram à direita. Contagens de destinatários serão estimadas localmente para os públicos globais e para listas de IDs válidos; o texto deixará explícito que o backend apenas confirma o início do processamento, não a entrega final.

Não serão criados push notifications, mensagens automáticas no Telegram, histórico de entregas ou novas métricas de servidor neste ciclo, pois não há armazenamento/API para sustentá-los com precisão. Isso evita apresentar dados ou notificações que não existem.

## Passos detalhados

1. Estender os seletores do CRM com cálculos puros para ativação, novos cadastros, usuários sem canais, blacklist e uma lista tipada de alertas operacionais.
2. Atualizar a visão geral com métricas adicionais, bloco de saúde/ativação e uma lista de alertas com CTA para filtrar ou abrir Usuários; incluir estados vazios claros.
3. Adicionar ao topo do admin uma central compacta de notificações derivadas, com contador, navegação por teclado, rótulos acessíveis e fechamento ao perder foco/pressionar Escape.
4. Reorganizar `AdminNoticeTab` em seções mínimas: alcance, mensagem, mídia, botões e revisão; adicionar resumo de destinatários e validações visíveis antes de abrir a confirmação.
5. Manter a prévia Telegram fixa apenas em telas largas, sem `dangerouslySetInnerHTML` adicional além do preview já sanitizado; garantir fallback de imagem e estados vazio, carregando, inválido e de limite de caracteres.
6. Criar estilos escopados ao admin com tokens existentes e breakpoints para celular, sem introduzir cores verdes ou superfícies translúcidas.
7. Executar typecheck e build do dashboard; verificar classes/semântica dos novos controles e revisar o resultado renderizado em desktop e celular, se o ambiente local puder iniciar com credenciais seguras.
8. Registrar a decisão de alertas derivados e atualizar a memória; mover este plano para `done` após a verificação.

## Riscos

- O número de destinatários é uma estimativa da base carregada; o servidor pode ter mudanças entre a abertura do painel e o processamento do broadcast.
- Não há confirmação de entrega por destinatário; a UI não deve afirmar que um disparo foi concluído, apenas iniciado.
- A base de dados pode ser grande, portanto os cálculos serão memorizados e lineares sobre os dados já carregados.
- A central in-app não substitui alertas por Telegram/e-mail. Esses canais exigiriam requisito de produto, persistência e API próprios.

## Impactos esperados

- Owner/admin passa a identificar rapidamente bloqueios, ausência de ativação e crescimento recente.
- Broadcast fica mais seguro de operar, mostrando o público e o que será enviado antes da confirmação.
- A interface ganha consistência visual com o restante do Minimalism UI e mantém responsividade.

## Compatibilidade

- Linux
- macOS
- Windows
- Docker
- CI/CD
- Desktop e celular

## Como testar

### Build

```bash
cd dashboard && npx tsc --noEmit && npm run build
```

### Testes

```bash
go test ./...
```

### Execução

```bash
make run
```

### Verificação visual

1. Abrir `/admin/dash` em largura desktop e celular.
2. Confirmar que alertas e contador refletem os dados reais da resposta de overview.
3. Alternar cada público do Broadcast, editar texto, mídia e botões; conferir contagem, preview, estados inválidos e confirmação.
4. Navegar a central de notificações por teclado e confirmar foco/fechamento.

## Rollback

Reverter apenas os arquivos de dashboard listados neste plano e restaurar o layout anterior do Broadcast; não haverá migração, mudança de banco ou alteração de Docker.

## Observações

- A implementação se limita a alertas dentro do painel. Caso o objetivo seja receber notificações no Telegram/e-mail, será necessário um plano separado para definir destinatários, limites, persistência, retries e preferência de entrega.
- Não modificar `docker-compose`, conforme orientação anterior de que ele é exclusivo de desenvolvimento.
