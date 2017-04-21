package models

type Order struct {
	OrderId  string      `json:"id" bson:"_id"`
	Location string      `json:"Location" bson:"Location"`
	Items    []Item      `json:"Items" bson:"Items"`
	Status   OrderStatus `json:"Status" bson:"Status"`
}

type Item struct {
	Name     string `json:"Name" bson:"Name"`
	Milk     string `json:"Milk" bson:"Milk"`
	Size     string `json:"Size" bson:"Size"`
	Quantity int    `json:"Quantity" bson:"Quantity"`
}

//Order Status
type OrderStatus string

var (
	PLACED    OrderStatus = "PLACED"
	PAID      OrderStatus = "PAID"
	PREPARING OrderStatus = "PREPARING"
	SERVED    OrderStatus = "SERVED"
	COLLECTED OrderStatus = "COLLECTED"
)
