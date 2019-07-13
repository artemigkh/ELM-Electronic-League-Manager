// import {BrowserModule} from '@angular/platform-browser';
// import {RouterModule} from '@angular/router';
// import {NgModule} from '@angular/core';
// import {HttpClientModule} from '@angular/common/http';
// import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
// import {
//     DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE,
//     MatButtonModule, MatButtonToggleModule, MatCardModule, MatCheckboxModule, MatChipsModule, MatDatepickerModule,
//     MatDialogModule,
//     MatDividerModule,
//     MatExpansionModule,
//     MatFormFieldModule, MatIconModule, MatIconRegistry,
//     MatInputModule, MatNativeDateModule, MatProgressSpinnerModule,
//     MatSelectModule, MatSlideToggleModule, MatSnackBarModule, MatSpinner, MatStepperModule,
//     MatTabsModule
// } from '@angular/material';
// import {MatTableModule} from '@angular/material/table';
//
// import {NavBar} from "./shared/navbar/navbar";
// import {AppComponent} from './app.component';
// import {HomeComponent} from "./home/home";
// import {StandingsComponent} from "./standings/standings";
// import {TeamsComponent} from "./teams/teams";
//
// import {LeagueService} from './httpServices/leagues.service';
//
// import {ELM_ROUTES} from './routes'
// import {ManageComponent} from "./manage/manage";
// import {ManageLeagueComponent} from "./manage/league/manage-league";
// import {ManageTeamPopup, ManageTeamsComponent} from "./manage/teams/manage-teams";
// import {ManagePermissionsComponent} from "./manage/permissions/manage-permissions";
// import {ManageDatesComponent} from "./manage/dates/manage-dates";
// import {ManagePlayersComponent} from "./manage/players/manage-players";
// import {ManageGamePopup, ManageGamesComponent, ReportGamePopup} from "./manage/games/manage-games";
// import {WarningPopup} from "./manage/warningPopup/warning-popup";
// import {LoginComponent} from "./login/login";
// import {FormsModule, ReactiveFormsModule} from "@angular/forms";
// import {TeamsService} from "./httpServices/teams.service";
// import {PlayersService} from "./httpServices/players.service";
// import {SignupComponent} from "./signup/signup";
// import {GameEntry} from "./shared/gameEntry/game-entry";
// import {GamesComponent} from "./games/games";
// import {GamesService} from "./httpServices/games.service";
// import {NgxMaterialTimepickerModule} from "ngx-material-timepicker";
// import {MAT_MOMENT_DATE_FORMATS, MatMomentDateModule, MomentDateAdapter} from "@angular/material-moment-adapter";
// import {ConfirmationComponent} from "./shared/confirmation/confirmation-component";
// import {TeamEntry} from "./shared/teamEntry/team-entry";
// import {TournamentRegistrationComponent} from "./tournamentRegistration/tournament-registration";
// import {
//     ManagePlayersPopup,
//     ManagePlayersTeamComponent
// } from "./manage/players/manage-players-team/manage-players-team";
// import {UserService} from "./httpServices/user.service";
// import {LeaguesComponent} from "./leagues/leagues";
// import {LeagueCreationComponent} from "./leagueCreation/league-creation";
// import {StatsComponent} from "./stats/stats";
// import {GenericPlayerEntry} from "./teams/playerEntry/generic-player-entry";
// import {PlayerEntryComponent, PlayerEntryDirective} from "./teams/playerEntry/player-entry";
// import {LeagueOfLegendsPlayerEntry} from "./teams/playerEntry/league-of-legends-player-entry";
// import {
//     ManagePlayersPopupLeagueOfLegends,
//     ManagePlayersTeamLeagueOfLegendsComponent
// } from "./manage/players/manage-players-team/league-of-legends/manage-players-team-league-of-legends";
// import {
//     ManagePlayersTeamContainerComponent,
//     ManagePlayersTeamDirective
// } from "./manage/players/manage-players-team/manage-players-team-loader";
// import {RulesComponent} from "./rules/rules";
// import { MarkdownModule } from 'ngx-markdown';
// import { LMarkdownEditorModule } from 'ngx-markdown-editor';
// import {ManageRulesComponent} from "./manage/rules/manage-rules";
// import {ManageAvailabilityPopup, ManageScheduleComponent} from "./manage/schedule/manage-schedule";
// import {StatsPageComponent, StatsPageDirective} from "./stats/statsPage/stats-page";
// import {GenericStatsPage} from "./stats/statsPage/generic-stats-page";
// import {LeagueOfLegendsStatsPage} from "./stats/statsPage/league-of-legends-stats-page";
// import {LeagueOfLegendsStatsService} from "./httpServices/league-of-legends-stats.service";
// import {ElmState} from "./shared/state/state.service";
//
// @NgModule({
//     declarations: [
//         AppComponent,
//         StandingsComponent,
//         HomeComponent,
//         TeamsComponent,
//         GamesComponent,
//         ManageComponent,
//         ManageLeagueComponent,
//         ManageTeamsComponent,
//         ManagePermissionsComponent,
//         ManageDatesComponent,
//         ManagePlayersComponent,
//         ManageGamesComponent,
//         ManageScheduleComponent,
//         LeaguesComponent,
//         StatsComponent,
//         LeagueCreationComponent,
//         ReportGamePopup,
//         ManageGamePopup,
//         ManageTeamPopup,
//         WarningPopup,
//         ManagePlayersPopup,
//         ConfirmationComponent,
//         LoginComponent,
//         SignupComponent,
//         TournamentRegistrationComponent,
//         ManagePlayersTeamComponent,
//         NavBar,
//         GameEntry,
//         TeamEntry,
//         GenericPlayerEntry,
//         GenericStatsPage,
//         PlayerEntryComponent,
//         StatsPageComponent,
//         PlayerEntryDirective,
//         LeagueOfLegendsPlayerEntry,
//         LeagueOfLegendsStatsPage,
//         ManagePlayersTeamDirective,
//         ManagePlayersTeamContainerComponent,
//         ManagePlayersTeamLeagueOfLegendsComponent,
//         ManagePlayersPopupLeagueOfLegends,
//         RulesComponent,
//         ManageRulesComponent,
//         ManageAvailabilityPopup,
//         StatsPageDirective
//     ],
//     imports: [
//         BrowserModule,
//         HttpClientModule,
//         BrowserAnimationsModule,
//         MatTabsModule,
//         MatTableModule,
//         MatButtonModule,
//         MatCardModule,
//         MatDividerModule,
//         MatButtonToggleModule,
//         MatExpansionModule,
//         MatFormFieldModule,
//         MatSelectModule,
//         MatInputModule,
//         MatCheckboxModule,
//         MatDatepickerModule,
//         MatNativeDateModule,
//         MatSlideToggleModule,
//         MatIconModule,
//         MatChipsModule,
//         MatDialogModule,
//         FormsModule,
//         ReactiveFormsModule,
//         MatProgressSpinnerModule,
//         MatSnackBarModule,
//         MatMomentDateModule,
//         MatStepperModule,
//         NgxMaterialTimepickerModule.forRoot(),
//         RouterModule.forRoot(ELM_ROUTES),
//         MarkdownModule.forRoot(),
//         LMarkdownEditorModule
//     ],
//     providers: [ElmState, LeagueService, TeamsService, PlayersService, GamesService, UserService, LeagueOfLegendsStatsService, MatIconRegistry],
//     bootstrap: [AppComponent],
//     entryComponents: [
//         ReportGamePopup,
//         ManageGamePopup,
//         ManageTeamPopup,
//         WarningPopup,
//         ManagePlayersPopup,
//         ConfirmationComponent,
//         GenericPlayerEntry,
//         LeagueOfLegendsPlayerEntry,
//         LeagueOfLegendsStatsPage,
//         ManagePlayersTeamComponent,
//         ManagePlayersTeamLeagueOfLegendsComponent,
//         ManagePlayersPopupLeagueOfLegends,
//         ManageAvailabilityPopup,
//         GenericStatsPage
//     ]
// })
import {
    MatButtonModule,
    MatButtonToggleModule, MatCardModule,
    MatCheckboxModule,
    MatChipsModule,
    MatDatepickerModule, MatDialogModule, MatDividerModule, MatExpansionModule, MatFormFieldModule, MatIconModule,
    MatIconRegistry, MatInputModule,
    MatNativeDateModule, MatProgressSpinnerModule, MatSelectModule,
    MatSlideToggleModule, MatSnackBarModule, MatStepperModule, MatTableModule, MatTabsModule
} from "@angular/material";
import {NgModule} from "@angular/core";
import {AppComponent} from "./app.component";
import {LeagueService} from "./httpServices/leagues.service";
import {ElmState} from "./shared/state/state.service";
import {LMarkdownEditorModule} from "ngx-markdown-editor";
import {RouterModule} from "@angular/router";
import {HomeComponent} from "./home/home";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {ELM_ROUTES} from "./routes";
import {MarkdownModule} from "ngx-markdown";
import {MatMomentDateModule} from "@angular/material-moment-adapter";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {BrowserModule} from "@angular/platform-browser";
import {NavBar} from "./shared/navbar/navbar";
import {HttpClientModule} from "@angular/common/http";
import {NgxMaterialTimepickerModule} from "ngx-material-timepicker";
import {UserService} from "./httpServices/user.service";
import {LoggerModule, NgxLoggerLevel} from "ngx-logger";
import {GamesService} from "./httpServices/games.service";
import {TeamsService} from "./httpServices/teams.service";
import {StandingsComponent} from "./standings/standings";
import {TeamEntry} from "./shared/teamEntry/team-entry";
import {GameEntry} from "./shared/gameEntry/game-entry";
import {TeamsComponent} from "./teams/teams";
import {PlayerEntryComponent, PlayerEntryDirective} from "./teams/playerEntry/player-entry";
import {GenericPlayerEntry} from "./teams/playerEntry/generic-player-entry";
import {GamesComponent} from "./games/games";
import {RulesComponent} from "./rules/rules";
import {ManageComponent} from "./manage/manage";
import {ManageLeagueComponent} from "./manage/league/manage-league";
import {EventDisplayerComponent} from "./shared/eventDisplayer/event-displayer";
import {EventDisplayerService} from "./shared/eventDisplayer/event-displayer.service";
import {ManageRulesComponent} from "./manage/rules/manage-rules";
import {ManageDatesComponent} from "./manage/dates/manage-dates";
import {ManageTeamPopup, ManageTeamsComponent} from "./manage/teams/manage-teams";
import {WarningPopup} from "./manage/warningPopup/warning-popup";
import {ManagePlayersComponent} from "./manage/players/manage-players";
import {ManagePlayerPopup, ManagePlayersTeamComponent} from "./manage/players/manage-players-team/manage-players-team";
import {
    ManagePlayersTeamDirective,
    ManagePlayersTeamContainerComponent
} from "./manage/players/manage-players-team/manage-players-team-loader";
import {SignupComponent} from "./signup/signup";
import {LoginComponent} from "./login/login";
import {LeagueCreationComponent} from "./leagueCreation/league-creation";
import {LeaguesComponent} from "./leagues/leagues";
import {ManageAvailabilityPopup, ManageScheduleComponent} from "./manage/schedule/manage-schedule";
import {ManageGamePopup, ManageGamesComponent, ReportGamePopup} from "./manage/games/manage-games";
import {ManagePermissionsComponent} from "./manage/permissions/manage-permissions";
import {LeagueOfLegendsPlayerEntry} from "./teams/playerEntry/league-of-legends-player-entry";
import {
    ManageLoLPlayerPopup,
    ManagePlayersTeamLeagueOfLegendsComponent
} from "./manage/players/manage-players-team/league-of-legends/manage-players-team-league-of-legends";

