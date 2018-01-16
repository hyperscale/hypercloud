import { Component, OnInit, OnDestroy, Input, ChangeDetectionStrategy, keyframes } from '@angular/core';
import { Service } from '../../entities/docker';
import { ServiceService } from '../../services/service.service';
import { Subscription } from 'rxjs/Subscription';
import { Subject } from 'rxjs/Subject';

interface Memory {
    usage: number;
    limit: number;
    percent: number;
}

interface Stats {
    memory: Memory;
}

@Component({
    selector: 'app-service-detail-metrics',
    templateUrl: './service-detail-metrics.component.html',
    styleUrls: ['./service-detail-metrics.component.less'],
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class ServiceDetailMetricsComponent implements OnInit, OnDestroy {
    @Input()
    service: Service;

    stats$ = new Subject<Stats>();

    metrics: { [key: string]: Stats } = {};

    private subscriptions: Subscription[] = [];

    constructor(private serviceService: ServiceService) {}

    aggregate() {
        const keys = Object.keys(this.metrics);
        const size = keys.length;

        console.log('keys:', keys);

        const stats: Stats = {
            memory: {
                usage: 0,
                limit: 0,
                percent: 0
            }
        };

        keys.forEach(id => {
            const task = this.metrics[id];
            stats.memory.usage += task.memory.usage;
            stats.memory.limit += task.memory.limit;
        });

        stats.memory.usage = (stats.memory.usage / size);
        stats.memory.limit = (stats.memory.limit / size);
        stats.memory.percent = ((stats.memory.usage / stats.memory.limit) * 100);

        console.log('Percent:', stats.memory.percent);

        this.stats$.next(stats);
    }

    ngOnInit() {
        this.subscriptions.push(
            this.serviceService.stats(this.service.ID).subscribe(event => {
                console.log('Stats:', event);

                this.metrics[event.id] = {
                    memory: {
                        usage: event.memory_stats.usage,
                        limit: event.memory_stats.limit,
                        percent: 0,
                    }
                };

                console.log('metrics:', this.metrics);

                this.aggregate();
            })
        );
    }

    ngOnDestroy() {
        this.subscriptions.forEach(subscription => subscription.unsubscribe());
        this.stats$.complete();
    }
}
