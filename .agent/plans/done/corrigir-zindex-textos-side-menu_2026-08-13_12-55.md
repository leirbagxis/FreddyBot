# Plano: Corrigir Z-Index, Rótulos de Texto e Bordas do SideMenu

## Pedido do usuário
1. **Z-Index Acima do Footer/TabBar**: O menu lateral e seu overlay devem ficar com `z-index` superior a qualquer footer/barra inferior (`z-[20000]`), para que a barra de navegação inferior não fique sobreposta ao menu.
2. **Remover "Meus" dos Rótulos**:
   - "Meus Templates" -> "Templates"
   - "Meus Agendamentos" -> "Agendamentos"
3. **Remover a Borda Direita da Gaveta**: Tirar a borda direita do container principal do menu (`border-r border-white/10`).
4. **Remover o Rodapé do Menu**: Tirar o texto "LegendasBOT Dashboard" da parte inferior do menu lateral.

## Objetivo
Garantir que o menu lateral seja renderizado com prioridade máxima na camada visual (`z-[20000]`), ajustando os rótulos de texto, removendo a borda externa direita e eliminando o rodapé do menu.

## Contexto atual
- `SideMenu.tsx` utiliza `z-50`, enquanto o `.bottom-nav` possui `z-1000`, fazendo com que a barra inferior do Telegram/Dashboard sobreponha o menu lateral.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Alterar `z-50` para `z-[20000]` na div raiz do `SideMenu`.
   - Remover a classe `border-r border-white/10` da gaveta lateral (`relative z-10 w-64 max-w-[75vw] h-full bg-background p-5 ...`).
   - Alterar rótulo "Meus Templates" para "Templates".
   - Alterar rótulo "Meus Agendamentos" para "Agendamentos".
   - Remover a div do rodapé `<div className="pt-3 border-t border-white/10 text-center ...">LegendasBOT Dashboard</div>`.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Ajustes pontuais de UI e z-index em `SideMenu.tsx`.

## Impactos esperados
- Menu lateral renderizado por cima de toda a interface (inclusive barra de navegação/footer), com nomes diretos "Templates" e "Agendamentos", sem borda lateral externa e sem rodapé.

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
