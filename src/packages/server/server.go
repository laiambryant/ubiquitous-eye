package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	"ubiquitous-eye/packages/controllers"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func RunServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      getHttpOtelHandler(),
	}
	srvErr := srv.ListenAndServe()
	if srvErr != nil && srvErr != http.ErrServerClosed {
		log.Fatal(srvErr)
	}
}

func getHttpOtelHandler() http.Handler {
	mux := http.NewServeMux()
	handlerSetup(mux)
	return getHandler(mux)
}

func handlerSetup(mux *http.ServeMux) {
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}
	handleFunc(controllers.ROOT, controllers.RootController)
}

func getHandler(mux *http.ServeMux) http.Handler {
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}
