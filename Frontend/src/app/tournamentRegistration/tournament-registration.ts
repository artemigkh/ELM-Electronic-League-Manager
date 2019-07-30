import {Component} from "@angular/core";
import {EventDisplayerService} from "../shared/eventDisplayer/event-displayer.service";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {TeamsService} from "../httpServices/teams.service";
import {getEmptyTeamWithRosters, Team, TeamCore, TeamWithRosters} from "../interfaces/Team";
import {Router} from "@angular/router";

@Component({
    selector: 'app-tournament-registration',
    templateUrl: './tournament-registration.html',
    styleUrls: ['./tournament-registration.scss']
})
export class TournamentRegistrationComponent {
    teamForm: FormGroup;
    team: TeamWithRosters;
    imagePreview: any;

    constructor(
        private router: Router,
        private teamsService: TeamsService,
        private _formBuilder: FormBuilder,
        private eventDisplayer: EventDisplayerService) {
        this.team = getEmptyTeamWithRosters();
        this.imagePreview = null;
    }

    newForm() {
        this.teamForm = this._formBuilder.group({
            'name': ['', [Validators.required, Validators.minLength(3), Validators.maxLength(25)],
                this.teamsService.validateTeamNameUniqueness(0)],
            'tag': ['', [Validators.required, Validators.minLength(3), Validators.maxLength(5)],
                this.teamsService.validateTagUniqueness(0)],
            'description': ['', Validators.maxLength(500)],
            'icon': null
        }, {updateOn: 'blur'});
        this.imagePreview = null;
    }

    ngOnInit() {
        this.newForm();
    }

    onSubmit() {
        this.teamsService.createTeamWithPlayers(
            new TeamCore(
                this.teamForm.value.name,
                this.teamForm.value.description,
                this.teamForm.value.tag
            ), this.teamForm.value.icon, this.team.mainRoster.concat(this.team.substituteRoster)).subscribe(
                createdTeamId => {
                    this.newForm();
                    this.eventDisplayer.displaySuccess("Team successfully registered");
                    this.router.navigate(["teams/", createdTeamId]);
                }, error => this.eventDisplayer.displayError(error)
        );
    }

    onFileChange(event) {
        if (event.target.files.length > 0) {
            let file: File = event.target.files[0];
            let reader = new FileReader();
            reader.onloadend = () => {
                this.imagePreview = reader.result;
                console.log(this.imagePreview);
                let encoded = reader.result.replace(/^data:(.*;base64,)?/, '');
                if ((encoded.length % 4) > 0) {
                    encoded += '='.repeat(4 - (encoded.length % 4));
                }
                this.teamForm.value.icon = encoded;
            };
            reader.readAsDataURL(file);
        }
    }

    removeImage() {
        this.imagePreview = null;
        this.teamForm.value.icon = null;
    }
}
