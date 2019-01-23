import { PipeTransform, Pipe } from '@angular/core';
import { Service } from '../entities/docker';
import { StackNamePipe } from './stack-name.pipe';

@Pipe({
    name: 'serviceName',
    pure: true
})
export class ServiceNamePipe implements PipeTransform {
    constructor(private stackName: StackNamePipe) {}

    transform(service: Service): string {
        if (!service) {
            return '';
        }

        const stack = this.stackName.transform(service);

        return service.Spec.Name.substr(stack.length + 1);
    }
}
