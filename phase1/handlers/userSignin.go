package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"phase1/config"
	"phase1/helper"
	"phase1/models"

	//"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func SignInController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateSigninRequest(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user credentials against MongoDB
	if !validateUser("patient", user.Email, user.Password) && !validateUser("doctor", user.Email, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Respond with a success message
	//log.Println(user.Email)
	helper.StoreUserInSession(w, r, user.Email)
	fmt.Fprintf(w, "User signed in successfully")

}

func validateUser(collectionName, email, password string) bool {
	collection := config.Client.Database("clinic").Collection(collectionName)
	filter := bson.D{{Key: "email", Value: email}, {Key: "password", Value: password}}

	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	return err == nil
}

func validateSigninRequest(user models.User) error { // check if required fields exist
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil

}