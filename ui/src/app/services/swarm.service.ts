import { Injectable } from '@angular/core';

import { Swarm } from '../entities/swarm';

import { ApiService } from './api.service';

@Injectable()
export class SwarmService {
    constructor(private apiService: ApiService) { }

    getInfo(): Promise<Swarm> {
        return this.apiService.get("/swarm").then(response => response.json() as Swarm);
    }
}
