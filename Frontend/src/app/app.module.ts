import {BrowserModule} from '@angular/platform-browser';
import {RouterModule} from '@angular/router';
import {NgModule} from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
    DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE,
    MatButtonModule, MatButtonToggleModule, MatCardModule, MatCheckboxModule, MatChipsModule, MatDatepickerModule,
    MatDialogModule,
    MatDividerModule,
    MatExpansionModule,
    MatFormFieldModule, MatIconModule, MatIconRegistry,
    MatInputModule, MatNativeDateModule, MatProgressSpinnerModule,
    MatSelectModule, MatSnackBarModule, MatSpinner, MatStepperModule,
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
import {WarningPopup} from "./manage/warningPopup/warning-popup";
import {LoginComponent} from "./login/login";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {TeamsService} from "./httpServices/teams.service";
import {PlayersService} from "./httpServices/players.service";
import {SignupComponent} from "./signup/signup";
import {GameEntry} from "./shared/gameEntry/game-entry";
import {GamesComponent} from "./games/games";
import {GamesService} from "./httpServices/games.service";
import {NgxMaterialTimepickerModule} from "ngx-material-timepicker";
import {MAT_MOMENT_DATE_FORMATS, MatMomentDateModule, MomentDateAdapter} from "@angular/material-moment-adapter";
import {ConfirmationComponent} from "./shared/confirmation/confirmation-component";
import {TeamEntry} from "./shared/teamEntry/team-entry";
import {TournamentRegistrationComponent} from "./tournamentRegistration/tournament-registration";
import {ManagePlayersPopup, ManagePlayersTeamComponent} from "./manage/players/manage-players-team/manage-players-team";
import {UserService} from "./httpServices/user.service";
import {LeaguesComponent} from "./leagues/leagues";

@NgModule({
    declarations: [
        AppComponent,
        StandingsComponent,
        HomeComponent,
        MatchHistoryComponent,
        TeamsComponent,
        MatchHistoryComponent,
        UpcomingGamesComponent,
        GamesComponent,
        ManageComponent,
        ManageLeagueComponent,
        ManageTeamsComponent,
        ManagePermissionsComponent,
        ManageDatesComponent,
        ManagePlayersComponent,
        ManageGamesComponent,
        LeaguesComponent,
        ReportGamePopup,
        ManageGamePopup,
        ManageTeamPopup,
        WarningPopup,
        ManagePlayersPopup,
        ConfirmationComponent,
        LoginComponent,
        SignupComponent,
        TournamentRegistrationComponent,
        ManagePlayersTeamComponent,
        NavBar,
        GameEntry,
        TeamEntry
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
        FormsModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
        MatMomentDateModule,
        MatStepperModule,
        NgxMaterialTimepickerModule.forRoot(),
        RouterModule.forRoot(ELM_ROUTES)
    ],
    providers: [LeagueService, TeamsService, PlayersService, GamesService, UserService, MatIconRegistry],
    bootstrap: [AppComponent],
    entryComponents: [
        ReportGamePopup,
        ManageGamePopup,
        ManageTeamPopup,
        WarningPopup,
        ManagePlayersPopup,
        ConfirmationComponent
    ]
})
export class AppModule {
    constructor(public matIconRegistry: MatIconRegistry) {
        matIconRegistry.registerFontClassAlias('fontawesome', 'fa');
    }
}
