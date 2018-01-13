import { Injectable } from '@angular/core';
import { URLSearchParams } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Service } from '../entities/service';
import { ApiService } from './api.service';

export interface ServicesQueryParams {
    stack_id?: string;
}

@Injectable()
export class ServiceService {
    constructor(private apiService: ApiService) {}

    public getService(id: string): Promise<Service> {
        return this.apiService.get(`/v1/services/${id}`)
        .then(response => response.json())
        .then(response => Object.assign(new Service(), response));
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
                return response.map(service => Object.assign(new Service(), service));
            });
    }

    public create(service: Service): Promise<Service> {
        return this.apiService.post('/v1/services', service)
            .then(response => response.json())
            .then(response => Object.assign(new Service(), response));
    }
}
