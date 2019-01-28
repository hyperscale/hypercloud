import { NgZone, Injectable } from '@angular/core';
import { HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
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

    public getService(id: string): Observable<Service> {
        return this.apiService.get<Service>(`/v1/services/${id}`);
    }

    public getServices(query?: ServicesQueryParams): Observable<Service[]> {
        const params = new HttpParams();

        if (query) {
            Object.keys(query).map((key) =>  {
                params.set(key, query[key]);
            });
        }

        let url = '/v1/services';

        if (params.keys().length > 0) {
            url += '?' + params.toString();
        }

        return this.apiService.get<Service[]>(url);
    }

    public create(service: ServiceRequest): Observable<Service> {
        return this.apiService.post<Service>('/v1/services', service);
    }

    public update(id: string, service: Service): Observable<Service> {
        return this.apiService.put<Service>(`/v1/services/${id}`, service);
    }

    public stats(id: string): Observable<StatsJSON> {
        return new Observable<StatsJSON>(obs => {
            const eventSource = new this.eventSource(this.apiService.getUrl(`/v1/services/${id}/stats`));

            eventSource.onmessage = (event) => {
                const data = JSON.parse(event.data) as StatsJSON;

                this.ngZone.run(() => obs.next(data));
            };

            return () => eventSource.close();
        });
    }
}
