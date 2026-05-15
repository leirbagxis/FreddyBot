# Plano: fix-nil-pointer-panic_2026-05-15_11-00.md

## Pedido do usuário
Corrigir o erro (panic: invalid memory address) que ocorre ao acessar a rota do canal no dashboard (`/api/channel/:id`).

## Objetivo
Garantir que a função `ToUserDTO` não tente desreferenciar um ponteiro nulo e assegurar que a relação `Owner` do canal seja corretamente carregada (hidratada) pelo GORM nas consultas.

## Contexto atual
- A rota `/api/channel/:id` chama `GetChannelByIDController`.
- O controller chama `dto.ToUserDTO(channel.Owner)`.
- Se `channel.Owner` for nulo, a aplicação sofre um panic pois `ToUserDTO` não faz a checagem.
- O método `ChannelRepository.GetChannelByID` usa `Joins("Owner")`, o qual pode não estar populando corretamente o struct `User` aninhado no SQLite após as mudanças recentes de performance.

## Arquivos analisados
- `internal/api/dto/mapper.go`
- `internal/api/controllers/channelController.go`
- `internal/database/repositories/channel.go`

## Arquivos que poderão ser modificados
- `internal/api/dto/mapper.go`
- `internal/database/repositories/channel.go`

## Estratégia de implementação

### 1. Proteção contra Nil Pointer no DTO Mapper
Adicionar checagens iniciais nas funções do `mapper.go`:
```go
func ToUserDTO(u *models.User) UserDTO {
    if u == nil {
        return UserDTO{}
    }
    // ...
}

func ToChannelDTO(c *models.Channel) ChannelDTO {
    if c == nil {
        return ChannelDTO{}
    }
    // ...
}
```

### 2. Garantir o Preload do Owner
No arquivo `internal/database/repositories/channel.go`, alterar a linha `Joins("Owner")` para `Preload("Owner")` no método `GetChannelByID`. Isso obriga o GORM a fazer uma consulta separada ou usar a inteligência de eager loading para garantir que o objeto em memória não venha com ponteiro vazio (se o dono existir).

## Passos detalhados
1. Alterar `mapper.go` para adicionar as condições de segurança.
2. Alterar `channel.go` (repository) mudando a estratégia de hidratação de `Owner` de `Joins` para `Preload`.

## Riscos
- Mudar de `Joins` para `Preload` executa uma query extra silenciosa no banco (GORM faz isso internamente). Como se trata da consulta de 1 canal específico (`GetChannelByID`), o impacto no banco é de ~1ms (desprezível), especialmente estando atrás do novo Cache L1 implementado.

## Impactos esperados
- Estabilidade da rota de Dashboard e fim do panic 500.

## Compatibilidade
- Linux, Docker

## Como testar

### Build
```bash
go build ./...
```

### Execução
Acessar o painel Web (`/dashboard/:id`) logado. A API deve retornar `200 OK` e o JSON com os dados do canal e do usuário.

## Rollback
Desfazer as alterações pontuais nos 2 arquivos.