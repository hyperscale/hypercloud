import { Component, OnInit } from '@angular/core';
import { Node, Swarm, Info } from '../../entities/docker';
import { NodeService } from '../../services/node.service';

@Component({
    selector: 'app-node-list',
    templateUrl: './node-list.component.html',
    styleUrls: ['./node-list.component.less']
})
export class NodeListComponent implements OnInit {
    public nodes: Node[] = [];
    public swarm: Swarm;
    public info: Info;
    public advertiseAddr: string;

    public _createNodeModalOpened = false;

    constructor(private nodeService: NodeService) { }

    ngOnInit() {
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

            /*
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
            */
        });
    }
}
