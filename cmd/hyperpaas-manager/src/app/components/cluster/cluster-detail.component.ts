import { Component, OnInit } from '@angular/core';
import { Node, Swarm, Info } from '../../entities/docker';
import { NodeService } from '../../services/node.service';

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

    constructor(private nodeService: NodeService) {
        this.nodes = [];

        this.cluster = {
            nodes: 0,
            cpu: 0,
            memory: 0,
            version: ''
        };
    }

    ngOnInit() {
        this.nodeService.getVersion().then(version => {
            this.cluster.version = version.ApiVersion;
        });

        this.nodeService.getInfo().then(info => {
            this.info = info;

            info.Swarm.RemoteManagers.forEach(item => {
                if (item.NodeID === info.Swarm.NodeID) {
                    this.advertiseAddr = item.Addr;
                }
            });
        });

        this.nodeService.getSwarm().then(swarm => this.swarm = swarm);

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
