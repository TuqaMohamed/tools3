
import { Injectable } from '@angular/core';
import { HttpClient,HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable ,tap} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DoctorService {
  private apiUrl = 'http://localhost:8080'; 

  constructor(private http: HttpClient) { }

  getDoctors(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/getslots`);
  }

  getDoctorSlots(doctorId: string): Observable<any[]> {
    const url = `${this.apiUrl}/getdrslot/${doctorId}/slots`;
    return this.http.get<any[]>(url);
  }

  setDoctorSchedule(slot: any,doctorId:string): Observable<any> {
    const url = `${this.apiUrl}/insertslot/${doctorId}`;  
    console.log('Request payload:', slot);
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<any>(url, slot);
  }

  getSlotId(date: string, time: string, doctorName: string): Observable<any> {
    const url = `${this.apiUrl}/slotid/${date}/${time}/${doctorName}`;
    return this.http.get(url);
  }
  
}
