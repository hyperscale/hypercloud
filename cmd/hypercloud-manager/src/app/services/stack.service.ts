import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Stack } from '../entities/docker';
import { ApiService } from './api.service';

@Injectable()
export class StackService {
    constructor(private apiService: ApiService) {}

    public getStacks(): Promise<Stack[]> {
        return this.apiService.get('/v1/stacks')
            .then(response => response.json())
            .then(response => {
                return response.map(stack => stack as Stack);
            });
    }

    public create(stack: Stack): Promise<Stack> {
        return this.apiService.post('/v1/stacks', stack)
            .then(response => response.json() as Stack);
    }
}
