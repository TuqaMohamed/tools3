// config.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {
  private configUrl = 'assets/config.json'; 
  private apiUrl: string='';

  constructor(private http: HttpClient) {}

  loadConfig(): Promise<any> {
    return this.http.get(this.configUrl).toPromise()
      .then((config: any) => {
        this.apiUrl = config.API_URL || 'http://localhost:8080';
        console.log('Constructed URL:', this.apiUrl);
      })
      .catch(error => {
        console.error('Error loading configuration:', error);
      });
  }
  

  getApiUrl(): string {
    return this.apiUrl;
  }
}
