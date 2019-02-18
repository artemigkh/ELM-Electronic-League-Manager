import {Component, ViewEncapsulation} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {LeagueService} from "../httpServices/leagues.service";
import {Markdown} from "../httpServices/api-return-schemas/markdown";

@Component({
    selector: 'app-rules',
    templateUrl: './rules.html',
    styleUrls: ['./rules.scss'],
    encapsulation: ViewEncapsulation.None
})
export class RulesComponent {
    markdown: string;
    constructor(private leagueService: LeagueService) {
        this.leagueService.getMarkdown().subscribe(
            (next: Markdown) => {
                this.markdown = next.markdown;
                console.log(this.markdown)
            }, error => {
                console.log(error);
            }
        )
    }
}
