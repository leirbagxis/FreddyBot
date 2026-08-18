# Plano: Corrigir Fundo Branco Persistente na Dashboard Admin

## Pedido do usuário
"mano o tema da dashboard admin ainda continua branco mano"

## Causa Raiz Identificada
No arquivo `dashboard/src/index.css` (linhas 3531-3585), a classe `.admin-layout-v2` redefinia localmente todas as variáveis de tema (`--bg: #f7f7f5;`, `--card: #ffffff;`, `--background: #f7f7f5;`, `--nav-bg: #ffffff;`).
Como o hook `useTheme.ts` define o atributo `data-theme="telegram"` na raiz da aplicação, as regras do `.admin-layout-v2` ignoravam o tema do Telegram e caíam nos valores padrão da classe que forçavam fundo branco e cinza claro.

## Objetivo
Remover as redefinições de variáveis legadas brancas de `.admin-layout-v2` em `index.css`, permitindo que a Dashboard Admin herde diretamente as variáveis globais do tema do Telegram (`var(--bg)` escuro e `var(--card)` translúcido `rgba(255, 255, 255, 0.05)`).

## Arquivos analisados
- `dashboard/src/index.css` (linhas 3531-3630)
- `dashboard/src/hooks/useTheme.ts`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Remover/limpar os blocos de variáveis de `.admin-layout-v2` nas linhas 3531-3630 que sobressaíam e forçavam fundo branco `#f7f7f5` e `#ffffff`.
   - Garantir que `.admin-layout-v2` respeite as variáveis globais de `:root` e `[data-theme="telegram"]` (`--bg`, `--card`, `--text`, `--hint`, `--accent`, `--border`).

2. **Validação**:
   - Compilar o frontend com `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Eliminação de override incorreto de CSS.

## Impactos esperados
- A Dashboard Admin exibirá imediatamente o fundo escuro/translúcido do tema do Telegram, eliminando qualquer tela ou card branco indesejado.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/index.css
```
