# Plano: remover-threejs-gsap-dashboard-v2_2026-05-18_15-00.md

## Pedido do usuário
Remover `three.js` e `gsap` da dashboard. Não quer nada de Canvas ou 2D/3D (BackgroundEffect), mas quer manter as animações de UI.

## Objetivo
Remover as bibliotecas pesadas e o efeito de fundo em WebGL, substituindo as animações de interface por CSS puro para otimização máxima.

## Contexto atual
- `three.js` renderiza um fundo de partículas em `BackgroundEffect.tsx`.
- `gsap` faz o efeito de "cascata" (stagger) na entrada de componentes.
- O CSS já possui algumas animações de entrada baseadas em `@keyframes`.

## Arquivos analisados
- `dashboard/package.json`
- `dashboard/src/App.tsx`
- `dashboard/src/components/BackgroundEffect.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/components/PermissionsCard.tsx`
- `dashboard/src/index.css`

## Arquivos que poderão ser modificados
- `dashboard/package.json`
- `dashboard/src/App.tsx`
- `dashboard/src/components/ButtonGrid.tsx`
- `dashboard/src/components/PermissionsCard.tsx`
- `dashboard/src/index.css`

## Arquivos que serão excluídos
- `dashboard/src/components/BackgroundEffect.tsx`

## Estratégia de implementação
1.  **Limpeza de Dependências:** Remover `three`, `@types/three` e `gsap`.
2.  **Remoção do Background:** Excluir o componente `BackgroundEffect` e sua referência no `App.tsx`.
3.  **Animações CSS:** 
    - No `index.css`, garantir que existam animações de entrada (`fade-in-up`, `fade-in-right`).
    - Nos componentes (`ButtonGrid`, `PermissionsCard`), substituir o `gsap.fromTo` por classes CSS e aplicar `animation-delay` via style inline nos itens filhos para manter o efeito de stagger.
4.  **Refatoração App.tsx:** Remover a lógica de animação de troca de abas feita com GSAP e usar transições CSS simples.

## Passos detalhados

1.  **Remover Dependências:**
    - Editar `dashboard/package.json` removendo as libs citadas.
    - `npm install` na pasta dashboard.

2.  **Remover BackgroundEffect:**
    - Deletar `dashboard/src/components/BackgroundEffect.tsx`.
    - No `App.tsx`, remover o import e a tag `<BackgroundEffect />`.

3.  **Implementar Stagger em CSS (PermissionsCard):**
    - Remover `useEffect` com GSAP.
    - Adicionar classe de animação aos itens `.perm-row`.
    - Aplicar `style={{ animationDelay: `${index * 0.05}s` }}`.

4.  **Implementar Stagger em CSS (ButtonGrid):**
    - Remover `useEffect` com GSAP.
    - Aplicar delay semelhante aos elementos do grid.

5.  **Ajustar Transição de Abas (App.tsx):**
    - Substituir o `gsap.fromTo('.tab-content-wrapper', ...)` por uma animação CSS disparada pela mudança de estado (ou simplesmente confiar no `entrance-slide` já existente no CSS).

## Riscos
- Sem o `clearProps: 'all'` do GSAP, é necessário garantir que as classes CSS não interfiram em estados futuros (ex: hover). Vou usar `animation-fill-mode: both`.

## Impactos esperados
- Dashboard muito mais leve e rápida.
- Remoção total de WebGL e Canvas do projeto.
- Manutenção da experiência visual fluida.

## Como testar
- Verificar se as abas e listas ainda entram com suavidade.
- Confirmar que não há erros de "module not found" no console.

## Rollback
- Reverter mudanças via git e reinstalar libs.
