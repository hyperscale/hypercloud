import { Injectable } from '@angular/core';
import { Node, Swarm, Info, VersionResponse, Task } from '../entities/docker';
import { ApiService } from './api.service';

@Injectable()
export class NodeService {
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

    public getTasks(): Promise<Task[]> {
        return this.apiService.get('/docker/tasks').then(response => response.json() as Task[]);
    }
}
