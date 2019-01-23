import { Component, Input, OnInit } from '@angular/core';
import * as _ from 'lodash';

@Component({
    selector: 'app-docker-state',
    templateUrl: './docker-state.component.html',
    styleUrls: ['./docker-state.component.less']
})
export class DockerStateComponent implements OnInit {
    @Input()
    state: string;

    @Input()
    type: string;

    default = 'info';

    definitions = {
        // @see: https://docs.docker.com/engine/swarm/how-swarm-mode-works/swarm-task-states/
        task: {
            running: 'success',
            shutdown: 'danger',
            accepted: 'info',
        },
        node: {
            ready: 'success'
        }
    };

    label: {
        style: string,
        name: string
    };

    constructor() {
        this.label = {
            style: '',
            name: ''
        };
    }

    ngOnInit() {
        this.label.style = _.get(this.definitions, `${this.type}.${this.state}`, this.default);
        this.label.name = this.state;
    }
}
