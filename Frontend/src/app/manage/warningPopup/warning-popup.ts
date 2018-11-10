
import {Component, Inject} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {ManageComponentInterface} from "../manage-component-interface";

class WarningPopupData {
    entity: string;
    name: string;
    Id: number;
    Id2: number;
    caller: ManageComponentInterface;
}

@Component({
    selector: 'warning-popup',
    templateUrl: 'warning-popup.html',
    styleUrls: ['./warning-popup.scss'],
})
export class WarningPopup {
    constructor(
        public dialogRef: MatDialogRef<WarningPopup>,
        @Inject(MAT_DIALOG_DATA) public data: WarningPopupData) {}

    onCancel(): void {
        this.dialogRef.close();
    }

    onConfirm(): void {
        if(this.data.Id2) {
            this.data.caller.notifyDelete(this.data.Id, this.data.Id2);
        } else {
            this.data.caller.notifyDelete(this.data.Id);
        }

        this.dialogRef.close();
    }
}
