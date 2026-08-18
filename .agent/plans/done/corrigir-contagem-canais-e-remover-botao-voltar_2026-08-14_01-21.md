# Plano: Corrigir Contagem de Canais do Usuário e Remover Botão de Voltar no Header

## Pedido do usuário
"Quando estou na dashboard de configuracao de canal /dashboard/:channelID, ele nao mostra a quantidade de canais corretas do usuario. Tbm remova aquele botao de voltar no header"

## Objetivo
1. Corrigir o carregamento da lista de canais do proprietário na rota de configuração de canal (`/dashboard/:channelID`), garantindo que o contador no menu lateral (`SideMenu`) exiba a quantidade exata de canais do usuário.
2. Remover o botão de voltar (`ArrowLeft`) presente na barra superior (`.top-bar`) do header da dashboard.

## Contexto atual
- Na rota `/dashboard/:channelID`, a API backend `/api/channel/:channelId` faz a busca do canal pelo repositório `GetChannelByID`, mas precarregava apenas `Preload("Owner")` sem precarregar `Preload("Owner.Channels")`. Por isso, `user.channels` retornava vazio, fazendo o SideMenu exibir 0 canais.
- No header (`App.tsx`), a `.top-bar` renderizava um botão com o ícone `ArrowLeft` acionando `handleBack` quando `isSpecificChannel` era `true`.

## Arquivos analisados
- `internal/database/repositories/channel.go`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- `internal/database/repositories/channel.go`
- `dashboard/src/App.tsx`

## Estratégia de implementação

1. **Backend (`internal/database/repositories/channel.go`)**:
   - Adicionar `.Preload("Owner.Channels")` na query `GetChannelByID` para que o objeto `channel.Owner` inclua a lista completa dos canais cadastrados pelo usuário.

2. **Frontend (`dashboard/src/App.tsx`)**:
   - Remover a renderização condicional do botão de voltar (`ArrowLeft`) na `.top-bar` do header.

3. **Validação**:
   - Executar `go build ./cmd/FreddyBot` para validar a compilação do Go.
   - Executar `npm run build` na pasta `dashboard` para validar a compilação do frontend.

## Riscos
- **Nenhum**: Ajuste de pré-carregamento no repositório de canais e remoção pontual do botão de voltar no header conforme solicitado.

## Impactos esperados
- O contador de canais no menu lateral exibirá o número correto de canais do usuário em qualquer tela, inclusive em `/dashboard/:channelID`.
- A barra superior do header ficará mais limpa, sem o botão de voltar.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
go build ./cmd/FreddyBot
cd dashboard && npm run build
```

## Rollback
```bash
git checkout internal/database/repositories/channel.go dashboard/src/App.tsx
```
