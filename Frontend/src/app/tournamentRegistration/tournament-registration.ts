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
import {isUndefined} from "util";

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
            tag: ['', Validators.required],
            description: '',
            icon: null
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

    onFileChange(event) {
        if(event.target.files.length > 0) {
            let file = event.target.files[0];
            this.firstFormGroup.value.icon = file;
        }
    }

    onAnimationEnd() {
        if(this.playerTransition) {
            let form = new FormData();
            form.append("name", this.firstFormGroup.value.name);
            form.append("tag", this.firstFormGroup.value.tag);
            form.append("description", this.firstFormGroup.value.description);
            form.append("icon", this.firstFormGroup.value.icon);
            if(isUndefined(this.team)) {
                this.teamsService.createNewTeamWithIcon(form).subscribe(
                    (next: Id) => {
                        let team = EmptyTeam();
                        team.id = next.id;
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
            } else {
                this.teamsService.updateTeamWithIcon(this.team.id, form).subscribe(
                    next => {}, error => {
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
                    });
            }
            this.playerTransition = false;
        }
    }

}
