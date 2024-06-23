import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = 'http://localhost:8000'; // URL del API Gateway

  constructor(private http: HttpClient) { }

  processData(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/process`, data);
  }
}

