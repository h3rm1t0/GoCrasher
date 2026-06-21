# GoCrash - In-Memory Concurrent Mutation Fuzzer

Um motor de *fuzzing* de mutação cega (Bit-Flipping) desenvolvido em Go, focado na descoberta de vulnerabilidades de corrupção de memória (*Memory Corruption*) em bibliotecas e binários compilados em C/C++. A ferramenta opera inteiramente na memória RAM via STDIN, garantindo alta performance de I/O e evasão de gargalos de disco.

## Objetivo
O projeto foi construído do zero para atuar na triagem agressiva de *parsers* complexos (como o ImageMagick e bibliotecas de processamento de documentos). O motor automatiza a injeção de anomalias matemáticas, gerencia o ciclo de vida do processo alvo para evitar travamentos (*hangs*) e intercepta sinais de Kernel para capturar execuções que resultam em `SIGSEGV` ou vazamentos detectados pelo `AddressSanitizer` (ASAN).

## Arquitetura e Funcionalidades
* **In-Memory Fuzzing (STDIN):** Os dados mutados são canalizados diretamente para a entrada padrão do alvo, multiplicando a velocidade de execução em comparação com abordagens baseadas em gravação de arquivos em disco.
* **Concorrência Massiva:** Utilização de *Goroutines* (*Worker Pool*) com sementes de entropia independentes, saturando múltiplos núcleos do processador para testes paralelos ininterruptos.
* **Detecção Dupla de Anomalias:** Captura nativa de falhas de segmentação (Sinal 11) integradas ao monitoramento da saída de erro (STDERR) para isolar micro-vazamentos estruturais identificados pelo ASAN.
* **Prevenção de Livelock/Hang:** Implementação de limite de tempo de execução via pacote `context`, assassinando processos alvos que entram em laços infinitos e liberando as *threads* do fuzzer.
* **Alertas em Tempo Real:** Integração nativa com a API do Telegram para notificação imediata (*Push*) no momento em que um *crash* explorável é salvo no disco.

## Setup e Ambiente de Testes
O *fuzzer* requer que o alvo seja preferencialmente compilado com instrumentação (`-fsanitize=address`) para maximizar a captura de erros *Off-by-One* e *Use-After-Free*.

**Compilação do Fuzzer:**
```bash
go build -o GoCrash GoCrasher.go
