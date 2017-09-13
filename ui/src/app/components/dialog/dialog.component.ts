import { Component, OnInit, Input } from '@angular/core';

@Component({
    selector: 'ui-dialog',
    templateUrl: './dialog.component.html',
    styleUrls: ['./dialog.component.css']
})
export class DialogComponent implements OnInit {
    @Input()
    width: number;

    @Input()
    top: number;

    @Input()
    opened: boolean;

    constructor() { }

    ngOnInit() {
    }

    open() {
        this.opened = true;
    }

    close() {
        this.opened = false;
    }
}
