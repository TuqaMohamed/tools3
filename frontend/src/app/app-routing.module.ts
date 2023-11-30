import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SigninComponent } from './signin/signin.component';
import { LoginComponent } from './login/login.component';
import { DoctorComponent } from './doctor/doctor.component';
import { PatientComponent } from './patient/patient.component';
const routes: Routes = [
  { path: 'signin', component: SigninComponent },
  { path: 'login', component: LoginComponent},
  { path: 'doctor/:id', component: DoctorComponent },
  { path: 'patient/:id', component: PatientComponent },
  { path: '', redirectTo: '/signin', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
