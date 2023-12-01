package helper

import (
	"log"
	"net/http"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("menna@2023"))

func StoreUserInSession(w http.ResponseWriter, r *http.Request ,email string) {
	session, err := store.Get(r, "Clinic-session")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		log.Printf("Error retrieving session: %v", err)
		return
	}

	session.Values["email"] = email
	//session.Values["id"] = id
	//log.Printf(email)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		log.Printf("Error saving session: %v", err)
		return
	}
}

func GetUserFromSession(r *http.Request) string {
	session, _ := store.Get(r, "Clinic-session")

	//id, okID := session.Values["id"].(string)
	email, okName := session.Values["email"].(string)
	// if !okID {
	// 	log.Println("id not found in session")
	// 	return "", ""
	// }

	if !okName {
		log.Println("name not found in session")
		return ""
	}

	return email
}