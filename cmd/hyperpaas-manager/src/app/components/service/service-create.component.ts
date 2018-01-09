import { Component, OnInit, OnDestroy, OnChanges, EventEmitter, Input, Output, SimpleChanges} from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Application } from '../../entities/application';

@Component({
    selector: 'app-service-create',
    templateUrl: './service-create.component.html',
    styleUrls: ['./service-create.component.less']
})
export class ServiceCreateComponent implements OnInit, OnDestroy {
    public application = new Application();

    public submitted = false;

    constructor() {
    }

    ngOnInit(): void {
    }

    ngOnDestroy(): void {
        this.submitted = false;
    }

    onSubmit() {
        this.submitted = true;
    }
}
