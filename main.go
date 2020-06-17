package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/Meexe/videocall/app"
	"github.com/Meexe/videocall/controllers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {

	users := &controllers.Connections{sync.Mutex{}, make(map[string]*websocket.Conn)}

	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", controllers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/login", controllers.LoginUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/ws/online", users.GetOnlineUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/echo", controllers.Echo).Methods("GET", "OPTIONS")
	router.Use(app.Authentication)

	// router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("http_port")
	if port == "" {
		port = "8000"
	}
	fmt.Printf("Server listening on port %s\n", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
