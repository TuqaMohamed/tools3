import { Injectable } from '@angular/core';
import { HttpClient,HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable ,tap} from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class PatientService {
  private apiUrl = 'http://localhost:8080'; 
  constructor(private http: HttpClient) { }
  
  reserveSlot(slotId: string, patientId: string): Observable<any> {
  const url=`${this.apiUrl}/reserveslot/${slotId}/${patientId}`;
    return this.http.post(url,{});
  }

  viewPatientSlots(patientID:any)
  {
    const url=`${this.apiUrl}/viewslots/${patientID}`;
    return this.http.get<any[]>(url,{});
  }

  deleteSlot(slotId: any): Observable<any> {
    const url = `${this.apiUrl}/delete/${slotId}`;
    return this.http.put(url, {});
  }

  updateReservedSlot(oldSlotId: string, newSlotId: string, patientId: string): Observable<any> {
    const url = `${this.apiUrl}/update-slot/${oldSlotId}/${newSlotId}/${patientId}`;
    return this.http.put(url, {});
  }
}
