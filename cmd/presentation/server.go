package presentation

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
)

func StartHttpServer(cfg *config.Config, gatewayService application.ProcessorService) {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mutation", AuthMiddleware(NewMutationHandler(gatewayService), cfg))

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		fmt.Println("Server crashed for some reason")
		os.Exit(1)
	}

}
