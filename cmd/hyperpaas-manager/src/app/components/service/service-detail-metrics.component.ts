import { Component, OnInit, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Service, StatsJSON, NetworkStats } from '../../entities/docker';
import { ServiceService } from '../../services/service.service';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';
import * as _ from 'lodash';

interface StatsMemory {
    usage: number;
    limit: number;
    percent: number;
}

interface StatsNetwork {
    rx: number;
    tx: number;
}

interface StatsCPU {
    percent: number;
}

interface Stats {
    date: string;
    memory: StatsMemory;
    cpu: StatsCPU;
    task: number;
    network: StatsNetwork;
}

@Component({
    selector: 'app-service-detail-metrics',
    templateUrl: './service-detail-metrics.component.html',
    styleUrls: ['./service-detail-metrics.component.less'],
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class ServiceDetailMetricsComponent implements OnInit, OnDestroy {
    service: Service;

    stats$ = new Subject<Stats>();

    metrics: { [key: string]: Stats } = {};

    private previousStats: Stats;

    private subscriptions: Subscription[] = [];

    constructor(private serviceService: ServiceService, private route: ActivatedRoute) {}

    aggregate() {
        const keys = Object.keys(this.metrics);
        const size = keys.length;

        const stats: Stats = {
            date: '',
            memory: {
                usage: 0,
                limit: 0,
                percent: 0
            },
            cpu: {
                percent: 0
            },
            task: size,
            network: {
                tx: 0,
                rx: 0
            }
        };

        keys.forEach(id => {
            const task = this.metrics[id];
            stats.memory.usage += task.memory.usage;
            stats.memory.limit += task.memory.limit;

            stats.cpu.percent += task.cpu.percent;

            stats.network.rx += task.network.rx;
            stats.network.tx += task.network.tx;
        });

        stats.memory.usage = (stats.memory.usage / size);
        stats.memory.limit = (stats.memory.limit / size);
        stats.memory.percent = ((stats.memory.usage / stats.memory.limit) * 100);

        stats.cpu.percent = (stats.cpu.percent / size);

        stats.network.rx = (stats.network.rx / size);
        stats.network.tx = (stats.network.tx / size);

        if (!_.isEqual(this.previousStats, stats)) {
            this.stats$.next(stats);
        }
    }

    ngOnInit() {
        this.service = this.route.snapshot.parent.data['service'];

        this.subscriptions.push(
            this.serviceService.stats(this.service.ID).subscribe(event => {
                console.log('Stats:', event);

                this.metrics[event.id] = {
                    date: event.read,
                    memory: {
                        usage: (event.memory_stats.usage - event.memory_stats.stats['cache']),
                        limit: event.memory_stats.limit,
                        percent: 0,
                    },
                    task: 0,
                    cpu: {
                        percent: this.calculateCPUPercent(event),
                    },
                    network: this.calculateNetwork(event.networks),
                };

                this.aggregate();
            })
        );
    }

    ngOnDestroy() {
        this.subscriptions.forEach(subscription => subscription.unsubscribe());
        this.stats$.complete();
    }

    private calculateCPUPercent(stats: StatsJSON): number {
        const cpuDelta = (stats.cpu_stats.cpu_usage.total_usage - stats.precpu_stats.cpu_usage.total_usage);
        const systemDelta = (stats.cpu_stats.system_cpu_usage - stats.precpu_stats.system_cpu_usage);

        let onlineCPUs = stats.cpu_stats.online_cpus;
        let cpuPercent = 0.0;

        if (onlineCPUs === 0) {
            onlineCPUs = stats.cpu_stats.cpu_usage.percpu_usage.length;
        }

        if (systemDelta > 0 && cpuDelta > 0) {
            cpuPercent = ((cpuDelta / systemDelta) * onlineCPUs * 100);
        }

        return cpuPercent;
    }

    private calculateNetwork(networks: { [key: string]: NetworkStats }): { tx: number, rx: number } {
        const stats = {
            rx: 0,
            tx: 0
        };

        Object.keys(networks).forEach(key => {
            const network = networks[key];

            stats.rx += network.rx_bytes;
            stats.tx += network.tx_bytes;
        });

        return stats;
    }
}
