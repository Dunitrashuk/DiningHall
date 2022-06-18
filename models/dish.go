package models

type Dish struct {
	DishId int    `json:"id"`
	Name            string `json:"name"`
	PreparationTime int    `json:"preparation-time"`
	Complexity      int    `json:"complexity"`
	CookingApparatus string `json:"cooking-apparatus"`
}
