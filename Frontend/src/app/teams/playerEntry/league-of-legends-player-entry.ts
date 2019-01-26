import {PlayerEntryInterface} from "./player-entry";
import {Component, Input} from "@angular/core";
import {LeagueOfLegendsPlayer} from "../../interfaces/Player";

@Component({
    template: `
        <span class="player-entry">
            <span class="rank">
                <div *ngIf="mainRoster" class="position">
                    <img [src]="getPositionIcon()">
                </div>
                <div class="game-identifier">
                    {{player.gameIdentifier}}
                </div>
                <div class="emblem">
                    <img [src]="getEmblem()">
                </div>
                <div>
                    {{getRankString()}}
                </div>
            </span>
        </span>
  `
})
export class LeagueOfLegendsPlayerEntry implements PlayerEntryInterface {
    @Input() player: LeagueOfLegendsPlayer;
    @Input() mainRoster: boolean = false;
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

    getPositionIcon(): string {
        return "assets/leagueOfLegends/" + this.player.position + "_Icon.png";
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

