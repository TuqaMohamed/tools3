import { Component } from '@angular/core';
import { LoginService } from './login.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  email: string = '';
  password: string = '';
  constructor(private loginService: LoginService, private router: Router) {}
  onSubmit() {
    this.loginService.login(this.email, this.password).subscribe(
      (response)=>{
        const id = response.user._id;
        const userType = response.user.userType;
        console.log(id);
        console.log(response);
        console.log(response.userType);
        console.log(this.email);
        if(userType == 'doctor'){
          console.log('Sign in successful');
          this.router.navigate(['/doctor',id]);
          this.resetForm();
          
        }else if (userType == 'patient'){
          console.log('Sign in successful');
          this.router.navigate(['/patient',id]);
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
    this.email = '';
    this.password = '';
  }
}
