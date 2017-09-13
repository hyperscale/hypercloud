import { Injectable } from '@angular/core';

import { Auth } from '../entities/auth';

const TOKEN_KEY = '__jwt__';

@Injectable()
export class AuthService {
    authenticated: boolean;

    constructor() {
        this.authenticated = false;

        this.autoconnect();
    }

    public autoconnect() {
        const token = this.getToken();

        if (token) {
            this.authenticated = true;
        }
    }

    public getToken(): string {
        let value = localStorage.getItem(TOKEN_KEY);
        if (value) {
            let auth = <Auth>JSON.parse(value);

            if (auth.expires <= Math.floor(Date.now()/1000)-20) {
                console.log('JWT has expired.');

                localStorage.removeItem(TOKEN_KEY);

                this.authenticated = false;

                return null;
            }

            return auth.token;
        }

        return null;
    }

    public authenticate(auth: Auth): Auth {
        localStorage.setItem(TOKEN_KEY, JSON.stringify(auth));
        this.authenticated = true;

        return auth;
    }

    public isAuthenticated(): boolean {
        this.getToken();

        return this.authenticated;
    }

    public logout(): void {
        this.authenticated = false;

        localStorage.removeItem(TOKEN_KEY);
    }
}
