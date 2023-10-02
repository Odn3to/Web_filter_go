package squid

import (
	"web-filter/models"

	"bufio"
	"os"
	"fmt"
	"strings"
	"encoding/json"
	"os/exec"
)

func removerPrefixoWWW(url string) string {
	// Verifique se a string começa com "www."
	if strings.HasPrefix(url, "www") {
		// Remova o prefixo "www." e retorne a string resultante
		return strings.Replace(url, "www", "", 1)
	}
	// Se não houver prefixo "www.", retorne a string original
	return url
}

func GetURL(wf models.WebFilter) (string, error) {
    var data map[string]string
    err := json.Unmarshal([]byte(wf.Data), &data)
    if err != nil {
        return "", err
    }
    url, found := data["url"]
    if !found {
        return "", nil
    }
    return url, nil
}

func ConfiguradorSquid(webFilters []models.WebFilter) {
    nomeArquivo := "/usr/local/squid/etc/webfilter/webfilter.conf"

    // Abre o arquivo no modo de escrita, cria se não existir
    arquivo, err := os.Create(nomeArquivo)
    if err != nil {
        fmt.Println("Erro ao abrir o arquivo:", err)
        return
    }
    defer arquivo.Close()

    // Cria um escritor bufio para escrever no arquivo
    escritor := bufio.NewWriter(arquivo)

    // Escreve no arquivo
    for i, webFilter := range webFilters {
        url, err := GetURL(webFilter)
        if err != nil {
            fmt.Println("Erro ao obter a URL:", err)
            return
        }
		dominio := removerPrefixoWWW(url)

		if(i == 0){
			_, err = escritor.WriteString("# ACL site bloqueio\nacl bloquear_sites dstdomain ")
			if err != nil {
				fmt.Println("Erro ao escrever no arquivo:", err)
				return
			}
		}

        _, err = escritor.WriteString(dominio + " ")
        if err != nil {
            fmt.Println("Erro ao escrever no arquivo:", err)
            return
        }
    }

	if(len(webFilters) > 0){
		_, err = escritor.WriteString("\n# Negue aos sites listados na ACL\nhttp_access deny bloquear_sites\n")
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}
	}

    // Certifique-se de que todos os dados sejam escritos no arquivo
    err = escritor.Flush()
    if err != nil {
        fmt.Println("Erro ao fazer flush no arquivo:", err)
        return
    }

    fmt.Println("Dados escritos com sucesso no arquivo.")

	//restart squid
	cmd := exec.Command("systemctl", "restart", "squid")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erro ao reiniciar o serviço do Squid:", err)
		return
	}
	fmt.Println("Saída do comando:", string(output))

}