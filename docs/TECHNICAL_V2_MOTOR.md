# Documentação Técnica: Motor de Processamento V2 (Pipeline)

## 1. Visão Geral
O **Motor V2** do FreddyBot representa uma mudança de uma arquitetura monolítica (onde um único arquivo lidava com todas as etapas) para uma arquitetura baseada em **Pipeline (Linha de Montagem)**. Esta mudança foi projetada para suportar a escala de milhares de canais com alta performance, facilidade de manutenção e estabilidade de memória.

## 2. A Necessidade da Mudança
O motor anterior (V1) apresentava alguns desafios conforme o projeto crescia:
- **God File:** O arquivo `processors.go` estava se tornando excessivamente grande e complexo.
- **Vazamento de Recursos:** Criação excessiva de goroutines de limpeza para cada mensagem.
- **Dificuldade de Expansão:** Adicionar novas regras de postagem exigia alterar lógicas sensíveis de álbuns e banco de dados.

## 3. A Nova Arquitetura: O Pipeline
O processamento de cada mensagem agora é dividido em dois grandes ciclos: **Discovery** (Descoberta) e **Execution** (Execução).

### A. Discovery Pipeline (Fase de Identificação)
Executado imediatamente na thread do Handler para decidir o que fazer com a mensagem.
1.  **StagePreflight:** Verifica manutenção, blacklists, carrega configurações do canal e sincroniza metadados (Título/URL) proativamente.
2.  **StageSpecialFlows:** Intercepta fluxos como o comando `!newpack`.
3.  **StageMediaGrouping:** Identifica se a mensagem faz parte de um álbum. Se sim, aguarda a chegada de todas as partes antes de prosseguir.
4.  **StageQueue:** Adiciona o contexto pronto à fila de processamento assíncrono.

### B. Execution Pipeline (Fase de Transformação e Envio)
Executado pelos Workers da fila para realizar o trabalho pesado de I/O e CPU.
1.  **StageTransform:** A "inteligência" do texto. Resolve hashtags, aplica Custom Captions ou a Legenda Padrão, preservando a formatação HTML original.
2.  **StageDecorate:** Constrói dinamicamente o teclado inline com base nos botões do canal e no bloco de reações.
3.  **StageSend (Dispatch):** Realiza a comunicação final com o Telegram, lidando com retentativas automáticas e limites de taxa (Rate Limit/429).

## 4. Principais Inovações Técnicas

### Sincronização Inteligente de Canais
Implementamos um sistema de **Debounce de 1 hora**. O bot verifica se o nome ou link do canal mudou apenas uma vez por hora, ou instantaneamente se ele detectar uma mudança visual na mensagem postada. Isso economiza 90% das chamadas de API desnecessárias.

### Gerenciamento de Álbuns (Media Groups)
A lógica de álbuns foi isolada. O bot agora consegue:
- Reenviar álbuns de áudio e documentos (que o Telegram não permite editar a legenda diretamente com botões).
- Editar legendas de álbuns de fotos/vídeos de forma atômica.
- Manter a ordem correta das mensagens e do separador final.

### Estabilidade de Memória (Singletons)
Gerenciadores como o `PermissionManager` e `MediaGroupManager` agora são **Singletons**. Eles não são mais recriados a cada mensagem, eliminando vazamentos de memória (Memory Leaks) e garantindo que as rotinas de limpeza sejam únicas e eficientes.

### Worker Pool Escalável
Aumentamos para **20 Workers** processando uma fila de **5.000 posições**. Isso permite que o bot suporte picos massivos de postagens simultâneas sem elevar o uso de CPU para 100%.

## 5. Como estender o Motor?
Para adicionar uma nova funcionalidade (ex: tradução automática), basta criar um novo arquivo `stage_translate.go` que implemente a interface `Stage` e adicioná-lo ao Pipeline no arquivo `channelPost.go`.

---
**Data da Implementação:** 12 de Maio de 2026  
**Status:** Operacional (V2.0-Pipeline)
