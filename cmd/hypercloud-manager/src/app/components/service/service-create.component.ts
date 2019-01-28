import { Component, OnInit, OnDestroy, OnChanges, EventEmitter, Input, Output, SimpleChanges} from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Router, ActivatedRoute, Params } from '@angular/router';
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

    public submitted = false;

    constructor(
        private stackService: StackService,
        private serviceService: ServiceService,
        private router: Router,
        private activatedRoute: ActivatedRoute
    ) {
    }

    ngOnInit(): void {
        this.activatedRoute.queryParams.subscribe((params: Params) => {
            console.log('StackID:', params['stack_id']);

            this.service.stack_id = params['stack_id'] || '';

            this.stackService.getStacks().subscribe(stacks => this.stacks = stacks);
        });
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

        this.serviceService.create(this.service).subscribe(service => {
            console.log('Service Response', service);

            this.onReset();

            this.router.navigate(['/service/', service.ID]);
        });
    }
}
