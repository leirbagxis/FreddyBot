# Plano: corrigir-persistencia-links-dinamicos_2026-05-16_16-30.md

## Pedido do usuário
O toggle de Links Dinâmicos não persiste após fechar e abrir o Mini App.

## Objetivo técnico
Garantir que os dados de Links Dinâmicos sejam incluídos na resposta da API enviada ao Dashboard, atualizando o DTO e o Mapeador.

## Contexto atual
Os campos foram adicionados ao modelo de banco de dados e a API de atualização funciona, mas os dados não são retornados no carregamento inicial do dashboard porque o DTO (`ChannelDTO`) e o mapeador (`ToChannelDTO`) foram esquecidos.

## Arquivos analisados
- `internal/api/dto/dto.go`
- `internal/api/dto/mapper.go`

## Arquivos que poderão ser modificados
- `internal/api/dto/dto.go`
- `internal/api/dto/mapper.go`

## Estratégia de implementação
1. **DTO:** Adicionar `DynamicLinks`, `DLBotButtons`, `DLBotCaptions` e `DLBotReactions` ao `ChannelDTO`.
2. **Mapper:** Atualizar `ToChannelDTO` para copiar esses valores do modelo `models.Channel` para o `dto.ChannelDTO`.

## Passos detalhados
1. Editar `internal/api/dto/dto.go` e adicionar os campos booleanos no `ChannelDTO`.
2. Editar `internal/api/dto/mapper.go` e atualizar a função `ToChannelDTO`.

## Riscos
- Nenhum risco identificado, apenas exposição de campos já existentes no banco.

## Impactos esperados
- O Dashboard passará a mostrar o estado correto (salvo) dos toggles ao carregar.

## Como testar

### Build
```bash
go build ./cmd/FreddyBot
```

### Verificação
1. Abrir o Dashboard.
2. Ativar Links Dinâmicos e sub-toggles.
3. Recarregar a página.
4. Os toggles devem permanecer ativos.

## Rollback
Reverter as alterações nos arquivos DTO e Mapper.
