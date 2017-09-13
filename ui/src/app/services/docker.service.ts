import { Injectable } from '@angular/core';

import { Version } from '../entities/version';

import { ApiService } from './api.service';

interface VersionResponse {
    ApiVersion: string;
    Arch: string;
    BuildTime: string;
    Experimental: boolean;
    GitCommit: string;
    GoVersion: string;
    KernelVersion: string;
    MinAPIVersion: string;
    Os: string;
    Version: string;
}

@Injectable()
export class DockerService {
    constructor(private apiService: ApiService) { }

    getVersion(): Promise<VersionResponse> {
        return this.apiService.get("/version").then(response => response.json() as VersionResponse);
    }
}
