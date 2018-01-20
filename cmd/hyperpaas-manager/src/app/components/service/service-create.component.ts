import { Component, OnInit, OnDestroy, OnChanges, EventEmitter, Input, Output, SimpleChanges} from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Stack, ServiceRequest } from '../../entities';
import { StackService, ServiceService } from '../../services';

@Component({
    selector: 'app-service-create',
    templateUrl: './service-create.component.html',
    styleUrls: ['./service-create.component.less']
})
export class ServiceCreateComponent implements OnInit, OnDestroy {
    public service: ServiceRequest = {
        stack_id: '',
        name: ''
    };
    public stacks: Stack[] = [];
    public stack = '';

    public submitted = false;

    constructor(
        private stackService: StackService,
        private serviceService: ServiceService
    ) {
    }

    ngOnInit(): void {
        this.stackService.getStacks().then(stacks => this.stacks = stacks);
    }

    ngOnDestroy(): void {
        this.submitted = false;
    }

    onReset() {
        this.service = {
            stack_id: '',
            name: ''
        };
        this.submitted = false;
    }

    onSubmit() {
        this.submitted = true;

        console.log('Service Request:', this.service);

        this.serviceService.create(this.service).then(service => {
            console.log('Service Response', service);

            this.onReset();
            // @todo: redirect
        });
    }

    get diagnostic() {
        return JSON.stringify(this.service, null, '  ');
    }
}
