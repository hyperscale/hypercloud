import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Rx';
import { WebSocketService } from './websocket.service';

@Injectable()
export class EventService {
    private eventSource: Subject<any> = new Subject<any>();

    public event = this.eventSource.asObservable();

    constructor(private websocketService: WebSocketService) {
        const protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

        const url = `${protocol}//${location.host}/ws/events`;

        this.websocketService.connect(url);

        this.websocketService.message$.subscribe((event: MessageEvent) => {
            this.eventSource.next(JSON.parse(event.data));
        });
    }

    close(): void {
        this.websocketService.close();
    }
}
