import {Component, Inject, OnInit} from "@angular/core";
import {Team, TeamCore} from "../../interfaces/Team";
import {ElmState} from "../../shared/state/state.service";
import {NGXLogger} from "ngx-logger";
import {EventDisplayerService} from "../../shared/eventDisplayer/event-displayer.service";
import {TeamsService} from "../../httpServices/teams.service";
import {UserWithPermissions} from "../../interfaces/User";
import {Action} from "../actions";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {WarningPopup, WarningPopupData} from "../warningPopup/warning-popup";

class TeamData {
    action: Action;
    team: Team;
    onSuccess: (team?: Team) => void;
}

@Component({
    selector: 'app-manage-teams',
    templateUrl: './manage-teams.html',
    styleUrls: ['./manage-teams.scss'],
})
export class ManageTeamsComponent implements OnInit {
    teams: Team[];
    user: UserWithPermissions;

    constructor(private state: ElmState,
                private log: NGXLogger,
                private teamsService: TeamsService,
                private eventDisplayer: EventDisplayerService,
                private dialog: MatDialog) {
    }

    ngOnInit(): void {
        this.state.subscribeUser(user => {
            this.user = user;
            this.getEditableTeamList();
        });
    }

    private getEditableTeamList() {
        this.teamsService.getLeagueTeams().subscribe(
            teams => {
                if (this.user.leaguePermissions.administrator || this.user.leaguePermissions.editTeams) {
                    this.teams = teams;
                } else {
                    this.teams = teams.filter(team => {
                        this.user.teamPermissions
                            .filter(t => t.administrator || t.information)
                            .map(t => t.teamId)
                            .includes(team.teamId);
                    });
                }
            }, error => this.eventDisplayer.displayError(error)
        );
    }

    private createTeam() {
        this.dialog.open(ManageTeamPopup, {
            data: <TeamData>{
                action: Action.Create,
                team: new Team(),
                onSuccess: (team => {
                    this.eventDisplayer.displaySuccess("Team Successfully Created");
                    this.teams = [...this.teams].concat(team);
                })
            },
            autoFocus: false, width: '500px'
        });
    }

    private editTeam(team: Team) {
        this.dialog.open(ManageTeamPopup, {
            data: <TeamData>{
                action: Action.Edit,
                team: team,
                onSuccess: () => this.eventDisplayer.displaySuccess("Team Successfully Updated")
            },
            autoFocus: false, width: '500px'
        });
    }

    private _deleteTeam(teamId: number) {
        this.teamsService.deleteTeam(teamId).subscribe(
            () => {
                this.eventDisplayer.displaySuccess("Team Successfully Deleted");
                this.teams = this.teams.filter(team => team.teamId != teamId);
            }, error => this.eventDisplayer.displayError(error)
        );
    }

    private deleteTeam(team: Team) {
        this.dialog.open(WarningPopup, {
            data: <WarningPopupData>{
                entity: "team",
                name: team.name,
                onAccept: () => this._deleteTeam(team.teamId)
            },
            autoFocus: false, width: '500px'
        });
    }
}

@Component({
    selector: 'manage-teams-popup',
    templateUrl: 'manage-teams-popup.html',
    styleUrls: ['./manage-teams-popup.scss'],
})
export class ManageTeamPopup {
    title: string;
    teamForm: FormGroup;

    constructor(
        @Inject(MAT_DIALOG_DATA) public data: TeamData,
        public dialogRef: MatDialogRef<ManageTeamPopup>,
        private log: NGXLogger,
        private teamsService: TeamsService,
        private formBuilder: FormBuilder) {
        this.title = this.data.action == Action.Create ? "Create New Team" : "Edit Team";
        this.teamForm = this.formBuilder.group({
            'name': [this.data.team.name, [Validators.required, Validators.minLength(3), Validators.maxLength(20)]],
            'tag': [this.data.team.tag, [Validators.required, Validators.minLength(3), Validators.maxLength(5)]],
            'description': [this.data.team.description, Validators.maxLength(500)],
            'icon': null
        });
    }

    onFileChange(event) {
        if(event.target.files.length > 0) {
            this.teamForm.value.icon = event.target.files[0];
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    saveTeam(): void {
        Object.keys(new TeamCore("", "", "")).forEach(k => this.data.team[k] = this.teamForm.value[k]);
        let form = new FormData();
        Object.keys(this.teamForm.value).forEach(k => form.append(k, this.teamForm.value[k]));
        if(this.data.action == Action.Create) {
            this.teamsService.createTeam(form).subscribe(
                res => {
                    this.data.team.teamId = res.teamId;
                    this.data.onSuccess(this.data.team);
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        } else if(this.data.action == Action.Edit) {
            this.teamsService.updateTeam(this.data.team.teamId, form).subscribe(
                () => {
                    this.data.onSuccess();
                    this.dialogRef.close();
                }, error => {
                    this.log.error(error);
                    this.dialogRef.close();
                }
            );
        }
    }
}
