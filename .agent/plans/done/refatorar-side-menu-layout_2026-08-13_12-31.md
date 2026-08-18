# Plano: Refatorar Layout do Menu Lateral (SideMenu.tsx)

## Pedido do usuário
1. **Remover o cabeçalho superior**: Tirar a palavra "MENU" e o ícone "X" de fechar.
2. **Largura mais compacta**: Deixar a largura lateral do drawer mais curta (de `w-80` para `w-64`).
3. **Refatorar a div "Minha Conta"**:
   - Remover o ícone de pessoa de dentro da div.
   - Deixar a div mais alta/ampla (maior padding).
   - Remover a cor preta da borda e aplicar um tom de branco suave (`border-white/15`).
   - Exibir o nome "Minha Conta" bem pequeno no canto superior esquerdo.
   - Exibir o número de canais do usuário bem grande e destacado no centro da div.

## Objetivo
Ajustar o visual do componente `SideMenu.tsx` tornando-o mais compacto na largura e refinando a área "Minha Conta" para dar destaque ao número de canais em fonte grande e centralizada.

## Contexto atual
- O componente `SideMenu.tsx` possui atualmente `w-80`, um cabeçalho com "Menu" e o ícone de fechar "X", e um card compacto de "Minha Conta" com um ícone de usuário.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Alterar container principal para `w-64 max-w-[75vw]`.
   - Remover o bloco com a palavra "Menu" e o botão `<X size={20} />`.
   - Reestruturar o container "Minha Conta":
     - `p-5 rounded-2xl bg-muted/30 border border-white/15 text-left`
     - Canto superior esquerdo: `<span className="text-[11px] font-bold text-muted-foreground uppercase tracking-wider">Minha Conta</span>`
     - Centro: `<div className="text-4xl font-extrabold text-foreground text-center my-2">{channelsCount}</div>` + `<p className="text-[11px] text-muted-foreground text-center">canais cadastrados</p>`
   - Manter o topo da foto circular + nome + ID ao lado, perfeitamente centralizados verticalmente.
   - Manter os links de navegação rápida com fechamento ao clicar no overlay/backdrop.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Mudança de estilização e layout no componente React `SideMenu.tsx`.

## Impactos esperados
- Menu lateral mais estreito e minimalista, sem botão X ou título "MENU", com um card "Minha Conta" visualmente impactante focado na contagem de canais.

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
