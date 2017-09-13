import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, Params} from '@angular/router';
import * as _ from 'lodash';

import { Task } from '../../entities/task';
import { Node } from '../../entities/node';
import { Service } from '../../entities/service';

import { TaskService } from '../../services/task.service';
import { NodeService } from '../../services/node.service';
import { ServiceService } from '../../services/service.service';

interface Env {
    name: string;
    value: string;
}

interface Label {
    name: string;
    value: string;
}

interface Constraint {
    name: string;
    operator: string;
    value: string;
}

@Component({
    selector: 'app-service-task',
    templateUrl: './service-task.component.html',
    styleUrls: ['./service-task.component.css']
})
export class ServiceTaskComponent implements OnInit {
    tasks: Task[];
    nodes: { [key: string]: Node };
    service: Service;
    envs: Env[];
    containerLabels: Label[];
    serviceLabels: Label[];
    constraints: Constraint[];

    constructor(
        private serviceService: ServiceService,
        private nodeService: NodeService,
        private taskService: TaskService,
        private activatedRoute: ActivatedRoute
    ) {
        this.tasks = [];
        this.nodes = {};
        this.service = null;
        this.envs = [];
        this.containerLabels = [];
        this.serviceLabels = [];
        this.constraints = [];
    }

    fetchTask(id) {
        this.nodeService.getNodes().then(nodes => {
            nodes.forEach((node) => {
                this.nodes[node.ID] = node;
            });

            return this.serviceService.getService(id);
        }).then(service => {
            this.service = service;

            this.envs = _.get(service, 'Spec.TaskTemplate.ContainerSpec.Env', []).map(item => {
                const s = item.split('=');

                return {
                    name: s[0],
                    value: s[1]
                };
            });

            this.constraints = <Constraint[]>ServiceService.translateConstraintsToKeyValue(
                _.get(service, 'Spec.TaskTemplate.Placement.Constraints', [])
            );

            Object.keys(_.get(service, 'Spec.TaskTemplate.ContainerSpec.Labels', {})).forEach(key => {
                this.containerLabels.push({
                    name: key,
                    value: service.Spec.Labels[key]
                });
            });

            Object.keys(_.get(service, 'Spec.Labels', {})).forEach(key => {
                this.serviceLabels.push({
                    name: key,
                    value: service.Spec.Labels[key]
                });
            });



            return this.taskService.getTasks(service.Spec.Name);
        }).then(tasks => {
            console.log(tasks);
            this.tasks = tasks;
        });
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            this.fetchTask(params['id']);
        });
    }

}
