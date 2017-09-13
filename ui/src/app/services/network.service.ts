import { Injectable } from '@angular/core';

import { Network } from '../entities/network';

import { ApiService } from './api.service';

@Injectable()
export class NetworkService {
    constructor(private apiService: ApiService) { }

    getNetworks(): Promise<Network[]> {
        return this.apiService.get("/networks").then(response => response.json() as Network[]);
    }

    addNetwork(network: Network): Promise<Network> {
        return this.apiService.post("/networks", network).then(response => response.json() as Network);
    }

    getNetwork(id: string): Promise<Network> {
        return this.apiService.get(`/networks/${id}`).then(response => response.json() as Network);
    }
}
