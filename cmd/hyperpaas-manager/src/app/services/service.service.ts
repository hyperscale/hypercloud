import { NgZone, Injectable } from '@angular/core';
import { URLSearchParams } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import { Service, StatsJSON, ServiceRequest } from '../entities';
import { ApiService } from './api.service';

export interface ServicesQueryParams {
    stack_id?: string;
}

@Injectable()
export class ServiceService {
    private eventSource: any = window['EventSource'];

    constructor(private apiService: ApiService, private ngZone: NgZone) {
    }

    public getService(id: string): Promise<Service> {
        return this.apiService.get(`/v1/services/${id}`)
        .then(response => response.json() as Service);
    }

    public getServices(query?: ServicesQueryParams): Promise<Service[]> {
        const params: URLSearchParams = new URLSearchParams();

        if (query) {
            Object.keys(query).map((key) =>  {
                params.set(key, query[key]);
            });
        }

        let url = '/v1/services';

        if (params.paramsMap.size > 0) {
            url += '?' + params.toString();
        }

        return this.apiService.get(url)
            .then(response => response.json())
            .then(response => {
                return response.map(service => service as Service);
            });
    }

    public create(service: ServiceRequest): Promise<Service> {
        return this.apiService.post('/v1/services', service)
            .then(response => response.json() as Service);
    }

    public stats(id: string): Observable<StatsJSON> {
        return new Observable<StatsJSON>(obs => {
            const eventSource = new this.eventSource(this.apiService.getUrl(`/v1/services/${id}/stats`));

            eventSource.onmessage = event => {
                const data = JSON.parse(event.data) as StatsJSON;

                this.ngZone.run(() => obs.next(data));
            };

            return () => eventSource.close();
        });
    }
}
