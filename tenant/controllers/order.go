package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

// OrderController represents the controller for operating on the Order resource
type OrderController struct {
	session *mgo.Session
}

// NewOrderController provides a reference to a OrderController with provided mongo session
func NewOrderController(s *mgo.Session) *OrderController {
	return &OrderController{s}
}

// CreateOrder creates a new Order
func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("inside createorder	")
	var o models.Order

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&o)

	// Add an Id, using uuid for
	o.OrderId = uuid.NewV4().String()

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
		return
	}

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(o)
}

//Delete Order deletes the order with specified order id
func (oc OrderController) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["id"]

	o := models.Order{}

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&o); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := oc.session.DB("test").C("Order").RemoveId(orderId); err != nil {
		fmt.Println("Could not find order - %s to delete", orderId)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

//Update Order updates the order with specified details
func (oc OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["id"]

	o := models.Order{}

	// Fetch order
	json.NewDecoder(r.Body).Decode(&o)

	if err := oc.session.DB("test").C("Order").UpdateId(orderId, &o); err != nil {
		w.WriteHeader(404)
		return
	}

	fmt.Println("Order Updated:", orderId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(o)

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

// GetOrders retrieves all the orders
func (oc OrderController) OrderPayment(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("inside createorder	")
	vars := mux.Vars(r)
	orderId := vars["id"]

	o := models.Order{}

	json.NewDecoder(r.Body).Decode(&o)

	if err := oc.session.DB("test").C("Order").UpdateId(orderId, &o); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write content-type, statuscode, payload
	fmt.Println("Order Status Updated: ", o.Status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(o)
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
