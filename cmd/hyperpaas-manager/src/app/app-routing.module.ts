import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';
import { CollectionComponent } from './components/collection/collection.component';
import { NodeListComponent } from './components/node/node-list.component';
import { ServiceListComponent } from './components/service/service-list.component';
import { ServiceCreateComponent } from './components/service/service-create.component';
import { ServiceDetailComponent } from './components/service/service-detail.component';
import { StackListComponent } from './components/stack/stack-list.component';
import { StackDetailComponent } from './components/stack/stack-detail.component';

const routes: Routes = [
    { path: '', redirectTo: 'dashboard', pathMatch: 'full'},
    { path: 'dashboard', component: DashboardComponent },
    { path: 'stack', component: StackListComponent },
    { path: 'stack/:id', component: StackDetailComponent },
    { path: 'service', component: ServiceListComponent },
    { path: 'service/create', component: ServiceCreateComponent },
    { path: 'service/:id', component: ServiceDetailComponent },
    { path: 'node', component: NodeListComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
