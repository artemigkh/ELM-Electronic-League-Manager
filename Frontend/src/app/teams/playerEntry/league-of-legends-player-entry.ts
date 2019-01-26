import {PlayerEntryInterface} from "./player-entry";
import {Component, Input} from "@angular/core";
import {LeagueOfLegendsPlayer} from "../../interfaces/Player";

@Component({
    template: `
        <span class="player-entry">
            <span class="rank">
                <div>
                    {{player.gameIdentifier}}
                </div>
                <div>
                    <img [src]="getEmblem()">
                </div>
                <div>
                    {{getRankString()}}
                </div>
            </span>
        </span>
  `, styles: ['img { max-height: 100px; max-width: 100px; }',
    '.player-entry {height: 100px;}']
})
export class LeagueOfLegendsPlayerEntry implements PlayerEntryInterface {
    @Input() player: LeagueOfLegendsPlayer;
    getEmblem(): string {
        if(this.player.tier.length > 0) {
            return "assets/leagueOfLegends/" +
                this.player.tier.substring(0, 1) +
                this.player.tier.substring(1).toLowerCase() +
                "_Emblem.png";
        } else {
            return "assets/leagueOfLegends/Unranked_Emblem.png";
        }
    }

    addRankNum(): string {
        if(this.player.tier == "MASTER" ||
            this.player.tier == "GRANDMASTER" ||
            this.player.tier == "CHALLENGER" ) {
            return ""
        } else {
            switch(this.player.rank) {
                case "I": {
                    return " 1";
                }
                case "II": {
                    return " 2";
                }
                case "III": {
                    return " 3";
                }
                case "IV": {
                    return " 4"
                }
                default: {
                    return ""
                }
            }
        }
    }

    getRankString(): string {
        if(this.player.tier.length > 0) {
            return this.player.tier.substring(0, 1) +
                this.player.tier.substring(1).toLowerCase() +
                this.addRankNum()
        } else {
            return "Unranked";
        }
    }
}

