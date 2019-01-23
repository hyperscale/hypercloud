import { NgZone, Injectable } from '@angular/core';
import { URLSearchParams } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';

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
            const params: URLSearchParams = new URLSearchParams();

            if (query) {
                Object.keys(query).map((key) =>  {
                    params.set(key, query[key]);
                });
            }

            let url = '/v1/events';

            if (params.paramsMap.size > 0) {
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
