package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/Meexe/videocall/app"
	"github.com/Meexe/videocall/controllers"
	"net/http"
	"os"
)

func main() {

	httpRouter := mux.NewRouter()
	httpRouter.HandleFunc("/api/user/new", controllers.CreateUser).Methods("POST")
	httpRouter.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	httpRouter.Use(app.JwtAuthentication) //attach JWT auth middleware

	// router.NotFoundHandler = app.NotFoundHandler

	httpPort := os.Getenv("http_port")
	if httpPort == "" {
		httpPort = "8000"
	}
	fmt.Printf("HTTP listening on port %s\n", httpPort)

	err := http.ListenAndServe(":" + httpPort, httpRouter)
	if err != nil {
		fmt.Print(err)
	}
}