import { Component } from '@angular/core';
import { SigninService } from './signin.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent {
  name: string = '';
  email: string = '';
  password: string = '';
  userType: string = '';

  constructor(private signinService: SigninService, private router: Router) {}
  onSubmit() {
    this.signinService.signUp(this.name, this.email, this.password, this.userType).subscribe(
      (response)=>{
        const id = response.id;
        console.log(id);
        console.log(response);
        console.log(this.userType);
        if(this.userType=== 'doctor'){
          console.log('SignUp successful');
          this.router.navigate(['./doctor',id]);
          this.resetForm();
          
        }else if (this.userType === 'patient'){
          console.log('SignUp successful');
          this.router.navigate(['./patient',id]);
          this.resetForm();
        }
      },
      (error) => {
        console.error('SignUp failed:', error);
        this.resetForm();
      }
    );
  }
  private resetForm() {
    this.name = '';
    this.email = '';
    this.password = '';
    this.userType = '';
    
  }
}
