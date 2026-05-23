# Plano: modo tema telegram dashboard

## Pedido do usuário
Adicionar ao toggle de modo claro/escuro um terceiro modo `Telegram`, que usa as informações de tema do Telegram WebApp para aplicar cores na dashboard.

## Objetivo
Permitir que o usuário alterne entre três modos de tema na dashboard: claro, escuro e Telegram. No modo Telegram, a dashboard deve usar `window.Telegram.WebApp.themeParams` e `colorScheme`, com fallback seguro fora do Telegram.

## Contexto atual
O hook `dashboard/src/hooks/useTheme.ts` possui apenas `Theme = 'light' | 'dark'`. Ele usa o tema do Telegram apenas como fallback inicial quando não existe preferência manual. O botão em `App.tsx` chama `toggleTheme` e alterna entre ícones de sol/lua. O CSS define variáveis para `[data-theme="light"]` e `[data-theme="dark"]`, mas não tem `[data-theme="telegram"]`.

## Arquivos analisados
- dashboard/src/hooks/useTheme.ts
- dashboard/src/App.tsx
- dashboard/src/index.css
- dashboard/src/types.ts

## Arquivos que poderão ser modificados
- dashboard/src/hooks/useTheme.ts
- dashboard/src/App.tsx
- dashboard/src/index.css

## Estratégia de implementação
Alterar o hook de tema para trabalhar com `Theme = 'light' | 'dark' | 'telegram'`. O botão passa a ciclar `light -> dark -> telegram -> light`. Quando o modo for `telegram`, o hook aplica `data-theme="telegram"` e injeta variáveis CSS baseadas em `Telegram.WebApp.themeParams`.

Como fora do Telegram `themeParams` pode não existir, o modo Telegram deve cair em cores derivadas de `colorScheme` ou nos defaults atuais. Também será necessário atualizar o header/background do Telegram usando as cores do próprio WebApp quando disponíveis.

## Passos detalhados

1. Alterar o tipo `Theme` em `useTheme.ts` para incluir `telegram`.
2. Ajustar leitura de `localStorage` para aceitar apenas `light`, `dark` ou `telegram`.
3. Criar helpers para validar tema, escolher esquema efetivo e aplicar/remover variáveis CSS do modo Telegram.
4. Quando `theme === 'telegram'`, aplicar variáveis como `--bg`, `--text`, `--hint`, `--link`, `--accent`, `--card`, `--nav-bg`, `--surface`, `--input-bg`, `--border` usando `themeParams`.
5. Quando sair do modo Telegram, remover overrides inline de variáveis para voltar ao CSS normal de light/dark.
6. Atualizar `toggleTheme` para ciclar entre os três modos e marcar preferência manual.
7. Atualizar o retorno do hook se necessário com `theme` e `toggleTheme`.
8. Atualizar `App.tsx` para exibir um ícone/label adequado no botão de tema: sol, lua ou ícone representando Telegram/sistema.
9. Adicionar CSS para `[data-theme="telegram"]` como fallback, reaproveitando valores próximos ao light/dark, e ajustar `.theme-switch` se precisar comportar label curto.
10. Rodar `npm run build`.
11. Rodar `git diff --check`.

## Riscos
- Baixo a médio: mudança em tema global pode afetar contraste se o Telegram enviar cores incompletas ou muito próximas.
- Fora do Telegram, o modo `telegram` precisará de fallback visual estável.
- O hook atual possui lógica automática por horário; ela deve continuar apenas quando não há preferência manual.

## Impactos esperados
- O toggle passa a ter três estados: claro, escuro e Telegram.
- Dentro do Telegram WebApp, o modo Telegram acompanha as cores do cliente do usuário.
- Fora do Telegram, o modo Telegram funciona com fallback sem quebrar a dashboard.

## Compatibilidade
- Linux
- macOS
- Windows
- Docker
- CI/CD

## Como testar

### Build
```bash
npm run build
```

### Testes
```bash
go test ./...
```

### Execução
```bash
npm run dev
```

Teste manual:
1. Abrir dashboard e clicar no botão de tema.
2. Confirmar ciclo claro -> escuro -> Telegram -> claro.
3. Confirmar persistência após recarregar.
4. Dentro do Telegram WebApp, confirmar que o modo Telegram usa `themeParams`.
5. Fora do Telegram, confirmar que o modo Telegram mantém fallback visual legível.

## Rollback
Reverter alterações em `useTheme.ts`, `App.tsx` e `index.css` para voltar ao toggle binário claro/escuro.

## Observações
Essa alteração é apenas frontend. Não exige alteração na API ou no bot.
