package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"phase1/config"
	"phase1/helper"
	"phase1/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Message string             `json:"message"`
	User    models.User        `json:"user,omitempty"`
	ID      primitive.ObjectID `json:"id,omitempty"`
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}

func SignUpController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	helper.PanicIfError(err)

	if err := validateSignupRequest(user); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userExists(user.Email, user.UserType) {
		fmt.Println("hello")
		respondWithError(w, "User already registered", http.StatusConflict)
		return
	}

	collectionName := "users"
	if user.UserType == "doctor" {
		collectionName = "doctor"
	} else if user.UserType == "patient" {
		collectionName = "patient"
	}

	collection := config.Client.Database("clinic").Collection(collectionName)
	result, err := collection.InsertOne(context.Background(), bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"userType": user.UserType,
	})
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		respondWithError(w, "Failed to retrieve inserted ID", http.StatusInternalServerError)
		return
	}

	// Set the generated ID in the user object
	user.ID = insertedID

	response := Response{
		Message: "User Signed up successfully",
		User:    user,
		ID:      user.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func userExists(email string, Usertype string) bool {

	collection := config.Client.Database("clinic").Collection(Usertype)
	filter := bson.D{{Key: "email", Value: email}}

	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	return err == nil
}

func validateSignupRequest(user models.User) error { // check if required fields exist
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	if user.UserType == "" {
		return errors.New("user type is required")
	}
	if user.Name == "" {
		return errors.New("name is required")
	}

	return nil

}
