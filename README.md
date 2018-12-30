# ELM  - electronic-league-manager
Go/PostGreSQL/Angular5/Angular Materials Sports League Manager

Easily create, manage, and view leagues of any sport

# Initial Setup

## Frontend
Install node/NPM, then angular CLI:

```
npm install -g @angular/cli
```



Download and install GoLang. Very highly recommend using JetBrains GoLand IDE.



Set local gopath to project root
```
export $GOPATH=<path to git repo>
```

create src folder and symlink it to Backend so that go compiler can find packages
```
mkdir src
ln -s <path to git repo>Backend <path to git repo>src
```

### Database
1. Download postgres
2. Run `Backend/Database/createTables.sql` against postgres
3. Change the parameters in `Backend/conf.json` to the credentials you set in postgres install

### Go Libraries
```
go get ./...
```

### Frontend
```
cd Frontend/
npm install
```

## How To Run

### Server
```
go build Backend/Server/main.go
./main
```

Or in GoLand, open main.go, right click, and 'run go build main.go'.
If errors with "cannot find find file" happen, set the working directory to the
root directory of the git repo.

#### Server Tests

##### Unit Tests
```
go test -v ./Backend/UnitTests/...
```

Or in GoLand, right click on the UnitTests folder in project view, and press Run->Go Test

##### Integration Tests

Note that running the integration tests requires the following environment:

* Database set up, `createTables.sql` ran against it
* Database permissions to allow the user in the connection string in `conf.json` to connect to the database
* Permission for the application to serve on the local loopback on port 8080
* Permission for the go testing executable to send requests to the local loopback on port 8080 on HTTP

```
go test -v ./Backend/IntegrationTests/...
```
Or in GoLand, right click on the IntegrationTests folder in project view, and press Run->Go Test

### Generate Endpoint Documentation
Download and install [http://apidocjs.com/](apiDoc)
Execute this command in the Backend directory of this project:
```
cd Backend/
apidoc
```
Open `doc/index.html` on the browser of your choice

### Frontend

```
cd Frontend/
ng serve
```

Then, navigate to `http://localhost:4200/` in your browser

## Development

### Go Backend

This repository strives to maintain readable, reusable, stable, and easily augmentable code.

As such, the style guide at https://golang.org/doc/effective_go.html is strictly enforced (GoLand automatically highlights
problems with these guidelines).

All PRs must come with unit tests that cover new functionality, and pass the full suite of tests.

All PRs must have files `go fmt`d

This project uses `github.com/stretchr/testify/mock` for testing. To generate mocks after changing an interface, run
```
go get github.com/vektra/mockery/.../
mockery -all
```

#

This project is licensed under the terms of the GPL-3.0 license