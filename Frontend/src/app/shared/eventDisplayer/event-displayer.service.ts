import {Injectable} from "@angular/core";
import {MatSnackBar} from "@angular/material";
import {EventDisplayerComponent} from "./event-displayer";
import {NGXLogger} from "ngx-logger";

@Injectable()
export class EventDisplayerService {
    constructor(private log: NGXLogger, public confirmation: MatSnackBar) {
    }

    public displaySuccess(message: string) {
        this.confirmation.openFromComponent(EventDisplayerComponent, {
            duration: 1250,
            panelClass: ['blue-snackbar'],
            data: message
        });
    }

    public displayError(error: any) {
        this.log.error(error);
        this.confirmation.openFromComponent(EventDisplayerComponent, {
            duration: 5000,
            panelClass: ['red-snackbar'],
            data: "Error: ".concat(error.status == 400 ? error.error.errorDescription : error.statusText)
        });
    }
}
