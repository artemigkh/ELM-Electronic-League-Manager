import {
    Component, ComponentFactoryResolver, Directive, Input, OnInit, ViewChild,
    ViewContainerRef
} from '@angular/core';
import {Player} from "../../interfaces/Player";
import {GenericPlayerEntry} from "./generic-player-entry";
import {LeagueOfLegendsPlayerEntry} from "./league-of-legends-player-entry";
import {LeagueService} from "../../httpServices/leagues.service";

@Directive({
    selector: '[player-entry-host]',
})
export class PlayerEntryDirective {
    constructor(public viewContainerRef: ViewContainerRef) { }
}

export interface PlayerEntryInterface {
    players: Player[];
    mainRoster: boolean;
}

const componentMapping: { [id: string] : any; } = {
    "genericsport": GenericPlayerEntry,
    "basketball": GenericPlayerEntry,
    "curling": GenericPlayerEntry,
    "football": GenericPlayerEntry,
    "hockey":GenericPlayerEntry,
    "rugby": GenericPlayerEntry,
    "soccer": GenericPlayerEntry,
    "volleyball": GenericPlayerEntry,
    "waterpolo": GenericPlayerEntry,
    "genericesport": GenericPlayerEntry,
    "csgo": GenericPlayerEntry,
    "leagueoflegends": LeagueOfLegendsPlayerEntry,
    "overwatch": GenericPlayerEntry
};
@Component({
    selector: 'player-entry-component',
    template: `<ng-template player-entry-host></ng-template>`
})
export class PlayerEntryComponent implements OnInit {
    @Input() players: Player[];
    @Input() mainRoster: boolean = false;
    @ViewChild(PlayerEntryDirective) playerEntryHost: PlayerEntryDirective;

    constructor(private componentFactoryResolver: ComponentFactoryResolver,
                private leagueService: LeagueService) { }

    ngOnInit() {
        this.loadComponent();
    }

    loadComponent() {
        console.log(this.leagueService.getGame());
        let componentFactory = this.componentFactoryResolver.
            resolveComponentFactory(componentMapping[this.leagueService.getGame()]);
        let viewContainerRef = this.playerEntryHost.viewContainerRef;
        let componentRef = viewContainerRef.createComponent(componentFactory);
        (<PlayerEntryInterface>componentRef.instance).players = this.players;
        (<PlayerEntryInterface>componentRef.instance).mainRoster = this.mainRoster;
    }
}
