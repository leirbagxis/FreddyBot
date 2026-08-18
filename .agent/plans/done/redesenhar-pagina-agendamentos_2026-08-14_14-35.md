# Plano: Redesenhar Página de Agendamentos (Schedule UI Refinement)

## Pedido do usuário
Redesenhar e refinar a página de agendamentos (`ScheduleTab.tsx` / visualização "Todos os Agendamentos") no padrão de Telegram Mini App do FreddyBot, melhorando a hierarquia de informações, visibilidade de status, contraste de datas, ações táticas de toque mobile e botão flutuante (FAB) de novo agendamento, com confirmação obrigatória de exclusão.

## Objetivo
Transformar o card de agendamento de uma linha genérica em uma estrutura visualmente polida, com cabeçalho limpo, seções de próxima postagem e intervalo bem divididas, botões de ação confortáveis para toque e confirmação modal ao deletar.

## Contexto atual
- O arquivo `dashboard/src/components/ScheduleTab.tsx` exibe os cards de agendamento em uma única linha horizontal amontoada, com status brutos em inglês (ex: `pending`, `paused`).
- Em `App.tsx`, o cabeçalho exibe `"Agendamentos de Todos os Canais"` sem descrição.
- O botão de excluir deleta o agendamento imediatamente sem pedir confirmação.
- Não há botão flutuante de criação nem formatação sutil dos rótulos de tempo.

## Arquivos analisados
- `dashboard/src/components/ScheduleTab.tsx` — Lista de agendamentos e card individual
- `dashboard/src/App.tsx` — Cabeçalho e contêiner da página "Todos os Agendamentos"
- `dashboard/src/components/ConfirmModal.tsx` — Modal de confirmação reutilizável

## Arquivos que poderão ser modificados
- `dashboard/src/components/ScheduleTab.tsx`
- `dashboard/src/App.tsx`

## Estratégia de implementação

### 1. Ajuste do Cabeçalho da Página (`App.tsx`)
- Substituir o texto `"Agendamentos de Todos os Canais"` por **`Todos os Agendamentos`**.
- Adicionar descrição mutada abaixo: `"Visualize e gerencie todos os agendamentos de postagens."`.
- Manter o botão de voltar à esquerda (`<ArrowLeft />`).

### 2. Seção de Agendamentos Ativos (`ScheduleTab.tsx`)
- Título da seção: **`Agendamentos Ativos`**.
- Contador discreto em badge arredondada à direita (ex: `Agendamentos Ativos` ... `1`).

### 3. Reorganização Estrutural do Card de Agendamento (`ScheduleCard`)
- **Cabeçalho do Card**:
  - Ícone de calendário em container quadrado arredondado com borda e fundo suave (`bg-accent/15 text-accent border border-accent/25`).
  - Nome do canal/agendamento como texto principal em destaque.
  - Tipo de agendamento em texto secundário mutado logo abaixo (`Intervalo`, `Único`, `Diário`, `Semanal`, `Fila`).
  - Badge compacta de status à direita com indicador de ponto (`● Pendente`, `● Ativo`, `● Pausado`, `● Concluído`, `● Falhou`).
  - Tradução completa dos status:
    - `pending` / `scheduled` → `Pendente`
    - `active` → `Ativo`
    - `paused` → `Pausado`
    - `completed` / `sent` → `Concluído`
    - `failed` / `error` / `cancelled` → `Falhou`

- **Linha Próxima Postagem**:
  - Ícone de relógio (`Clock`).
  - Rótulo mutado `"Próxima postagem"`.
  - Data/hora em destaque com contraste forte (ex: `"Hoje, 11:32"` ou `"15/08/2026 às 14:00"`).

- **Linha Intervalo / Janela**:
  - Ícone de repetição (`Repeat`).
  - Texto principal: `"A cada 5 minutos"`.
  - Sub-texto: `"Janela: 09:05 – 12:00"`.

- **Linha Divisória & Área de Ações (Rodapé)**:
  - Divisor sutil (`border-t border-border/60 my-3`).
  - Botões de ação com **Ícone + Texto** com áreas de toque mobile confortáveis:
    - `✎ Editar` (ícone azul, hover suave)
    - `📌 Fixar` / `📌 Desfixar` (ícone amarelo/laranja)
    - `Ⅱ Pausar` / `▶ Retomar` (ícone amarelo/laranja ou verde)
    - `🗑 Excluir` (ícone vermelho)

- **Confirmação de Exclusão**:
  - Utilizar o `ConfirmModal` para solicitar confirmação antes de remover qualquer agendamento (`"Tem certeza que deseja excluir este agendamento?"`).

### 4. Botão Flutuante de Novo Agendamento (FAB)
- Botão circular flutuante em Azul Telegram (`bg-[#2481cc] hover:bg-[#1f70b2] text-white size-14 rounded-full shadow-xl flex items-center justify-center font-bold text-2xl fixed bottom-6 right-6 z-40`).
- Respeita `env(safe-area-inset-bottom)` e posicionamento em relação à barra de navegação bottom.
- Ao ser clicado, exibe um modal amigável para criação de novos agendamentos no bot.

## Passos detalhados
1. Atualizar o cabeçalho em `App.tsx` com o novo título, descrição e layout.
2. Reformular o componente `ScheduleTab.tsx` implementando a nova hierarquia de cards, badge de status em português, formatador de datas ("Hoje, HH:mm"), botões de ação com ícone + rótulo e modal de confirmação de exclusão.
3. Adicionar o FAB azul de novo agendamento.
4. Testar o build com `cd dashboard && npm run build` para garantir que o TypeScript e o bundle estão 100% corretos.

## Riscos
- Nenhum risco ao backend. As chamadas de API (`fetchMySchedules`, `deleteSchedule`, `updateScheduleStatus`, `updateScheduleTime`) continuam idênticas.

## Impactos esperados
- Interface 100% alinhada ao Telegram Mini App FreddyBot.
- Hierarquia de informações limpa e fácil de ler no celular.
- Ações táticas de toque sem risco de exclusão acidental.

## Compatibilidade
- Linux, macOS, Windows, Mobile (iOS/Android via Telegram WebApp)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Restaurar os arquivos `ScheduleTab.tsx` e `App.tsx`.
