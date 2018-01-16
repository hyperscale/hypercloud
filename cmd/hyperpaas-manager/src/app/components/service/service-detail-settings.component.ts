import { Component, OnInit, Input, OnChanges, SimpleChanges, SimpleChange } from '@angular/core';
import { Service } from '../../entities/docker';

@Component({
    selector: 'app-service-detail-settings',
    templateUrl: './service-detail-settings.component.html',
    styleUrls: ['./service-detail-settings.component.less']
})
export class ServiceDetailSettingsComponent implements OnInit, OnChanges {
    @Input()
    service: Service;

    hosts: string[] = [];

    envs: { key: string, value: string }[] = [];

    constructor() {}

    ngOnInit() {

    }

    ngOnChanges(changes: SimpleChanges) {
        const serviceChange: SimpleChange = changes.service;

        if (serviceChange.currentValue) {
            const service = (<Service>serviceChange.currentValue);

            const rule = service.Spec.Labels['traefik.frontend.rule'];

            this.hosts = rule.substring(5).split(',');

            this.envs = service.Spec.TaskTemplate.ContainerSpec.Env.map(item => {
                const parts = item.split('=');

                return {
                    key: parts[0],
                    value: parts[1],
                };
            });
        }
    }
}
