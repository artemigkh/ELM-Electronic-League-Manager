import {Component, Inject, ViewContainerRef} from "@angular/core";
import {MAT_SNACK_BAR_DATA, MatDialogRef} from "@angular/material";
import {WarningPopup} from "../../manage/warningPopup/warning-popup";

export class ConfirmationData {
    message: string
}

@Component({
    selector: 'app-confirmation-component',
    templateUrl: './confirmation-component.html',
    styles: [`
        span {
            text-align: center;
        }
    `],
})
export class ConfirmationComponent {
    constructor(
        public containerRef: ViewContainerRef,
        @Inject(MAT_SNACK_BAR_DATA) public data: ConfirmationData) {}
}
