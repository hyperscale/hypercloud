import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { Task } from '../entities/task';
//import { Log } from '../entities/log';

import { ApiService } from './api.service';

@Injectable()
export class TaskService {
    constructor(private apiService: ApiService) { }

    /*getLogs(id: string): Promise<Log[]> {
        return new Promise((resolve, reject) => {
            let protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

            var ws = new WebSocket(`${protocol}//${location.host}/ws/container/${id}/log`);
            ws.addEventListener('error', () => reject("error"));
            ws.addEventListener('close', (event) => reject(event.reason))
            ws.addEventListener('message', (event) => resolve(JSON.parse(event.data)), false);
        })
    }*/

    getTasks(service_name: string): Promise<Task[]> {
        let filters = {
            service: [
                service_name
            ]
        };

        return this.apiService.get(`/tasks?filters=${JSON.stringify(filters)}`)
            .then(response => response.json() as Task[])
            .then(tasks => {
                return tasks.map(task => {
                    task.Name = `${service_name}.${task.Slot}`;
                    task.ContainerID = `${task.Name}.${task.ID}`

                    return task;
                })
            });
    }
}
