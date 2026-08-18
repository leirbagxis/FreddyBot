# Plano: Refatorar Header da Dashboard de Canais (Transparente + Sem Título/Usuário/Avatar)

## Pedido do usuário
Na dashboard de configuração de canais (`/dashboard/:channelID` e listagem):
- Remover os títulos (nome do usuário e "Visão Geral"/"Meus Canais").
- Remover o ícone com a foto do perfil do usuário (`top-avatar`).
- Manter apenas o ícone de trocar o tema (alinhado à direita) e o botão de voltar (quando em um canal específico).
- Deixar a top-bar transparente com efeito blur (`backdrop-filter: blur(12px)`), sem borda inferior nem sombra, como se ela pertencesse ao body.

## Objetivo
Simplificar o header em `App.tsx` e `index.css` removendo avatar, nome do usuário e subtítulos, deixando apenas o botão de ação de tema e o botão de voltar em um layout limpo, transparente e integrado ao body.

## Contexto atual
- A `.top-bar` em `App.tsx` exibe o botão voltar, avatar (`top-avatar`), nome (`{displayName}`), rótulo da rota ("Visão Geral", "Meus Canais") e botão de tema (`theme-switch`).
- No CSS (`index.css`), `.top-bar` possui `background: var(--nav-bg)`, `border-bottom: 1px solid var(--border)` e `box-shadow`.

## Arquivos analisados
- `dashboard/src/App.tsx` (estrutura HTML do header)
- `dashboard/src/index.css` (estilos de `.top-bar`)

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

1. **`dashboard/src/App.tsx`**:
   - No bloco `<div className="top-bar animate-stagger-in">`:
     - Remover `<div className="top-avatar">...</div>` (foto de perfil).
     - Remover `<div className="min-w-0 flex-1">...</div>` (nome do usuário + "Visão Geral"/"Meus Canais").
     - Manter o botão de voltar (`isSpecificChannel && <button ... />`) se estiver na página de canal `/dashboard/:channelID`.
     - Posicionar o botão de trocar tema (`theme-switch`) com `ml-auto` alinhado à direita.

2. **`dashboard/src/index.css`**:
   - Em `.top-bar`:
     - Alterar `background: transparent`.
     - Manter `-webkit-backdrop-filter: blur(12px)` e `backdrop-filter: blur(12px)`.
     - Remover `border-bottom` (ou `border-bottom: none`) e `box-shadow: none` para mesclar o header perfeitamente com o body.
     - Ajustar `padding: 12px 16px` para um visual mais clean.

3. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Modificar `<div className="top-bar">` em `dashboard/src/App.tsx`.
2. Modificar regra de estilo `.top-bar` em `dashboard/src/index.css`.
3. Executar `npm run build` em `dashboard`.

## Riscos
- **Baixo**: Alteração puramente visual na UI.

## Impactos esperados
- Header completamente limpo, transparente, integrado ao body, sem avatar/foto de perfil ou títulos, mantendo apenas as ações essenciais.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/App.tsx
git checkout dashboard/src/index.css
```
