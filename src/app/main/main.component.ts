import { Component } from '@angular/core';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css']
})
export class MainComponent {
  inputData: any = {};
  results: any;

  constructor(private apiService: ApiService) { }

  onSubmit() {
    this.apiService.processData(this.inputData).subscribe(
      response => {
        this.results = response;
      },
      error => {
        console.error('Error al procesar los datos', error);
      }
    );
  }
}

