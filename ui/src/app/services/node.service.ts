import { Injectable } from '@angular/core';

import { Node } from '../entities/node';

import { ApiService } from './api.service';

@Injectable()
export class NodeService {
    constructor(private apiService: ApiService) { }

    getNodes(): Promise<Node[]> {
        return this.apiService.get('/nodes').then(response => response.json() as Node[]);
    }

    getNode(id: string): Promise<Node> {
        return this.apiService.get(`/nodes/${id}`).then(response => response.json() as Node);
    }
}
