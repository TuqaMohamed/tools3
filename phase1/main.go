package main

// 	"context"
// 	"encoding/json"
// 	"log"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
import (
	"phase1/config"
	"phase1/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Starting the application...")

	config.ConnectDB()
	defer config.Client.Disconnect(config.Context)
	o := mux.NewRouter()
	// o.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
	// 	email := r.FormValue("email")
	// 	helper.StoreUserInSession(w, r, email)
	// })

	// Routes for doctors
	o.HandleFunc("/signup", handlers.SignUpController).Methods("POST")

	o.HandleFunc("/signin", handlers.SignInController).Methods("POST")
	o.HandleFunc("/insertslot", handlers.SetScheduleHandler).Methods("POST")
	o.HandleFunc("/getslots", handlers.GetDoctors).Methods("GET")
	o.HandleFunc("/getdrslot/{name}/slots", handlers.GetDoctorSlotsByName).Methods("GET")
	o.HandleFunc("/reserveslot/{name}/slot", handlers.ReserveSlot).Methods("POST")
	o.HandleFunc("/viewslots/{name}", handlers.GetPatientReservations).Methods("GET")
	o.HandleFunc("/delete/{id}", handlers.DeleteSlot).Methods("PUT")
	o.HandleFunc("/updateslot/{id}", handlers.UpdateReservedSlot).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", o))
}