# Plano: Alterar Ícone de "Meus Canais" no SideMenu

## Pedido do usuário
Substituir o ícone atual de "Meus Canais" (`Hash`) por um ícone mais adequado de canal/transmissão (ex: `Radio` ou `Tv`).

## Objetivo
Atualizar a importação do `lucide-react` em `SideMenu.tsx` substituindo o ícone `#` (`Hash`) pelo ícone de sinal de transmissão (`Radio`).

## Contexto atual
- O item "Meus Canais" em `SideMenu.tsx` está utilizando `<Hash size={18} />`.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Substituir `Hash` por `Radio` nas importações de `lucide-react`.
   - Substituir `<Hash size={18} className="text-foreground/80 shrink-0" />` por `<Radio size={18} className="text-foreground/80 shrink-0" />`.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Substituição pontual de um ícone em `SideMenu.tsx`.

## Impactos esperados
- Ícone moderno de sinal/transmissão (`Radio`) representando os canais do usuário no menu lateral.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/SideMenu.tsx
```
