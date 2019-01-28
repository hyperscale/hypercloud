import { Component, OnInit } from '@angular/core';
import { Node, Swarm, Info } from '../../entities';
import { DockerService } from '../../services';

interface ClusterInfo {
    cpu: number;
    memory: number;
    nodes: number;
    version: string;
}

@Component({
    selector: 'app-cluster-detail',
    templateUrl: './cluster-detail.component.html',
    styleUrls: ['./cluster-detail.component.less']
})
export class ClusterDetailComponent implements OnInit {
    public cluster: ClusterInfo;
    public nodes: Node[];
    public swarm: Swarm;
    public info: Info;
    public advertiseAddr: string;

    public _createNodeModalOpened = false;

    constructor(private dockerService: DockerService) {
        this.nodes = [];

        this.cluster = {
            nodes: 0,
            cpu: 0,
            memory: 0,
            version: ''
        };
    }

    ngOnInit() {
        this.dockerService.getVersion().subscribe(version => {
            this.cluster.version = version.ApiVersion;
        });

        this.dockerService.getInfo().subscribe(info => {
            this.info = info;

            info.Swarm.RemoteManagers.forEach(item => {
                if (item.NodeID === info.Swarm.NodeID) {
                    this.advertiseAddr = item.Addr;
                }
            });
        });

        this.dockerService.getSwarm().subscribe(swarm => this.swarm = swarm);

        this.dockerService.getNodes().subscribe(nodes => {
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
