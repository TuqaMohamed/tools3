package handlers

import (
	"phase1/config"
	"phase1/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPatientReservations(w http.ResponseWriter, r *http.Request) {
	//client := config.Client.getmo
	slotsCollection := config.Client.Database("clinic").Collection("slots")

	// Get patient email from request parameters
	patientEmail := mux.Vars(r)["name"]

	// Query the slots collection for reservations with the patient's email
	cursor, err := slotsCollection.Find(r.Context(), bson.M{"patientEmail": patientEmail})
	if err != nil {
		http.Error(w, "Failed to retrieve patient reservations", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	var patientReservations []models.Slot
	if err := cursor.All(r.Context(), &patientReservations); err != nil {
		http.Error(w, "Error decoding patient reservations", http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the reservations as JSON and write it to the response
	json.NewEncoder(w).Encode(patientReservations)
}
