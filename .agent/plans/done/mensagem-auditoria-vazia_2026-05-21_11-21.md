# Plano: mensagem-auditoria-vazia

## Pedido do usuário
Melhorar a Dashboard Admin para mostrar uma mensagem bonita quando a auditoria do XavolaBot nao encontra nenhum canal.

## Objetivo
Criar um estado vazio mais claro e visualmente consistente na aba de Auditoria, alem de manter um feedback amigavel em caso de erro real na varredura.

## Contexto atual
A aba `AdminAuditTab` executa `fetchAuditCheckBot()` e armazena o retorno em `results`. Quando `results.length === 0`, hoje aparece apenas um card simples com texto curto e icone. O usuario quer uma mensagem de erro/estado vazio mais bonita quando nenhum canal com `@XavolaBot` for encontrado.

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/api.ts`
- `dashboard/src/types.ts`
- `dashboard/src/index.css`

## Arquivos que poderao ser modificados
- `dashboard/src/components/AdminAuditTab.tsx`

## Estrategia de implementacao
Manter a mudanca concentrada no componente da aba de auditoria, usando os estilos ja existentes (`card`, `section-icon`, variaveis CSS e lucide icons). Evitar alterar API/backend porque o backend ja retorna sucesso com lista vazia; o problema e de apresentacao no frontend.

## Passos detalhados

1. Ajustar os imports de icones em `AdminAuditTab.tsx` para usar icones apropriados ao estado vazio.
2. Normalizar o retorno de `fetchAuditCheckBot()` para tratar `res.data` ausente como lista vazia.
3. Melhorar o toast de lista vazia com texto mais natural.
4. Substituir o card simples de `results.length === 0` por um estado vazio mais rico, com icone, titulo, texto explicativo e pequenos indicadores do que foi verificado.
5. Garantir que o estado vazio so apareca depois da varredura terminar.
6. Rodar build do dashboard com `npm run build`.

## Riscos
- Pequeno risco de destoar do visual atual se o estado vazio usar estilo muito diferente.
- Pequeno risco de `res.data` vir em formato inesperado; a normalizacao reduz esse risco.

## Impactos esperados
- Quando nenhum canal com `@XavolaBot` for encontrado, o admin vera uma confirmacao clara e bonita.
- Erros reais da API continuam aparecendo via toast de erro.
- Sem impacto esperado no backend, banco, bot ou rotas.

## Compatibilidade
- Linux: compativel
- macOS: compativel
- Windows: compativel
- Docker: sem impacto
- CI/CD: apenas build frontend deve ser afetado

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
```bash
go test ./...
```

Observacao: na varredura anterior, `go test ./...` falhou neste ambiente por instalacao Go inconsistente (`go: no such tool "vet"`), nao por erro do codigo.

### Execucao
```bash
make dev
```

Abrir a Dashboard Admin, ir na aba Auditoria e executar a varredura em um cenario sem canais com `@XavolaBot`.

## Rollback
Reverter as alteracoes em `dashboard/src/components/AdminAuditTab.tsx`.

## Observacoes
- Nao sera feita alteracao destrutiva.
- Nao sera alterado o contrato da API.
- O arquivo `dashboard/src/index.css` ja possui variaveis e classes suficientes, entao a primeira opcao e nao adicionar CSS global.
