package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"./controllers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// Get a UserController instance
	oc := controllers.NewOrderController(getSession())

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

}

func getSession() (s *mgo.Session) {
	// Connect to local mongodb
	s, _ = mgo.Dial("mongodb://localhost")
	return s
}

/*
func YourHandler(w http.ResponseWriter, r *http.Request) {
w.Write([]byte("Gorilla!\n"))
}

func main() {
r := mux.NewRouter()
// Routes consist of a path and a handler function.
r.HandleFunc("/", YourHandler)

// Bind to a port and pass our router in
log.Fatal(http.ListenAndServe(":8000", r))
}
*/
