import {HttpHeaders} from "@angular/common/http";

export const httpOptions = {
    withCredentials: true,
    headers: new HttpHeaders({
        'Content-Type': 'application/json'
    })
};
