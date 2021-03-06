import { Component, OnInit } from '@angular/core';
import { Stack } from '../../entities';
import { StackService } from '../../services';

@Component({
    selector: 'app-stack-list',
    templateUrl: './stack-list.component.html',
    styleUrls: ['./stack-list.component.less']
})
export class StackListComponent implements OnInit {
    public stacks: Stack[] = [];

    public stack: Stack;

    public submitted = false;

    public _createStackModalOpened = false;

    constructor(private stackService: StackService) {
        this.stack = {
            Name: '',
            Services: 0,
        };
    }

    ngOnInit() {
        this.stackService.getStacks().subscribe(stacks => this.stacks = stacks);
    }

    onSubmit() {
        this.submitted = true;
        console.log('Stack Request:', this.stack);

        this.stackService.create(this.stack).subscribe(stack => {
            console.log('Stack Response:', stack);

            this.stacks.push(stack);

            this.onResetStack();
        });
    }

    onResetStack() {
        this._createStackModalOpened = false;
        this.stack = {
            Name: '',
            Services: 0,
        };
        this.submitted = false;
    }
}
