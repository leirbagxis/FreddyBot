# Plano: Adicionar Separadores Verticais Explícitos (|) entre os Botões de Ação

## Pedido do usuário
O usuário notou que os separadores verticais (`|`) entre os botões de ação (Editar, Fixar, Pausar, Excluir) não ficaram visíveis ou destacados o suficiente. Ele solicita a inclusão clara e visível dessas linhas divisórias.

## Objetivo
Inserir separadores verticais explícitos (`<div className="w-px h-7 bg-border/80 self-center shrink-0" />`) entre os 4 botões de ação do rodapé do card de agendamento em `ScheduleTab.tsx`.

## Contexto atual
- `divide-x divide-border/60` no container grid flexível em Tailwind escuro nem sempre gera contraste visível nos navegadores mobile.
- Substituir por elementos divisores estáticos explícitos (`<div className="w-px h-7 bg-border/80 self-center shrink-0" />`) entre os botões dentro de um flexbox/grid garante visibilidade total em qualquer tema.

## Arquivos analisados
- `dashboard/src/components/ScheduleTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/ScheduleTab.tsx`

## Estratégia de implementação
Substituir o rodapé de botões em `ScheduleTab.tsx` por uma estrutura flexível com separadores verticais reais:

```tsx
<div className="border-t border-border/60 pt-2.5 mt-1 flex items-center justify-between">
  <button type="button" className="...">
    <Edit3 size={16} />
    <span>Editar</span>
  </button>

  <div className="w-px h-7 bg-border/80 shrink-0" />

  <button type="button" className="...">
    <Pin size={16} />
    <span>Fixar</span>
  </button>

  <div className="w-px h-7 bg-border/80 shrink-0" />

  <button type="button" className="...">
    <Pause size={16} />
    <span>Pausar</span>
  </button>

  <div className="w-px h-7 bg-border/80 shrink-0" />

  <button type="button" className="...">
    <Trash2 size={16} />
    <span>Excluir</span>
  </button>
</div>
```

## Passos detalhados
1. Atualizar `ScheduleTab.tsx` inserindo os separadores verticais explícitos `w-px h-7 bg-border/80`.
2. Executar `cd dashboard && npm run build` para validar a compilação.

## Riscos
- Nenhum risco.

## Impactos esperados
- Linhas verticais separadoras `|` visíveis e destacadas entre todos os botões de ação.

## Compatibilidade
- Linux, macOS, Windows, Mobile (iOS/Android)

## Como testar

### Build
```bash
cd dashboard && npm run build
```

## Rollback
Reverter as alterações em `ScheduleTab.tsx`.
