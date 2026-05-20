# Plano: Corrigir Erro de Foreign Key no Postgres ao Deletar Canal

## Pedido do usuário
O bot apresentou erro de violação de chave estrangeira (`fk_channels_buttons`) ao tentar excluir um canal no PostgreSQL.

## Objetivo técnico
Garantir que a exclusão de um canal e seus dados dependentes (botões, separadores, legendas) funcione corretamente no PostgreSQL, contornando a rigidez das restrições de integridade que o SQLite às vezes ignora ou trata de forma diferente.

## Contexto atual
No SQLite, o `CASCADE DELETE` costuma funcionar apenas com as tags do GORM se habilitado via PRAGMA. No PostgreSQL, a restrição de integridade é verificada no momento da execução. O erro `ERROR: update or delete on table "channels" violates foreign key constraint` indica que existem registros na tabela `buttons` apontando para o canal que está sendo deletado, e o banco impediu a ação para não deixar dados órfãos.

## Arquivos analisados
- `internal/database/repositories/channel.go`: O método `DeleteChannelWithRelations` está tentando deletar o canal diretamente, confiando que o banco fará o cascade.

## Arquivos que poderão ser modificados
- `internal/database/repositories/channel.go`

## Estratégia de implementação
Em vez de confiar apenas na configuração do banco de dados (que pode variar dependendo de como as tabelas foram criadas originalmente no Postgres), vamos implementar uma **limpeza manual em transação** dentro do repositório. 

Isso garante compatibilidade total entre SQLite e Postgres:
1. Abrir uma transação.
2. Deletar manualmente todos os registros das tabelas filhas (`buttons`, `separators`, `custom_captions`, `default_captions`) que possuem o `owner_channel_id` do canal alvo.
3. Deletar o canal por último.
4. Confirmar a transação.

## Passos detalhados

1. **Modificar `internal/database/repositories/channel.go`**:
    - Reescrever o método `DeleteChannelWithRelations`.
    - Usar `r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })`.
    - Dentro da transação, executar `Delete` nas tabelas:
        - `models.Button`
        - `models.Separator`
        - `models.CustomCaption` (e suas `CustomCaptionButton`)
        - `models.DefaultCaption` (e suas permissões)
    - Por fim, deletar o `models.Channel`.

## Riscos
- **Nenhum.** Como a operação será feita dentro de uma transação, se qualquer deleção falhar, nada será alterado no banco, mantendo a integridade.

## Impactos esperados
- A remoção de canais passará a funcionar 100% no PostgreSQL.
- Limpeza garantida de dados órfãos em qualquer banco de dados.

## Como testar

### Build
```bash
go build -v ./cmd/FreddyBot/...
```

### Testes
1. Remover o bot de um canal (gatilho automático de limpeza).
2. Tentar remover um canal manualmente via Dashboard Admin.
3. Verificar se o erro 23503 do Postgres desapareceu.
4. Conferir no banco se os botões do canal deletado também sumiram.

## Rollback
Reverter as mudanças no arquivo `internal/database/repositories/channel.go`.
