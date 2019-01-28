import { Component, OnInit } from '@angular/core';
import { Service, Task, Node } from '../../entities';
import { ActivatedRoute } from '@angular/router';
import { ServiceService, DockerService, DockerFilterParams } from '../../services';
import { switchMap } from 'rxjs/operators';

@Component({
    selector: 'app-service-detail-scalability',
    templateUrl: './service-detail-scalability.component.html',
    styleUrls: ['./service-detail-scalability.component.less']
})
export class ServiceDetailScalabilityComponent implements OnInit {
    service: Service;
    tasks: Task[] = [];
    nodes: { [key: string]: Node } = {};

    autoscale = {
        enable: false,
        min: 1,
        max: 5,
    };

    constructor(private serviceService: ServiceService, private dockerService: DockerService, private route: ActivatedRoute) { }

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];

        Object.keys(this.service.Spec.Labels).forEach(label => {
            const value = this.service.Spec.Labels[label];

            switch (label) {
                case 'com.hypercloud.autoscale.enable':
                    this.autoscale.enable = value === 'true';
                    break;
                case 'com.hypercloud.autoscale.min':
                    this.autoscale.min = parseInt(value, 10) || 1;
                    break;
                case 'com.hypercloud.autoscale.max':
                    this.autoscale.max = parseInt(value, 10) || 5;
                    break;
            }
        });

        this.dockerService.getNodes().pipe(
            switchMap(nodes => {
                nodes.forEach(node => this.nodes[node.ID] = node);

                return this.dockerService.getTasks({
                    filters: new DockerFilterParams({
                        service: this.service.Spec.Name,
                        desired_state: ['running', 'accepted']
                    })
                });
            })
        ).subscribe(tasks => {
            this.tasks = tasks;
        });
    }

    onScaleSubmit() {
        console.log('Scale', this.autoscale);

        this.service.Spec.Labels['com.hypercloud.autoscale.enable'] = this.autoscale.enable ? 'true' : 'false';
        this.service.Spec.Labels['com.hypercloud.autoscale.min'] = `${this.autoscale.min}`;
        this.service.Spec.Labels['com.hypercloud.autoscale.max'] = `${this.autoscale.max}`;

        if (this.autoscale.enable) {
            this.service.Spec.Mode.Replicated.Replicas = this.autoscale.min;
        } else {
            this.service.Spec.Mode.Replicated.Replicas = 1;
        }

        console.log('Service Request:', this.service);

        this.serviceService.update(this.service.ID, this.service).subscribe(service => {
            console.log('Service Response:', service);

            this.route.snapshot.parent.data['service'] = service;
        });
    }
}
