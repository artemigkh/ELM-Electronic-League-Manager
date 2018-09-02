import { Component } from '@angular/core';
import {LeagueService} from '../httpServices/leagues.service';

export interface StandingsEntry {
  name: string;
  tag: string;
  wins: number;
  losses: number;
}

@Component({
  selector: 'app-standings',
  templateUrl: './standings.component.html',
  styleUrls: ['./standings.component.css']
})
export class StandingsComponent {
  displayedColumns: string[] = ['Name', 'Tag', 'Wins', 'Losses'];
  dataSource: StandingsEntry[];
  constructor(private leagueService: LeagueService) {
      this.leagueService.getTeamSummary().subscribe(
        teamSummary => {
          console.log('success');
          console.log(teamSummary);
          this.dataSource = teamSummary;
        },
        error => {console.log('error'); console.log(error); });
    }
}
