package main

import (
	"fmt"
	"log"
	"net/http"
	"project/userlogin/connection"
	"project/userlogin/library"
	"project/userlogin/login"

	"github.com/gorilla/mux"
)

func init() {
	//mongodb connection
	err := connection.ConnectDB()
	if err != nil {
		log.Fatalf("Error in client connection")
	}
}
func main() {
	//router creation
	server := mux.NewRouter().StrictSlash(true)
	server.HandleFunc("/register", login.Register).Methods("POST")
	server.HandleFunc("/login", login.Login).Methods("POST")
	server.Handle("/userinfo", library.AuthenticateMiddleware(http.HandlerFunc(login.UserInfo))).Methods("GET")
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", server)
}
