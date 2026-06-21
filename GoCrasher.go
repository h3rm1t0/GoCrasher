//O motor alvo pode ser modificado de acordo com a lógica da aplicação que é objeto de pesquisa, se desejar o código possúi uma função que habilita integração com bots do telegram.

package main

import (
    "bytes"
    "context"
    "crypto/md5"
    crand "crypto/rand"
    "fmt"
    "math/big"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "syscall"
    "time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#@!%&$*()<>çÇ"

const mensagem = "[INFO] Sucesso em um crash contra magick."

func telegram(mensagem string, tipoCrash string, dir string) {
    mensagem = fmt.Sprintf("%s\nTipo:%s\nArquivo:%s", tipoCrash, dir)
    var telegram_key string = "<INSIRA SUA API KEY DO TELEGRAM AQUI>"
    var chat_id string = "<INSIRA O ID DO SEU CHAT NO TELEGRAM>"
    ritmo := time.NewTicker(time.Millisecond / time.Duration(2500))
    defer ritmo.Stop()

    api_url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegram_key)

    dados := url.Values{}
    dados.Set("chat_id", chat_id)
    dados.Set("text", mensagem)
    dados.Set("parse_mode", "Markdown")

    <-ritmo.C
    resp, err := http.PostForm(api_url, dados)
    if err != nil {
        fmt.Printf(momento()+" - [ERRO] Falha ao conectar na API do Telegram: %v\n", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == 200 {
        fmt.Printf("%s", momento()+" - [INFO] Notificação enviada para o celular com sucesso!\n")
    } else {
        fmt.Printf(momento()+" - [ERRO] Telegram retornou status HTTP %d\n", resp.StatusCode)
    }

}

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

    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
    defer cancel()

    cmd := exec.CommandContext(ctx, "./"+alvo, "svg:-", "null:")
    cmd.Env = append(os.Environ(), "MAGICK_CONFIGURE_PATH=/home/h3rm1t0/Documentos/Testes/ImageMagick/config")

    cmd.Stdin = bytes.NewReader(corrompido)

    var bufferErros bytes.Buffer
    cmd.Stderr = &bufferErros

    err_pg := cmd.Run()

    if ctx.Err() == context.DeadlineExceeded {
        return
    }

    crashDetectado := false
    tipoCrash := ""

    if err_pg != nil {
        saidaErroStr := bufferErros.String()

        if strings.Contains(saidaErroStr, "AddressSanitizer:") {
            crashDetectado = true
            tipoCrash = "ASAN_LEAK"
        } else if exitErr, ok := err_pg.(*exec.ExitError); ok {
            if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
                if status.Signaled() && status.Signal() == syscall.SIGSEGV {
                    crashDetectado = true
                    tipoCrash = "SIGSEGV"
                }
            }
        }

        if crashDetectado {
            string_aleatoria := GerarStringRandomica(5)
            hash := md5.Sum([]byte(string_aleatoria))
            nome_arq_final := fmt.Sprintf("%x.ini", hash)
            dir := filepath.Join("Crashs", nome_arq_final)
            err := os.WriteFile(dir, corrompido, 0644)
            if err == nil {
                fmt.Printf("%s - [!! %s !!] Payload fatal salvo em: %s\n", momento(), tipoCrash, dir)
            }
            telegram(mensagem, tipoCrash, dir)
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
