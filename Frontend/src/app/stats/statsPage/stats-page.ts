import {
    Component,
    ComponentFactoryResolver,
    Directive, OnChanges,
    OnInit,
    SimpleChanges,
    Type,
    ViewChild,
    ViewContainerRef
} from "@angular/core";
import {GenericStatsPage} from "./generic-stats-page";
import {LeagueOfLegendsStatsPage} from "./league-of-legends-stats-page";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";

@Directive({
    selector: '[stats-page-host]',
})
export class StatsPageDirective {
    constructor(public viewContainerRef: ViewContainerRef) {
    }
}

export interface StatsPageInterface {

}

@Component({
    selector: 'stats-page-component',
    template: `
        <ng-template stats-page-host></ng-template>`
})
export class StatsPageComponent implements OnInit, OnChanges {
    @ViewChild(StatsPageDirective) statsPageHost: StatsPageDirective;

    game: string;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private componentFactoryResolver: ComponentFactoryResolver) {
    }

    ngOnInit() {
        this.state.subscribeLeague(league => {
            this.game = league.game;
            this.loadComponent()
        });
    }

    ngOnChanges(changes: SimpleChanges): void {
        if (this.game) {
            this.loadComponent();
        }
    }

    loadComponent() {
        this.log.debug("Received league with game = " + this.game);

        let componentFactory = this.componentFactoryResolver.resolveComponentFactory(StatsPageComponent.getComponent(this.game));
        this.statsPageHost.viewContainerRef.createComponent(componentFactory);
    }

    static getComponent(game: string): Type<StatsPageInterface> {
        if (game == "leagueoflegends") {
            return LeagueOfLegendsStatsPage;
        } else {
            return GenericStatsPage;
        }
    }
}
