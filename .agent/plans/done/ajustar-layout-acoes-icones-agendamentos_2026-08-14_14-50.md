# Plano: Ajuste no Layout de Ações e Ícones da Página de Agendamentos

## Pedido do usuário
1. Reorganizar os botões de ação do rodapé do card de agendamento com o **ícone em cima** e o **nome da função embaixo**, separados visualmente por divisores verticais (`|`).
2. Mudar a cor dos ícones da seção **Próxima Postagem** (`Clock`) e **Janela de Ciclo** (`RotateCw`) para um **azul mais escuro**.
3. Adicionar/reforçar um **divisor horizontal claro** entre as seções de informações de horário do card.

## Objetivo
Refinar o card de agendamento para alinhar exatamente à preferência visual solicitada: ícone acima do texto nas ações com separadores verticais e ícones de tempo com destaque em tom azul elegante.

## Contexto atual
- Em `dashboard/src/components/ScheduleTab.tsx`, as ações atualmente estão dispostas em linha horizontal (ícone ao lado do texto).
- Os ícones de `Clock` e `RotateCw` estão em caixas neutras cinzas (`bg-muted/30 text-muted-foreground`).

## Arquivos analisados
- `dashboard/src/components/ScheduleTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/ScheduleTab.tsx`

## Estratégia de implementação

### 1. Botões de Ação no Rodapé
- Alterar o layout dos botões para orientação vertical: `flex flex-col items-center justify-center gap-1 py-2 text-[10px] font-semibold`.
- Manter o container da barra de ações em `grid grid-cols-4 divide-x divide-border/60 pt-2 border-t border-border/60` para que cada botão fique perfeitamente centralizado com a linha vertical `|` separando-os.

### 2. Ícones em Azul Mais Escuro
- Atualizar os containers de ícone de **Próxima postagem** (`Clock`) e **Intervalo de ciclo** (`RotateCw`) para um tom de azul mais forte e elegante: `bg-blue-500/15 text-blue-500 border border-blue-500/25` (ou azul escuro `#1f70b2`).

### 3. Divisores Horizontais
- Garantir uma linha divisória horizontal limpa entre o bloco de **Próxima Postagem** e o bloco de **Intervalo de Ciclo** (`divide-y divide-border/60`).

## Passos detalhados
1. Atualizar o componente `ScheduleTab.tsx` com os estilos de ícone azul escuro e botões de ação em formato de ícone no topo + legenda abaixo com dividor `|`.
2. Executar `cd dashboard && npm run build` para validar.

## Riscos
- Nenhum risco. Ajuste puramente estético de frontend.

## Impactos esperados
- Ações no rodapé mais fáceis de tocar e identificar visualmente.
- Melhor contraste nos ícones de tempo e divisão limpa das seções.

## Compatibilidade
- Linux, macOS, Windows, Mobile (iOS/Android)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Reverter as alterações em `ScheduleTab.tsx`.
