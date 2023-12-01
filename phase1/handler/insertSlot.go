package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"phase1/config"
	"phase1/helper"
	"phase1/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := helper.GetUserFromSession(r)
	fmt.Println(email)
	doctorID := mux.Vars(r)["id"]
	var input struct {
		Slot models.Slot `json:"slot"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Reached before user validation")
	objectID, err := primitive.ObjectIDFromHex(doctorID)
	doctor, err := validatedoctor(objectID)
	if err != nil {
		log.Println("Error validating doctor:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println("User validation completed")
	// if !userExists {
	// 	fmt.Println("doctor doesn't exist")
	// 	http.Error(w, "doctor not found", http.StatusNotFound)
	// 	return
	// }

	// Set the doctor's schedule by inserting the time slot into MongoDB
	slot := models.Slot{
		Date:     input.Slot.Date,
		Time:     input.Slot.Time,
		IsBooked: false,
	}

	// Insert the slot into the doctor_schedule collection
	err = insertSlot(doctor.Name, doctor.ID, input.Slot)
	if err != nil {
		log.Println("Error inserting slot:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":           true,
		"message":           "slot inserted successfully",
		"insertedSlot":      slot,
		"doctorInformation": doctor,
	}

	json.NewEncoder(w).Encode(response)
}

func validatedoctor(id primitive.ObjectID) (*models.User, error) {
	collection := config.Client.Database("clinic").Collection("doctor")
	filter := bson.D{{Key: "_id", Value: id}}

	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		// User not found
		return nil, nil
	} else if err != nil {
		// Other errors
		return nil, err
	}

	// User found
	return &existingUser, nil
}

func insertSlot(name string, id primitive.ObjectID, slot models.Slot) error {

	collection := config.Client.Database("clinic").Collection("slots")

	// Check if the slot already exists for the given date and hour
	existingSlot, err := findSlot(name, slot.Date, slot.Time)
	if err != nil {
		return err
	}
	if existingSlot != nil {
		return errors.New("slot already reserved") // Slot already exists, you may handle this case as needed
	}
	// Insert the slot into the collection
	_, err = collection.InsertOne(context.Background(), bson.M{
		"doctorName":   name,
		"doctorID":     id,
		"date":         slot.Date,
		"time":         slot.Time,
		"isBooked":     false,
		"patientEmail": "",
	})
	if err != nil {
		log.Println("Error inserting slot", err)
		return err
	}

	return err
}

func findSlot(name string, date string, time string) (*models.Slot, error) {
	collection := config.Client.Database("clinic").Collection("slots")
	filter := bson.D{
		{Key: "doctorName", Value: name},
		{Key: "date", Value: date},
		{Key: "time", Value: time},
	}

	var slot models.Slot
	err := collection.FindOne(context.Background(), filter).Decode(&slot)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		log.Println("Error checking slot avalability")
		return nil, err
	}
	return &slot, nil
}
