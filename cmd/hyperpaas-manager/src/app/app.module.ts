import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './components/app/app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';
import { CollectionComponent } from './components/collection/collection.component';
import { NodeComponent } from './components/node/node.component';
import { ServiceCreateComponent } from './components/service/service-create.component';
import { ApplicationNameValidatorDirective } from './directives/application-name-validator.directive';


@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        CollectionComponent,
        NodeComponent,
        ServiceCreateComponent,
        ApplicationNameValidatorDirective
    ],
    imports: [
        BrowserModule,
        FormsModule,
        AppRoutingModule,
        ClarityModule
    ],
    providers: [],
    bootstrap: [AppComponent]
})
export class AppModule { }
