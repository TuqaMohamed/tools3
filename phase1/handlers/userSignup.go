package handlers

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
)

func SignUpController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	helper.PanicIfError(err)

	if err := validateSignupRequest(user); err != nil { //check if required fields exist
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userExists(user.Email, user.UserType) { //check if user is already registered
		http.Error(w, "User already registered", http.StatusConflict)
		return
	}

	collectionName := "users"
	if user.UserType == "doctor" {
		collectionName = "doctor"
	} else if user.UserType == "patient" {
		collectionName = "patient"
	}

	collection := config.Client.Database("clinic").Collection(collectionName)
	_, err = collection.InsertOne(context.Background(), bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)

	fmt.Fprintf(w, "User Signed up successfully")

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
