import {Component, ViewEncapsulation} from "@angular/core";
import {LeagueInformation} from "../../interfaces/LeagueInformation";
import {MatSnackBar} from "@angular/material";
import {LeagueService} from "../../httpServices/leagues.service";
import {esportsDef, physicalSportsDef} from "../../shared/sports.defs";
import {ConfirmationComponent} from "../../shared/confirmation/confirmation-component";
import {Markdown} from "../../httpServices/api-return-schemas/markdown";

@Component({
    selector: 'app-manage-rules',
    templateUrl: './manage-rules.html',
    styleUrls: ['./manage-rules.scss'],
})
export class ManageRulesComponent {
    markdown: string;
    constructor(public confirmation: MatSnackBar, private leagueService: LeagueService) {
        this.leagueService.getMarkdown().subscribe(
            (next: Markdown) => {
                this.markdown = next.markdown;
                console.log(this.markdown)
            }, error => {
                console.log(error);
            }
        )
    }

    updateAtServer() {
        this.leagueService.setMarkdown(this.markdown).subscribe(
            next => {
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 1250,
                    panelClass: ['blue-snackbar'],
                    data: {
                        message: "Rules and Information Successfully Updated"
                    }
                });
            }, error => {
                console.log(error);
                this.confirmation.openFromComponent(ConfirmationComponent, {
                    duration: 2000,
                    panelClass: ['red-snackbar'],
                    data: {
                        message: "Update Failed"
                    }
                });
            }
        );

    }
}
