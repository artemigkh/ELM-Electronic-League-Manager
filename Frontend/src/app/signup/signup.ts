import {Component} from "@angular/core";
import {Router} from "@angular/router";
import {UserService} from "../httpServices/user.service";
import {NGXLogger} from "ngx-logger";

@Component({
    selector: 'app-signup',
    templateUrl: './signup.html',
    styleUrls: ['./signup.scss']
})
export class SignupComponent {
    email: string;
    password: string;

    constructor(private log: NGXLogger,
                private router: Router,
                private userService: UserService) {
    }

    signup(): void {
        this.userService.signup(this.email, this.password).subscribe(
            () => this.router.navigate(["login"]),
            error => this.log.error(error)
        );
    }
}
