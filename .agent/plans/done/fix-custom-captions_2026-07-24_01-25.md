# Plano: Correção e Ativação das Custom Captions (Legendas Dinâmicas)

## Pedido do usuário
Ativar e finalizar a feature de "Custom Captions" que já existe parcialmente no projeto. O usuário quer postar uma mídia com um comando (ex: `#promo`), o bot deve remover o comando, inserir a legenda configurada e adicionar botões específicos.

## Objetivo
- Corrigir os bugs de backend que impedem a detecção correta de comandos (atualmente só detecta a primeira hashtag).
- Criar a interface visual (Dashboard UI) para que o usuário possa criar, editar e excluir seus comandos e botões.
- Garantir o correto funcionamento da hierarquia: botões do comando devem sobrescrever botões padrão do canal.

## Contexto atual
- O banco de dados (`CustomCaption` e `CustomCaptionButton`) já existe e as migrations estão corretas.
- As rotas da API REST (`/api/channel/:channelId/custom-captions`) já estão criadas no `routes.go` e os controllers existem.
- O pipeline de envio já tenta buscar o Custom Caption, mas falha devido à lógica limitada de extração de hashtags (`extractHashtag` só pega a primeira).
- **A UI da Dashboard não existe**, impossibilitando o usuário de usar a funcionalidade.

## Arquivos analisados
- `internal/database/models/models.go` — Structs da Custom Caption.
- `internal/telegram/events/channelPost/utils_v2.go` — Função `extractHashtag` e `findCustomCaption`.
- `internal/telegram/events/channelPost/stage_transform_telego.go` — Aplicação da legenda.
- `internal/telegram/events/channelPost/stage_decorate_telego.go` — Extração duplicada da hashtag para botões.
- `internal/api/routes/routes.go` — Rotas de API existentes.

## Arquivos que poderão ser modificados
- `internal/telegram/events/channelPost/utils_v2.go`
- `internal/telegram/events/channelPost/stage_transform_telego.go`
- `internal/telegram/events/channelPost/stage_decorate_telego.go`
- `dashboard/src/App.tsx`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`

## Arquivos que serão criados
- `dashboard/src/components/CustomCaptionsCard.tsx` — Novo componente de UI.

## Estratégia de implementação

### 1. Correções no Backend (Pipeline)

**Passo 1.1: Melhorar extração de hashtags (`utils_v2.go`)**
Substituir a função `extractHashtag` atual (que só pega a 1ª hashtag) por `extractAllHashtags`:
```go
func extractAllHashtags(text string) []string {
    matches := hashtagRegex.FindAllStringSubmatch(text, -1)
    var hashtags []string
    for _, match := range matches {
        if len(match) > 1 {
            hashtags = append(hashtags, match[1]) // Sem a hash
        }
    }
    return hashtags
}
```

**Passo 1.2: Ajustar a busca pelo comando (`utils_v2.go`)**
```go
func findCustomCaptionFromText(channel *dbmodels.Channel, text string) (*dbmodels.CustomCaption, string) {
    hashtags := extractAllHashtags(text)
    for _, ht := range hashtags {
        for _, cc := range channel.CustomCaptions {
            // cc.Code pode estar com ou sem '#' salvo no banco. 
            cleanCode := strings.TrimPrefix(cc.Code, "#")
            if strings.EqualFold(cleanCode, ht) {
                return &cc, ht
            }
        }
    }
    return nil, ""
}
```
*Dessa forma, retorna o objeto e qual foi a hashtag exata ("ht") que disparou.*

**Passo 1.3: Atualizar o Pipeline (`stage_transform_telego.go` e `stage_decorate_telego.go`)**
- Em ambos os arquivos, chamar `findCustomCaptionFromText`.
- Se encontrado, remover **somente** a hashtag encontrada usando `strings.Replace(text, "#"+foundHashtag, "", 1)`.

### 2. Integração Front-end (Dashboard)

**Passo 2.1: Types e API Calls (`types.ts` e `api.ts`)**
- Assegurar que `CustomCaption` e `CustomCaptionButton` estão corretos no `types.ts`.
- Adicionar chamadas no `api.ts`:
  - `createCustomCaption(channelId, data)`
  - `updateCustomCaption(channelId, captionId, data)`
  - `deleteCustomCaption(channelId, captionId)`
  - Mesma coisa para os Botões atrelados à Caption.

**Passo 2.2: Criar o `CustomCaptionsCard.tsx`**
- Layout similar ao `ReactionsCard.tsx` ou `ButtonsCard.tsx`.
- Estado: listagem de "Comandos Ativos".
- Formulário para Novo Comando:
  - Input `Code`: (ex: `promo`)
  - Textarea `Caption`: (O texto da legenda)
  - Toggle `LinkPreview`
- Sistema de botões dinâmicos (idêntico ao `ButtonsCard.tsx`), mas atrelado especificamente ao comando selecionado.

**Passo 2.3: Inserir na Dashboard (`App.tsx`)**
- Adicionar a tab "Comandos" ou "Custom Captions" ao array `BASE_TABS`.
- Renderizar o `<CustomCaptionsCard channelId={channel.id} />` quando a tab for selecionada.

## Riscos
- O texto original (além da hashtag) pode ficar colado com a legenda customizada se não houver quebra de linha tratada. O backend já faz a formatação original + legenda custom, mas precisa revisar o espaçamento.
- Dupla extração de hashtag: Atualmente `transform` e `decorate` extraem a hashtag de forma independente. É importante que ambos achem a mesma hashtag se o texto tiver mais de uma.

## Impactos esperados
- Usuários poderão configurar comandos automáticos (ex: `#oferta`) que alteram instantaneamente a legenda e os botões postados.
- Nenhum impacto negativo na versão gratuita se a feature for liberada para todos, porém, sendo uma funcionalidade poderosa, pode facilmente ser limitada a contas Premium no futuro (ex: Free = 1 comando, Premium = infinito).

## Compatibilidade
- Suporte normal para Telegram.
- UI compatível com o WebApp existente.

## Rollback
- Remover a tab `Comandos` da Dashboard e reverter `utils_v2.go` para a versão `extractHashtag` simples.

## Como testar

### Dashboard
1. Rodar `npm run dev` na pasta dashboard.
2. Acessar a nova aba "Comandos".
3. Criar o comando `teste`, adicionar a legenda `Este é um teste.` e criar o botão `[🔗 Google | https://google.com]`.

### Bot
1. Postar uma foto no canal com a legenda: `Olha isso aqui! #teste`.
2. O bot deve enviar a foto, remover a palavra `#teste`, adicionar a frase `Este é um teste.` e anexar o botão do Google.
3. Postar com duas hashtags: `#foto #teste`. O bot deve processar normalmente, removendo apenas `#teste`.
