import { Component, OnInit, ViewChild } from '@angular/core';

import { DialogComponent } from '../dialog/dialog.component';

import { Stack } from '../../entities/stack';

import { StackService } from '../../services/stack.service';


@Component({
    selector: 'app-stack',
    templateUrl: './stack.component.html',
    styleUrls: ['./stack.component.css']
})
export class StackComponent implements OnInit {

    stacks: Stack[];

    stack: Stack;

    @ViewChild(DialogComponent)
    private dialog: DialogComponent;

    constructor(private stackService: StackService) {
        this.stack = {};
        this.stacks = [];
    }

    onFetchStacks() {
        this.stackService.getStacks().then(stacks => {
            console.log(stacks);
            this.stacks = stacks;
        });
    }

    ngOnInit() {
        this.onFetchStacks();
    }

    onAdd() {
        this.stackService.addStack(this.stack).then(stack => {
            this.stack = {};

            console.log('Stack:', stack);

            this.onFetchStacks();

            this.dialog.close();
        });
    }
}
