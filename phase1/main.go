package main

// 	"context"
// 	"encoding/json"
// 	"log"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
import (
	"fmt"
	"log"
	"net/http"
	"phase1/config"
	"phase1/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Starting the application...")

	config.ConnectDB("mongo-container")
	defer config.Client.Disconnect(config.Context)
	o := mux.NewRouter()

	o.HandleFunc("/signup", handler.SignUpController).Methods("POST")
	o.HandleFunc("/signin", handler.SignInController).Methods("POST")
	o.HandleFunc("/insertslot/{id}", handler.SetScheduleHandler).Methods("POST")
	o.HandleFunc("/getslots", handler.GetDoctors).Methods("GET")
	o.HandleFunc("/getdrslot/{id}/slots", handler.GetDoctorSlotsByName).Methods("GET")
	o.HandleFunc("/reserveslot/{slotID}/{patientID}", handler.ReserveSlot).Methods("POST")
	o.HandleFunc("/slotid/{date}/{time}/{doctorName}", handler.SlotIDHandler).Methods("GET")
	o.HandleFunc("/viewslots/{patientID}", handler.GetPatientReservations).Methods("GET")
	o.HandleFunc("/delete/{id}", handler.DeleteSlot).Methods("PUT")
	o.HandleFunc("/update-slot/{oldID}/{newID}/{patientID}", handler.UpdateReservedSlot).Methods("PUT")

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(o)))
}
