package structs

type Table struct {
	Id             int
	State          string // free, WO (waiting to order), WS (waiting to be served)
	OrderId string
}