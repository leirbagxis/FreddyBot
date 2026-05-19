# Plano: Revelar Link de Acesso e Implementar HUD de Telemetria

## Pedido do usuário
"tique aquele REVEAL ACCESS LINK e aquele link."
Interpretado como: remover a necessidade de "clicar para revelar" o link (sempre mostrar) e possivelmente finalizar a transição para o design HUD que estava incompleta.

## Objetivo técnico
1. Remover o efeito de "raspadinha" (`security-scratch`) que esconde o link de convite.
2. Implementar o design de **HUD de Telemetria** (estilo McLaren/Pit Wall) no card de identidade, conforme planejado anteriormente mas não totalmente integrado.

## Contexto atual
- O arquivo `DashboardInicioTab.tsx` ainda contém as classes `.security-scratch` e `.scratch-overlay`.
- O CSS correspondente a essas classes foi removido do `index.css` em uma tarefa anterior, tornando o componente visualmente inconsistente ou "quebrado".
- O link de convite exige uma interação desnecessária para ser visualizado.

## Arquivos analisados
- `dashboard/src/components/DashboardInicioTab.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

**1. Componente React (`DashboardInicioTab.tsx`):**
- Substituir a estrutura do `security-scratch` por um painel de telemetria fixo.
- O link será exibido diretamente em um bloco de "SECURE UPLINK".
- Adicionar um indicador visual de status ("SYSTEM ACTIVE").

**2. Estilos CSS (`index.css`):**
- Adicionar as classes para `.telemetry-hud`, `.telemetry-header`, `.telemetry-body` e `.telemetry-footer`.
- Incluir a animação de pulso para o LED de status.
- Manter a consistência com o tema "Sentri Soft" (Violeta/Indigo) recém aplicado, usando o roxo como cor do LED e bordas.

## Passos detalhados

1. **Modificar `DashboardInicioTab.tsx`:**
   - Remover a div `.security-scratch` e seu `onClick`.
   - Inserir a nova estrutura `.telemetry-hud`.
   - Adicionar o ícone `Cpu` para reforçar a estética técnica.

2. **Modificar `index.css`:**
   - Adicionar estilos para o HUD (borda lateral colorida, fontes mono, espaçamento técnico).
   - Criar a animação `@keyframes pulse-dot`.

## Riscos
- O link pode ser longo; usarei `truncate` e `overflow: hidden` para garantir que não quebre o card em telas menores.

## Como testar
- Abrir a dashboard na aba Início.
- O link de convite deve estar visível imediatamente, sem a necessidade de clicar em "REVEAL".
- O visual deve parecer um painel técnico de monitoramento.

## Rollback
- Reverter as alterações nos dois arquivos mencionados.
