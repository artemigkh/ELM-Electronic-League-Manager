import {Component, Inject, ViewContainerRef} from "@angular/core";
import {MAT_SNACK_BAR_DATA} from "@angular/material";

@Component({
    selector: 'app-event-displayer-component',
    templateUrl: './event-displayer.html',
    styles: [`
        span {
            text-align: center;
        }
    `],
})
export class EventDisplayerComponent {
    constructor(
        public containerRef: ViewContainerRef,
        @Inject(MAT_SNACK_BAR_DATA) public data: String) {}
}
