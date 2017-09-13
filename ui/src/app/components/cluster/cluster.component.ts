import { Component, OnInit } from '@angular/core';

import { Node } from '../../entities/node';

import { NodeService } from '../../services/node.service';
import { DockerService } from '../../services/docker.service';

interface ClusterInfo {
    cpu: number;
    memory: number;
    nodes: number;
    version: string;
}

@Component({
    selector: 'app-cluster',
    templateUrl: './cluster.component.html',
    styleUrls: ['./cluster.component.css']
})
export class ClusterComponent implements OnInit {
    cluster: ClusterInfo;
    nodes: Node[];

    constructor(private nodeService: NodeService, private dockerService: DockerService) {
        this.nodes = [];
        this.cluster = {
            nodes: 0,
            cpu: 0,
            memory: 0,
            version: ''
        };
    }

    ngOnInit() {
        this.dockerService.getVersion().then(version => {
            this.cluster.version = version.ApiVersion;
        });

        this.nodeService.getNodes().then(nodes => {
            console.log(nodes);
            this.nodes = nodes;

            let cpu = 0;
            let memory = 0;

            nodes.forEach(node => {
                cpu += node.Description.Resources.NanoCPUs;
                memory += node.Description.Resources.MemoryBytes;
            });

            this.cluster = Object.assign(this.cluster, {
                nodes: nodes.length,
                cpu: cpu,
                memory: memory
            });
        });
    }
}
