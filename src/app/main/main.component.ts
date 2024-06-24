import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { JsonPipe, NgIf } from '@angular/common';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-main',
  standalone: true,
  imports: [FormsModule, JsonPipe, NgIf],
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
