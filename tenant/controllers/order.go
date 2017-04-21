package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
	var o models.Order

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&o)

	// Add an Id, using uuid for
	o.OrderId = uuid.NewV4().String()
	var links models.Links
	links.Payment = "http://localhost:8080/order/" + o.OrderId + "/pay"
	links.Order = "http://localhost:8080/order/" + o.OrderId

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

	o := models.Order{}

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

	o := models.Order{}

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

	orderFromJson := models.Order{}
	orderFromDb := models.Order{}
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

	var orders []models.Order

	iter := oc.session.DB("test").C("Order").Find(nil).Iter()
	result := models.Order{}
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
	order := models.Order{}
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

	if order.Status == "PAID" || order.Status == "PREPARING" || order.Status == "SERVED" || order.Status == "COLLECTED" {
		w.WriteHeader(400)
		data := `{"status":"error","message":"Order payment rejected "}`
		json.NewEncoder(w).Encode(data)
		return
	}

	//code to update status to paid goes here

	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"Status": "PAID", "Message": "Payment Accepted"}})
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$unset": bson.M{"Links.Payment": ""}})

	//CHANGE THE ORDER PROCESSIGN STATUS
	fmt.Println("Order Status Updated: ", order.Status)
	time.Sleep(10000)
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"Status": "PREPARING"}})
	time.Sleep(10000)
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"Status": "SERVED"}})
	time.Sleep(10000)
	oc.session.DB("test").C("Order").UpdateId(orderId, bson.M{"$set": bson.M{"Status": "COLLECTED"}})

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

// DeleteOrders deletes all the orders
/*func (oc OrderController) DeleteOrders(w http.ResponseWriter, r *http.Request) {

	var orders []models.Order



	oc.session.DB("test").C("Order").Remove(models.Order{})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&orders)

}*/

//ping resource function
func (oc OrderController) PingOrderResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinging Order Resource")
}
