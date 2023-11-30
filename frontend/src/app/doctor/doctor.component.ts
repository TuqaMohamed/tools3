import { Component, OnInit } from '@angular/core';
import { DoctorService } from './doctor.service';
import { ActivatedRoute } from '@angular/router';
@Component({
  selector: 'app-doctor',
  templateUrl: './doctor.component.html',
  styleUrls: ['./doctor.component.css']
})
export class DoctorComponent implements OnInit{
  appointments: any[] = [];
  showMessagesDropdown: boolean = false;
  doctorMessages: string[] = [
  'ReservationCreated',
  'ReservationUpdated',
  'ReservationCancelled'
]; 
  id: any;
  newAppointment: { date: string, time: string } = { date: '', time: '' };
  constructor(private doctorService: DoctorService,private route: ActivatedRoute) { }
  ngOnInit() {
    this.route.params.subscribe(params=>
      {this.id= params['id'];}
      );
      this.loadDoctorSlots();
  }
  loadDoctorSlots() {
    const doctorId = this.id;
    console.log(doctorId);
    this.doctorService.getDoctorSlots(this.id).subscribe(
      (response) => {
        console.log('Slots successfully retrieved:', response);
        this.appointments = response;
      },
      (error) => {
        console.error('Error fetching doctor slots:', error);
      }
    );
  }
  toggleMessagesDropdown() 
  {
    this.showMessagesDropdown = !this.showMessagesDropdown;
  }
  
  fetchDoctorMessages() 
  {
    this.doctorMessages = ['Message 1', 'Message 2', 'Message 3'];
  }

  editAppointment(appointment: any) 
  {
    console.log('Edit appointment:', appointment);
  }

  cancelAppointment(appointment: any) 
  {
    const index = this.appointments.indexOf(appointment);
    if (index !== -1) {
      this.appointments.splice(index, 1);
    }
  }
  addAppointmentAndSetSchedule() {
    const parsedTime = this.parseTimeString(this.newAppointment.time);
    if (parsedTime) {
      const newAppointment = {
        slot: {
          date: this.newAppointment.date,
          time: this.formatTime(parsedTime)
        }
      };
      this.doctorService.setDoctorSchedule(newAppointment, this.id).subscribe(
        (response) => {
          console.log('Appointment added and schedule set successfully', response);
          this.loadDoctorSlots();
        },
        (error) => {
          console.error('Error adding appointment and setting schedule:', error);
        }
      );
  
      this.newAppointment = { date: '', time: '' };
    } else {
      console.error('Invalid time format');
    }
  }
  
  private formatTime(date: Date): string {
    const hours = date.getHours();
    const minutes = date.getMinutes();
    const ampm = hours >= 12 ? 'PM' : 'AM';
    const formattedHours = hours % 12 === 0 ? 12 : hours % 12;
    const formattedMinutes = minutes < 10 ? '0' + minutes : minutes;
  
    return `${formattedHours}:${formattedMinutes}${ampm}`;
  }
  
    private parseTimeString(timeString: string): Date | null 
  {
    const [hoursStr, minutesStr] = timeString.split(':');
    const hours = parseInt(hoursStr, 10);
    const minutes = parseInt(minutesStr, 10);
  
    if (isNaN(hours) || isNaN(minutes)) {
      console.error('Invalid time format. Please enter a valid time.');
      return null;
    }
    const date = new Date();
    date.setHours(hours, minutes, 0, 0);
    return date;
  }
  
}
