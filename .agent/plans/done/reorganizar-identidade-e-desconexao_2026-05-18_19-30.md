# Plano: Simplificação e Reorganização do Card de Identidade

## Pedido do usuário
1. Manter a cor azul (Sentri Soft).
2. Remover o bloco de link (Telemetry HUD / Access Link).
3. Colocar o botão "Desconectar Bot" dentro do card de identificação, abaixo das informações do usuário.
4. Reajustar a mensagem de saudação que está "fora da div" (integrar ou alinhar melhor).

## Objetivo técnico
- Integrar a saudação e o botão de desconectar dentro do card principal de identidade para um visual mais compacto e direto.
- Eliminar o HUD de telemetria e o link de acesso.
- Ajustar o CSS para que essa nova estrutura interna fique harmoniosa.

## Arquivos analisados
- `dashboard/src/components/DashboardInicioTab.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

**1. Componente React (`DashboardInicioTab.tsx`):**
- Remover o `welcome-card` separado.
- Mover a lógica de `getGreeting` e `getGreetingEmoji` para dentro do card principal.
- Substituir o `telemetry-hud` pelo botão de desconexão.
- Remover o botão de desconexão antigo que ficava no final da página.

**2. Estilos CSS (`index.css`):**
- Criar uma classe `.card-disconnect-btn` para o botão quando estiver dentro do card (mais discreto, sem bordas externas pesadas).
- Ajustar o espaçamento interno do card para acomodar a saudação no topo.

## Passos detalhados

1. **Modificar `DashboardInicioTab.tsx`**:
   - Unificar a saudação e as informações do usuário no primeiro `.card`.
   - Inserir o botão "Desconectar Bot" logo abaixo do ID/Nome.
   - Limpar o código morto (HUD e botão antigo).

2. **Modificar `index.css`**:
   - Ajustar estilos para que o botão de desconectar dentro do card pareça um item de ação integrado, não um botão gigante isolado.

## Riscos
Nenhum. É uma reorganização de layout solicitada para simplificação.

## Como testar
- Verificar se a saudação aparece no topo do card de identidade.
- Verificar se o botão de desconectar está logo abaixo do nome do usuário.
- Confirmar que não há mais nenhum link de acesso ou HUD visível.
