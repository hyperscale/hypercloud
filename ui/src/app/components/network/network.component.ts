import { Component, OnInit, ViewChild } from '@angular/core';

import { DialogComponent } from '../dialog/dialog.component';

import { Network } from '../../entities/network';

import { NetworkService } from '../../services/network.service';

@Component({
    selector: 'app-network',
    templateUrl: './network.component.html',
    styleUrls: ['./network.component.css']
})
export class NetworkComponent implements OnInit {

    networks: Network[];

    network: Network;

    _openCreate: boolean;

    constructor(private networkService: NetworkService) {
        this.networks = [];

        this._openCreate = false;
    }

    onFetchNetworks() {
        this.networkService.getNetworks().then(networks => {
            console.log(networks);
            this.networks = networks;
        });
    }

    ngOnInit() {
        this.onFetchNetworks();
    }

    onAddNetwork() {
        this.networkService.addNetwork(this.network).then(network => {
            this.network = {};

            console.log('Network:', network);

            this.onFetchNetworks();
        });
    }

    onCancel() {
        this._openCreate = false;
    }

    onOpenCreateNetwork() {
        this._openCreate = true;
    }

}
