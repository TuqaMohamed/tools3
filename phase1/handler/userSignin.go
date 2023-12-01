package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"phase1/config"
	"phase1/helper"
	"phase1/models"

	//"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type response struct {
	Message  string             `json:"message"`
	User     models.User        `json:"user,omitempty"`
	ID       primitive.ObjectID `json:"id,omitempty"`
	UserType string             `json:"userType" bson:"userType"`
	Email    string             `json:"email" bson:"email"`
}

func RespondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
func SignInController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateSigninRequest(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user credentials against MongoDB
	// if !validateUser("patient", user.Email, user.Password) && !validateUser("doctor", user.Email, user.Password) {
	// 	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	// 	return
	// }

	// // Retrieve the user ID from the database
	// userID, err := getUserID(user.UserType, user.Email, user.Password)
	// fmt.Println(userID)
	// if err != nil {
	// 	http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
	// 	return
	// }

	exists, user := validateUser(user.Email, user.Password)

	if exists {
		fmt.Println("User exists:", user)
	} else {
		fmt.Println("User does not exist")
	}

	// _ , err := getUser(user.UserType, user.Email, user.Password)
	// if err != nil {
	//     http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	//     return
	// }

	response := response{
		Message:  "User Signed in successfully",
		User:     user,
		UserType: user.UserType,
		Email:    user.Email,
		ID:       user.ID,
	}

	fmt.Println(response.UserType)
	helper.StoreUserInSession(w, r, user.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
func getUser(collectionName, email, password string) (models.User, error) {
	collection := config.Client.Database("clinic").Collection(collectionName)
	filter := bson.D{{Key: "email", Value: email}, {Key: "password", Value: password}}

	var existingUser models.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err != nil {
		return models.User{}, err
	}

	return existingUser, nil
}

func validateUser(email, password string) (bool, models.User) {
	// Specify the two collections to check
	collectionNames := []string{"doctor", "patient"}

	// Iterate over the collections and check for the user
	for _, collectionName := range collectionNames {
		collection := config.Client.Database("clinic").Collection(collectionName)
		filter := bson.D{{Key: "email", Value: email}, {Key: "password", Value: password}}

		var existingUser models.User
		err := collection.FindOne(context.Background(), filter).Decode(&existingUser)

		if err == nil {
			// User found in one of the collections, return the result
			return true, existingUser
		}
	}

	// User not found in any collection
	return false, models.User{}
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
