import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class SigninService {

  private apiUrl = 'http://localhost:8080'; 

  constructor(private http: HttpClient) {}

  signUp(Name: string, Email: string, Password: string,UserType: string): Observable<any> {
    const user = { Name, Email, Password ,UserType };
    console.log('Request payload:', user); // Log the payload
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });

   return this.http.post(`${this.apiUrl}/signup`, user);
}
}
