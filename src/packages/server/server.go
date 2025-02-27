package server

import (
	"fmt"
	"net/http"
	"ubiquitous-eye/packages/controllers"
)

func RunServer() {
	mux := http.NewServeMux()
	configureEndpoints(mux)
	fmt.Printf("Listening on port 8080...\n")
	http.ListenAndServe(":8080", mux)
}

func configureEndpoints(mux *http.ServeMux) {
	fmt.Printf("Registering endpoints\n")

	mux.HandleFunc(controllers.ROOT, controllers.RootController)

	fmt.Printf("Endpoints registered\n")
}
