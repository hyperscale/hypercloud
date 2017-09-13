import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';

import { Network, NetworkContainer } from '../../entities/network';

import { NetworkService } from '../../services/network.service';

interface Option {
    name: string;
    value: string;
}

@Component({
    selector: 'app-network-view',
    templateUrl: './network-view.component.html',
    styleUrls: ['./network-view.component.css']
})
export class NetworkViewComponent implements OnInit {
    network: Network;
    options: Option[];
    containers: NetworkContainer[];

    constructor(private networkService: NetworkService, private activatedRoute: ActivatedRoute) {
        this.network = {};
        this.options = [];
        this.containers = [];
    }

    private fetchNetwork(id: string) {
        this.networkService.getNetwork(id).then(network => {
            this.network = network;

            Object.keys(network.Options).forEach(key => {
                this.options.push({
                    name: key,
                    value: network.Options[key]
                });
            });

            Object.keys(network.Containers).forEach(key => {
                this.containers.push(network.Containers[key]);
            });
        });
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            console.log('ID:', params['id']);

            this.fetchNetwork(params['id']);
        });
    }
}
