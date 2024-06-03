package handlers

import (
	"candy-server/csrc"
	"candy-server/restapi/operations"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"os"
)

// Prices for candies
var Prices = map[string]int64{
	"CE": 10,
	"AA": 25,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}

// BuyCandyHandler is the custom handler for buying candy
func BuyCandyHandler(params operations.BuyCandyParams) middleware.Responder {
	order := params.Order
	price, validCandyType := Prices[*order.CandyType]
	cCount := order.CandyCount

	if !validCandyType {
		return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{
			Error: "Invalid candyType",
		})
	}

	if *cCount <= 0 {
		return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{
			Error: "Invalid candyCount",
		})
	}

	if *order.Money <= 0 {
		return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{
			Error: "Invalid money amount",
		})
	}

	totalCost := price * *cCount
	if totalCost > *order.Money {
		return operations.NewBuyCandyPaymentRequired().WithPayload(&operations.BuyCandyPaymentRequiredBody{
			Error: fmt.Sprintf("You need %d more money", totalCost-*order.Money),
		})
	}

	change := *order.Money - totalCost
	thanksMessage := csrc.AskCow("Thank you!") // Использование функции AskCow для генерации сообщения

	return operations.NewBuyCandyCreated().WithPayload(&operations.BuyCandyCreatedBody{
		Change: change,
		Thanks: thanksMessage,
	})
}

// ConfigureTLS configures the TLS settings
func ConfigureTLS(tlsConfig *tls.Config) {
	fmt.Println(os.Getwd())
	certFile := "../candy.tld/cert.pem"
	keyFile := "../candy.tld/key.pem"
	caFile := "../minica.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate and key: %v", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.ClientCAs = caCertPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
}
