import { Component, OnInit, OnDestroy, OnChanges, EventEmitter, Input, Output, SimpleChanges} from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Service, Stack } from '../../entities/docker';
import { StackService } from '../../services/stack.service';
import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service-create',
    templateUrl: './service-create.component.html',
    styleUrls: ['./service-create.component.less']
})
export class ServiceCreateComponent implements OnInit, OnDestroy {
    public service: Service;
    public stacks: Stack[] = [];

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
        this.service = null;
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
