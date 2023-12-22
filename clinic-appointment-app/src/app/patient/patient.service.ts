import { Injectable } from '@angular/core';
import { HttpClient,HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable ,tap} from 'rxjs';
import Config from "src/config.json";
@Injectable({
  providedIn: 'root'
})
export class PatientService {

  private apiUrl =Config.API_URL ; 
  constructor(private http: HttpClient) { }
  
  reserveSlot(slotId: string, patientId: string): Observable<any> {
  const url=`${this.apiUrl}/reserveslot/${slotId}/${patientId}`;
  console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.post(url,{});
  }

  viewPatientSlots(patientID:any)
  {
    const url=`${this.apiUrl}/viewslots/${patientID}`;
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.get<any[]>(url,{});
  }

  deleteSlot(slotId: any): Observable<any> {
    const url = `${this.apiUrl}/delete/${slotId}`;
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.put(url, {});
  }

  updateReservedSlot(oldSlotId: string, newSlotId: string, patientId: string): Observable<any> {
    const url = `${this.apiUrl}/update-slot/${oldSlotId}/${newSlotId}/${patientId}`;
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.put(url, {});
  }
}
