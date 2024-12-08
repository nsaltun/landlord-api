package httpwrap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nsaltun/landlord-api/pkg/middlewares"
	"github.com/spf13/viper"
)

func NewHttpServer(mux *http.ServeMux) *http.Server {
	vi := viper.New()
	vi.AutomaticEnv()
	vi.SetDefault("READ_TIMEOUT_IN_SECONDS", time.Second*10)
	vi.SetDefault("WRITE_TIMEOUT_IN_SECONDS", time.Second*10)
	vi.SetDefault("IDLE_TIMEOUT_IN_SECONDS", time.Second*10)
	vi.SetDefault("HOST_ADDRESS", "localhost")
	vi.SetDefault("PORT", 8080)
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", vi.GetString("HOST_ADDRESS"), vi.GetInt("PORT")),
		Handler:      http.Handler(middlewares.CorsMiddleware(mux)),
		ReadTimeout:  time.Duration(vi.GetInt("READ_TIMEOUT_IN_SECONDS")) * time.Second,
		WriteTimeout: time.Duration(vi.GetInt("WRITE_TIMEOUT_IN_SECONDS")) * time.Second,
		IdleTimeout:  time.Duration(vi.GetInt("IDLE_TIMEOUT_IN_SECONDS")) * time.Second,
	}
}
