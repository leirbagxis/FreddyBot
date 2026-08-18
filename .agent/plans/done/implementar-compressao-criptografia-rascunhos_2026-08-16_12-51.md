# Plano: implementar-compressao-criptografia-rascunhos_2026-08-16_12-51

## Pedido do usuário
Implementar o pipeline de **Compactação GZIP + Criptografia AES-256-GCM** (*Compress-Then-Encrypt*) no armazenamento de rascunhos (`UserPostTemplate`), garantindo economia de até 80% do espaço de banco de dados e proteção total dos dados contra vazamentos.

## Objetivo
1. Criar o pacote `pkg/crypto` com funções de compactação `gzip` + cifragem `AES-256-GCM`.
2. Integrar a compactação e criptografia no `UserPostTemplateService`, cifrando os rascunhos antes de salvar no banco de dados e decifrando ao ler.
3. Manter **compatibilidade retroativa** total (caso haja rascunhos salvos anteriormente sem criptografia, o sistema identifica e lê em texto claro sem erros).

## Arquivos analisados
- `internal/core/services/user_post_template.go`
- `pkg/config/config.go`
- `internal/container/appContainer.go`

## Arquivos que serão criados / modificados
- `pkg/crypto/crypto.go` *(novo)*
- `pkg/crypto/crypto_test.go` *(novo)*
- `internal/core/services/user_post_template.go`

## Estratégia de implementação

1. **Pacote `pkg/crypto`**:
   - `CompressAndEncrypt(data []byte, secretKey string) (string, error)`:
     1. Compacta o array de bytes via `gzip`.
     2. Cifra o conteúdo via `AES-256-GCM` usando a chave secreta.
     3. Codifica em Base64 prefixado com `"enc:v1:"`.
   - `DecryptAndDecompress(encodedStr string, secretKey string) ([]byte, error)`:
     1. Verifica se inicia com `"enc:v1:"`. Se não iniciar, entende como dado legado em texto claro (JSON antigo) e retorna diretamente.
     2. Decodifica Base64.
     3. Decifra com `AES-256-GCM`.
     4. Descompacta via `gzip`.

2. **Integração no `UserPostTemplateService`**:
   - `SaveTemplate`: executa `crypto.CompressAndEncrypt([]byte(templateData), secretKey)` e armazena a string cifrada.
   - `GetTemplateByID` / `ListTemplates`: executa `crypto.DecryptAndDecompress(tpl.TemplateData, secretKey)` reconstruindo o JSON original.

3. **Validação & Testes**:
   - Testes unitários do pacote `crypto_test.go` (verificar redução de tamanho e ciclo de cifragem/decifragem).
   - Executar suíte geral (`go test ./...`).
   - Compilar o executável (`go build ./cmd/FreddyBot`).

## Passos detalhados
1. Salvar o plano em `.agent/plans/pending/implementar-compressao-criptografia-rascunhos_2026-08-16_12-51.md`.
2. Apresentar o resumo ao usuário e solicitar aprovação explícita.
3. Criar `pkg/crypto/crypto.go` e `pkg/crypto/crypto_test.go`.
4. Atualizar `internal/core/services/user_post_template.go`.
5. Rodar testes e compilar (`go test ./...`, `go build ./cmd/FreddyBot`).

## Riscos
- Mínimo. Mantém 100% de compatibilidade retroativa com rascunhos já existentes.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD.

## Como testar

### Testes Go
```bash
go test ./pkg/crypto/... ./internal/...
go build ./cmd/FreddyBot
```

## Rollback
`git checkout internal/core/services/user_post_template.go && rm -rf pkg/crypto`
