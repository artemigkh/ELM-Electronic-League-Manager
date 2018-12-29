import {Component} from "@angular/core";
import {UserService} from "../../httpServices/user.service";

@Component({
    selector: 'app-navbar',
    templateUrl: './navbar.html',
    styleUrls: ['./navbar.scss']
})
export class NavBar {
    loggedIn: boolean;
    constructor(private userService: UserService) {
        this.userService.registerNavBar(this);
        this.userService.checkIfLoggedIn().subscribe(
            next => {this.loggedIn = next}
        )
    }

    logout(): void {
        this.userService.logout().subscribe(next=>{this.loggedIn = false;})
    }

    public notifyLogin() {
        console.log("notify logged in");
        this.loggedIn = true;
    }
}
