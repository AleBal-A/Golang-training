// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"candy-server/restapi/handlers"
	"candy-server/restapi/operations"
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"net/http"
)

//go:generate swagger generate server --target ../../candy-server --name CandyServer --spec ../api.yaml --principal interface{}

func configureFlags(api *operations.CandyServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CandyServerAPI) http.Handler {
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Используем пользовательский обработчик
	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(handlers.BuyCandyHandler)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts
func configureTLS(tlsConfig *tls.Config) {
	handlers.ConfigureTLS(tlsConfig)
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

//func BuyCandy(params operations.BuyCandyParams) (ret middleware.Responder) {
//	return
//}
