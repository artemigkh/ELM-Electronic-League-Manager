import {PlayerEntryInterface} from "./player-entry";
import {Component, Input} from "@angular/core";
import {Player} from "../../interfaces/Player";

@Component({
    templateUrl: './generic-player-entry.html',
    styleUrls: ['./generic-player-entry.scss'],
})
export class GenericPlayerEntry implements PlayerEntryInterface {
    @Input() players: Player[];
    @Input() displayAsMainRoster: boolean = false;
}
