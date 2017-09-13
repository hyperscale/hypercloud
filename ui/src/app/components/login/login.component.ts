import { Component, OnInit } from '@angular/core';

import { Login } from '../../entities/login';
import { UserService } from '../../services/user.service';
import { Router } from '@angular/router';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    login: Login;

    constructor(private userService: UserService, private router: Router) {
        this.login = {
            username: '',
            password: ''
        };
    }



    ngOnInit() {
    }

    public onSubmit() {
        this.userService.authenticate(this.login).then(auth => {
            this.login = {
                username: '',
                password: ''
            };

            console.log('Auth:', auth);

            this.router.navigate(['dashboard']);
        }).catch(reason => {
            console.error(reason);
        });
    }
}
