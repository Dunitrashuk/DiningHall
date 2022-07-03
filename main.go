package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Dunitrashuk/DiningHall/config"
	"github.com/Dunitrashuk/DiningHall/structs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var tables []structs.Table

func getHall(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hall Server is Listening on port 8082")
	fmt.Fprintf(w, "Hall Server is Listening on port 8082")
}

func getDish(w http.ResponseWriter, r *http.Request) {
	var dish models.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Dish %d received. Name: %s\n", dish.Dish_id, dish.Name)
	//fmt.Println("Dishes:", ordersDone)
}

func sendDishes() {
	time.Sleep(2 * time.Second)
	for i := 1; i <= 10; i++ {
		sendDish(i)
		time.Sleep(1 * time.Second)
	}
}

func sendDish(index int) {
	data := config.GetDish(index)
	jsonData, errMarshall := json.Marshal(data)
	if errMarshall != nil {
		log.Fatal(errMarshall)
	}
	resp, err := http.Post("http://"+config.GetKitchenAddr()+"/order", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Dish %d sent to kitchen. Status: %d\n", data.Dish_id, resp.StatusCode)
}

func hallServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", getHall).Methods("GET")
	myRouter.HandleFunc("/distribution", getDish).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+config.GetHallPort(), myRouter))
}

//function to create tables
func createTables() {
	for i := 0; i < config.NrOfTables(); i++ {
		table := structs.Table{
			i,
			"free",
			0,
		}
		tables = append(tables, table)
	}
}

func main() {
	go sendDishes()
	hallServer()
}

