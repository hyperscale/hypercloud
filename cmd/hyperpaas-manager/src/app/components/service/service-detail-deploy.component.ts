import { Component, OnInit } from '@angular/core';
import { Service } from '../../entities';
import { ActivatedRoute } from '@angular/router';

@Component({
    selector: 'app-service-detail-deploy',
    templateUrl: './service-detail-deploy.component.html',
    styleUrls: ['./service-detail-deploy.component.less']
})
export class ServiceDetailDeployComponent implements OnInit {
    service: Service;

    constructor(private route: ActivatedRoute) { }

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];
    }
}
