package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Slot struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date         string             `json:"date" bson:"date"`
	DoctorName   string             `json:"doctorName" bson:"doctorName"`
	IsBooked     bool               `json:"isBooked" bson:"isBooked"`
	PatientName  string             `json:"patientName,omitempty" bson:"patientName,omitempty"`
	Time         string             `json:"time" bson:"time"`
	PatientID    primitive.ObjectID `json:"patientID,omitempty" bson:"patientID,omitempty"`
}
