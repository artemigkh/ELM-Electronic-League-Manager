import {Component, Input} from "@angular/core";
import {StatsPageInterface} from "./stats-page";

@Component({
    templateUrl: './generic-stats-page.html',
})
export class GenericStatsPage implements StatsPageInterface {}
