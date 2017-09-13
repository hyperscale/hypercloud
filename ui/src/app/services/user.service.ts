import { Injectable } from '@angular/core';

import { User } from '../entities/user';
import { Auth } from '../entities/auth';
import { Login } from '../entities/login';

import { ApiService } from './api.service';
import { AuthService } from './auth.service';

const TOKEN_KEY = '__jwt__';

@Injectable()
export class UserService {
    constructor(private apiService: ApiService, private authService: AuthService) {
    }

    public authenticate(login: Login): Promise<Auth> {
        return this.apiService.post("/authenticate", login)
            .then(response => response.json() as Auth)
            .then(auth => {
                return this.authService.authenticate(auth);
            });
    }
}
