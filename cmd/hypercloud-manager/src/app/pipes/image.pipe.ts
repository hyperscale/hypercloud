import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'image'
})
export class ImagePipe implements PipeTransform {
    transform(value: any, format?: any): any {
        if (!value) {
            return value;
        }

        switch (format) {
            case 'short':
                return value.split('@')[0];

            default:
                return value;
        }
    }
}
