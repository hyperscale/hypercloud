import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';

import { Node } from '../../entities/node';

import { NodeService } from '../../services/node.service';
import { DockerService } from '../../services/docker.service';

interface Label {
    name: string;
    value: string;
}

@Component({
    selector: 'app-cluster-node',
    templateUrl: './cluster-node.component.html',
    styleUrls: ['./cluster-node.component.css']
})
export class ClusterNodeComponent implements OnInit {
    node: Node;
    labels: Label[];

    constructor(private nodeService: NodeService, private dockerService: DockerService, private activatedRoute: ActivatedRoute) {
        this.labels = [];
    }

    private fetchNode(id: string) {
        this.nodeService.getNode(id).then(node => {
            this.node = node;

            Object.keys(node.Spec.Labels).forEach(key => {
                this.labels.push({
                    name: key,
                    value: node.Spec.Labels[key]
                });
            });
        });
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            console.log('ID:', params['id']);

            this.fetchNode(params['id']);
        });
    }
}
