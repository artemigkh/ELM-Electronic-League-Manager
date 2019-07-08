import {Component} from "@angular/core";
import {Router} from "@angular/router";
import {UserService} from "../httpServices/user.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-login',
    templateUrl: './login.html',
    styleUrls: ['./login.scss']
})
export class LoginComponent {
    email: string;
    password: string;

    constructor(private log: NGXLogger,
                private router: Router,
                private userService: UserService) {
    }

    login(): void {
        this.userService.login(this.email, this.password).subscribe(
            () => this.router.navigate(["login"]),
            error => this.log.error(error)
        );
    }
}
