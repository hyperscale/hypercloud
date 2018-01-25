import { Injectable } from '@angular/core';
import { URLSearchParams } from '@angular/http';
import { Node, Swarm, Info, VersionResponse, Task } from '../entities/docker';
import { ApiService } from './api.service';
import * as _ from 'lodash';

export class DockerFilterParams {
    containers: { [key: string]: string[] };

    constructor(obj?: { [key: string]: any }) {
        this.containers = {};

        if (obj) {
            Object.keys(obj).forEach(key => {
                const value = obj[key];

                if (_.isArray(value)) {
                    value.forEach(val => this.add(key, val));
                } else if (_.isString(value)) {
                    this.set(key, value);
                }
            });
        }
    }

    private normalizeKey(key: string): string {
        return key.replace('_', '-').toLowerCase();
    }

    public set(key: string, value: string) {
        this.containers[this.normalizeKey(key)] = [value];
    }

    public add(key: string, value: string) {
        key = this.normalizeKey(key);

        if (!this.containers[key]) {
            this.containers[key] = [];
        }

        this.containers[key].push(value);
    }

    public toString(): string {
        return JSON.stringify(this.containers);
    }
}

export interface TasksQueryParams {
    filters?: DockerFilterParams;
}

@Injectable()
export class DockerService {
    constructor(private apiService: ApiService) { }

    public getVersion(): Promise<VersionResponse> {
        return this.apiService.get('/docker/version').then(response => response.json() as VersionResponse);
    }

    public getInfo(): Promise<Info> {
        return this.apiService.get('/docker/info').then(response => response.json() as Info);
    }

    public getSwarm(): Promise<Swarm> {
        return this.apiService.get('/docker/swarm').then(response => response.json() as Swarm);
    }

    public getNodes(): Promise<Node[]> {
        return this.apiService.get('/docker/nodes').then(response => response.json() as Node[]);
    }

    public getNode(id: string): Promise<Node> {
        return this.apiService.get(`/docker/nodes/${id}`).then(response => response.json() as Node);
    }

    public getTasks(query?: TasksQueryParams): Promise<Task[]> {
        const params: URLSearchParams = new URLSearchParams();

        if (query) {
            Object.keys(query).map((key) =>  {
                let value = query[key];

                if (key === 'filters') {
                    value = (<DockerFilterParams>value).toString();
                }

                params.set(key, value);
            });
        }

        let url = '/docker/tasks';

        if (params.paramsMap.size > 0) {
            url += '?' + params.toString();
        }

        return this.apiService.get(url)
            .then(response => response.json() as Task[]);
    }
}
