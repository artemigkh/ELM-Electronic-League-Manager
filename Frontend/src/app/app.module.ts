import {BrowserModule} from '@angular/platform-browser';
import {RouterModule} from '@angular/router';
import {NgModule} from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
    MatButtonModule, MatButtonToggleModule, MatCardModule, MatCheckboxModule, MatChipsModule, MatDatepickerModule,
    MatDialogModule,
    MatDividerModule,
    MatExpansionModule,
    MatFormFieldModule, MatIconModule, MatIconRegistry,
    MatInputModule, MatNativeDateModule,
    MatSelectModule,
    MatTabsModule
} from '@angular/material';
import {MatTableModule} from '@angular/material/table';

import {NavBar} from "./shared/navbar/navbar";
import {AppComponent} from './app.component';
import {HomeComponent} from "./home/home";
import {MatchHistoryComponent} from "./matchHistory/match-history";
import {StandingsComponent} from "./standings/standings";
import {TeamsComponent} from "./teams/teams";
import {UpcomingGamesComponent} from "./upcomingGames/upcoming-games";

import {LeagueService} from './httpServices/leagues.service';

import {ELM_ROUTES} from './routes'
import {ManageComponent} from "./manage/manage";
import {ManageLeagueComponent} from "./manage/league/manage-league";
import {ManageTeamPopup, ManageTeamsComponent} from "./manage/teams/manage-teams";
import {ManagePermissionsComponent} from "./manage/permissions/manage-permissions";
import {ManageDatesComponent} from "./manage/dates/manage-dates";
import {ManagePlayersComponent} from "./manage/players/manage-players";
import {ManageGamePopup, ManageGamesComponent, ReportGamePopup} from "./manage/games/manage-games";
import {AmazingTimePickerModule} from "amazing-time-picker";
import {WarningPopup} from "./manage/warningPopup/warning-popup";

@NgModule({
    declarations: [
        AppComponent,
        StandingsComponent,
        HomeComponent,
        MatchHistoryComponent,
        TeamsComponent,
        MatchHistoryComponent,
        UpcomingGamesComponent,
        ManageComponent,
        ManageLeagueComponent,
        ManageTeamsComponent,
        ManagePermissionsComponent,
        ManageDatesComponent,
        ManagePlayersComponent,
        ManageGamesComponent,
        ReportGamePopup,
        ManageGamePopup,
        ManageTeamPopup,
        WarningPopup,
        NavBar
    ],
    imports: [
        BrowserModule,
        HttpClientModule,
        BrowserAnimationsModule,
        MatTabsModule,
        MatTableModule,
        MatButtonModule,
        MatCardModule,
        MatDividerModule,
        MatButtonToggleModule,
        MatExpansionModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
        MatCheckboxModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatIconModule,
        MatChipsModule,
        MatDialogModule,
        AmazingTimePickerModule,
        RouterModule.forRoot(ELM_ROUTES)
    ],
    providers: [LeagueService, MatIconRegistry],
    bootstrap: [AppComponent],
    entryComponents: [ReportGamePopup, ManageGamePopup, ManageTeamPopup, WarningPopup]
})
export class AppModule {
    constructor(public matIconRegistry: MatIconRegistry) {
        matIconRegistry.registerFontClassAlias('fontawesome', 'fa');
    }
}
