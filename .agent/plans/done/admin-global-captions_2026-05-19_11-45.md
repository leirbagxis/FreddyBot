# Plano: Configurações Globais de Legendas no Painel Admin

## Pedido do usuário
- Criar campos na Dashboard Admin para configurar as legendas globais: "Legenda Padrão" e "Legenda de Novo Pack".
- Definir o valor inicial da "Legenda Padrão" como: `🐈‍⠀៹ [t.me/legendasbot](https://t.me/botusername)  ‹` (ajustado para suportar marcações de variável posteriormente).
- Essas legendas devem ser usadas como valores iniciais quando um novo canal for adicionado ao sistema.

## Objetivo
1. Estender a tabela `ServerConfig` para armazenar `GlobalDefaultCaption` e `GlobalNewPackCaption`.
2. Atualizar as rotas e controladores do backend para receber e enviar esses novos campos.
3. Modificar o componente frontend `AdminConfigTab.tsx` para exibir dois `RichTextEditor` permitindo a edição dessas legendas.
4. Conectar a criação de novos canais (`CreateChannelWithDefaults`) a essas configurações globais.

## Contexto atual
- O painel de administrador (`AdminConfigTab.tsx`) exibe atualmente apenas configurações booleanas (Maintence e ForceJoin).
- A tabela `ServerConfig` já existe e armazena o estado global do bot.
- O fluxo de adição de canais que corrigimos recentemente envia strings vazias `""` para as legendas. Precisamos conectar essas funções.

## Arquivos analisados
- `internal/database/models/models.go`
- `internal/core/services/server.go`
- `internal/api/controllers/adminController/configController.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/components/AdminConfigTab.tsx`
- `internal/telegram/handlers/events/addChannel/addChannel.go`

## Arquivos que poderão ser modificados
- `internal/database/models/models.go`
- `internal/core/services/server.go`
- `internal/api/controllers/adminController/configController.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/components/AdminConfigTab.tsx`
- `internal/telegram/handlers/events/addChannel/addChannel.go`
- `internal/core/services/channels.go`

## Estratégia de implementação

1. **Banco de Dados:**
   - Adicionar `GlobalDefaultCaption` e `GlobalNewPackCaption` à struct `ServerConfig` em `models.go`.

2. **Backend API:**
   - Em `ServerService.UpdateConfig`, adicionar os novos parâmetros.
   - Em `ConfigController.UpdateConfig`, atualizar a validação do `body` e passar os novos parâmetros para o serviço.

3. **Frontend:**
   - Atualizar a interface `ServerConfig` em `types.ts`.
   - Atualizar a função `updateServerConfig` em `api.ts`.
   - Modificar `AdminConfigTab.tsx` para importar o componente `RichTextEditor` e criar seções para edição das legendas, com botões "Salvar" individuais ou acoplados ao botão principal, aproveitando o estado de salvamento.

4. **Integração na Adição do Canal:**
   - Em `AskAddChannelHandlerTelego` e `AddYesHandlerTelego`, será necessário obter a instância do `ServerService` (do contêiner) para buscar as configurações.
   - Ler os valores globais e passá-los para `c.ChannelService.CreateChannelWithDefaults`.

## Passos detalhados
1. Atualizar o `ServerConfig` (models.go, server.go, configController.go).
2. Atualizar o frontend (types.ts, api.ts, AdminConfigTab.tsx).
3. Modificar a injeção do contêiner em `addChannel.go` para acessar o `ServerConfig` e injetar as strings globais durante o cadastro do canal.

## Riscos
- Migração do banco: como o GORM já faz auto-migrate, novos campos de texto serão criados automaticamente, mas podem ficar nulos inicialmente. Teremos que tratar caso sejam nulos e definir um "default" via código se aplicável.

## Impactos esperados
- O administrador poderá trocar a legenda inicial de todos os novos canais diretamente pela dashboard.
- A consistência do bot é aumentada removendo valores "hardcoded" no banco ou no YAML (para novos canais).

## Como testar
1. Abrir a Dashboard e ir na aba "Admin" -> "Configurações".
2. Editar as legendas globais e salvar.
3. Dar reload na página para confirmar se os dados persistiram.
4. Adicionar um novo canal no bot.
5. Ir nas configurações desse canal e verificar se a legenda padrão importada é a mesma que foi definida no painel Admin.

## Rollback
- Reverter as alterações nos arquivos modificados usando o histórico do git.
