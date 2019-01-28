import { Injectable } from '@angular/core';
import { HttpParams } from '@angular/common/http';
import { Node, Swarm, Info, VersionResponse, Task } from '../entities/docker';
import { ApiService } from './api.service';
import * as _ from 'lodash';
import { Observable } from 'rxjs';

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

    public getVersion(): Observable<VersionResponse> {
        return this.apiService.get<VersionResponse>('/docker/version');
    }

    public getInfo(): Observable<Info> {
        return this.apiService.get<Info>('/docker/info');
    }

    public getSwarm(): Observable<Swarm> {
        return this.apiService.get<Swarm>('/docker/swarm');
    }

    public getNodes(): Observable<Node[]> {
        return this.apiService.get<Node[]>('/docker/nodes');
    }

    public getNode(id: string): Observable<Node> {
        return this.apiService.get<Node>(`/docker/nodes/${id}`);
    }

    public getTasks(query?: TasksQueryParams): Observable<Task[]> {
        const params = new HttpParams();

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

        if (params.keys().length > 0) {
            url += '?' + params.toString();
        }

        return this.apiService.get<Task[]>(url);
    }
}
