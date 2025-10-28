package presentation

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
)

func StartHttpServer(cfg *config.Config, gatewayService application.ProcessorService) {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mutation", AuthMiddleware(NewMutationHandler(gatewayService), cfg))
	mux.HandleFunc("POST /generateJWT", Authenticate)

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		log.Println("Server stopped")
		os.Exit(1)
	}

}
