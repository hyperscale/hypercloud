import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Stack } from '../entities/stack';
import { ApiService } from './api.service';

@Injectable()
export class StackService {
    constructor(private apiService: ApiService) {}

    public getStack(id: string): Promise<Stack> {
        return this.apiService.get(`/v1/stacks/${id}`)
        .then(response => response.json())
        .then(response => Object.assign(new Stack(), response));
    }

    public getStacks(): Promise<Stack[]> {
        return this.apiService.get('/v1/stacks')
            .then(response => response.json())
            .then(response => {
                return response.map(stack => Object.assign(new Stack(), stack));
            });
    }

    public create(stack: Stack): Promise<Stack> {
        return this.apiService.post('/v1/stacks', stack)
            .then(response => response.json())
            .then(response => Object.assign(new Stack(), response));
    }
}
