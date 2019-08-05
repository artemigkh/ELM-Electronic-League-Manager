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

    isTeamWinner(teamNum: number): boolean {
        return this.game.winnerId == this.game.team1.teamId && teamNum == 1 ||
                this.game.winnerId == this.game.team2.teamId && teamNum == 2;
    }

    gameResultText(teamNum: number): string {
        if (this.game.winnerId == this.game.team1.teamId && teamNum == 1 ||
            this.game.winnerId == this.game.team2.teamId && teamNum == 2) {
            return "VICTORY";
        } else {
            return "DEFEAT";
        }
    }
}
