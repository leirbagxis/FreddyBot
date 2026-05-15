# Plano: refatoracao-nucleo-unico_2026-05-13_12-45.md

## Pedido do usuário
Refazer a parte de conexão, repositórios e serviços para unificar o Bot e a API, mantendo as tabelas intactas, mas removendo duplicidade.

## Objetivo
Estabelecer uma arquitetura de "Núcleo Único" onde o `AppContainer` fornece serviços compartilhados que garantem consistência de dados, lógica de negócio e cache tanto para o Telegram quanto para o Dashboard.

## Contexto atual
- `AppContainer` expõe repositórios e serviços de forma desorganizada.
- API acessa o banco (GORM) diretamente em muitos serviços.
- Bot acessa repositórios diretamente, pulando validações que existem na API.
- Cache (Redis) é invalidado de forma inconsistente entre as duas interfaces.

## Arquivos analisados
- `internal/container/appContainer.go`
- `internal/database/repositories/`
- `internal/api/service/`
- `internal/database/models/models.go`

## Arquivos que poderão ser modificados
- `internal/container/appContainer.go`
- `internal/database/repositories/user.go` (e outros)
- `internal/api/service/` (Será transformado em `internal/core/services`)
- `internal/telegram/callbacks/` e `events/`
- `internal/api/controllers/`

## Estratégia de implementação
1. **Consolidação do DB (Repositórios):** Mover TODO o código GORM (`db.Where`, `db.Create`, etc.) para os arquivos em `internal/database/repositories`. Ninguém fora dos repositórios toca no objeto `gorm.DB`.
2. **Serviços de Domínio (Core):** Transformar os atuais serviços da API em serviços de domínio que aceitam contextos e IDs, independentes de protocolo (HTTP ou Telegram).
3. **Injeção de Dependência Limpa:** O `AppContainer` será o dono dos Repositórios e dos Serviços. Ele injeta os repositórios nos serviços durante a inicialização.
4. **Remoção de Especialização Prematura:** Eliminar as pastas `adminRepositories` e `adminService`, movendo a lógica para os repositórios/serviços de `User` e `Channel`.

## Passos detalhados

1. **Fase 1: Preparação do Container**
   - Ajustar `AppContainer` para carregar todos os repositórios primeiro.
   - Inicializar os Serviços passando os repositórios necessários.

2. **Fase 2: Repositórios (Refatoração Cirúrgica)**
   - Mover consultas de Admin (Listar usuários, etc) para `UserRepository`.
   - Garantir que métodos de escrita (Create/Update/Delete) já incluam a chamada de `InvalidateCache`.
   - Remover `pkg/parser` de dentro dos repositórios.

3. **Fase 3: Unificação de Serviços**
   - Refatorar `ButtonsService`, `ChannelService`, etc., para que não dependam de tipos específicos do Gin/HTTP.
   - Implementar validações de negócio dentro destes serviços.

4. **Fase 4: Migração de Interface**
   - Atualizar os controladores da API para chamar os novos serviços.
   - Atualizar os handlers do Bot para chamar os serviços em vez dos repositórios (quando houver lógica de negócio envolvida).

5. **Fase 5: Limpeza Técnica**
   - Deletar pastas duplicadas.
   - Padronizar logs de erro e sucesso.

## Riscos
- Incompatibilidade de tipos de retorno (necessário usar structs de domínio limpas).
- Tempo de refatoração pode ser extenso devido ao volume de arquivos.

## Impactos esperados
- **Integridade:** Cache sempre sincronizado.
- **Velocidade de Desenvolvimento:** Regras de negócio alteradas em um único lugar.
- **Código Limpo:** Fim do "espaguete" de acesso a dados.

## Compatibilidade
- Linux, Docker, Go 1.24.

## Como testar
- Build do servidor: `make build-server`.
- Testar fluxo completo de criação de botão via Bot e verificar se reflete na API imediatamente.
- Testar exclusão de canal via Dashboard e verificar se o Bot limpa o estado.

## Rollback
- Git checkout para o estado anterior (recomenda-se criar uma branch `refactor/single-core`).

## Observações
O banco de dados (SQLite/Postgres) não sofrerá migrações de schema, apenas o acesso ao driver será reorganizado.
