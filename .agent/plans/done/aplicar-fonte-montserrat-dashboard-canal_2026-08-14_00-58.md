# Plano: Aplicar Fonte Montserrat na Dashboard do Canal (/dashboard/:channelID)

## Pedido do usuário
"vamos utilizar a fonte na dashboard de configuracao de canal /dashboard/:channelID https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap"

## Objetivo
Configurar a fonte **Montserrat** (`'Montserrat', sans-serif`) como a fonte primária de todo o dashboard e da tela de configuração do canal (`/dashboard/:channelID`), abrangendo todos os textos, títulos, botões, formulários e elementos visuais.

## Contexto atual
- O link de importação do Google Fonts para a Montserrat já está presente no `index.html` e no topo de `index.css`.
- Em `index.css`, a regra `html, body, #root` estava usando `'Geist Variable'` e os títulos `h1-h6` usavam `'Plus Jakarta Sans'`.

## Arquivos analisados
- `dashboard/src/index.css`
- `dashboard/index.html`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`

## Estratégia de implementação

1. **`dashboard/src/index.css`**:
   - Atualizar a regra `html, body, #root`:
     - `font-family: 'Montserrat', sans-serif;`
   - Atualizar a regra `h1, h2, h3, h4, h5, h6`:
     - `font-family: 'Montserrat', sans-serif;`
   - Garantir que todos os componentes e textos da tela de configuração do canal herdem Montserrat limpa.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Riscos
- **Nenhum**: Alteração tipográfica puramente estética conforme solicitado pelo usuário.

## Impactos esperados
- Toda a dashboard de configuração do canal (`/dashboard/:channelID`) renderizará uniformemente com a tipografia Montserrat, proporcionando visual moderno, legível e coeso com o menu lateral.

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
