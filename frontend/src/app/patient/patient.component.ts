import { Component, OnInit } from '@angular/core';
import { DoctorService } from '../doctor/doctor.service';
import { PatientService } from './patient.service';
import { ActivatedRoute } from '@angular/router';

interface Doctor {
  id: string;
  name: string;
}

interface Slot {
  time: string;
  date:string;
  id:string;
}
@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.css']
})
export class PatientComponent implements OnInit{
  selectedDoctor: any;
  doctors: Doctor[] = [];
  appointments: any[] = [];
  newAppointment: { date: string, time: string, slot?: Slot } = { date: '', time: ''};
  availableSlots: Slot[] = [];
  editedAppointment: any | null = null;
  selectedSlot: any ;
  patientID :any;
  editedDoctor: any;
  editedSlot: any;
  editedAvailableSlots: any[] = [];
  oldID:any;
  
  constructor(private doctorService: DoctorService , private route: ActivatedRoute,private patientService:PatientService) {
    this.loadDoctors();
    this.route.params.subscribe(params => {
      this.patientID = params['id'];
    });
    console.log(this.patientID);
  }

  ngOnInit() {
    
      this.loadPatientSlots();
  }
loadDoctors() {
    this.doctorService.getDoctors().subscribe((doctors: Doctor[]) => {
    this.doctors = doctors;
    console.log('Doctors:', this.doctors); 
  });
  }

  loadDoctorSlots(isEditing: boolean) {
    const selectedDoctor = isEditing ? this.editedDoctor : this.selectedDoctor;
  
    if (selectedDoctor) {
      this.doctorService.getDoctorSlots(selectedDoctor._id).subscribe((slots: Slot[]) => {
        if (isEditing) {
          this.editedAvailableSlots = slots;
        } else {
          this.availableSlots = slots;
        }
      });
    }
  }
  
  loadPatientSlots() {
    console.log(this.patientID)
    
    this.patientService.viewPatientSlots(this.patientID).subscribe((slots: Slot[]) => {
      this.appointments = slots;
      console.log(this.appointments);
      });
  }
  cancelAppointment(appointment: any) {
    const slotId = appointment.slotID;
    console.log(slotId);
      this.patientService.deleteSlot(slotId).subscribe(
      (response) => {
        console.log('Slot canceled successfully:', response);
        this.loadPatientSlots();
      },
      (error) => {
        console.error('Error canceling slot:', error);
      }
    );
  }
  editAppointment(appointment: any) {
    this.editedAppointment = { ...appointment };
    this.oldID=appointment.slotID;
    console.log("old id:",this.oldID);
  }
  saveEditedAppointment() {
    console.log("old slot ", this.oldID);
    console.log("new slot id",this.editedSlot.id);
    console.log("patient id ",this.patientID)
      const oldSlotId = this.oldID; 
      const newSlotId = this.editedSlot.id; 
      const patientId = this.patientID; 
      this.patientService.updateReservedSlot(oldSlotId, newSlotId, patientId).subscribe(
        () => {
          const index = this.appointments.findIndex(appointment => appointment.id === this.editedAppointment.id);
          if (index !== -1) {
            this.appointments[index] = { ...this.editedAppointment };
            this.editedAppointment = null;
          }
          this.loadPatientSlots();
        },
        (error) => {
          console.error('Error updating reserved slot:', error);
        }
      );
    }
reserveAppointment() {
  console.log("inside the reserve");
    console.log(this.selectedSlot.id);
    console.log(this.patientID)
    this.patientService.reserveSlot(this.selectedSlot.id, this.patientID)
            .subscribe(
              (reservationResponse: any) => {
                this.loadPatientSlots();
                console.log('Slot reserved successfully:', reservationResponse);
              },
              (reservationError: any) => {
                console.error('Error reserving slot:', reservationError);
              }
            );
}
}