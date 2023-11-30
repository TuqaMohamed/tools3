package handlers

import (
	"phase1/config"
	"phase1/helper"
	"phase1/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// GetDoctors returns a list of doctors' names
func GetDoctors(w http.ResponseWriter, r *http.Request) {
	collection := config.Client.Database("clinic").Collection("doctor")
	//log.Print(collection)
	var doctors []models.User
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve doctors ", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(r.Context())
	for cursor.Next(r.Context()) {
		var doctor models.User
		if err := cursor.Decode(&doctor); err != nil {
			http.Error(w, "Failed to decode doctor data", http.StatusInternalServerError)
			return
		}
		doctors = append(doctors, doctor)
	}

	var names []string
	for _, doctor := range doctors {
		names = append(names, doctor.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

// GetDoctorSlotsByName returns the available slots for a specific doctor by name
func GetDoctorSlotsByName(w http.ResponseWriter, r *http.Request) {

	doctorsCollection := config.Client.Database("clinic").Collection("doctor")
	slotsCollection := config.Client.Database("clinic").Collection("slots")
	//filter := bson.D{{Key: "name", Value: doctorsCollection.Name()}}

	// Get the doctor's name from the request parameters
	doctorName := mux.Vars(r)["name"]
	//log.Printf(mux.Vars(r)["name"])

	// Find the doctor by name
	var doctor models.User
	err := doctorsCollection.FindOne(r.Context(), bson.M{"name": doctorName}).Decode(&doctor)
	log.Printf(doctor.Name)
	if err != nil {
		http.Error(w, "Doctor not found", http.StatusNotFound)
		return
	}

	cursor, err := slotsCollection.Find(r.Context(), bson.M{"doctorName": doctorName, "isBooked": false})
	if err != nil {
		http.Error(w, "Failed to retrieve doctor slots", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(r.Context())
	var formattedSlots []map[string]string
	for cursor.Next(r.Context()) {
		var slot models.Slot
		if err := cursor.Decode(&slot); err != nil {
			http.Error(w, "Error decoding slot data", http.StatusInternalServerError)
			return
		}
		formattedSlots = append(formattedSlots, map[string]string{
			"date": slot.Date,
			"time": slot.Time,
		})
	}
	//log.Printf(formattedSlots[string][string])

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the slots as JSON and write it to the response
	json.NewEncoder(w).Encode(formattedSlots)
}

func ReserveSlot(w http.ResponseWriter, r *http.Request) {

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	slotsCollection := config.Client.Database("clinic").Collection("slots")
	email := helper.GetUserFromSession(r)
	date := request["date"]
	time := request["time"]

	// Get the doctor's name, date, and time from the URL parameters
	doctorName := mux.Vars(r)["name"]
	//date := mux.Vars(r)["date"]
	//time := mux.Vars(r)["time"]

	// Find the doctor by name
	var doctor models.User
	doctorsCollection := config.Client.Database("clinic").Collection("doctor")
	ero := doctorsCollection.FindOne(r.Context(), bson.M{"name": doctorName}).Decode(&doctor)
	if ero != nil {
		http.Error(w, "Doctor not found", http.StatusNotFound)
		return
	}

	// Find the slot for the specified date, time, and doctor ID
	var slot models.Slot
	err = slotsCollection.FindOne(r.Context(), bson.M{"date": date, "isBooked": false, "time": time}).Decode(&slot)
	if err != nil {
		http.Error(w, "Slot not found or already reserved", http.StatusNotFound)
		return
	}

	// Update the slot to mark it as reserved and add the patient name
	var updatedSlot models.Slot
	err = slotsCollection.FindOneAndUpdate(
		r.Context(),
		bson.M{"_id": slot.ID},
		bson.M{
			"$set": bson.M{"isBooked": true, "patientEmail": email},
		},
	).Decode(&updatedSlot)

	if err != nil {
		http.Error(w, "Failed to reserve slot or update patient name", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Slot reserved successfully"})
}