@NgModule({
    declarations: [
        AppComponent,
        NavBar,
        SignupComponent,
        LoginComponent,
        LeaguesComponent,
        LeagueCreationComponent,
        HomeComponent,
        StandingsComponent,
        TeamEntry,
        GameEntry,
        TeamsComponent,
        PlayerEntryDirective,
        GenericPlayerEntry,
        LeagueOfLegendsPlayerEntry,
        PlayerEntryComponent,
        GamesComponent,
        RulesComponent,
        ManageComponent,
        ManageLeagueComponent,
        EventDisplayerComponent,
        ManageRulesComponent,
        ManageDatesComponent,
        ManageTeamsComponent,
        ManagePlayersComponent,
        ManagePlayersTeamComponent,
        ManagePlayersTeamLeagueOfLegendsComponent,
        ManagePlayersTeamDirective,
        ManagePlayersTeamContainerComponent,
        ManageGamesComponent,
        ManageScheduleComponent,
        ManagePermissionsComponent,
        ManageTeamPopup,
        ManagePlayerPopup,
        ManageLoLPlayerPopup,
        ManageAvailabilityPopup,
        ManageGamePopup,
        ReportGamePopup,
        WarningPopup
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
        MatSlideToggleModule,
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
        RouterModule.forRoot(ELM_ROUTES),
        MarkdownModule.forRoot(),
        LMarkdownEditorModule,
        LoggerModule.forRoot({level: NgxLoggerLevel.DEBUG})
    ],
    providers: [ElmState,
        EventDisplayerService,
        UserService,
        LeagueService,
        GamesService,
        TeamsService,
        MatIconRegistry
    ],
    bootstrap: [AppComponent],
    entryComponents: [
        GenericPlayerEntry,
        LeagueOfLegendsPlayerEntry,
        ManagePlayersTeamComponent,
        ManagePlayersTeamLeagueOfLegendsComponent,
        EventDisplayerComponent,
        ManageTeamPopup,
        ManagePlayerPopup,
        ManageLoLPlayerPopup,
        ManageAvailabilityPopup,
        ManageGamePopup,
        ReportGamePopup,
        WarningPopup
    ]
})

export class AppModule {
    constructor(public matIconRegistry: MatIconRegistry) {
        matIconRegistry.registerFontClassAlias('fontawesome', 'fa');
    }
}
