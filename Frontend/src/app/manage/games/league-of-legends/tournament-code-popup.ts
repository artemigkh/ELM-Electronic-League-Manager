import {Component, Inject} from "@angular/core";
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE, MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {NGXLogger} from "ngx-logger";
import {GamesService} from "../../../httpServices/games.service";
import {Game} from "../../../interfaces/Game";
import {EventDisplayerService} from "../../../shared/eventDisplayer/event-displayer.service";

@Component({
    selector: 'tournament-code-popup',
    templateUrl: 'tournament-code-popup.html',
    styleUrls: ['./tournament-code-popup.scss'],
})
export class TournamentCodePopup {
    code: string;
    constructor(
        public dialogRef: MatDialogRef<TournamentCodePopup>,
        private log: NGXLogger,
        private eventDisplayer: EventDisplayerService,
        private gamesService: GamesService,
        @Inject(MAT_DIALOG_DATA) public data: {game: Game}) {
        this.code = '';
        this.gamesService.getLoLTournamentCode(data.game.gameId, false).subscribe(
            res => this.code = res.tournamentCode,
            error => {
                this.eventDisplayer.displayError(error);
                this.dialogRef.close();
            }
        );
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    generateNewTournamentCode(): void {
        this.gamesService.getLoLTournamentCode(this.data.game.gameId, true).subscribe(
            res => {
                this.code = res.tournamentCode;
                this.eventDisplayer.displaySuccess("Successfully generated new tournament code");
            },
            error => {
                this.eventDisplayer.displayError(error);
                this.dialogRef.close();
            }
        );
    }
}
