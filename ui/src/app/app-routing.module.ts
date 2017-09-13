import { NgModule }              from '@angular/core';
import { RouterModule, Routes }  from '@angular/router';

import { DashboardComponent }   from './components/dashboard/dashboard.component';
import { NetworkComponent } from './components/network/network.component';
import { NetworkViewComponent } from './components/network/network-view.component';
import { ServiceComponent } from './components/service/service.component';
import { NotFoundComponent } from './components/not-found/not-found.component';
import { ClusterComponent } from './components/cluster/cluster.component';
import { ClusterNodeComponent } from './components/cluster/cluster-node.component';
import { ImageComponent } from './components/image/image.component';
import { ServiceTaskComponent } from './components/service/service-task.component';
import { RegistryComponent } from './components/registry/registry.component';
import { ServiceLogComponent } from './components/service/service-log.component';
import { StackComponent } from './components/stack/stack.component';
import { LoginComponent } from './components/login/login.component';

import { AuthGuard } from './services/auth-guard.service';

const appRoutes: Routes = [
    { path: 'login', component: LoginComponent },
    { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'network', component: NetworkComponent, canActivate: [AuthGuard] },
    { path: 'network/:id', component: NetworkViewComponent, canActivate: [AuthGuard] },
    { path: 'service', component: ServiceComponent, canActivate: [AuthGuard] },
    { path: 'service/:id/task', component: ServiceTaskComponent, canActivate: [AuthGuard] },
    { path: 'service/:id/log', component: ServiceLogComponent, canActivate: [AuthGuard] },
    { path: 'cluster', component: ClusterComponent, canActivate: [AuthGuard] },
    { path: 'cluster/node/:id', component: ClusterNodeComponent, canActivate: [AuthGuard] },
    { path: 'image', component: ImageComponent, canActivate: [AuthGuard] },
    { path: 'registry', component: RegistryComponent, canActivate: [AuthGuard] },
    { path: 'stack', component: StackComponent, canActivate: [AuthGuard] },
    { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
    { path: '**', component: NotFoundComponent }
];

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule {}
