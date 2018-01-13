import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'containerPort'
})
export class ContainerPortPipe implements PipeTransform {
    transform(value: any, args?: any): any {
        if (!value) {
            return '--';
        }

        console.log('containerPort Value:', value);

        return null;
    }

}
