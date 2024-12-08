package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsaltun/landlord-api/internal/handler"
	"github.com/nsaltun/landlord-api/pkg/httpwrap"
	"github.com/nsaltun/landlord-api/pkg/logging"
	"github.com/nsaltun/landlord-api/pkg/middlewares"
)

func main() {
	logging.InitSlog()

	httpHandler := handler.NewLandCalculationHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /land/calculate", middlewares.MiddlewareRunner(httpHandler.HandleFieldCalculation, middlewares.LoggingMiddleware))
	server := httpwrap.NewHttpServer(mux)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server stopped! %v", err)
		}
	}()

	log.Println("server is running on ", server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("server stopped gracefully")
}
