# Plano: Deleção em Massa na Auditoria e Melhorias Visuais

## Pedido do usuário
1. Adicionar um botão de "Excluir Todos os Canais" ao lado do nome do usuário na aba de Auditoria, para remover em massa os canais onde o @XavolaBot foi detectado.
2. Melhorar o distanciamento visual entre os blocos de usuários na lista de auditoria, pois estão muito próximos.

## Objetivo técnico
1. Implementar um endpoint de API para deleção em massa de canais de um usuário.
2. Atualizar o componente `AdminAuditTab` para incluir o botão de ação e chamar a API.
3. Ajustar o layout CSS/Tailwind para proporcionar uma separação visual mais clara entre os grupos de usuários.

## Contexto atual
A aba de Auditoria identifica canais problemáticos, mas exige que o administrador entre em cada um para desconectar o bot. A deleção em massa agilizará a limpeza do sistema.

## Arquivos analisados
- `internal/api/controllers/adminController/auditController.go`
- `dashboard/src/components/AdminAuditTab.tsx`
- `internal/api/routes/routes.go`

## Arquivos que poderão ser modificados
- `internal/api/controllers/adminController/auditController.go`
- `internal/api/routes/routes.go`
- `dashboard/src/components/AdminAuditTab.tsx`
- `dashboard/src/api.ts`

## Estratégia de implementação

### 1. Backend (Go)
- No `AuditController`, criar o método `BulkDeleteUserChannels`.
- Este método receberá um JSON com uma lista de IDs de canais e o ID do usuário dono.
- Utilizará uma transação ou loop seguro chamando o `ChannelService.DisconnectChannel` para garantir que a saída do canal e a limpeza do banco ocorram para cada item.

### 2. Frontend (React)
- Em `api.ts`, adicionar `bulkDeleteChannels`.
- Em `AdminAuditTab.tsx`:
    - Adicionar um estado de carregamento específico por usuário (opcional, ou global simplificado).
    - Adicionar o botão "Limpar Canais" ao lado do nome do usuário.
    - Implementar a lógica de confirmação antes de deletar.
    - Aumentar o `margin-bottom` e adicionar uma borda ou separador entre os blocos de usuários.

## Passos detalhados

1. **Modificar `internal/api/controllers/adminController/auditController.go`**:
    - Adicionar struct `BulkDeleteRequest`.
    - Implementar `BulkDeleteUserChannels`.

2. **Registrar Rota**:
    - Em `internal/api/routes/routes.go`, adicionar `adminRoute.POST("/audit/bulk-delete", auditController.BulkDeleteUserChannels)`.

3. **Atualizar Dashboard**:
    - `api.ts`: Adicionar `bulkDeleteChannels`.
    - `AdminAuditTab.tsx`: 
        - Refatorar o loop de renderização para usar `gap-8` ou `space-y-8`.
        - Adicionar o botão de delete com ícone `Trash2`.
        - Adicionar modal de confirmação.

## Riscos
- **Deleção Acidental**: O administrador pode clicar por engano.
  - *Mitigação:* Usar `window.confirm` ou o componente `ConfirmModal` existente.
- **Rate Limit do Telegram**: Se um usuário tiver centenas de canais, o bot tentando sair de todos de uma vez pode ser bloqueado.
  - *Mitigação:* Adicionar um pequeno delay entre as saídas se necessário (sleep de 100ms).

## Impactos esperados
- Limpeza massiva de canais obsoletos com poucos cliques.
- Dashboard mais legível e profissional.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/... && cd dashboard && npm run build
```

### Testes
1. Realizar a auditoria.
2. Clicar no botão de excluir ao lado de um usuário com múltiplos canais.
3. Confirmar a ação.
4. Verificar se os canais sumiram da lista de auditoria e do banco de dados.

## Rollback
Reverter as mudanças nos controladores e no componente de auditoria.
