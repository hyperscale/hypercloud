import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs/Rx';

import { Service } from '../entities/service';

import { ApiService } from './api.service';

@Injectable()
export class ServiceService {
    constructor(private apiService: ApiService) {}

    getServices(): Promise<Service[]> {
        return this.apiService.get("/services").then(response => response.json() as Service[]);
    }

    getService(id: string): Promise<Service> {
        return this.apiService.get(`/services/${id}`).then(response => response.json() as Service);
    }

    static translateConstraintsToKeyValue(constraints: string[]): { [key: string]: any}[] {
        function getOperator(constraint) {
            var indexEquals = constraint.indexOf('==');
            if (indexEquals >= 0) {
                return [indexEquals, '=='];
            }

            return [constraint.indexOf('!='), '!='];
        }

        if (constraints) {
            var keyValueConstraints = [];
            constraints.forEach((constraint) => {
                const operatorIndices = getOperator(constraint);

                keyValueConstraints.push({
                    name: constraint.slice(0, operatorIndices[0]),
                    value: constraint.slice(operatorIndices[0] + 2),
                    operator: operatorIndices[1]
                });
            });

            return keyValueConstraints;
        }

        return [];
    }
}
