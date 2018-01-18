import { Component, OnInit } from '@angular/core';
import { Service, ServiceInfo } from '../../entities/docker';
import { ServiceService } from '../../services/service.service';
import { NodeService } from '../../services/node.service';
import { unescape } from 'querystring';

@Component({
    selector: 'app-service-list',
    templateUrl: './service-list.component.html',
    styleUrls: ['./service-list.component.less']
})
export class ServiceListComponent implements OnInit {
    public services: Service[] = [];
    public serviceInfo: { [key: string]: ServiceInfo };

    constructor(private serviceService: ServiceService, private nodeService: NodeService) {
        this.serviceInfo = {};
    }

    ngOnInit() {
        const running: { [key: string]: number } = {};
        const tasksNoShutdown: { [key: string]: number } = {};
        const activeNodes: { [key: string]: boolean } = {};

        this.serviceService.getServices().then(services => {
            this.services = services;

            return services;
        }).then(services => {
            return this.nodeService.getNodes();
        }).then(nodes => {
            nodes.forEach(node => {
                if (node.Status.State !== 'down') {
                    activeNodes[node.ID] = true;
                }
            });

            return nodes;
        }).then(nodes => {
            return this.nodeService.getTasks();
        }).then(tasks => {
            tasks.forEach(task => {
                if (task.DesiredState !== 'shutdown') {
                    if (!tasksNoShutdown.hasOwnProperty(task.ServiceID)) {
                        tasksNoShutdown[task.ServiceID] = 0;
                    }

                    tasksNoShutdown[task.ServiceID]++;
                }

                if (activeNodes.hasOwnProperty(task.NodeID) && task.Status.State === 'running') {
                    if (!running.hasOwnProperty(task.ServiceID)) {
                        running[task.ServiceID] = 0;
                    }

                    running[task.ServiceID]++;
                }

                this.services.forEach(service => {
                    if (service.Spec.Mode.Replicated !== undefined && service.Spec.Mode.Replicated.Replicas !== undefined) {
                        this.serviceInfo[service.ID] = {
                            Mode: 'replicated',
                            Replicas: {
                                Running: running[service.ID],
                                Desired: service.Spec.Mode.Replicated.Replicas,
                            }
                        };
                    } else if (service.Spec.Mode.Global !== unescape) {
                        this.serviceInfo[service.ID] = {
                            Mode: 'global',
                            Replicas: {
                                Running: running[service.ID],
                                Desired: tasksNoShutdown[service.ID],
                            }
                        };
                    }
                });

            });


        });
    }
}
