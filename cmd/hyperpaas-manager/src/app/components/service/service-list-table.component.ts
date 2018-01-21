import { Component, OnInit, Input } from '@angular/core';
import { Service, ServiceInfo } from '../../entities';
import { ServiceService, DockerService } from '../../services';

@Component({
    selector: 'app-service-list-table',
    templateUrl: './service-list-table.component.html',
    styleUrls: ['./service-list-table.component.less']
})
export class ServiceListTableComponent implements OnInit {
    @Input()
    public services: Service[];

    @Input()
    public stack?: string;

    public serviceInfo: { [key: string]: ServiceInfo };

    constructor(private serviceService: ServiceService, private dockerService: DockerService) {
        this.serviceInfo = {};
    }

    ngOnInit() {
        const running: { [key: string]: number } = {};
        const tasksNoShutdown: { [key: string]: number } = {};
        const activeNodes: { [key: string]: boolean } = {};

        this.dockerService.getNodes().then(nodes => {
            nodes.forEach(node => {
                if (node.Status.State !== 'down') {
                    activeNodes[node.ID] = true;
                }
            });

            return nodes;
        }).then(nodes => {
            return this.dockerService.getTasks();
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
                                Running: running[service.ID] || 0,
                                Desired: service.Spec.Mode.Replicated.Replicas,
                            }
                        };
                    } else if (service.Spec.Mode.Global !== undefined) {
                        this.serviceInfo[service.ID] = {
                            Mode: 'global',
                            Replicas: {
                                Running: running[service.ID] || 0,
                                Desired: tasksNoShutdown[service.ID],
                            }
                        };
                    }
                });

            });

        });
    }
}
