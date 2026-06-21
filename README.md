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

**Aviso Educacional**
Este projeto é um artefato técnico desenvolvido estritamente para estudos avançados em engenharia reversa e pesquisa independente de vulnerabilidades (Vulnerability Research). A ferramenta deve ser utilizada exclusivamente em ambientes de laboratório locais, targets sob contrato de Pentest e programas de Bug Bounty com escopo aberto. O autor não se responsabiliza pelo uso do motor contra infraestruturas onde não haja autorização explícita.

## Setup e Ambiente de Testes
OS do Ambiente: Ubuntu 24.04.4 LTS x86_64
Go: go version go1.22.2 linux/amd64
O *fuzzer* requer que o alvo seja preferencialmente compilado com instrumentação (`-fsanitize=address`) para maximizar a captura de erros *Off-by-One* e *Use-After-Free*.

## Execução Padrão:
O fuzzer exige a semente mínima (corpus) e o caminho do binário alvo. Exemplo: ./GoCrash seed.svg magick

Nota: Variáveis de ambiente específicas do alvo (ex: MAGICK_CONFIGURE_PATH) já estão injetadas de forma autônoma no subprocesso do motor.

## Compilação do Fuzzer:
```bash
go build -o GoCrash GoCrasher.go
utilizada exclusivamente em ambientes de laboratório locais, targets sob contrato de Pentest e programas de Bug Bounty com escopo aberto. O autor não se responsabiliza pelo uso do motor contra infraestruturas onde não haja autorização explícita.
