import { Injectable } from '@angular/core';
import { Node, Swarm, Info } from '../entities/docker';
import { ApiService } from './api.service';

@Injectable()
export class NodeService {
    constructor(private apiService: ApiService) { }

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
}
