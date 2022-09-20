package structs

type Order struct {
	Order_Id string `json:"order_id"`
	Table_Id int `json:"table_id"`
	Items []int `json:"items"`
	Priority int `json:"priority"`
	Max_Wait int `json:"max_wait"`
	Pickup_Time int `json:"pickup_time"`
	Waiter_Id int `json:"waiter_id"`
}