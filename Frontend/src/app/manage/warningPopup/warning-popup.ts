import {Component, Inject} from "@angular/core";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";

export interface WarningPopupData {
    entity: string;
    name: string;
    onAccept: () => void;
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
        this.data.onAccept();
        this.dialogRef.close();
    }
}
