import { Directive, Input } from '@angular/core';
import { AbstractControl, NG_VALIDATORS, Validator, ValidatorFn, Validators } from '@angular/forms';

export function applicationNameValidator(): ValidatorFn {
    return (control: AbstractControl): {[key: string]: any} => {
        const isValid = /^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$/.test(control.value);

        return isValid ? null : {'applicationName': {value: control.value}};
    };
}

@Directive({
    selector: '[appApplicationNameValidator]',
    providers: [
        {
            provide: NG_VALIDATORS,
            useExisting: ApplicationNameValidatorDirective,
            multi: true
        }
    ]
})
export class ApplicationNameValidatorDirective implements Validator  {
    validate(control: AbstractControl): {[key: string]: any} {
        return applicationNameValidator()(control);
    }
}
