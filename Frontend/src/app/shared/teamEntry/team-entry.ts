import {Team} from "../../interfaces/Team";
import {Component, Input} from "@angular/core";

@Component({
    selector: 'app-team-entry',
    templateUrl: './team-entry.html',
    styleUrls: ['./team-entry.scss']
})

export class TeamEntry {
    @Input() team: Team;
    @Input() position: Number;
    @Input() height: Number;
}
