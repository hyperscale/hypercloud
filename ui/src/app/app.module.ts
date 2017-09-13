import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { ClarityModule } from "clarity-angular";

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './components/app/app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { NetworkComponent } from './components/network/network.component';
import { NetworkViewComponent } from './components/network/network-view.component';
import { ServiceComponent } from './components/service/service.component';
import { ServiceTaskComponent } from './components/service/service-task.component';
import { NotFoundComponent } from './components/not-found/not-found.component';
import { PanelComponent } from './components/panel/panel.component';
import { ClusterComponent } from './components/cluster/cluster.component';
import { ClusterNodeComponent } from './components/cluster/cluster-node.component';
import { ImageComponent } from './components/image/image.component';
import { RegistryComponent } from './components/registry/registry.component';
import { ServiceLogComponent } from './components/service/service-log.component';
import { DialogComponent } from './components/dialog/dialog.component';
import { ConsoleComponent } from './components/console/console.component';
import { StackComponent } from './components/stack/stack.component';
import { LoginComponent } from './components/login/login.component';

import { NetworkService } from './services/network.service';
import { ApiService } from './services/api.service';
import { ServiceService } from './services/service.service';
import { ImageService } from './services/image.service';
import { SwarmService } from './services/swarm.service';
import { NodeService } from './services/node.service';
import { TaskService } from './services/task.service';
import { StackService } from './services/stack.service';
import { RegistryService } from './services/registry.service';
import { UserService } from './services/user.service';
import { AuthGuard } from './services/auth-guard.service';
import { AuthService } from './services/auth.service';
import { DockerService } from './services/docker.service';

import { TruncatePipe } from './pipes/truncate.pipe';
import { SizePipe } from './pipes/size.pipe';
import { ImagePipe } from './pipes/image.pipe';
import { ContainerPortPipe } from './pipes/container-port.pipe';



@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        NetworkComponent,
        NetworkViewComponent,
        ServiceComponent,
        NotFoundComponent,
        TruncatePipe,
        PanelComponent,
        ClusterComponent,
        ClusterNodeComponent,
        ImageComponent,
        ServiceTaskComponent,
        SizePipe,
        RegistryComponent,
        DialogComponent,
        ServiceLogComponent,
        ConsoleComponent,
        StackComponent,
        LoginComponent,
        ImagePipe,
        ContainerPortPipe
    ],
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        AppRoutingModule,
        BrowserAnimationsModule,
        ClarityModule.forRoot()
    ],
    providers: [
        NetworkService,
        ApiService,
        ServiceService,
        NodeService,
        ImageService,
        SwarmService,
        TaskService,
        StackService,
        RegistryService,
        UserService,
        AuthGuard,
        AuthService,
        DockerService
    ],
    schemas: [ NO_ERRORS_SCHEMA ],
    bootstrap: [AppComponent]
})
export class AppModule { }
