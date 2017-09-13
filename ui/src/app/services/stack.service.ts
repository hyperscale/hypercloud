import { Injectable } from '@angular/core';

import { Stack } from '../entities/stack';

import { ApiService } from './api.service';


@Injectable()
export class StackService {
    constructor(private apiService: ApiService) { }

    getStacks(): Promise<Stack[]> {
        return this.apiService.get("/stacks").then(response => response.json() as Stack[]);
    }

    addStack(stack: Stack): Promise<Stack> {
        return this.apiService.post("/stacks", stack).then(response => response.json() as Stack);
    }
}
