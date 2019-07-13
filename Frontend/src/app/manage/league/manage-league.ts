import {Component, OnInit} from "@angular/core";
import {League} from "../../interfaces/League";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {LeagueService} from "../../httpServices/leagues.service";
import {Option} from "../../interfaces/UI";
import {eSportsDef, physicalSportsDef} from "../../shared/lookup.defs";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";

@Component({
    selector: 'app-manage-league',
    templateUrl: './manage-league.html',
    styleUrls: ['./manage-league.scss'],
})
export class ManageLeagueComponent implements OnInit{
    league: League;
    physicalSports: Option[];
    eSports: Option[];

    constructor(private state: ElmState,
                private log: NGXLogger,
                private leagueService: LeagueService,
                private eventDisplayer: EventDisplayerService) {
    }

    ngOnInit(): void {
        this.state.subscribeLeague(league => this.league = league);
        this.physicalSports = Object.entries(physicalSportsDef).map(o => <Option>{value: o[0], display: o[1]});
        this.eSports = Object.entries(eSportsDef).map(o => <Option>{value: o[0], display: o[1]});
    }

    updateAtServer() {
        this.leagueService.updateLeagueInformation(this.league).subscribe(
            () => this.eventDisplayer.displaySuccess("League Information Successfully Updated"),
            error => this.eventDisplayer.displayError(error)
        );
    }
}
