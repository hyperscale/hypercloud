import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, of } from 'rxjs';
import { catchError, switchMap } from 'rxjs/operators';

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

export interface RequestOptions {
    headers?: HttpHeaders;
    observe?: 'body';
    params?: HttpParams | {
        [param: string]: string | string[];
    };
    reportProgress?: boolean;
    responseType?: 'json';
    withCredentials?: boolean;
}

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

    constructor(private http: HttpClient) {
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

    private handleError(error: HttpErrorResponse) {
        if (error.error instanceof ErrorEvent) {
          // A client-side or network error occurred. Handle it accordingly.
          console.error('An error occurred:', error.error.message);
        } else {
          // The backend returned an unsuccessful response code.
          // The response body may contain clues as to what went wrong,
          console.error(
            `Backend returned code ${error.status}, ` +
            `body was: ${error.error}`);
        }

        // return an observable with a user-facing error message
        return throwError('Something bad happened; please try again later.');
    }

    refreshToken(): Observable<OAuthToken> {
        const token = this.getToken();

        return this.post<OAuthToken>('/v1/oauth/token', {
            grant_type: 'refresh_token',
            refresh_token: token.refresh_token
        }, {
            headers: new HttpHeaders({
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authorization': this.getAuthorizationBasic()
            }),
        }).pipe(
            catchError((error: HttpErrorResponse) => {
                this.removeToken();

                return this.handleError(error);
            }),
            switchMap(response => {
                this.setToken(response);

                return of(response);
            })
        );
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

    getOptions(options?: RequestOptions): RequestOptions {
        const headers = new HttpHeaders({
            'Content-Type': 'application/json',
            'Authorization': this.getAuthorization()
        });

        if (!options) {
            options = {};
        }

        if (options.headers) {
            headers.keys().forEach((key) => {
                if (!options.headers.has(key)) {
                    options.headers.set(key, headers.getAll(key));
                }
            });
        } else {
            options.headers = headers;
        }

        options.headers.keys().forEach((key) => {
            options.headers.getAll(key).forEach((val) => {
                if ((<any>val) === false) {
                    options.headers.delete(key);
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

    processBody(body: any, options: RequestOptions): any {
        if (options.headers.get('Content-Type') === 'application/x-www-form-urlencoded') {
            return serialize(body);
        }

        return body;
    }

    /**
     * Performs a request with `get` http method.
     */
    get<T>(path: string, options?: RequestOptions): Observable<T> {
        return this.http.get<T>(this.getUrl(path), this.getOptions(options));
    }

    /**
     * Performs a request with `post` http method.
     */
    post<T>(path: string, body: any, options?: RequestOptions): Observable<T> {
        options = this.getOptions(options);

        return this.http.post<T>(this.getUrl(path), this.processBody(body, options), options);
    }

    /**
     * Performs a request with `put` http method.
     */
    put<T>(path: string, body: any, options?: RequestOptions): Observable<T> {
        options = this.getOptions(options);

        return this.http.put<T>(this.getUrl(path), this.processBody(body, options), options);
    }

    /**
     * Performs a request with `delete` http method.
     */
    delete<T>(path: string, options?: RequestOptions): Observable<T> {
        return this.http.delete<T>(this.getUrl(path), this.getOptions(options));
    }
}
