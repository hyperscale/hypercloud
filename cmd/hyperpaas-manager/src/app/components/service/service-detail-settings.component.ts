import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Service } from '../../entities';

@Component({
    selector: 'app-service-detail-settings',
    templateUrl: './service-detail-settings.component.html',
    styleUrls: ['./service-detail-settings.component.less']
})
export class ServiceDetailSettingsComponent implements OnInit {
    service: Service;

    hosts: string[] = [];

    envs: { key: string, value: string }[] = [];

    constructor(private route: ActivatedRoute) {}

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];

        this.processService();
    }

    private processService() {
        const rule = this.service.Spec.Labels['traefik.frontend.rule'];

        this.hosts = rule.substring(5).split(',');

        this.envs = this.service.Spec.TaskTemplate.ContainerSpec.Env.map(item => {
            const parts = item.split('=');

            return {
                key: parts[0],
                value: parts[1],
            };
        });
    }

/*
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
*/
}
