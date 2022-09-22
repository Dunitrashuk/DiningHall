package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Dunitrashuk/DiningHall/config"
	"github.com/Dunitrashuk/DiningHall/structs"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex
var tables []structs.Table
var orderList []structs.Order

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

func sendOrder(order structs.Order) {
	data := order
	jsonData, errMarshall := json.Marshal(data)
	if errMarshall != nil {
		log.Fatal(errMarshall)
	}
	resp, err := http.Post("http://"+config.GetKitchenAddr()+"/order", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order %d sent to kitchen. Status: %d\n", data.Order_Id, resp.StatusCode)
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
			Id:    i,
			State: "free",
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
		fmt.Printf("Table %d: %s\n", i, tables[i].State)
		// wait for about 1 min to start occupation
		time.Sleep(time.Duration(rand.Intn(300)+800) * time.Millisecond)
		go occupy(i)
	}
}

func printTables() {
	time.Sleep(time.Duration(rand.Intn(100)+4900) * time.Millisecond)
	fmt.Println()
	for i := 0; i < config.NrOfTables(); i++ {
		fmt.Printf("%+v\n", tables[i])
	}
}

func generateOrder(waiterId int, tableId int) structs.Order{
	var items []int
	maxWait := 0

	 for i := 0; i < rand.Intn(5) + 1; i++ {
		 items = append(items, config.GetDish(rand.Intn(9) + 1).Dish_id)
	 }

	for _, dishId := range items {
		preparationTime := config.GetDish(dishId - 1).Preparation_time
		if maxWait < preparationTime {
			maxWait = preparationTime
		}
	}

	order := structs.Order{
		Order_Id: uuid.New().String(),
		Table_Id: tableId,
		Items: items,
		Priority: rand.Intn(5) + 1,
		Max_Wait: int(float32(maxWait) * 1.3),
		Pickup_Time: int(time.Now().Unix()),
		Waiter_Id: waiterId,
	}
	return order
}

func createWaiters() {
	for i := 1; i <= config.NrOfWaiters() + 1; i++ {
		go waiter(i)
		fmt.Printf("Waiter %d created\n", i)
	}
}


func waiter(waiterId int) {
	for {
		time.Sleep(time.Duration(rand.Intn(100)+100) * time.Millisecond)
		mutex.Lock() //lock mutex in order to access the shared resource tables
		for i := 0; i < config.NrOfTables(); i++ {
			if tables[i].State == "WO" {
				order := generateOrder(waiterId, i)
				orderList = append(orderList, order)
				tables[i].OrderId = order.Order_Id
				fmt.Printf("%+v\n", order)
				sendOrder(order)
				tables[i].State = "WS"
			}
		}
		mutex.Unlock()
	}
}

func main() {
	//go sendDishes()
	createTables()
	createWaiters()
	occupyTables()
	printTables()
	hallServer()
}
