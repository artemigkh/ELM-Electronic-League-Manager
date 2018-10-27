
import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";

@Component({
    selector: 'warning-popup',
    templateUrl: 'warning-popup.html',
    styleUrls: ['./warning-popup.scss'],
})
export class WarningPopup {
    constructor(
        public dialogRef: MatDialogRef<WarningPopup>,
        @Inject(MAT_DIALOG_DATA) public data: Object,
        private leagueService: LeagueService) {

    }

    onNoClick(): void {
        this.dialogRef.close();
    }
}
