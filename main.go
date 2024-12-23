package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// QueryWhoisServer realiza uma consulta WHOIS para um domínio em um servidor específico
func QueryWhoisServer(server, domain string) (string, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:43", server), 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to WHOIS server %s: %v", server, err)
	}
	defer conn.Close()

	// Envia o domínio para o servidor WHOIS
	fmt.Fprintf(conn, "%s\r\n", domain)

	// Lê a resposta
	var response strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		response.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from WHOIS server %s: %v", server, err)
	}

	return response.String(), nil
}

// GetWhoisServerForDomain verifica o servidor WHOIS correto para um domínio
func GetWhoisServerForDomain(domain string) (string, error) {
	response, err := QueryWhoisServer("whois.iana.org", domain)
	if err != nil {
		return "", fmt.Errorf("failed to query IANA WHOIS: %v", err)
	}

	// Procura pela linha 'refer:' para obter o servidor WHOIS correto
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "refer:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "refer:")), nil
		}
	}

	return "", fmt.Errorf("could not find refer server for domain: %s", domain)
}

// GetDomainWhoisData realiza a consulta WHOIS completa
func GetDomainWhoisData(domain string) (string, error) {
	// Busca o servidor correto
	whoisServer, err := GetWhoisServerForDomain(domain)

	fmt.Printf("Servidor utilizado: %s", whoisServer)
	
	if err != nil {
		return "", fmt.Errorf("error finding WHOIS server: %v", err)
	}

	fmt.Printf("Using WHOIS server: %s\n", whoisServer)

	// Consulta o servidor WHOIS correto
	whoisData, err := QueryWhoisServer(whoisServer, domain)
	if err != nil {
		return "", fmt.Errorf("failed to query WHOIS server %s: %v", whoisServer, err)
	}

	return whoisData, nil
}

func main() {
	domain := "mcdonalds.com.br"
	whoisData, err := GetDomainWhoisData(domain)
	if err != nil {
		fmt.Println("Erro ao buscar dados WHOIS:", err)
		return
	}

	fmt.Println("Dados WHOIS para", domain)
	fmt.Println(whoisData)
}