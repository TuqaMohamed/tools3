import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import Config from "src/config.json";
@Injectable({
  providedIn: 'root'
})
export class SigninService {
  private apiUrl =Config.API_URL ; 
  constructor(private http: HttpClient) {}

  signUp(Name: string, Email: string, Password: string, UserType: string): Observable<any> {
    const user = { Name, Email, Password, UserType };
    console.log('Request payload:', user);
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    console.log('Constructed URL:', `${this.apiUrl}/signup`);
    return this.http.post(`${this.apiUrl}/signup`, user);
  }
}
