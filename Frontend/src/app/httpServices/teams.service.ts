import {Injectable} from "@angular/core";
import {HttpClient} from "@angular/common/http";
import {httpOptions} from "./http-options";
import {Observable} from "rxjs/Rx";

@Injectable()
export class TeamsService {
    constructor(private http: HttpClient) {}

    public createNewTeam(name: string, tag: string): Observable<Object> {
        return this.http.post('http://localhost:8080/api/teams/', {
            name: name,
            tag: tag
        }, httpOptions)
    }

    public updateTeam(id: number, name: string, tag: string) {
        return this.http.put('http://localhost:8080/api/teams/updateTeam/' + id, {
            name: name,
            tag: tag
        }, httpOptions)
    }

    public deleteTeam(id: number): Observable<Object> {
        return this.http.delete('http://localhost:8080/api/teams/removeTeam/' + id, httpOptions)
    }
}
