package handler

import (
	"encoding/json"
	"net/http"
	"phase1/config"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseMessage is a struct to represent the response message
type ResponseMessage struct {
	Message string `json:"message"`
}

// DeleteSlot deletes a slot based on its ID
func DeleteSlot(w http.ResponseWriter, r *http.Request) {

	slotsCollection := config.Client.Database("clinic").Collection("slots")

	// Get the slot ID from the URL parameters
	slotID := mux.Vars(r)["id"]

	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		http.Error(w, "Invalid slot ID", http.StatusBadRequest)
		return
	}

	// Update the 'available' parameter to true
	result, err := slotsCollection.UpdateOne(
		r.Context(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"isBooked": false, "patientName": "", "patientID": primitive.NilObjectID}},
	)
	if err != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	// Check if the slot was found and updated
	if result.ModifiedCount == 0 {
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	// Respond with a success message as JSON
	responseMessage := ResponseMessage{
		Message: "Reservation deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseMessage)
}
