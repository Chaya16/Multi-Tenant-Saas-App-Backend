package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
func NewOrderController(s *mgo.Session) *OrderController {
	return &OrderController{s}
}

// CreateOrder creates a new Order
func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Stub an order to be populated from the body
	//o := models.Order{}
	fmt.Println("inside createorder	")
	var o models.Order

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&o)

	// Add an Id
	//o.OrderId = bson.NewObjectId()
	o.OrderId = uuid.NewV4().String()

	// Write the user to mongo
	oc.session.DB("test").C("Order").Insert(&o)

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	fmt.Println("New Order Created, Order ID:", o.OrderId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	//fmt.Fprintf(w, "%s", uj)
	json.NewEncoder(w).Encode(o)
}

// GetOrder retrieves an individual order
func (oc OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {

	// Grab id
	vars := mux.Vars(r)
	orderId := vars["id"]

	// Verify id is ObjectId, otherwise bail
	/*if !bson.IsObjectIdHex(orderId) {
		w.WriteHeader(404)
		return
	}
	*/
	// Grab id
	//	oid := bson.ObjectIdHex(orderId)

	// Stub user
	o := models.Order{}

	// Fetch order
	if err := oc.session.DB("test").C("Order").FindId(orderId).One(&o); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	//uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(o)

	//fmt.Fprintf(w, "%s", uj)
}

//Delete Order deletes the order with specified order id
func (oc OrderController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["orderId"]

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(orderId) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(orderId)
	// Remove from database
	if err := oc.session.DB("test").C("Order").RemoveId(oid); err != nil {
		fmt.Println("Could not find order - %s to delete", orderId)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

//Update Order updates the order with specified deails that come in the json format
func (oc OrderController) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	orderId := vars["orderId"]

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(orderId) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(orderId)
	// Remove from database
	if err := oc.session.DB("test").C("Order").RemoveId(oid); err != nil {
		fmt.Println("Could not find order - %s to delete", orderId)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

//ping resource function
func (oc OrderController) PingOrderResource(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pinging Order Resource")
}
