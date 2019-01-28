import { PipeTransform, Pipe } from '@angular/core';
import { Service } from '../entities/docker';

@Pipe({
    name: 'stackName',
    pure: true
})
export class StackNamePipe implements PipeTransform {
    transform(service?: Service): string {
        if (!service) {
            return '';
        }

        return service.Spec.Labels['com.docker.stack.namespace'] || '';
    }
}
