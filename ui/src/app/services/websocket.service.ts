import { Injectable, NgZone } from '@angular/core';
import { Subject } from 'rxjs/Subject';

export enum WebSocketState {
    Connecting = 0,
    Connected,
    Disconnecting,
    Disconnected,
}

export enum WebSocketCloseCode {
    Normal          = 1000,
    GoingAway       = 1001,
    ProtocolError   = 1002,
    Unsupported     = 1003,
    NoStatus        = 1005,
    Abnormal        = 1006,
    TooLarge        = 1009
}

@Injectable()
export class WebSocketService {
    private pingTimer: number;
    private pongTimer: number;
    private ws: WebSocket;
    private state: WebSocketState;
    private url: string;

    private stateSource: Subject<WebSocketState> = new Subject<WebSocketState>();
    private messageSource: Subject<MessageEvent> = new Subject<MessageEvent>();

    public state$ = this.stateSource.asObservable();
    public message$ = this.messageSource.asObservable();



    constructor() {
        this.state = WebSocketState.Disconnected;

        window.addEventListener('offline', () => {
            this.close();
        }, false);

        window.addEventListener('online', () => {
            if (this.state === WebSocketState.Disconnected) {
                setTimeout(() => this.connect(), 200);
            }
        }, false);
    }

    public changeState(state: WebSocketState) {
        this.state = state;
        this.stateSource.next(state);
    }

    public connect(url?: string) {
        this.changeState(WebSocketState.Connecting);

        if (url) {
            this.url = url;
        }

        this.ws = new WebSocket(this.url);

        this.ws.addEventListener('open', (event) => {
            this.changeState(WebSocketState.Connected);

            this.pingTimer = window.setInterval(() => {
                this.send('{"type":"ping"}');

                this.pongTimer = window.setTimeout(() => {
                    this.close();
                }, 10000);
            }, 30000);
        });

        this.ws.addEventListener('error', (event) => {
            this.messageSource.error(event);
        });

        this.ws.addEventListener('close', (event: CloseEvent) => {
            this.ws = null;

            this.changeState(WebSocketState.Disconnected);

            this.messageSource.complete();
        });

        this.ws.addEventListener('message', (event: MessageEvent) => {
            if (event.data === '{"type":"pong"}') {
                window.clearTimeout(this.pongTimer);
            } else {
                this.messageSource.next(event);
            }
        });
    }

    public send(data: string) {
        if (this.state !== WebSocketState.Connected) {
            return;
        }

        this.ws.send(data);
    }

    public close() {
        this.changeState(WebSocketState.Disconnecting);

        window.clearInterval(this.pingTimer);
        window.clearTimeout(this.pongTimer);

        this.ws.close(WebSocketCloseCode.Normal);
    }
}
