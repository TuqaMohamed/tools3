package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"phase1/config"

	//"phase1/helper"
	"phase1/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name" bson:"name"`
}

// GetDoctors returns a list of doctors' names
func GetDoctors(w http.ResponseWriter, r *http.Request) {
	collection := config.Client.Database("clinic").Collection("doctor")
	//log.Print(collection)
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve doctors ", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var doctorInfoList []DoctorInfo
	for cursor.Next(r.Context()) {
		var doctor models.User
		if err := cursor.Decode(&doctor); err != nil {
			http.Error(w, "Failed to decode doctor data", http.StatusInternalServerError)
			return
		}
		doctorInfo := DoctorInfo{
			ID:   doctor.ID.Hex(),
			Name: doctor.Name,
		}
		doctorInfoList = append(doctorInfoList, doctorInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctorInfoList)
}

func GetDoctorSlotsByName(w http.ResponseWriter, r *http.Request) {

	doctorsCollection := config.Client.Database("clinic").Collection("doctor")
	slotsCollection := config.Client.Database("clinic").Collection("slots")
	//filter := bson.D{{Key: "name", Value: doctorsCollection.Name()}}

	// Get the doctor's name from t
	doctorID := mux.Vars(r)["id"]
	objectID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		http.Error(w, "Invalid ObjectID", http.StatusBadRequest)
		return
	}
	//log.Printf(mux.Vars(r)["name"])

	// Find the doctor by name
	var doctor models.User
	err = doctorsCollection.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&doctor)
	log.Printf(doctor.Name)
	if err != nil {
		http.Error(w, "Doctor not found", http.StatusNotFound)
		return
	}

	// Find the slots for the doctor by ID
	//var doctorSlots []models.Slot
	cursor, err := slotsCollection.Find(r.Context(), bson.M{"doctorName": doctor.Name, "isBooked": false})
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
			"id": slot.ID.Hex(),
		})
	}
	//log.Printf(formattedSlots[string][string])

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the slots as JSON and write it to the response
	json.NewEncoder(w).Encode(formattedSlots)
}

func ReserveSlot(w http.ResponseWriter, r *http.Request) {
	// Extract slot ID and patient ID from URL parameters
	slotID := mux.Vars(r)["slotID"]
	patientID := mux.Vars(r)["patientID"]
	log.Printf(slotID)
	log.Printf(patientID)

	slotsCollection := config.Client.Database("clinic").Collection("slots")
	patientsCollection := config.Client.Database("clinic").Collection("patient")

	// Find the slot by ID
	objectID, err := primitive.ObjectIDFromHex(slotID)
	if err != nil {
		http.Error(w, "Invalid Slot ID", http.StatusBadRequest)
		return
	}

	var slot models.Slot
	err = slotsCollection.FindOne(r.Context(), bson.M{"_id": objectID, "isBooked": false}).Decode(&slot)
	if err != nil {
		http.Error(w, "Slot not found or already reserved", http.StatusNotFound)
		return
	}

	// Find the patient by ID
	var patient models.User
	objectID1, err := primitive.ObjectIDFromHex(patientID)
	if err != nil {
		http.Error(w, "Invalid Slot ID", http.StatusBadRequest)
		return
	}
	err = patientsCollection.FindOne(r.Context(), bson.M{"_id": objectID1}).Decode(&patient)
	if err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	// Update the slot to mark it as reserved and add the patient information
	var updatedSlot models.Slot
	err = slotsCollection.FindOneAndUpdate(
		r.Context(),
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"isBooked":    true,
				"patientID":   patient.ID,
				"patientName": patient.Name,
			},
		},
	).Decode(&updatedSlot)

	if err != nil {
		http.Error(w, "Failed to reserve slot or update patient information", http.StatusInternalServerError)
		log.Printf("Error updating slot: %v", err)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Slot reserved successfully"})
}
func SlotIDHandler(w http.ResponseWriter, r *http.Request) {
	doctorName := mux.Vars(r)["doctorName"]
	date := mux.Vars(r)["date"]
	time := mux.Vars(r)["time"]

	collection := config.Client.Database("clinic").Collection("slots")

	filter := bson.M{
		"date":       date,
		"time":       time,
		"doctorName": doctorName,
	}

	var result models.Slot
	err := collection.FindOne(r.Context(), filter).Decode(&result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Respond with the slot ID as JSON
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"slotID": result.ID.Hex()}
	json.NewEncoder(w).Encode(response)
}
