import {
    Component,
    ComponentFactoryResolver,
    Directive,
    Input, OnChanges,
    OnInit, SimpleChanges, Type,
    ViewChild,
    ViewContainerRef
} from "@angular/core";
import {Player} from "../../interfaces/Player";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {GenericPlayerEntry} from "./generic-player-entry";
import {LeagueOfLegendsPlayerEntry} from "./league-of-legends-player-entry";

@Directive({
    selector: '[player-entry-host]',
})
export class PlayerEntryDirective {
    constructor(public viewContainerRef: ViewContainerRef) {
    }
}

export interface PlayerEntryInterface {
    displayAsMainRoster: boolean;
    players: Player[];
}

@Component({
    selector: 'player-entry-component',
    template: `
        <ng-template player-entry-host></ng-template>`
})
export class PlayerEntryComponent implements OnInit, OnChanges {
    @Input() displayAsMainRoster: boolean;
    @Input() players: Player[];
    @ViewChild(PlayerEntryDirective) playerEntryHost: PlayerEntryDirective;

    game: string;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private componentFactoryResolver: ComponentFactoryResolver) {
    }

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
        this.log.debug(this.players);
        this.playerEntryHost.viewContainerRef.clear();
        let componentFactory = this.componentFactoryResolver.resolveComponentFactory(PlayerEntryComponent.getComponent(this.game));
        let componentRef = this.playerEntryHost.viewContainerRef.createComponent(componentFactory);
        (<PlayerEntryInterface>componentRef.instance).displayAsMainRoster = this.displayAsMainRoster;
        (<PlayerEntryInterface>componentRef.instance).players = this.players;
    }

    static getComponent(game: string): Type<PlayerEntryInterface> {
        if (game == "leagueoflegends") {
            return LeagueOfLegendsPlayerEntry;
        } else {
            return GenericPlayerEntry;
        }
    }
}
