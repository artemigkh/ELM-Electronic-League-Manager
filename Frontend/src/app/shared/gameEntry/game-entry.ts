import {Component, Input} from "@angular/core";
import {Game} from "../../interfaces/Game";

@Component({
    selector: 'app-game-entry',
    templateUrl: './game-entry.html',
    styleUrls: ['./game-entry.scss']
})

export class GameEntry {
    @Input() game: Game;
    @Input() compact: Boolean;
    @Input() tentative: Boolean = false;

    teamNameClass(teamNum: number): string {
        let toReturn = "team-name";
        if (!this.game.complete) {
            return toReturn;
        } else {
            if (this.game.winnerId == this.game.team1.teamId && teamNum == 1 ||
                this.game.winnerId == this.game.team2.teamId && teamNum == 2 ) {
                return toReturn + " victory";
            } else {
                return toReturn + " defeat";
            }
        }
    }

    gameResultClass(teamNum: number): string {
        let c = "";
        if (this.game.winnerId == this.game.team1.teamId && teamNum == 1 ||
            this.game.winnerId == this.game.team2.teamId && teamNum == 2 ) {
            c = "victory";
        } else {
            c = "defeat";
        }

        if (!this.game.complete) {
            c += " ignore";
        }

        return c;
    }

    gameResultText(teamNum: number): string {
        if (this.game.winnerId == this.game.team1.teamId && teamNum == 1 ||
            this.game.winnerId == this.game.team2.teamId && teamNum == 2 ) {
            return "VICTORY";
        } else {
            return "DEFEAT";
        }
    }
}
