import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { Stack } from '../../entities/stack';
import { Service } from '../../entities/service';
import { StackService } from '../../services/stack.service';
import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-stack-detail',
    templateUrl: './stack-detail.component.html',
    styleUrls: ['./stack-detail.component.less']
})
export class StackDetailComponent implements OnInit {
    public stacks: Stack[] = [];
    public stack: Stack = {};
    public services: Service[] = [];

    constructor(
        private stackService: StackService,
        private serviceService: ServiceService,
        private activatedRoute: ActivatedRoute
    ) { }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            console.log('ID:', params['id']);

            this.fetchStack(params['id']);
            this.fetchServices(params['id']);
        });
    }

    private fetchStack(id: string) {
        this.stackService.getStack(id).then(stack => this.stack = stack);
    }

    private fetchServices(id: string) {
        this.serviceService.getServices({
            stack_id: id
        }).then(services => this.services = services);
    }


}
