import {
    Component,
    ComponentFactoryResolver,
    Directive,
    Input, OnChanges,
    OnInit, SimpleChanges, Type,
    ViewChild,
    ViewContainerRef
} from "@angular/core";
import {TeamWithRosters} from "../../../interfaces/Team";
import {ManagePlayersTeamComponent} from "./manage-players-team";
import {ManageComponentInterface} from "../../manage-component-interface";
import {ElmState} from "../../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {ManagePlayersTeamLeagueOfLegendsComponent} from "./league-of-legends/manage-players-team-league-of-legends";

@Directive({
    selector: '[manage-players-team-host]',
})
export class ManagePlayersTeamDirective {
    constructor(public viewContainerRef: ViewContainerRef) { }
}

@Component({
    selector: 'manage-players-team-container',
    template: `<ng-template manage-players-team-host></ng-template>`
})
export class ManagePlayersTeamContainerComponent implements OnInit, OnChanges  {
    @Input() team: TeamWithRosters;
    @Input() sendToServer: boolean = true;
    @ViewChild(ManagePlayersTeamDirective) managePlayersTeamHost: ManagePlayersTeamDirective;

    game: string;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private componentFactoryResolver: ComponentFactoryResolver) { }

    ngOnInit() {
        this.state.subscribeLeague(league => {this.game = league.game; this.loadComponent()});
    }

    ngOnChanges(changes: SimpleChanges): void {
        if (this.game) {
            this.loadComponent();
        }
    }

    loadComponent() {
        this.log.debug("Received league with game = " + this.game);
        this.log.debug("Team in manage player team loader: ", this.team);
        this.managePlayersTeamHost.viewContainerRef.clear();
        let componentFactory = this.componentFactoryResolver.resolveComponentFactory(ManagePlayersTeamContainerComponent.getComponent(this.game));
        let componentRef = this.managePlayersTeamHost.viewContainerRef.createComponent(componentFactory);

        (<ManagePlayersTeamComponent>componentRef.instance).sendToServer = this.sendToServer;
        (<ManagePlayersTeamComponent>componentRef.instance).setTeam(this.team);
    }

    static getComponent(game: string): Type<ManageComponentInterface> {
        if (game == "leagueoflegends") {
            return ManagePlayersTeamLeagueOfLegendsComponent;
        } else {
            return ManagePlayersTeamComponent;
        }
    }
}


