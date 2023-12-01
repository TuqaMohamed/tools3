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

// UpdateReservedSlot updates a reserved slot for a patient

func UpdateReservedSlot(w http.ResponseWriter, r *http.Request) {
    slotsCollection := config.Client.Database("clinic").Collection("slots")

    // Extract old slot ID, new slot ID, and patient ID from URL parameters
    oldSlotID := mux.Vars(r)["oldID"]
    newSlotID := mux.Vars(r)["newID"]
    patientID := mux.Vars(r)["patientID"]

    // Convert string IDs to primitive.ObjectID
    oldObjectID, err := primitive.ObjectIDFromHex(oldSlotID)
    if err != nil {
        http.Error(w, "Invalid Old Slot ID", http.StatusBadRequest)
        return
    }

    newObjectID, err := primitive.ObjectIDFromHex(newSlotID)
    if err != nil {
        http.Error(w, "Invalid New Slot ID", http.StatusBadRequest)
        return
    }

    patientObjectID, err := primitive.ObjectIDFromHex(patientID)
    if err != nil {
        http.Error(w, "Invalid Patient ID", http.StatusBadRequest)
        return
    }

    var existingSlot models.Slot
    err = slotsCollection.FindOne(r.Context(), bson.M{"_id": oldObjectID}).Decode(&existingSlot)
    if err != nil {
        http.Error(w, "Old Slot not found or reserved", http.StatusNotFound)
        return
    }

    var newSlot models.Slot
    err = slotsCollection.FindOne(r.Context(), bson.M{"_id": newObjectID}).Decode(&newSlot)
    if err != nil {
        http.Error(w, "New Slot not found or reserved", http.StatusNotFound)
        return
    }

    // Fetch patient information using patient ID
    patientsCollection := config.Client.Database("clinic").Collection("patient")
    var patient models.User
    err = patientsCollection.FindOne(r.Context(), bson.M{"_id": patientObjectID}).Decode(&patient)
    if err != nil {
        http.Error(w, "Patient not found", http.StatusNotFound)
        return
    }

    // if existingSlot.IsBooked {
    //     http.Error(w, "Cannot edit an unreserved appointment", http.StatusBadRequest)
    //     return
    // } else if newSlot.IsBooked {
    //     http.Error(w, "Cannot book a reserved slot", http.StatusBadRequest)
    //     return
    // }

    // Clear old slot
    existingSlot.IsBooked = false
    existingSlot.PatientID = primitive.NilObjectID
    existingSlot.PatientName = ""

    _, err = slotsCollection.ReplaceOne(r.Context(), bson.M{"_id": oldObjectID}, existingSlot)
    if err != nil {
        http.Error(w, "Failed to update existing slot", http.StatusInternalServerError)
        return
    }

    // Book new slot
    newSlot.IsBooked = true
    newSlot.PatientID = patientObjectID
    newSlot.PatientName = patient.Name
    // You can update other fields as needed

    _, err = slotsCollection.ReplaceOne(r.Context(), bson.M{"_id": newObjectID}, newSlot)
    if err != nil {
        http.Error(w, "Failed to update new slot", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Slot updated successfully"})
}
