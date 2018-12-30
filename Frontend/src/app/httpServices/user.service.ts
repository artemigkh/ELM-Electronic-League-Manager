import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Rx";
import {User} from "../interfaces/User";
import {httpOptions} from "./http-options";
import {Id} from "./api-return-schemas/id";
import {NavBar} from "../shared/navbar/navbar";
import {HttpClient} from "@angular/common/http";

@Injectable()
export class UserService {
    user: User;
    navBar: NavBar;

    constructor(private http: HttpClient) {
        this.user = null;
        this.navBar = null;
    }

    public login(email: string, password: string): Observable<User> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/login', {
                email: email,
                password: password
            }, httpOptions).subscribe(
                (next: Id) => {
                    console.log(this.navBar);
                    this.navBar.notifyLogin();
                    this.user = {
                        id: next.id,
                        email: email
                    };
                    observer.next(this.user)
                }, error => {
                    observer.error(error);
                }
            )
        })
    }

    public signup(email: string, password: string): Observable<boolean> {
        return new Observable(observer => {
            this.http.post('http://localhost:8080/api/users/', {
                email: email,
                password: password
            }, httpOptions).subscribe(
                next => {observer.next(true);},
                error => {observer.next(false);}
            )
        })
    }

    public logout(): Observable<Object> {
        return this.http.post('http://localhost:8080/logout', httpOptions);
    }

    public checkIfLoggedIn(): Observable<boolean> {
        return new Observable(observer => {
            this.http.get('http://localhost:8080/api/users/profile', httpOptions).subscribe(
                next => {observer.next(true);},
                error => {observer.next(false);}
            )
        });
    }

    public getCurrentUser() {
        return this.user;
    }

    public registerNavBar(navBar: NavBar) {
        this.navBar = navBar;
    }
}
