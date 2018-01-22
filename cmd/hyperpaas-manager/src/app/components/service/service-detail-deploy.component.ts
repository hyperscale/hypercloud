import { Component, OnInit } from '@angular/core';
import { Service } from '../../entities';
import { ActivatedRoute } from '@angular/router';
import { ServiceService } from '../../services';

@Component({
    selector: 'app-service-detail-deploy',
    templateUrl: './service-detail-deploy.component.html',
    styleUrls: ['./service-detail-deploy.component.less']
})
export class ServiceDetailDeployComponent implements OnInit {
    service: Service;

    deploy_method: string;

    docker = {
        image: ''
    };

    constructor(private serviceService: ServiceService, private route: ActivatedRoute) { }

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];

        this.deploy_method = this.service.Spec.Labels['com.hyperpaas.deploy.method'] || '';

        switch (this.deploy_method) {
            case 'docker':
                this.docker.image = this.service.Spec.TaskTemplate.ContainerSpec.Image;
        }
    }

    onDeployDockerSubmit() {
        console.log('Docker Image Ref:', this.docker.image);

        this.service.Spec.Labels['com.hyperpaas.deploy.method'] = 'docker';
        this.service.Spec.Labels['com.docker.stack.image'] = this.docker.image;
        this.service.Spec.TaskTemplate.ContainerSpec.Image = this.docker.image;

        console.log('Service Request:', this.service);

        this.serviceService.update(this.service.ID, this.service).then(service => {
            console.log('Service Response:', service);

            this.route.snapshot.parent.data['service'] = service;
        });
    }
}
