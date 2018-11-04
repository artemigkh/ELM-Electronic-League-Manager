import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {MatDialog} from '@angular/material'
import {LeagueService} from "../httpServices/leagues.service";

@Component({
    selector: 'app-signup',
    templateUrl: './signup.html',
    styleUrls: ['./signup.scss']
})
export class SignupComponent implements OnInit {
    constructor(private router: Router, private leagueService: LeagueService) { }
    email: string;
    password: string;
    ngOnInit() {
    }
    signup() : void {
        this.leagueService.signup(this.email, this.password).subscribe(
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

