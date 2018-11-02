import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {MatDialog} from '@angular/material'
import {LeagueService} from "../httpServices/leagues.service";
import {User} from "../interfaces/User";

@Component({
    selector: 'app-login',
    templateUrl: './login.html',
    styleUrls: ['./login.scss']
})
export class LoginComponent implements OnInit {
    constructor(private router: Router, private leagueService: LeagueService) { }
    email: string;
    password: string;
    ngOnInit() {
    }
    login() : void {
        this.leagueService.login(this.email, this.password).subscribe(
            (next: User) => {
                console.log("logged in with user with id ", next.id)
                this.router.navigate([""]);
            }, error => {
                console.log("error");
                alert("Incorrect email or password");
            }
        );
    }
}

