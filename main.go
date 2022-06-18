package main

import (
	"fmt"
	"github.com/Dunitrashuk/DiningHall/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getHall(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Hall")
	fmt.Fprintf(w, "Welcome to the Hall!")
}

func hallServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", getHall).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+config.GetHallPort(), myRouter))
}

func main() {
	hallServer();
}

