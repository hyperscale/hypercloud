import { Component, OnInit, ViewChild } from '@angular/core';

import { Registry } from '../../entities/registry';

import { RegistryService } from '../../services/registry.service';

@Component({
    selector: 'app-registry',
    templateUrl: './registry.component.html',
    styleUrls: ['./registry.component.css']
})
export class RegistryComponent implements OnInit {

    registries: Registry[];

    registry: Registry;

    _openCreate: boolean;

    constructor(private registryService: RegistryService) {
        this.registries = [];
        this.registry = {};

        this._openCreate = false;
    }

    onFetchRegistries() {
        this.registryService.getRegistries().then(registries => {
            console.log(registries);
            this.registries = registries;
        });
    }

    ngOnInit() {
        this.onFetchRegistries();
    }

    onAddRegistry() {
        this.registryService.addRegistry(this.registry).then(registry => {
            this.registry = {};

            console.log('Registry:', registry);

            this.onFetchRegistries();
        });
    }

    onCancel() {
        this._openCreate = false;
    }

    onOpenCreateRegistry() {
        this._openCreate = true;
    }
}
