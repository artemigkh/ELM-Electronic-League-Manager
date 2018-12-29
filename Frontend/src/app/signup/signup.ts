import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {MatDialog} from '@angular/material'
import {LeagueService} from "../httpServices/leagues.service";
import {UserService} from "../httpServices/user.service";

@Component({
    selector: 'app-signup',
    templateUrl: './signup.html',
    styleUrls: ['./signup.scss']
})
export class SignupComponent implements OnInit {
    constructor(private router: Router, private userService: UserService) { }
    email: string;
    password: string;
    ngOnInit() {
    }
    signup() : void {
        this.userService.signup(this.email, this.password).subscribe(
            next => {
                console.log("sign up successful");
                this.router.navigate(["login"]);
            }, error => {
                console.log("error");
                alert("sign up failed");
            }
        );
    }
}

