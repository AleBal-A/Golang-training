package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Параметры командной строки
	candyType := flag.String("k", "", "Type of candy (two-letter abbreviation)")
	candyCount := flag.Int("c", 0, "Count of candy to buy")
	money := flag.Int("m", 0, "Amount of money you gave to the machine")
	serverAddr := flag.String("addr", "https://candy.tld:3333", "Server address")
	caFile := "../minica.pem"
	clientCert := "../candy.tld/cert.pem"
	clientKey := "../candy.tld/key.pem"
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification")
	flag.Parse()

	if *candyType == "" || *candyCount <= 0 || *money <= 0 {
		log.Fatalf("Invalid input parameters")
	}

	// Загрузка клиентского сертификата и ключа
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Загрузка CA сертификата
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка TLS конфигурации
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: *insecure, // Пропуск проверки сертификата, если включено
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Формирование URL и тела запроса
	url := fmt.Sprintf("%s/buy_candy", *serverAddr)
	body := fmt.Sprintf(`{"money": %d, "candyType": "%s", "candyCount": %d}`, *money, *candyType, *candyCount)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Printf("%s\n", responseBody)
}
