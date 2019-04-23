import {Component, ComponentFactoryResolver, Directive, OnInit, ViewChild, ViewContainerRef} from "@angular/core";
import {Player} from "../../interfaces/Player";
import {GenericPlayerEntry} from "../../teams/playerEntry/generic-player-entry";
import {LeagueOfLegendsPlayerEntry} from "../../teams/playerEntry/league-of-legends-player-entry";
import {GenericStatsPage} from "./generic-stats-page";
import {LeagueService} from "../../httpServices/leagues.service";
import {PlayerEntryDirective, PlayerEntryInterface} from "../../teams/playerEntry/player-entry";
import {LeagueOfLegendsStatsPage} from "./league-of-legends-stats-page";

@Directive({
    selector: '[stats-page-host]',
})
export class StatsPageDirective {
    constructor(public viewContainerRef: ViewContainerRef) { }
}

export interface StatsPageInterface {

}

const componentMapping: { [id: string] : any; } = {
    "genericsport": GenericStatsPage,
    "basketball": GenericStatsPage,
    "curling": GenericStatsPage,
    "football": GenericStatsPage,
    "hockey":GenericStatsPage,
    "rugby": GenericStatsPage,
    "soccer": GenericStatsPage,
    "volleyball": GenericStatsPage,
    "waterpolo": GenericStatsPage,
    "genericesport": GenericStatsPage,
    "csgo": GenericStatsPage,
    "leagueoflegends": LeagueOfLegendsStatsPage,
    "overwatch": GenericStatsPage
};

@Component({
    selector: 'stats-page-component',
    template: `<ng-template stats-page-host></ng-template>`
})
export class StatsPageComponent implements OnInit {
    @ViewChild(StatsPageDirective) statsPageHost: StatsPageDirective;

    constructor(private componentFactoryResolver: ComponentFactoryResolver,
                private leagueService: LeagueService) { }

    ngOnInit() {
        this.loadComponent();
    }

    loadComponent() {
        console.log(this.leagueService.getGame());
        let componentFactory = this.componentFactoryResolver.
        resolveComponentFactory(componentMapping[this.leagueService.getGame()]);
        let viewContainerRef = this.statsPageHost.viewContainerRef;
        let componentRef = viewContainerRef.createComponent(componentFactory);
    }
}
