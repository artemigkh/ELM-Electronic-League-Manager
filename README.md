# ELM  - electronic-league-manager
Go/PostGreSQL/Angular5/Angular Materials Sports League Manager

Easily create, manage, and view leagues of any sport

# Initial Setup

## Backend

### Database

If you are familiar with postgres
* Create a new database, with unique user if desired
* Run `Backend/Database/createTables.sql` against new db to generate tables
* Modify `Backend/conf.json` with created user, password, and database name

If you are not familiar with postgres
* Download and install postgresql and pgAdmin4. Packages for all major OS can be found on the postgresql website
* Set postgres user password: On ubuntu this can be done with
```
sudo -u postgres psql postgres
\password postgres
```

* Open pgAdmin
* If no servers appear, create a new one: Name can be anything, under connection
host=localhost, username=postgres, password=<what you set it to earlier>
* Opening the server, right click on databases and create new database with name `elmdb` and owner `postgres`
* Click on newly created database, then tools->query tool
* Paste in contents of `Backend/Database/createTables.sql` and execute (F5)
* Modify `Backend/conf.json` with dbUser="postgres, password=<what you set it to earlier>, and dbName=elmdb

### Server
Download and install Go following instructions from golang.org/doc/install

Set local gopath to project root - recommended to add this to ~/.profile or equivalent to make it persist through reboots
```
export GOPATH=<path to git repo>/Backend
```

Get all dependencies
```
cd <path to git repo>/Backend
go get -u ./...
```

Check that it compiles with
```
go build src/Server/main.go
```

## Frontend
Install node/NPM, then angular CLI:

```
cd <path git repo>/Frontend
npm install
npm install -g @angular/cli
```

Test that it works with

```
ng serve
```

And ensure that it compiles successfully


## How To Run
### Server
```
cd <path to git repo>
go build Backend/src/Server/main.go
./main
```

### Frontend
```
cd <path to git repo>/Frontend
ng serve
```

## Backend Tests

### Unit Tests
```
go test -v ./Backend/src/UnitTests/...
```

### Test Data Generation / Integration Tests

The integration tests that start the server and test both the endpoints and database access is also used to
generate mock league data for use of testing the API and the frontend

Note that running the integration tests requires the following environment:

* Database set up, `createTables.sql` ran against it
* Database permissions to allow the user in the connection string in `conf.json` to connect to the database
* Permission for the application to serve on the local loopback on port 8080
* Permission for the go testing executable to send requests to the local loopback on port 8080 on HTTP

```
go test -v ./Backend/IntegrationTests/...
```

## Generate Endpoint Documentation
Download and install [http://apidocjs.com/](apiDoc)
Execute this command in the Backend directory of this project:
```
cd Backend/src/Server
apidoc
```
Open `doc/index.html` on the browser of your choice

## Development

### Go Backend

All PRs must come with unit tests that cover new functionality and pass the full suite of tests

This project uses `github.com/stretchr/testify/mock` for testing. To generate mocks after changing an interface, run
```
go get github.com/vektra/mockery/.../
cd Backend/src/Server
mockery -all -output ../mocks
```

This project is licensed under the terms of the GPL-3.0 license