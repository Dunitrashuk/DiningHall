package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Dunitrashuk/DiningHall/config"
	"github.com/Dunitrashuk/DiningHall/structs"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex
var tables []structs.Table

func getHall(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hall Server is Listening on port 8082")
	fmt.Fprintf(w, "Hall Server is Listening on port 8082")
}

func getDish(w http.ResponseWriter, r *http.Request) {
	var dish structs.Dish
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
	for i := 1; i <= 5; i++ {
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

func occupy(table int) {
	for {
		// wait 2-3 min to occupy after table became free
		if tables[table].State == "free" {
			time.Sleep(time.Duration(rand.Intn(1000)+2000) * time.Millisecond)
			tables[table].State = "WO"
			fmt.Printf("Table %d: %s\n", table, tables[table].State)
		}
	}
}

func occupyTables() {
	for i := 0; i < config.NrOfTables(); i++ {
		// wait for about 1 min to start occupation
		time.Sleep(time.Duration(rand.Intn(300)+800) * time.Millisecond)
		go occupy(i)
	}
}

func printTables() {
	time.Sleep(time.Duration(rand.Intn(100)+4900) * time.Millisecond)
	fmt.Println()
	for i := 0; i < config.NrOfTables(); i++ {
		fmt.Printf("Table %d: %s\n", i, tables[i].State)
	}
}

func main() {
	go sendDishes()
	// fmt.Println(config.NrOfTables())
	// createTables()
	// occupyTables()
	// printTables()
	hallServer()
}
