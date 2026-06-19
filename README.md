## Go-Mutation-Fuzzer
Uma prova de conceito (PoC) desenvolvida em Go para demonstrar a descoberta automatizada de vulnerabilidades de corrupção de memória (Memory Corruption) em processos nativos. A ferramenta atua como um motor de Fuzzing baseado em mutação cega, utilizando slices de bytes, concorrência agressiva via Goroutines e monitoramento de chamadas de sistema (syscall) para capturar falhas críticas de segmentação (SIGSEGV) em tempo real.
No caso do cenário de testes, a mutação tinha como objetivo subverter as validações de entrada de um parser vulnerável para provocar um Buffer Overflow e identificar o payload exato responsável pelo travamento.

## Objetivo
Compreender a fundo a arquitetura de Fuzzers modernos, o comportamento do Kernel do Linux diante de violações de limites de memória e o gerenciamento de concorrência massiva (I/O e processamento) na automatização da descoberta de Zero-Days e falhas lógicas em bibliotecas C/C++.

## Aviso Educacional
Este projeto foi desenvolvido estritamente para fins de estudo em engenharia reversa, segurança de sistemas operacionais, pesquisa de vulnerabilidades e desenvolvimento de ferramentas ofensivas (Toolsmithing). A técnica demonstrada deve ser utilizada apenas em ambientes controlados, laboratórios locais e aplicações onde se possui autorização explícita para testes de estresse, não devendo ser aplicada em infraestruturas de terceiros.

## Arquitetura e Funcionamento

    Ingestão de Semente (Corpus):
    Leitura isolada de um arquivo estruturalmente válido para a memória RAM, servindo como molde intocável para preservar metadados e contornar validações superficiais do alvo.

    Motor de Mutação Lógica:
    Aplicação do algoritmo de Bit-Flipping matemático (XOR 0xFF) em índices aleatórios gerenciados por fontes independentes de entropia, garantindo a corrupção de uma porcentagem exata do byte array.

    Concorrência Agressiva (Worker Pool):
    Orquestração de múltiplas threads simultâneas para geração e gravação de arquivos temporários mutantes em disco, saturando os núcleos do processador de forma paralela e otimizada.

    Monitoramento de Sinais (Syscall Interception):
    Invocação do binário alvo via exec.Command e acoplamento de rotinas de análise para diferenciar saídas de erro comuns (código 1) de interrupções críticas do Kernel, filtrando exclusivamente travamentos oriundos do Sinal 11 (SIGSEGV).

    Crash Triage Automatizado:
    Isolamento e salvamento em formato hash do payload mutante exato que causou o Segmentation Fault, permitindo a posterior reprodução e depuração no GDB para desenvolvimento do exploit (RCE).

## Setup de Ambiente de Testes
- OS: Ubuntu 24.04.4 LTS
- Kernel: 6.7.10-060710-generic
- Linguagem & Compilador: Go version 1.22+ (linux/amd64)
- Alvo de Teste: Binário customizado em C (com vulnerabilidade intencional de Stack Overflow)
- Arquitetura: Alvo: x86_64 | Host: x86_64
- Build & Flags: go build -o fuzzer fuzzer.go
- OS Defense: ASLR Ligado (Address Space Layout Randomization) - mitigação irrelevante na fase de Fuzzing, a ser contornada na fase de exploração.
