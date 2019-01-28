import {PlayerEntryInterface} from "./player-entry";
import {Component, Input} from "@angular/core";
import {LeagueOfLegendsPlayer} from "../../interfaces/Player";

@Component({
    templateUrl: './league-of-legends-player-entry.html',
    styleUrls: ['./league-of-legends-player-entry.scss'],
})
export class LeagueOfLegendsPlayerEntry implements PlayerEntryInterface {
    @Input() players: LeagueOfLegendsPlayer[];
    @Input() mainRoster: boolean = false;
    getEmblem(player): string {
        if(player.tier.length > 0) {
            return "assets/leagueOfLegends/" +
                player.tier.substring(0, 1) +
                player.tier.substring(1).toLowerCase() +
                "_Emblem.png";
        } else {
            return "assets/leagueOfLegends/Unranked_Emblem.png";
        }
    }

    getPositionIcon(player: LeagueOfLegendsPlayer): string {
        return "assets/leagueOfLegends/" + player.position + "_Icon.png";
    }

    addRankNum(player: LeagueOfLegendsPlayer): string {
        if(player.tier == "MASTER" ||
            player.tier == "GRANDMASTER" ||
            player.tier == "CHALLENGER" ) {
            return ""
        } else {
            switch(player.rank) {
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

    getRankString(player: LeagueOfLegendsPlayer): string {
        if(player.tier.length > 0) {
            return player.tier.substring(0, 1) +
                player.tier.substring(1).toLowerCase() +
                this.addRankNum(player)
        } else {
            return "Unranked";
        }
    }
}

