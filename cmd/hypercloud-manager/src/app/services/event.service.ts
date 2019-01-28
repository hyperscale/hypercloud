import { NgZone, Injectable } from '@angular/core';
import { HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Message } from '../entities/docker';
import { ApiService } from './api.service';

export interface EventsQueryParams {
    stack_id?: string;
}

@Injectable()
export class EventService {
    private eventSource: any = window['EventSource'];

    constructor(private apiService: ApiService, private ngZone: NgZone) {
    }

    public events(query?: EventsQueryParams): Observable<Message> {
        console.log('EventService::events()');

        return new Observable<Message>(obs => {
            const params = new HttpParams();

            if (query) {
                Object.keys(query).map((key) =>  {
                    params.set(key, query[key]);
                });
            }

            let url = '/v1/events';

            if (params.keys().length > 0) {
                url += '?' + params.toString();
            }

            const eventSource = new this.eventSource(this.apiService.getUrl(url));

            eventSource.onmessage = (event) => {
                console.log('Message:', event.data);
                const data = JSON.parse(event.data) as Message;

                this.ngZone.run(() => obs.next(data));
            };

            return () => eventSource.close();
        });
    }
}
