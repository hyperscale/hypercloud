import { Component, OnInit, OnDestroy } from '@angular/core';
import {Router, ActivatedRoute, Params} from '@angular/router';

import { Log } from '../../entities/log';

import { ServiceService } from '../../services/service.service';
import { WebSocketService } from '../../services/websocket.service';

@Component({
    selector: 'app-service-log',
    templateUrl: './service-log.component.html',
    styleUrls: ['./service-log.component.css'],
    providers: [
        WebSocketService,
    ]
})
export class ServiceLogComponent implements OnInit, OnDestroy {
    logs: Log[];

    constructor(
        private serviceService: ServiceService,
        private websocketService: WebSocketService,
        private activatedRoute: ActivatedRoute
    ) {
        this.logs = [];
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            const service = params['name'];

            const protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

            const url = `${protocol}//${location.host}/ws/services/${service}/logs`;

            this.websocketService.connect(url);

            this.websocketService.message$.subscribe((event: MessageEvent) => {
                let log = JSON.parse(event.data) as Log;

                console.log(log);
                this.logs.push(log);
            });

        })
    }

    ngOnDestroy() {
        this.websocketService.close();
    }
}
