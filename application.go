package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	mgo "gopkg.in/mgo.v2"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

type Order struct {
	OrderId  string `json:"id" bson:"_id"`
	Location string `json:"location" bson:"location"`
	Items    []Item `json:"items" bson:"items"`
	Status   string `json:"status" bson:"status"`
	Message  string `json:"message" bson:"message"`
	Links    Links  `json:"links" bson:"links"`
}

type Links struct {
	Payment string `json:"payment,omitempty"`
	Order   string `json:"order,omitempty"`
}

type Item struct {
	Name     string `json:"name" bson:"Name"`
	Milk     string `json:"milk" bson:"Milk"`
	Size     string `json:"size" bson:"Size"`
	Quantity int    `json:"qty" bson:"Quantity"`
}

// OrderController represents the controller for operating on the Order resource
type OrderController struct {
	session *mgo.Session
}

// NewOrderController provides a reference to a OrderController with provided mongo session
func NewOrderController(mgoSession *mgo.Session) *OrderController {
	return &OrderController{mgoSession}
}

// CreateOrder creates a new Order
func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("inside createorder	")
	var o Order

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&o)

	// Add an Id, using uuid for
	o.OrderId = uuid.NewV4().String()
	var links Links
	links.Payment = "http://localhost:8080/v1/starbucks/order/" + o.OrderId + "/pay"
	links.Order = "http://localhost:8080/v1/starbucks/order/" + o.OrderId

	o.Links = links
	o.Status = "PLACED"
	o.Message = "Order has been placed"
	// Write the user to mongo
	oc.session.DB("test").C("Order").Insert(&o)

	// Write content-type, statuscode, payload
	fmt.Println("New Order Created, Order ID:", o.OrderId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(o)
}

// GetOrder retrieves an individual order
func (oc OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {

	// Grab id
	vars := mux.Vars(r)
	orderId := vars["id"]

	o := Order{}

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&o); err != nil {
		w.WriteHeader(404)
		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(o)
}

// Delete Order deletes the order with specified order id
func (oc OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["id"]

	o := Order{}

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&o); err != nil {
		w.WriteHeader(404)
		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}

	//check for status and then delete
	if o.Status == "PAID" || o.Status == "PREPARING" || o.Status == "SERVED" || o.Status == "COLLECTED" {
		//fmt.Println("Order cannot be updated after payment has been made")
		data := `{"status":"error","message":"Order cannot be deleted after payment has been made"}`
		json.NewEncoder(w).Encode(data)
		return
		//http.Error(w, "Order cannot be updated after payment has been made", 400)
	}
	if err := oc.session.DB("test").C("Order").RemoveId(orderId); err != nil {
		fmt.Println("Could not find order - %s to delete", orderId)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
	data := `{"status":"success","message":"Order has been deleted"}`
	json.NewEncoder(w).Encode(data)
}

//Update Order updates the order with specified details
func (oc OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["id"]

	orderFromJson := Order{}
	orderFromDb := Order{}
	// Fetch order
	json.NewDecoder(r.Body).Decode(&orderFromJson)
	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&orderFromDb); err != nil {
		w.WriteHeader(404)
		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}

	if orderFromDb.Status == "PAID" || orderFromDb.Status == "PREPARING" || orderFromDb.Status == "SERVED" || orderFromDb.Status == "COLLECTED" {
		//fmt.Println("Order cannot be updated after payment has been made")
		w.WriteHeader(400)
		data := `{"status":"error","message":"Order cannot be updated after payment has been made "}`
		json.NewEncoder(w).Encode(data)
		return

	}

	if err := oc.session.DB("test").C("Order").UpdateId(orderId, &orderFromJson); err != nil {
		w.WriteHeader(404)
		return
	}

	fmt.Println("Order Updated:", orderId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(orderFromJson)

}

// GetOrders retrieves all the orders
func (oc OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {

	var orders []Order

	iter := oc.session.DB("test").C("Order").Find(nil).Iter()
	result := Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&orders)

}

// OrderPayment handles functionality after payment has been made for an order
func (oc OrderController) OrderPayment(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("inside createorder	")
	vars := mux.Vars(r)
	orderId := vars["id"]
	order := Order{}
	fmt.Println("pay order")
	fmt.Println(orderId)

	//call get order method
	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&order); err != nil {
		w.WriteHeader(404)
		return
	}

	json.NewDecoder(r.Body).Decode(&order)
	//order.Status :="PAID"
	fmt.Println(order)

	/*if order.Status == "PAID" || order.Status == "PREPARING" || order.Status == "SERVED" || order.Status == "COLLECTED" {
		w.WriteHeader(400)
		data := `{"status":"error","message":"Order payment rejected "}`
		json.NewEncoder(w).Encode(data)
		return
	}*/

	//code to update status to paid goes here

	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PAID", "message": "Payment Accepted"}})
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$unset": bson.M{"Links.Payment": ""}})
	fmt.Println("Order Status Updated: ", order.Status)
	runtime.GOMAXPROCS(1)
	time.Sleep(1000 * time.Millisecond)
	go changeStatusToPreparing(orderId, oc)
	time.Sleep(100 * time.Millisecond)
	go changeStatusToServing(orderId, oc)
	time.Sleep(100 * time.Millisecond)
	go changeStatusToCollected(orderId, oc)

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&order); err != nil {
		w.WriteHeader(404)
		data := `{"status":"error","message":"Order not found"}`
		json.NewEncoder(w).Encode(data)
		return
	}
	// to stop displaying payment after clicking on order pay(since payment set to omit empty)
	order.Links.Payment = ""
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(order)
}

func changeStatusToPreparing(orderId string, oc OrderController) {
	fmt.Println("preparing")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "PREPARING"}})

}

func changeStatusToServing(orderId string, oc OrderController) {
	fmt.Println("serving")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "SERVED"}})

}

func changeStatusToCollected(orderId string, oc OrderController) {
	fmt.Println("collected")
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"status": "COLLECTED"}})

}

//ping resource function
func (oc OrderController) PingOrderResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinging Order Resource")
}

func main() {

	r := mux.NewRouter()

	// Get a UserController instance
	oc := NewOrderController(getSession())

	r.HandleFunc("/v1/starbucks/order", oc.CreateOrder).Methods("POST")
	r.HandleFunc("/v1/starbucks/order/{id}", oc.GetOrder).Methods("GET")
	r.HandleFunc("/v1/starbucks/orders", oc.GetOrders).Methods("GET")
	r.HandleFunc("/v1/starbucks/order/{id}", oc.DeleteOrder).Methods("DELETE")
	r.HandleFunc("/v1/starbucks/order/{id}", oc.UpdateOrder).Methods("PUT")
	r.HandleFunc("/v1/starbucks/order/{id}/pay", oc.OrderPayment).Methods("POST")
	r.HandleFunc("/v1/starbucks/ping", oc.PingOrderResource)
	r.Handle("/", r)

	fmt.Println("serving on port 8080")
	http.ListenAndServe(":8080", r)
	//go changeDrinkStatus()
}

func getSession() (s *mgo.Session) {
	// Connect to local mongodb
	s, _ = mgo.Dial("mongodb://54.153.71.97")
	return s
}
