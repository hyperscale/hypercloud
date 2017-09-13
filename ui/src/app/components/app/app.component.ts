import { Component } from '@angular/core';
import { Router } from '@angular/router';

import 'clarity-icons';
import 'clarity-icons/shapes/all-shapes';

import { AuthService } from '../../services/auth.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {
    collapsible = true;
    collapsed = false;

    constructor(private authService: AuthService, private router: Router) {

    }

    public onLogout(event: Event): void {
        event.preventDefault();

        this.authService.logout();

        this.router.navigate(['login']);
    }
}
