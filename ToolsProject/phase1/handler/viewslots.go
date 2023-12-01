package handler

import (
	"encoding/json"
	"net/http"
	"phase1/config"
	"phase1/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SlotInfo struct {
	Date       string             `json:"date"`
	Time       string             `json:"time"`
	DoctorName string             `json:"doctorName"`
	SlotID     primitive.ObjectID `json:"slotID"`
}

func GetPatientReservations(w http.ResponseWriter, r *http.Request) {
	slotsCollection := config.Client.Database("clinic").Collection("slots")

	patientID := mux.Vars(r)["patientID"]
	objectID1, err := primitive.ObjectIDFromHex(patientID)

	// Query the slots collection for reservations with the patient's id
	cursor, err := slotsCollection.Find(r.Context(), bson.M{"patientID": objectID1})
	if err != nil {
		http.Error(w, "Failed to retrieve patient reservations", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	// Create a slice to store the simplified slot information
	var patientReservations []SlotInfo

	// Iterate over the slots and populate the simplified structure
	for cursor.Next(r.Context()) {
		var slot models.Slot
		if err := cursor.Decode(&slot); err != nil {
			http.Error(w, "Error decoding patient reservations", http.StatusInternalServerError)
			return
		}

		// Append the simplified slot information to the result
		patientReservations = append(patientReservations, SlotInfo{
			Date:       slot.Date,
			Time:       slot.Time,
			DoctorName: slot.DoctorName,
			SlotID:     slot.ID,
		})
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the simplified reservations as JSON and write it to the response
	json.NewEncoder(w).Encode(patientReservations)
}
