import {
    Component,
    ComponentFactoryResolver,
    Directive,
    Input,
    OnInit,
    ViewChild,
    ViewContainerRef
} from "@angular/core";
import {ManagePlayersTeamLeagueOfLegendsComponent} from "./league-of-legends/manage-players-team-league-of-legends";
import {Team} from "../../../interfaces/Team";
import {LeagueService} from "../../../httpServices/leagues.service";
import {ManagePlayersTeamComponent} from "./manage-players-team";


@Directive({
    selector: '[manage-players-team-host]',
})
export class ManagePlayersTeamDirective {
    constructor(public viewContainerRef: ViewContainerRef) { }
}


const mptComponentMapping: { [id: string] : any; } = {
    "genericsport": ManagePlayersTeamComponent,
    "basketball": ManagePlayersTeamComponent,
    "curling": ManagePlayersTeamComponent,
    "football": ManagePlayersTeamComponent,
    "hockey":ManagePlayersTeamComponent,
    "rugby": ManagePlayersTeamComponent,
    "soccer": ManagePlayersTeamComponent,
    "volleyball": ManagePlayersTeamComponent,
    "waterpolo": ManagePlayersTeamComponent,
    "genericesport": ManagePlayersTeamComponent,
    "csgo": ManagePlayersTeamComponent,
    "leagueoflegends": ManagePlayersTeamLeagueOfLegendsComponent,
    "overwatch": ManagePlayersTeamComponent
};

@Component({
    selector: 'manage-players-team-container',
    template: `<ng-template manage-players-team-host></ng-template>`
})
export class ManagePlayersTeamContainerComponent implements OnInit {
    @Input() team: Team;
    @ViewChild(ManagePlayersTeamDirective) managePlayersTeamHost: ManagePlayersTeamDirective;

    constructor(private componentFactoryResolver: ComponentFactoryResolver,
                private leagueService: LeagueService) { }

    ngOnInit() {
        this.loadComponent();
    }

    loadComponent() {
        console.log(this.leagueService.getGame());
        let componentFactory = this.componentFactoryResolver.
        resolveComponentFactory(mptComponentMapping[this.leagueService.getGame()]);
        let viewContainerRef = this.managePlayersTeamHost.viewContainerRef;
        let componentRef = viewContainerRef.createComponent(componentFactory);
        (<ManagePlayersTeamComponent>componentRef.instance).team = this.team;
    }
}


