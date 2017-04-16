package models

type Order struct {
	OrderId  string  `json:"id" bson:"_id"`
	Location string  `json:"Location" bson:"Location"`
	Drinks   []Drink `json:"Drinks" bson:"Drinks"`
}

type Drink struct {
	Name     string `json:"Name" bson:"Name"`
	Milk     string `json:"Milk" bson:"Milk"`
	Size     string `json:"Size" bson:"Size"`
	Quantity int    `json:"Quantity" bson:"Quantity"`
}
