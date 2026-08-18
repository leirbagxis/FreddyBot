# Plano: Desacoplar Agendamentos dos Canais para a Dashboard Inicial

## Pedido do usuário
Desacoplar o Agendamento das páginas individuais de configuração dos canais. Ele deve ficar localizado na dashboard inicial que exibe todos os canais, com uma div/seção "Agendamentos" onde o usuário possa visualizar os agendamentos de todos os seus canais.

## Objetivo
1. Criar um card/seção "Agendamentos" na Dashboard Inicial (`isChannels` / home) que carrega os agendamentos globais de todos os canais do usuário.
2. Permitir que o usuário visualize e gerencie todos os seus agendamentos consolidados a partir da dashboard inicial.
3. Conectar a opção "Agendamentos" do `SideMenu` diretamente a essa visão global de agendamentos.
4. Desacoplar a aba "Agendamentos" do menu inferior dos canais individuais (`BASE_TABS`).

## Contexto atual
- Atualmente, a aba `Agendamentos` fica dentro da página de configuração de um canal específico (`/dashboard/:channelID`).
- A API `fetchMySchedules()` já retorna a lista de todos os agendamentos do usuário autenticado para todos os seus canais.
- O componente `ScheduleTab.tsx` pode ser flexibilizado para filtrar por `channelId` quando informado, ou exibir todos os agendamentos quando não informado.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/components/ScheduleTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/components/ScheduleTab.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/ScheduleTab.tsx`**:
   - Tornar a prop `channelId` opcional (`channelId?: number`).
   - Se `channelId` não for fornecido, carregar todos os agendamentos do usuário via `fetchMySchedules()` sem filtrar por canal.
   - Exibir a identificação do canal (título do canal) em destaque em cada item quando em modo global.
   - Atualizar a mensagem de estado vazio ("Nenhum agendamento encontrado nos seus canais").

2. **`dashboard/src/App.tsx`**:
   - Adicionar o estado `showSchedulesModal` (visão dedicada a agendamentos globais na dashboard principal).
   - Na visão inicial (`isChannels`), adicionar a div/card de ação "Agendamentos" entre as ações principais:
     - Título: "Agendamentos"
     - Descrição: "Visualize e gerencie os agendamentos de todos os seus canais"
     - Ao clicar: abre a visão com `<ScheduleTab />` sem filtro de canal.
   - Atualizar `BASE_TABS` para remover a aba "Agendamentos" das páginas de configuração individual do canal.

3. **`dashboard/src/components/SideMenu.tsx`**:
   - Atualizar a propriedade `onNavigateChannels` ou adicionar `onOpenSchedules` para que o clique no item "Agendamentos" abra diretamente a visão global de agendamentos na dashboard inicial.

4. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/ScheduleTab.tsx`.
2. Editar `dashboard/src/components/SideMenu.tsx`.
3. Editar `dashboard/src/App.tsx`.
4. Executar `npm run build`.

## Riscos
- **Baixo**: Reorganização de UI e desacoplamento de rotas/componentes já existentes.

## Impactos esperados
- Os agendamentos ficam centralizados na dashboard inicial para todos os canais, tornando o gerenciamento mais rápido e prático.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/ScheduleTab.tsx
git checkout dashboard/src/components/SideMenu.tsx
git checkout dashboard/src/App.tsx
```
