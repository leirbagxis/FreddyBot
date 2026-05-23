# Plano: corrigir exibicao caption admin

## Pedido do usuário
Corrigir a dashboard admin porque a default caption configurada não aparece lá depois de setada.

## Objetivo
Garantir que os campos de legendas globais da dashboard admin reflitam o estado salvo/retornado pela API, especialmente após clicar em salvar.

## Contexto atual
A tela `AdminConfigTab` carrega `globalDefaultCaption` e `globalNewPackCaption` no `useEffect` inicial. No `handleSave`, após a API retornar sucesso, o componente atualiza `config` e os campos do PostBuilder fixo, mas não atualiza os estados locais `globalDefault` e `globalNewPack`. Isso pode deixar os editores exibindo valor antigo/vazio mesmo quando a API retornou a configuração atualizada.

O backend já expõe `globalDefaultCaption` e `globalNewPackCaption` no modelo `ServerConfig`, recebe esses campos no controller admin e usa a legenda global ao criar novos canais.

## Arquivos analisados
- dashboard/src/components/AdminConfigTab.tsx
- dashboard/src/api.ts
- dashboard/src/types.ts
- dashboard/src/components/RichTextEditor.tsx
- internal/api/controllers/adminController/configController.go
- internal/database/models/models.go
- internal/telegram/handlers/events/addChannel/addChannel.go
- internal/core/services/channels.go

## Arquivos que poderão ser modificados
- dashboard/src/components/AdminConfigTab.tsx

## Estratégia de implementação
Atualizar o `handleSave` para, quando receber `serverData`, sincronizar também `globalDefault` e `globalNewPack` com os valores retornados pela API. Isso mantém o editor controlado coerente após salvar, sem alterar contrato da API nem comportamento de criação de canais.

## Passos detalhados

1. Alterar `dashboard/src/components/AdminConfigTab.tsx` no bloco de sucesso do `handleSave`.
2. Após `setConfig(serverData)`, chamar `setGlobalDefault(serverData.globalDefaultCaption || '')`.
3. Também chamar `setGlobalNewPack(serverData.globalNewPackCaption || '')`.
4. Manter a sincronização já existente dos campos do PostBuilder fixo.
5. Rodar build da dashboard.
6. Rodar checagem de diff.

## Riscos
- Baixo risco: alteração restrita ao frontend admin.
- Se a API retornar texto transformado ou sanitizado, o editor passará a exibir esse texto retornado, que é o comportamento esperado depois de salvar.

## Impactos esperados
- A default caption salva passa a aparecer imediatamente no campo da dashboard admin.
- A legenda de novo pack também fica coerente após salvar.
- Nenhuma mudança na API ou no uso da legenda ao adicionar canais.

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

## Rollback
Reverter a alteração em `dashboard/src/components/AdminConfigTab.tsx` removendo a sincronização de `globalDefault` e `globalNewPack` dentro do sucesso do `handleSave`.

## Observações
A legenda global é usada como valor inicial para novos canais vinculados; ela não retroage automaticamente para canais que já existem.
