import { Component, OnInit, ViewChild } from '@angular/core';
import { Wizard } from "clarity-angular";

import { Service } from '../../entities/service';

import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service',
    templateUrl: './service.component.html',
    styleUrls: ['./service.component.css']
})
export class ServiceComponent implements OnInit {

    @ViewChild("createWizard")
    createWizard: Wizard;

    _openCreateWizard: boolean = false;

    services: Service[];

    service: Service;

    constructor(private serviceService: ServiceService) {
        this.services = [];
        this.service = {};
    }

    ngOnInit() {
        this.serviceService.getServices().then(services => {
            console.log(services);
            this.services = services;
        });
    }

    onOpenCreateWizard() {
        this._openCreateWizard = !this._openCreateWizard;
    }


    onCreateSubmit() {

    }
}
