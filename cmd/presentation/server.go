package presentation

import (
	"fmt"
	"net/http"
	"os"
)

func StartHttpServer() {

	mux := http.NewServeMux()

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		fmt.Println("Server crashed for some reason")
		os.Exit(1)
	}

}
