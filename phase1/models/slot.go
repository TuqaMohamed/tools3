package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Slot struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date         string             `json:"date" bson:"date"`
	DoctorName   string             `json:"doctorName" bson:"doctorName"`
	IsBooked     bool               `json:"isBooked" bson:"isBooked"`
	PatientEmail string             `json:"patientEmail,omitempty" bson:"patientEmail,omitempty"`
	Time         string             `json:"time" bson:"time"`
}
