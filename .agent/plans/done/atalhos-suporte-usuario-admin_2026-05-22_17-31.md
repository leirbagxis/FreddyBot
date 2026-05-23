# Plano: atalhos suporte usuario admin

## Pedido do usuário
Adicionar na dashboard admin um botão no detalhe de um usuário selecionado para enviar mensagem de suporte a ele. Também fazer com que, nos resultados do CheckBot/Auditoria, clicar no usuário leve para a aba/detalhe desse usuário.

## Objetivo
Reaproveitar o fluxo existente de disparo manual individual da dashboard admin, que já usa o título de suporte, e criar atalhos de navegação entre auditoria, detalhe do usuário e aba de disparo.

## Contexto atual
- A aba de disparo manual já possui o público-alvo `Individual`.
- Quando o alvo é `single`, o preview mostra `# 📨 MENSAGEM DO SUPORTE`.
- O backend também adiciona esse cabeçalho ao enviar mensagem individual.
- A aba de usuários permite selecionar um usuário e abrir o detalhe.
- A aba de auditoria mostra usuários encontrados no CheckBot, mas atualmente só os canais são clicáveis.
- O estado da aba admin fica em `App.tsx` e é passado para `AdminDashboard.tsx`.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/components/AdminNoticeTab.tsx`
- `.agent/plans/`

## Arquivos que poderão ser modificados
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`

## Estratégia de implementação
Adicionar callbacks controlados pelo `App.tsx` para evitar dessincronização entre a aba ativa do pai e a aba local do `AdminDashboard`.

Para mensagem de suporte:
- Criar um atalho que define o alvo do disparo como `single`.
- Preencher `noticeTargetId` com o ID do usuário selecionado.
- Navegar para a aba `notice`.
- Manter o texto/imagem/botões atuais do formulário para evitar apagar rascunhos do admin sem aviso.

Para auditoria:
- Passar um callback para `AdminAuditTab`.
- Transformar a área de identificação do usuário do resultado em um botão clicável.
- Ao clicar, selecionar o usuário e navegar para a aba `users`, abrindo o detalhe dele.

## Passos detalhados

1. Em `App.tsx`, criar callbacks para:
   - abrir detalhe de usuário na aba `users`;
   - abrir a aba `notice` com alvo `single` e ID do usuário preenchido.
2. Passar esses callbacks para `AdminDashboard`.
3. Em `AdminDashboard.tsx`, adicionar as novas props.
4. No detalhe do usuário, adicionar um botão "Mensagem de Suporte" com ícone apropriado.
5. Em `AdminDashboard.tsx`, repassar o callback de abrir usuário para `AdminAuditTab`.
6. Em `AdminAuditTab.tsx`, adicionar prop de navegação para usuário.
7. Fazer a área do usuário no resultado da auditoria ser clicável sem interferir no botão "Limpar Tudo".
8. Rodar build do dashboard.
9. Rodar `git diff --check`.

## Riscos
- Se o usuário do resultado da auditoria não estiver presente na lista carregada em `adminData.users`, a navegação pode ir para a aba de usuários sem encontrar o detalhe.
- A aba de disparo preservará rascunhos existentes; isso evita perda de texto, mas pode exigir que o admin revise o conteúdo antes de enviar.
- Mudanças são restritas ao frontend, sem alteração de contrato da API.

## Impactos esperados
- Admin consegue iniciar uma mensagem individual de suporte direto do detalhe do usuário.
- Admin consegue abrir rapidamente o detalhe do usuário a partir dos resultados do CheckBot.
- O título de suporte permanece centralizado no fluxo já existente de disparo individual.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
cd dashboard
npm run build
```

### Testes
```bash
git diff --check
```

### Execução
```bash
cd dashboard
npm run dev
```

## Rollback
Reverter as alterações nos arquivos:
- `dashboard/src/App.tsx`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/AdminAuditTab.tsx`

## Observações
- Não será necessário alterar o backend, porque o cabeçalho `MENSAGEM DO SUPORTE` já é aplicado no envio individual.
- O botão novo deve apenas preparar o formulário de disparo, não enviar automaticamente.
