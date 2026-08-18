# Plano: Reverter Alterações da Dashboard (Rollback)

## Pedido do usuário
Reverter todas as alterações visuais feitas na pasta `dashboard/` e restaurar o estado original dos arquivos da dashboard.

## Objetivo
Restaurar completamente a dashboard (`dashboard/src/`) para o estado anterior no Git (`git checkout dashboard/`), mantendo as correções de backend Go no `postBuilder.go`.

## Contexto atual
- A dashboard passou por uma refatoração Anti-AI-Slop.
- O usuário pediu para reverter toda a dashboard ao estado original.

## Arquivos analisados
- Todos os arquivos da pasta `dashboard/` listados no `git status`.
- `internal/telegram/handlers/events/postBuilder/postBuilder.go` (SERÁ PRESERVADO).

## Arquivos que poderão ser modificados (revertidos)
- `dashboard/` (todos os arquivos modificados na pasta `dashboard/`)

## Estratégia de implementação
1. Executar `git checkout dashboard/` para restaurar todos os arquivos da dashboard para o commit mais recente (`origin/feat/mtproto-custom-emoji-separator`).
2. Executar `npm run build` na pasta `dashboard/` para re-compilar a versão original do bundle e confirmar a integridade.

## Passos detalhados
1. Rodar `git checkout dashboard/`.
2. Rodar `cd dashboard && npm run build`.
3. Verificar a restauração completa.

## Riscos
- NENHUM. Restauração direta do repositório Git.

## Impactos esperados
- A dashboard retornará exatamente ao design e visual anteriores.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
git status
cd dashboard && npm run build
```

## Rollback
<não aplicável, esta ação JÁ É o rollback>
