import { Component, OnInit } from '@angular/core';
import { Service } from '../../entities';
import { ActivatedRoute } from '@angular/router';

@Component({
    selector: 'app-service-detail-overview',
    templateUrl: './service-detail-overview.component.html',
    styleUrls: ['./service-detail-overview.component.less']
})
export class ServiceDetailOverviewComponent implements OnInit {
    service: Service;

    constructor(private route: ActivatedRoute) { }

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];
    }
}
