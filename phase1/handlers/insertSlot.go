package handlers

import (
	"phase1/config"
	"phase1/helper"
	"phase1/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := helper.GetUserFromSession(r)
	//log.Println(email)
	var input struct {
		Slot models.Slot `json:"slot"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Find the doctor in the collection

	userExists, doctor, err := validatedoctor("doctor", email)
	if err != nil {
		log.Fatal(err)
	}

	if !userExists {
		fmt.Println("doctor doesn't exist")
	}

	// Set the doctor's schedule by inserting the time slot into MongoDB
	slot := models.Slot{
		Date:         input.Slot.Date,
		Time:         input.Slot.Time,
		IsBooked:     false,
		PatientEmail: "",
	}

	// Insert the slot into the doctor_schedule collection
	err = insertSlot(doctor.Name, doctor.ID, input.Slot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(slot)
	fmt.Fprintf(w, "slot inserted successfully")
}

func validatedoctor(collectionName, email string) (bool, *models.User, error) {
	collection := config.Client.Database("clinic").Collection(collectionName)
	filter := bson.D{{Key: "email", Value: email}}

	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		// User not found
		return false, nil, nil
	} else if err != nil {
		// Other errors
		return false, nil, err
	}

	// User found
	return true, &existingUser, nil
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
