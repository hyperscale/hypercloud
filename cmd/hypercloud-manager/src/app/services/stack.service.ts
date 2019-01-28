import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Stack } from '../entities/docker';
import { ApiService } from './api.service';

@Injectable()
export class StackService {
    constructor(private apiService: ApiService) {}

    public getStacks(): Observable<Stack[]> {
        return this.apiService.get<Stack[]>('/v1/stacks');
    }

    public create(stack: Stack): Observable<Stack> {
        return this.apiService.post<Stack>('/v1/stacks', stack);
    }
}
