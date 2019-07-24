import {Component, OnInit} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-manage-rules',
    templateUrl: './manage-rules.html',
    styleUrls: ['./manage-rules.scss'],
})
export class ManageRulesComponent implements OnInit{
    markdown: string;
    constructor(private leagueService: LeagueService,
                private log: NGXLogger,
                private eventDisplayer: EventDisplayerService) {
    }

    ngOnInit(): void {
        this.leagueService.getLeagueMarkdown().subscribe(
            md => this.markdown = md.markdown,
            error => this.log.error(error)
        )
    }

    updateAtServer() {
        this.leagueService.setLeagueMarkdown({markdown: this.markdown}).subscribe(
            success => this.eventDisplayer.displaySuccess("Rules and Information Successfully Updated"),
            error => this.eventDisplayer.displayError(error)
        );
    }
}
