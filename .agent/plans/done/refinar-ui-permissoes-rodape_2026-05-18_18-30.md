# Plano: Refinamento de UI e Permissões

## Pedido do usuário
1. Adicionar "Link Preview" às permissões de mensagens.
2. Simplificar o texto dos "Links Dinâmicos" (está muito grande).
3. Ajustar o rodapé (bottom-nav): mover um pouco mais para cima e deixá-lo maior.

## Objetivo técnico
1. Ajustar o componente `PermissionsCard` para garantir que `linkPreview` seja exibido.
2. Editar o texto de descrição em `App.tsx`.
3. Ajustar o CSS da `.bottom-nav` em `index.css`.

## Estratégia de implementação

**1. Permissões de Mensagem:**
- No componente `PermissionsCard.tsx`, alterarei a lógica de filtragem para que `linkPreview` apareça mesmo que não esteja presente no objeto de permissão inicial (garantindo que o usuário possa ativá-lo).
- Ajustar `App.tsx` para garantir que o estado inicial de `linkPreview` seja considerado.

**2. Simplificação de Texto:**
- Localizar a linha 988 de `App.tsx`.
- Alterar "Transforma !Nome e !URL ou links embutidos em botões automaticamente" para algo como "Transforma links em botões automaticamente".

**3. Ajuste do Rodapé (TabBar):**
- No arquivo `index.css`, classe `.bottom-nav`:
  - Alterar `bottom` de `16px` para `32px` (mais para cima).
  - Alterar `max-width` de `320px` para `400px` (maior).
  - Aumentar `padding` de `4px` para `8px`.
  - Aumentar `font-size` dos itens para `11px` ou `12px` se necessário.

## Passos detalhados
1.  **Modificar `dashboard/src/components/PermissionsCard.tsx`**:
    - Alterar a linha `const available = fields.filter(f => f.key in permission);` para permitir `linkPreview` explicitamente se o título for relacionado a Mensagens.
2.  **Modificar `dashboard/src/App.tsx`**:
    - Simplificar o texto da seção de Links Dinâmicos.
3.  **Modificar `dashboard/src/index.css`**:
    - Ajustar `.bottom-nav` para ser maior e ficar mais alto.

## Riscos
- Mover o rodapé muito para cima pode sobrepor conteúdo importante. Usarei um valor moderado (`32px`).
- Aumentar o rodapé pode fazer com que ele ocupe muito espaço em telas muito pequenas.

## Como testar
- Verificar se "Link Preview" agora aparece em Permissões de Mensagem.
- Verificar o novo texto dos Links Dinâmicos.
- Visualizar o novo posicionamento e tamanho do rodapé.
