import {PlayerEntryInterface} from "./player-entry";
import {Component, Input} from "@angular/core";
import {Player} from "../../interfaces/Player";

@Component({
    template: `
        <span class="player-entry">
            <span class="spacer"></span>
            <span class="name">
                {{player.name}}
            </span>
            <span class="rank">
                {{player.gameIdentifier}}
            </span>
        </span>
  `
})
export class GenericPlayerEntry implements PlayerEntryInterface {
    @Input() player: Player;
    @Input() mainRoster: boolean = false;
}
