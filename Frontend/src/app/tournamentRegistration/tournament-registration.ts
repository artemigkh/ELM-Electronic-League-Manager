import {LeagueService} from "../httpServices/leagues.service";
import {Component, ViewChild} from "@angular/core";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {EmptyTeam, Team} from "../interfaces/Team";
import {StepperSelectionEvent} from "@angular/cdk/typings/stepper";
import {CdkStepper} from "@angular/cdk/stepper";
import {MatSnackBar, MatStep, MatVerticalStepper} from "@angular/material";
import {TeamsService} from "../httpServices/teams.service";
import {Id} from "../httpServices/api-return-schemas/id";
import {ConfirmationComponent} from "../shared/confirmation/confirmation-component";

@Component({
    selector: 'app-tournament-registration',
    templateUrl: './tournament-registration.html',
    styleUrls: ['./tournament-registration.scss']
})
export class TournamentRegistrationComponent {
    @ViewChild('stepper') stepper: MatVerticalStepper;
    name: string;
    tag: string;
    description: string;
    firstFormGroup: FormGroup;
    secondFormGroup: FormGroup;
    team: Team;
    playerTransition: boolean;

    constructor(private leagueService: LeagueService,
                private teamsService: TeamsService,
                private _formBuilder: FormBuilder,
                public confirmation: MatSnackBar) {
        this.playerTransition = false;
    }

    ngOnInit() {
        this.firstFormGroup = this._formBuilder.group({
            name: ['', Validators.required],
            tag: ['', Validators.required]
        });
        this.secondFormGroup = this._formBuilder.group({
            secondCtrl: ['', Validators.required]
        });
    }

    onNext(e: StepperSelectionEvent) {
        if(e.selectedIndex == 1) {
            this.playerTransition = true;
            console.log("went to player information");
        }
    }

    onAnimationEnd() {
        if(this.playerTransition) {
            this.playerTransition = false;
            console.log(this.firstFormGroup.controls.name.value);
            console.log(this.firstFormGroup.controls.tag.value);
            this.teamsService.createNewTeam(
                this.firstFormGroup.controls.name.value,
                this.firstFormGroup.controls.tag.value).subscribe(
                (next: Id) => {
                    let team = EmptyTeam();
                    this.team.id = next.id;
                    this.team = team;
                }, error => {
                    let message = ": " + JSON.stringify(error.error);
                    if(error.error.error == "nameInUse") {
                        message = ": Name Is Already In Use"
                    } else if(error.error.error == "tagInUse") {
                        message = ": Tag Is Already In Use"
                    }
                    this.confirmation.openFromComponent(ConfirmationComponent, {
                        duration: 2000,
                        panelClass: ['red-snackbar'],
                        data: {
                            message: "Team Creation Failed" + message
                        }
                    });
                    this.stepper.previous();
                }
            );

        }
    }

}
