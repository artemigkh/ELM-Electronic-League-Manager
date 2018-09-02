import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatTabsModule} from '@angular/material';
import {MatTableModule} from '@angular/material/table';

import { AppComponent } from './app.component';
import { StandingsComponent } from '../standings/standings.component';

import {LeagueService} from '../httpServices/leagues.service';

@NgModule({
  declarations: [
    AppComponent,
    StandingsComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTabsModule,
    MatTableModule
  ],
  providers: [LeagueService],
  bootstrap: [AppComponent]
})
export class AppModule { }
