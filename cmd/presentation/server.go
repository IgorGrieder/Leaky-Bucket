package presentation

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/IgorGrieder/Leaky-Bucket/cmd/docs"
	"github.com/IgorGrieder/Leaky-Bucket/internal/application"
	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/swaggo/http-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a leaky bucket project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host leakyBucket
// @BasePath /v2
func StartHttpServer(cfg *config.Config, gatewayService application.ProcessorService, authService application.AuthService) {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mutation", AuthMiddleware(NewMutationHandler(gatewayService), authService, cfg))
	mux.HandleFunc("POST /generateJWT", Authenticate(authService))
	mux.HandleFunc("GET /swagger", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		log.Println("Server stopped")
		os.Exit(1)
	}

}
