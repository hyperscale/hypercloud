import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { HttpModule } from '@angular/http';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './components/app/app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { CollectionComponent } from './components/collection/collection.component';
import { ClusterDetailComponent } from './components/cluster/cluster-detail.component';
import { ServiceListComponent } from './components/service/service-list.component';
import { ServiceCreateComponent } from './components/service/service-create.component';
import { ServiceDetailComponent } from './components/service/service-detail.component';
import { ServiceDetailOverviewComponent } from './components/service/service-detail-overview.component';
import { ServiceDetailMetricsComponent } from './components/service/service-detail-metrics.component';
import { ServiceDetailSettingsComponent } from './components/service/service-detail-settings.component';
import { ApplicationNameValidatorDirective } from './directives/application-name-validator.directive';
import { StackListComponent } from './components/stack/stack-list.component';
import { StackDetailComponent } from './components/stack/stack-detail.component';

import { ApiService } from './services/api.service';
import { StackService } from './services/stack.service';
import { ServiceService } from './services/service.service';
import { NodeService } from './services/node.service';
import { EventService } from './services/event.service';

import { TruncatePipe } from './pipes/truncate.pipe';
import { SizePipe } from './pipes/size.pipe';
import { ImagePipe } from './pipes/image.pipe';
import { ContainerPortPipe } from './pipes/container-port.pipe';
import { StackNamePipe } from './pipes/stack-name.pipe';
import { ServiceNamePipe } from './pipes/service-name.pipe';

import { ServiceResolver } from './resolvers/service.resolver';

@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        CollectionComponent,
        ClusterDetailComponent,
        ServiceListComponent,
        ServiceCreateComponent,
        ServiceDetailComponent,
        ServiceDetailOverviewComponent,
        ServiceDetailMetricsComponent,
        ServiceDetailSettingsComponent,
        ApplicationNameValidatorDirective,
        StackListComponent,
        StackDetailComponent,
        TruncatePipe,
        SizePipe,
        ImagePipe,
        ContainerPortPipe,
        StackNamePipe,
        ServiceNamePipe
    ],
    imports: [
        BrowserModule,
        BrowserAnimationsModule,
        FormsModule,
        AppRoutingModule,
        ClarityModule,
        HttpModule
    ],
    providers: [
        StackNamePipe,
        ApiService,
        StackService,
        ServiceService,
        NodeService,
        EventService,
        ServiceResolver
    ],
    bootstrap: [AppComponent]
})
export class AppModule { }
