import { Component, OnInit } from '@angular/core';
import { Service } from '../../entities/service';
import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service-list',
    templateUrl: './service-list.component.html',
    styleUrls: ['./service-list.component.less']
})
export class ServiceListComponent implements OnInit {
    public services: Service[] = [];

    constructor(private serviceService: ServiceService) { }

    ngOnInit() {
        this.serviceService.getServices().then(services => this.services = services);
    }
}
