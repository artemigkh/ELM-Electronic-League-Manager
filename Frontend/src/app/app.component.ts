import { Component } from '@angular/core';
import {LeagueService} from '../httpServices/leagues.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  constructor(private leagueService: LeagueService) {
    this.leagueService.setActiveLeague(63).
      subscribe(
      success => {console.log('success'); console.log(success); },
      error => {console.log('error'); console.log(error); });
  }
}
