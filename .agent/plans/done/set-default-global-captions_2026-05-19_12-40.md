# Plano: set-default-global-captions

## Pedido do usuário
O usuário solicitou que, ao iniciar o bot pela primeira vez, os campos `GlobalDefaultCaption` e `GlobalNewPackCaption` sejam preenchidos com valores padrão específicos. Esses valores representam templates para legendas de mídia e montagem de pacotes, respectivamente.

## Objetivo
Atualizar a função de inicialização do banco de dados para injetar os templates fornecidos pelo usuário na criação inicial do registro `ServerConfig`, substituindo as strings vazias atuais.

## Contexto atual
- Em `internal/database/database.go`, a função `initServerConfig` cria o primeiro registro de configuração (ID 1) usando `FirstOrCreate`.
- Atualmente, as propriedades `GlobalDefaultCaption` e `GlobalNewPackCaption` estão sendo inicializadas com `""`.

## Arquivos analisados
- `internal/database/database.go`

## Arquivos que poderão ser modificados
- `internal/database/database.go`

## Estratégia de implementação
1. **Modificar `initServerConfig`**: Alterar a struct estática passada para `FirstOrCreate` no arquivo `database.go`.
2. Substituir `GlobalDefaultCaption: ""` pelo template:
   `"🐈‍⬛ ៹ [t.me/legendasbot](https://t.me/usernamebot)  ‹"` (ajustado com os caracteres corretos do prompt).
3. Substituir `GlobalNewPackCaption: ""` pelo template formatado (preservando as quebras de linha usando crases \` no Go):
   ```
   ╔═━──━═༻✧༺═━──━═╗

           𖦹⁠⁠⁠ ࣪ ⭑ ᥫ᭡
           (｡•́︿•̀｡)っ✧.*ೃ༄
           ˗ˏˋ [$name]($link) ⋆｡˚ ☁︎
               彡♡ ₊˚

   ⋆｡˚ ❀ @LegendasBrBot ☽⁺₊

   ╚═━──━═༻✧༺═━──━═╝
   ```
4. **Nota sobre Banco Existente:** Como a função usa `FirstOrCreate`, se o registro de ID 1 já existir no banco de dados (o que é o caso da nossa instância atual), ele **não** será sobrescrito. Para testar em ambiente de desenvolvimento local, o banco `FreddyBot.db` pode precisar ser deletado, ou o usuário precisará aplicar esses valores manualmente via Dashboard uma vez. No entanto, para novas instalações, a regra funcionará como esperado.

## Passos detalhados
1. Abrir `internal/database/database.go`.
2. Localizar `initServerConfig(db *gorm.DB) error`.
3. Substituir as strings vazias nos campos `GlobalDefaultCaption` e `GlobalNewPackCaption` pelos blocos de texto (usando \` \` para strings multilinhas).

## Riscos
- Nenhum risco arquitetural. Apenas risco de formatação de string (escapes no Go), que será mitigado pelo uso de raw string literals (`).

## Impactos esperados
- Novas implantações ou inicializações com banco limpo receberão os templates automáticos para "Legenda Padrão" e "Legenda de Novo Pacote", que serão copiados para os novos canais.

## Compatibilidade
- Linux
- macOS
- Windows

## Como testar

### Build
`go build -o tmp/FreddyBot ./cmd/FreddyBot/`

### Execução
1. Apagar/renomear o banco de dados local (`FreddyBot.db`).
2. Iniciar o bot (`./tmp/FreddyBot`).
3. Verificar a Dashboard Administrativa para confirmar se os valores apareceram nos campos globais.

## Rollback
Restaurar `internal/database/database.go` para o estado anterior.