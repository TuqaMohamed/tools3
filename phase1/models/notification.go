package models

type ReservationEvent struct {
	DoctorID  string `json:"doctorId" bson:"doctorId"`
	PatientID string `json:"patientId" bson:"patientId"`
	Operation string `json:"Operation" bson:"operation"`
}
