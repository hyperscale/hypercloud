import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { Stack } from '../../entities/stack';
import { Service } from '../../entities/docker';
import { StackService } from '../../services/stack.service';
import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service-detail',
    templateUrl: './service-detail.component.html',
    styleUrls: ['./service-detail.component.less']
})
export class ServiceDetailComponent implements OnInit {
    public stack: Stack = {};
    public service: Service;

    constructor(
        private stackService: StackService,
        private serviceService: ServiceService,
        private activatedRoute: ActivatedRoute
    ) { }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            console.log('ID:', params['id']);

            this.fetchService(params['id']).then(service => {
                this.stack = {
                    name: service.Spec.Labels['com.docker.stack.namespace'],
                };
            });
        });
    }

    private fetchService(id: string) {
        return this.serviceService.getService(id).then(service => this.service = service);
    }
}
