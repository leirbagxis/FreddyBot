# Plano: subs-toggles-links-dinamicos_2026-05-16_16-00.md

## Pedido do usuário
Adicionar 3 sub-toggles (opções) na dashboard para os "Links Dinâmicos":
1. Habilitar/desabilitar botões do bot na postagem.
2. Habilitar/desabilitar legendas do bot na postagem.
3. Habilitar/desabilitar reações do bot na postagem.
Essas regras só devem ser aplicadas caso um botão dinâmico seja detectado na mensagem.

## Objetivo técnico
Expandir o modelo de banco de dados para suportar as configurações granulares dos Links Dinâmicos e aplicar essas condicionais no pipeline do Telegram (especificamente no `StageTransform`).

## Contexto atual
Atualmente, o `DynamicLinks` é um toggle booleano simples. Se ativado, ele extrai os links e os adiciona à mensagem sem interferir nas configurações padrão do canal (botões extras, legendas e reações são aplicados normalmente).

## Arquivos analisados
- `internal/database/models/models.go`
- `internal/database/repositories/channel.go`
- `internal/api/controllers/permissionController.go`
- `internal/telegram/events/channelPost/stage_transform.go`
- `dashboard/src/types.ts`
- `dashboard/src/api.ts`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- **Backend:**
    - `models.go` (Adição de `DLBotButtons`, `DLBotCaptions`, `DLBotReactions`)
    - `repositories/channel.go` e `services/channels.go` (Atualizar assinatura do `UpdateDynamicLinks` ou criar novo update)
    - `permissionController.go` (Receber as novas configurações no JSON)
    - `stage_transform.go` (Aplicar as regras condicionais caso botões dinâmicos sejam extraídos)
- **Frontend:**
    - `types.ts` (Adicionar os 3 booleanos no Channel)
    - `api.ts` (Enviar o objeto completo no PUT)
    - `App.tsx` (Renderizar os 3 sub-toggles visualmente de forma indentada/aninhada abaixo do toggle principal de Links Dinâmicos)

## Estratégia de implementação
1. **Banco de Dados:**
    - Adicionar os campos booleanos com padrão `true` para manter a retrocompatibilidade.
2. **Backend API:**
    - Alterar o `UpdateDynamicLinksController` para processar um JSON com todas as 4 propriedades.
3. **Pipeline do Telegram:**
    - No `StageTransform`, registrar uma flag se os botões dinâmicos foram extraídos.
    - Se a flag for verdadeira:
        - Omitir `dbCaption` se `DLBotCaptions == false`.
        - Não fundir a lista de botões da DB se `DLBotButtons == false`.
        - Desativar a flag de reações (`CanAddReactions = false`) se `DLBotReactions == false`.
4. **Dashboard:**
    - Renderizar os 3 toggles sob o card "Links Dinâmicos". Eles só devem estar visíveis ou habilitados se o "Links Dinâmicos" principal estiver ativo.

## Passos detalhados
1. Atualizar modelo GORM.
2. Atualizar Repositório e Controller na API.
3. Implementar a lógica condicional no `StageTransform`.
4. Refatorar a interface no `App.tsx`.

## Riscos
- O banco de dados precisará criar colunas novas. O GORM fará isso automaticamente via AutoMigrate.

## Impactos esperados
- Maior granularidade para donos de canais que desejam que mensagens com links dinâmicos fiquem isoladas e mais limpas, sem a assinatura padrão do bot.

## Como testar
1. Ativar Links Dinâmicos e desativar os 3 sub-toggles.
2. Enviar mensagem com padrão `!Botão\n!Link`.
3. Validar se a mensagem final contém APENAS o botão dinâmico.

## Rollback
Desativar via Dashboard ou reverter os arquivos afetados.