package handlers

import (
	"encoding/json"

	"net/http"
	"phase1/config"
	"phase1/helper"

	"phase1/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateReservedSlot updates a reserved slot for a patient
func UpdateReservedSlot(w http.ResponseWriter, r *http.Request) {
	slotsCollection := config.Client.Database("clinic").Collection("slots")
	email := helper.GetUserFromSession(r)
	// Get the slot ID from the URL parameters
	slotID := mux.Vars(r)["id"]

	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		http.Error(w, "Invalid slot ID", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the updated slot information
	var updatedSlot models.Slot
	if err := json.NewDecoder(r.Body).Decode(&updatedSlot); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Find the slot by ID and check if it is reserved by the patient
	var existingSlot models.Slot
	err = slotsCollection.FindOne(
		r.Context(),
		bson.M{"_id": objID, "patientEmail": email},
	).Decode(&existingSlot)
	if err != nil {
		http.Error(w, "Slot not found or not reserved by the patient", http.StatusNotFound)
		return
	}

	// Check if updating the doctor only
	if updatedSlot.DoctorName != "" && updatedSlot.Time == "" && updatedSlot.Date == "" {
		// Find another available slot with the same date and time
		var availableSlot models.Slot
		err = slotsCollection.FindOne(
			r.Context(),
			bson.M{
				"doctorName": updatedSlot.DoctorName,
				"time":       existingSlot.Time,
				"date":       existingSlot.Date,
				"isBooked":   false,
			},
		).Decode(&availableSlot)

		if err != nil {
			http.Error(w, "No available slot with the same date and time found", http.StatusNotFound)
			return
		}

		// Reserve the new slot with the updated doctor information
		_, err = slotsCollection.UpdateOne(
			r.Context(),
			bson.M{"_id": availableSlot.ID},
			bson.M{"$set": bson.M{"isBooked": true, "patientEmail": email}},
		)
		if err != nil {
			http.Error(w, "Failed to update doctor", http.StatusInternalServerError)
			return
		}

		// Mark the old slot as available
		_, err = slotsCollection.UpdateOne(
			r.Context(),
			bson.M{"_id": objID},
			bson.M{"$set": bson.M{"isBooked": false, "patientEmail": ""}},
		)
		if err != nil {
			http.Error(w, "Failed to mark old slot as available", http.StatusInternalServerError)
			return
		}
	} else if updatedSlot.DoctorName != "" && updatedSlot.Time != "" && updatedSlot.Date != "" {
		var updatedSlot models.Slot
		err = slotsCollection.FindOneAndUpdate(
			r.Context(),
			bson.M{"_id": objID},
			bson.M{
				"$set": bson.M{"isBooked": false, "patientEmail": ""},
			},
		).Decode(&updatedSlot)

		if err != nil {
			http.Error(w, "Failed to reserve slot or update patient name", http.StatusInternalServerError)
			return
		}
		var doctor models.User
		doctorsCollection := config.Client.Database("clinic").Collection("doctor")
		ero := doctorsCollection.FindOne(r.Context(), bson.M{"name": updatedSlot.DoctorName}).Decode(&doctor)
		if ero != nil {
			http.Error(w, "Doctor not found", http.StatusNotFound)
			return
		}

		// Find the slot for the specified date, time, and doctor ID
		var slot models.Slot
		err = slotsCollection.FindOne(r.Context(), bson.M{"date": updatedSlot.Date, "isBooked": false, "time": updatedSlot.Time, "doctorName": updatedSlot.DoctorName}).Decode(&slot)
		if err != nil {
			http.Error(w, "Slot not found or already reserved", http.StatusNotFound)
			return
		}

		// Update the slot to mark it as reserved and add the patient name
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

	// If none of the above conditions are met, it's an invalid update
	http.Error(w, "Invalid update fields", http.StatusBadRequest)
}
