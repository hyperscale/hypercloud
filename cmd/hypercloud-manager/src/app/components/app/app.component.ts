import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { EventService } from '../../services';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.less']
})
export class AppComponent implements OnInit, OnDestroy {
    private subscriptions: Subscription[] = [];

    constructor(private eventService: EventService) {}

    ngOnInit() {
        console.log('AppComponent::ngOnInit');

        this.subscriptions.push(
            this.eventService.events().subscribe(event => {
                console.log('Event:', event);
            })
        );
    }

    ngOnDestroy() {
        console.log('AppComponent::ngOnDestroy');

        this.subscriptions.forEach(subscription => subscription.unsubscribe());
    }
}
