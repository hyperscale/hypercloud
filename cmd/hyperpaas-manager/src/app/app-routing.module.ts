import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';
import { CollectionComponent } from './components/collection/collection.component';
import { ClusterDetailComponent } from './components/cluster/cluster-detail.component';
import { ServiceListComponent } from './components/service/service-list.component';
import { ServiceCreateComponent } from './components/service/service-create.component';
import { ServiceDetailComponent } from './components/service/service-detail.component';
import { ServiceDetailOverviewComponent } from './components/service/service-detail-overview.component';
import { ServiceDetailMetricsComponent } from './components/service/service-detail-metrics.component';
import { ServiceDetailSettingsComponent } from './components/service/service-detail-settings.component';
import { ServiceDetailDeployComponent } from './components/service/service-detail-deploy.component';
import { StackListComponent } from './components/stack/stack-list.component';
import { StackDetailComponent } from './components/stack/stack-detail.component';

import { ServiceResolver } from './resolvers/service.resolver';

const routes: Routes = [
    { path: '', redirectTo: 'dashboard', pathMatch: 'full'},
    { path: 'dashboard', component: DashboardComponent },
    { path: 'stack', component: StackListComponent },
    { path: 'stack/:id', component: StackDetailComponent },
    { path: 'service', component: ServiceListComponent },
    { path: 'service/create', component: ServiceCreateComponent },
    {
        path: 'service/:id',
        component: ServiceDetailComponent,
        resolve: {
            'service': ServiceResolver,
        },
        children: [
            { path: '', redirectTo: 'overview', pathMatch: 'full' },
            { path: 'overview', component: ServiceDetailOverviewComponent },
            { path: 'deploy', component: ServiceDetailDeployComponent },
            { path: 'metrics', component: ServiceDetailMetricsComponent },
            { path: 'settings', component: ServiceDetailSettingsComponent },
        ]
    },
    { path: 'cluster', component: ClusterDetailComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
