import {Component} from "@angular/core";
import {LeagueService} from "../../httpServices/leagues.service";

@Component({
    selector: 'app-navbar',
    templateUrl: './navbar.html',
    styleUrls: ['./navbar.scss']
})
export class NavBar {
    loggedIn: boolean;
    constructor(private leagueService: LeagueService) {
        this.leagueService.registerNavBar(this);
        this.leagueService.checkIfLoggedIn().subscribe(
            next => {this.loggedIn = next}
        )
    }

    logout(): void {
        this.leagueService.logout().subscribe(next=>{this.loggedIn = false;})
    }

    public notifyLogin() {
        console.log("notify logged in");
        this.loggedIn = true;
    }
}
