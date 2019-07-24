import {Injectable, OnInit} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {EmptyUser, User, UserCreationInformation, UserId, UserWithPermissions} from "../interfaces/User";
import {httpOptions} from "./http-options";
import {HttpClient} from "@angular/common/http";
import {ElmState} from "../shared/state/state.service";
import {NGXLogger} from "ngx-logger";

@Injectable()
export class UserService{
    constructor(private state: ElmState,
                private log: NGXLogger,
                private http: HttpClient) {
        this.state.subscribeChanges(() => {this.updateUserPermissions()});
    }

    public login(email: string, password: string): Observable<User> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/login', {
                email: email,
                password: password
            }, httpOptions).subscribe(
                (user: User) => {
                    this.log.debug("User successfully logged in", user);
                    this.state.setUser(user);
                    observer.next(user)
                }, error => {
                    this.log.error("Failed to log in", error);
                    observer.error(error);
                }
            )
        });
    }

    public signup(email: string, password: string): Observable<UserId> {
        return this.http.post<UserId> ('http://localhost:8080/api/v1/users', {
            email: email,
            password: password
        }, httpOptions);
    }

    public logout(): void {
        this.http.post('http://localhost:8080/logout', httpOptions).subscribe(
            next => {this.state.setUserWithPermissions(EmptyUser())},
            error => {this.log.error(error)}
        );
    }

    public getCurrentUser(): Observable<Object> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/v1/users', httpOptions).subscribe(
                (user: User) => {
                    this.state.setUser(user);
                    observer.next(user);
                }, error => {
                    this.log.error(error);
                    observer.error(error);
                }
            )
        });
    }

    private updateUserPermissions() {
        this.log.debug("Received request to update user permissions because state change");
        this.http.get('http://localhost:8080/api/v1/users/leaguePermissions', httpOptions).subscribe(
            (user: UserWithPermissions) => {
                this.state.setUserWithPermissions(user);
            }, error => {
                this.log.error(error);
            }
        )
    }
}
