package models

type Order struct {
	OrderId  string `json:"id" bson:"_id"`
	Location string `json:"Location" bson:"Location"`
	Items    []Item `json:"Items" bson:"Items"`
	Status   string `json:"Status" bson:"Status"`
	Message  string `json:"Message" bson:"Message"`
	Links    Links  `json:"Links" bson:"Links"`
}

type Links struct {
	Payment string `json:"payment,omitempty"`
	Order   string `json:"order,omitempty"`
}

type Item struct {
	Name     string `json:"Name" bson:"Name"`
	Milk     string `json:"Milk" bson:"Milk"`
	Size     string `json:"Size" bson:"Size"`
	Quantity int    `json:"Quantity" bson:"Quantity"`
}
