import { Injectable } from '@angular/core';
import { Resolve, ActivatedRouteSnapshot } from '@angular/router';
import { ServiceService } from '../services/service.service';
import { Service } from '../entities/docker';

@Injectable()
export class ServiceResolver implements Resolve<Service> {
    constructor(private serviceService: ServiceService) {}

    resolve(route: ActivatedRouteSnapshot) {
        return this.serviceService.getService(route.paramMap.get('id'));
    }
}
