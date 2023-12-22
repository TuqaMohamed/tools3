
import { Injectable } from '@angular/core';
import { HttpClient,HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable ,tap} from 'rxjs';
import { environment } from 'src/environments/environment';
import Config from "src/config.json";

@Injectable({
  providedIn: 'root'
})
export class DoctorService {
  private apiUrl=Config.API_URL;

  constructor(private http: HttpClient) { }

  getDoctors(): Observable<any[]> {
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.get<any[]>(`${this.apiUrl}/getslots`);
  }

  getDoctorSlots(doctorId: string): Observable<any[]> {
    const url = `${this.apiUrl}/getdrslot/${doctorId}/slots`;
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.get<any[]>(url);
  }

  setDoctorSchedule(slot: any,doctorId:string): Observable<any> {
    const url = `${this.apiUrl}/insertslot/${doctorId}`;  
    console.log('Request payload:', slot);
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.post<any>(url, slot);
  }

  getSlotId(date: string, time: string, doctorName: string): Observable<any> {
    const url = `${this.apiUrl}/slotid/${date}/${time}/${doctorName}`;
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.get(url);
  }
  
}
