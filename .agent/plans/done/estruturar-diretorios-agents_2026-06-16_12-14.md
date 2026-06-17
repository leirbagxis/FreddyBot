# Plano: estruturar-diretorios-agents

## Pedido do usuário
Usuário pediu para ler e executar tudo que o AGENTS.md especifica.

## Objetivo
Garantir que a estrutura de diretórios `.agent/` esteja completa conforme especificado no AGENTS.md, com todos os subdiretórios e arquivos necessários.

## Contexto atual
- `.agent/` existe parcialmente:
  - `plans/done/` ✅ (com 100+ planos concluídos)
  - `plans/pending/` ❌ ausente
  - `plans/approved/` ❌ ausente
  - `memory/` ❌ ausente, `memory.md` não existe
  - `context.md` ✅
  - `decisions.md` ✅ (com decisões registradas)
  - `skills/` diretório extra não especificado no AGENTS.md

## Arquivos analisados
- AGENTS.md (regras e estrutura esperada)
- .agent/context.md
- .agent/decisions.md
- .agent/ (estrutura atual)

## Arquivos que poderão ser modificados
- .agent/memory/memory.md (criação)
- .agent/plans/pending/ (criação do diretório)
- .agent/plans/approved/ (criação do diretório)

## Estratégia de implementação
1. Criar diretórios faltantes: `plans/pending/`, `plans/approved/`, `memory/`
2. Extrair informações da arquitectura do `context.md` e `decisions.md` para popular `memory/memory.md`
3. O diretório `skills/` será mantido pois contém skills carregáveis pelo sistema

## Passos detalhados
1. Criar `mkdir -p .agent/plans/pending .agent/plans/approved .agent/memory`
2. Criar `.agent/memory/memory.md` com resumo da arquitetura, convenções, decisões técnicas e problemas conhecidos extraídos de context.md e decisions.md
3. Salvar plano atual em pending

## Riscos
- Nenhum risco significativo (apenas criação de diretórios e arquivos de documentação)
- Não há alteração de código fonte

## Impactos esperados
- Estrutura `.agent/` 100% compatível com AGENTS.md
- Memória persistente disponível para consultas futuras do agente
- Fluxo de planos (pending → approved → done) funcional

## Compatibilidade
- Linux ✅
- macOS ✅
- Windows ✅
- Docker ✅
- CI/CD ✅

## Como testar
```bash
ls -la .agent/plans/pending .agent/plans/approved .agent/memory/memory.md
```

## Rollback
```bash
rm -rf .agent/plans/pending .agent/plans/approved .agent/memory
```

## Observações
- Diretório `.agent/skills/` existe mas não está no spec do AGENTS.md; será mantido pois contém skills carregáveis
