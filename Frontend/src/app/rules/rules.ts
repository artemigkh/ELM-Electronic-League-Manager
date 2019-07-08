import {Component, OnInit, ViewEncapsulation} from "@angular/core";
import {LeagueService} from "../httpServices/leagues.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-rules',
    templateUrl: './rules.html',
    styleUrls: ['./rules.scss'],
    encapsulation: ViewEncapsulation.None
})
export class RulesComponent implements OnInit{
    markdown: string;
    constructor(private log: NGXLogger, private leagueService: LeagueService) {
        this.markdown = "";
    }

    ngOnInit(): void {
        this.leagueService.getLeagueMarkdown().subscribe(
            md => this.markdown = md.markdown,
            error => this.log.error(error)
        )
    }
}
