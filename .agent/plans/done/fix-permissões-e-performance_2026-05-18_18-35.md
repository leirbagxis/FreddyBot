# Plano: Correção de Permissões e Otimização de Performance

## Pedido do usuário
1. O "Link Preview" está aparecendo nas "Permissões de Botões" e não deveria. Deve ficar apenas nas "Permissões de Mensagem".
2. A dashboard está "travando demais" (problemas de performance/lag).

## Objetivo técnico
1. Ajustar a lógica no `PermissionsCard.tsx` para forçar a remoção de `linkPreview` se não for o card de Mensagem.
2. Otimizar o CSS (`index.css`) reduzindo a carga de renderização (CSS repaints e GPU memory) causada pelo excesso de `backdrop-filter` (glassmorphism) e animações pesadas.

## Estratégia de implementação

**1. Correção de Permissões (`PermissionsCard.tsx`):**
- A lógica atual permite `f.key in permission`. Se a API retornar `linkPreview` dentro de `buttonsPermission` por acidente (ou por herança do BD), ele será renderizado.
- Vou forçar a filtragem:
  ```tsx
  const isMessagePerm = title.toLowerCase().includes('mensagem');
  const available = fields.filter(f => {
    if (f.key === 'linkPreview' && !isMessagePerm) return false; // Bloqueia em Botões
    return (f.key in permission) || (isMessagePerm && f.key === 'linkPreview');
  });
  ```

**2. Otimização de Performance (`index.css`):**
A lentidão geralmente ocorre devido a muitos elementos com `backdrop-filter: blur()` e animações rodando simultaneamente, forçando a GPU no mobile.
- **Reduzir o desfoque (Blur):** Diminuir os valores de `blur(16px)` para `blur(8px)` ou `10px` nos cards e barras de navegação.
- **Remover animações infinitas desnecessárias:** O efeito `card-float` no access-card antigo pode ainda estar sendo chamado (vou limpar se sobrou algo), e vou garantir que não há `will-change` excessivo.
- **Otimizar os Cards:**
  - Reduzir `backdrop-filter: blur(12px)` para `blur(6px)`.
  - Melhorar as transições (remover `all 0.2s` indiscriminado e especificar apenas o que muda, como `background-color`, `transform`).

## Passos detalhados
1.  **Atualizar `PermissionsCard.tsx`:** Inserir a condição de bloqueio para `linkPreview` em cards que não sejam de mensagens.
2.  **Atualizar `index.css`:**
    - Ajustar `.card`, `.top-bar`, `.bottom-nav` reduzindo `blur` e `saturate`.
    - Revisar a classe `.animate-stagger-in` e outras animações de entrada para garantir que estão otimizadas (usando `transform` e `opacity` sem acionar layout recalculation).
    - Verificar e remover classes não utilizadas que possam estar pesando.

## Riscos
A redução do blur fará com que o efeito "vidro" fique um pouco menos "fosco", mostrando mais nitidamente o que está atrás, mas isso é uma troca comum para garantir fluidez (60fps) em dispositivos móveis. A estética geral será preservada.