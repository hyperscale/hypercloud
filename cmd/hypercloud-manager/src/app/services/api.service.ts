import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptionsArgs, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/toPromise';

import { environment } from '../../environments/environment';

const serialize = (obj: any) => {
    const str: any[] = [];
    for (const p in obj) {
        if (obj.hasOwnProperty(p)) {
            str.push(`${encodeURIComponent(p)}=${encodeURIComponent(obj[p])}`);
        }
    }

    return str.join('&');
};

export interface OAuthToken {
    access_token: string;
    refresh_token: string;
    expires_in: number;
}

const TOKEN_KEY = 'token';

@Injectable()
export class ApiService {

    private base = environment.api.url;

    private client_id = environment.api.client_id;

    private timeoutId: number;

    constructor(private http: Http) {
    }

    getToken(): OAuthToken {
        const tokenData = localStorage.getItem(TOKEN_KEY);

        if (tokenData === null) {
            return null;
        }

        return JSON.parse(tokenData) as OAuthToken;
    }

    setToken(token: OAuthToken): void {
        console.log('setToken', {token});

        localStorage.setItem(TOKEN_KEY, JSON.stringify(token));

        const now = Math.floor(Date.now() / 1000);

        if (this.timeoutId) {
            console.log('clearTimeout');

            window.clearTimeout(this.timeoutId);
        }

        const ttl = ((token.expires_in - 60) - now);

        console.log('access_token expire in ', {ttl});

        this.timeoutId = window.setTimeout(() => {
            console.log('refreshToken');

            this.refreshToken();
        }, ttl * 1000);
    }

    removeToken(): void {
        localStorage.removeItem(TOKEN_KEY);

        if (this.timeoutId) {
            window.clearTimeout(this.timeoutId);
        }
    }

    isFrechToken(): boolean {
        const now = Math.round(Date.now() / 1000);
        const token = this.getToken();

        if (token && token.expires_in > now) {
            return true;
        }

        return false;
    }

    refreshToken(): Promise<OAuthToken> {
        const token = this.getToken();

        return this.post('/v1/oauth/token', {
            grant_type: 'refresh_token',
            refresh_token: token.refresh_token
        }, {
            headers: new Headers({
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authorization': this.getAuthorizationBasic()
            }),
        }).then((response) => {
            const tok = response.json() as OAuthToken;

            this.setToken(tok);

            return tok;
        }).catch((reason) => {
            this.removeToken();

            return reason;
        });
    }

    getAuthorizationBasic(): string {
        const basic = `${this.client_id}:`;

        return `Basic ${btoa(basic)}`;
    }

    getAuthorization(): string {
        const token = this.getToken();

        if (token) {
            return `Bearer ${token.access_token}`;
        }

        return this.getAuthorizationBasic();
    }

    getOptions(options?: RequestOptionsArgs): RequestOptionsArgs {
        const headers = new Headers({
            'Content-Type': 'application/json',
            'Authorization': this.getAuthorization()
        });

        if (!options) {
            options = {};
        }

        if (options.headers) {
            headers.forEach((values, name) => {
                if (!options.headers.has(name)) {
                    options.headers.set(name, values);
                }
            });
        } else {
            options.headers = headers;
        }

        options.headers.forEach((values, name) => {
            values.forEach(val => {
                if ((<any>val) === false) {
                    options.headers.delete(name);
                }
            });
        });

        return options;
    }

    getUrl(path: string): string {
        if (path.charAt(0) === '/') {
            path = path.substr(1, path.length);
        }

        return `${this.base}/${path}`;
    }

    processBody(body: any, options: RequestOptionsArgs): any {
        if (options.headers.get('Content-Type') === 'application/x-www-form-urlencoded') {
            return serialize(body);
        }

        return body;
    }

    /**
     * Performs a request with `get` http method.
     */
    get(path: string, options?: RequestOptionsArgs): Promise<Response> {
        return this.get$(path, options).toPromise();
    }

    /**
     * Performs a request with `get` http method.
     */
    get$(path: string, options?: RequestOptionsArgs): Observable<Response> {
        return this.http.get(this.getUrl(path), this.getOptions(options));
    }

    /**
     * Performs a request with `post` http method.
     */
    post(path: string, body: any, options?: RequestOptionsArgs): Promise<Response> {
        return this.post$(path, body, options).toPromise();
    }

    /**
     * Performs a request with `post` http method.
     */
    post$(path: string, body: any, options?: RequestOptionsArgs): Observable<Response> {
        options = this.getOptions(options);

        return this.http.post(this.getUrl(path), this.processBody(body, options), options);
    }

    /**
     * Performs a request with `put` http method.
     */
    put(path: string, body: any, options?: RequestOptionsArgs): Promise<Response> {
        return this.put$(path, body, options).toPromise();
    }

     /**
     * Performs a request with `put` http method.
     */
    put$(path: string, body: any, options?: RequestOptionsArgs): Observable<Response> {
        options = this.getOptions(options);

        return this.http.put(this.getUrl(path), this.processBody(body, options), options);
    }

    /**
     * Performs a request with `delete` http method.
     */
    delete(path: string, options?: RequestOptionsArgs): Promise<Response> {
        return this.delete$(path, options).toPromise();
    }

    /**
     * Performs a request with `delete` http method.
     */
    delete$(path: string, options?: RequestOptionsArgs): Observable<Response> {
        return this.http.delete(this.getUrl(path), this.getOptions(options));
    }
}
