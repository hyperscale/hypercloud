import { Component, OnInit, DoCheck, Input, ViewChild, ElementRef } from '@angular/core';

import { Log } from '../../entities/log';

@Component({
    selector: 'console',
    templateUrl: './console.component.html',
    styleUrls: ['./console.component.css']
})
export class ConsoleComponent implements OnInit, DoCheck {
    @Input()
    public items: Log[];

    @Input()
    public follow: boolean;

    @ViewChild('console')
    private $console: ElementRef;

    constructor() {
        this.follow = true;
    }

    ngOnInit() {
        this.$console.nativeElement.addEventListener('scroll', (event: Event) => {
            const height = ((<any>event.target).scrollHeight - (<any>event.target).offsetHeight);

            if (this.follow === true && ((<any>event.target).scrollTop < height)) {
                console.log('Follow disabled');
                this.follow = false;
            } else if (this.follow === false && ((<any>event.target).scrollTop === height)) {
                console.log('Follow enabled');
                this.follow = true;
            }
        }, false);
    }

    ngDoCheck() {
        if (this.follow) {
            setTimeout(() => {
                this.$console.nativeElement.scrollTop = this.$console.nativeElement.scrollHeight;
            }, 100);
        }
    }
}
