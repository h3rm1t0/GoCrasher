package main

import (
	"crypto/md5"
	crand "crypto/rand" // Alias para evitar colisão com math/rand
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#@!%&$*()<>çÇ"

func GerarStringRandomica(i int) string {
	var sb strings.Builder
	sb.Grow(i)
	for n := 0; n <= i; n++ {
		num, _ := crand.Int(crand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[num.Int64()])
	}
	return sb.String()
}

func Crasher(alvo string, corrompido []byte) {
	os.MkdirAll("Crashs", 0755)
	ext := ".ini"

	corromp_arq := GerarStringRandomica(5) + ext
	err := os.WriteFile(corromp_arq, corrompido, 0644)
	if err != nil {
		fmt.Printf("%s - Falha ao salvar arquivo temporário\n", momento())
		return
	}
	defer os.Remove(corromp_arq)

	cmd := exec.Command("./"+alvo, corromp_arq)
	err_pg := cmd.Run()

	if err_pg != nil {
		if exitErr, ok := err_pg.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() && status.Signal() == syscall.SIGSEGV {
					string_aleatoria := GerarStringRandomica(5)
					hash := md5.Sum([]byte(string_aleatoria))
					nome_arq_final := fmt.Sprintf("%x%s", hash, ext)
					dir := filepath.Join("Crashs", nome_arq_final)
					err = os.WriteFile(dir, corrompido, 0644)
					if err == nil {
						fmt.Printf("%s - [CRASH DETECTADO] Payload salvo em: %s\n", momento(), dir)
					}
				}
			}
		}
	}
}

func Mutacao(seed []byte, rng *rand.Rand) []byte {
	copia := make([]byte, len(seed))
	copy(copia, seed)
	qnt_mut := len(copia) / 100
	if qnt_mut == 0 {
		qnt_mut = 1
	}
	for mut := 0; mut < qnt_mut; mut++ {
		indice := rng.Intn(len(copia))
		copia[indice] ^= 0xFF
	}
	return copia
}

func Worker(id int, alvo string, seed []byte) {
	fonte := rand.NewSource(time.Now().UnixNano() + int64(id))
	rng := rand.New(fonte)

	for {
		corrompido := Mutacao(seed, rng)
		Crasher(alvo, corrompido)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("%s - Informe os argumentos corretos.\n", momento())
		fmt.Println("Ex.: ./GoCrash <CAMINHO_DA_SEED> <NOME_DO_PROCESSO_ALVO>")
		return
	}

	alvo := os.Args[2]
	seed, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("%s - Falha ao abrir a seed.\n", momento())
		return
	}
	
	numWorkers := 10
	fmt.Printf("%s - [START] Iniciando Motor de Fuzzing...\n", momento())
	fmt.Printf("%s - Alvo: %s | Semente: %s | Threads Ativas: %d\n", momento(), alvo, os.Args[1], numWorkers)
	fmt.Printf("%s - Pressione CTRL+C para abortar.\n", momento())
	fmt.Println(strings.Repeat("-", 60))

	bloqueio := make(chan struct{})

	for i := 1; i <= numWorkers; i++ {
		go Worker(i, alvo, seed)
	}
	<-bloqueio 
}

func momento() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
