# Decisões Arquiteturais

## Decisão: Unificação do Núcleo via Core Services
### Data
2026-05-13

### Contexto
O projeto tinha lógica de negócio e acesso a dados duplicados entre os repositórios, controladores da API e handlers do Bot. Repositórios continham lógica complexa que dificultava testes e reutilização.

### Decisão tomada
Implementar uma camada de `Core Services` em `internal/core/services`. Toda a lógica de negócio e acesso a dados (GORM) deve passar por essa camada. Controladores da API e Handlers do Bot tornam-se "cascas" finas que apenas validam entrada/saída e chamam os serviços.

### Motivo
- **DRY (Don't Repeat Yourself):** Reutilização de lógica entre API e Bot.
- **Testabilidade:** Lógica isolada em serviços é mais fácil de testar.
- **Manutenibilidade:** Mudanças na regra de negócio são feitas em um único lugar.
- **Padronização:** Respostas da API via Generics `APIResponse[T]`.

### Impacto
- Removido acesso direto aos repositórios do `AppContainer`.
- Handlers do Bot e Middlewares migrados para usar Serviços.
- API refatorada para usar controladores com serviços e respostas padronizadas.

## Decisão: Mapeamento de Mensagens Inline para Sessões de Postagem
### Data
2026-05-13

### Contexto
Usuários criam postagens no Post Builder e as compartilham via modo inline. Ao votar nessas mensagens, o Telegram não fornece o teclado original, impedindo a atualização visual dos contadores de votos.

### Decisão tomada
Implementar um handler de `ChosenInlineResult` para mapear o `inline_message_id` gerado pelo Telegram para o ID da sessão da postagem no Redis.

### Motivo
Permite reconstruir o teclado original no momento do voto, possibilitando a atualização visual dos contadores enquanto mantém o `CallbackData` limpo (`vote:emoji`).

### Impacto
- Necessário ativar `Inline Feedback` no BotFather.
- Dependência de Redis para o mapeamento temporário (24h).

## Decisão: Uso de MatchFunc Customizado para ChosenInlineResult
### Data
2026-05-13

### Contexto
A biblioteca `go-telegram/bot` v1.19.0 não possui a constante `bot.HandlerTypeChosenInlineResult`, impossibilitando o uso de `RegisterHandler` padrão para esse tipo de update.

### Decisão tomada
Utilizar `RegisterHandlerMatchFunc` com uma função de match manual (`matchChosenInlineResult`) que verifica a presença do campo `ChosenInlineResult` no objeto `models.Update`.

### Motivo
Contornar a limitação da biblioteca sem a necessidade de forks ou atualização imediata da dependência, mantendo a funcionalidade de mapeamento de mensagens inline.

### Impacto
- Substituição do registro padrão por match manual em `internal/telegram/events/loader.go`.
- Código permanece compatível com versões futuras caso a constante seja adicionada.
