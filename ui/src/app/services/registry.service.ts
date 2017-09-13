import { Injectable } from '@angular/core';

import { Registry } from '../entities/registry';

import { ApiService } from './api.service';

@Injectable()
export class RegistryService {
    constructor(private apiService: ApiService) { }

    getRegistries(): Promise<Registry[]> {
        return this.apiService.get("/registries").then(response => response.json() as Registry[]);
    }

    addRegistry(registry: Registry): Promise<Registry> {
        return this.apiService.post("/registries", registry).then(response => response.json() as Registry);
    }
}
