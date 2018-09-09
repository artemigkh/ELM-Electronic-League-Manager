import { Component } from '@angular/core';
import {LeagueService} from './httpServices/leagues.service';

@Component({
  selector: 'elm-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  constructor(private leagueService: LeagueService) {
    this.leagueService.setActiveLeague(2).
      subscribe(
      success => {console.log('success'); console.log(success); },
      error => {console.log('error'); console.log(error); });
  }
}
