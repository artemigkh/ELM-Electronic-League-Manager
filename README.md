# ELM  - electronic-league-manager
Go/PostGreSQL/Angular5/Angular Materials Sports League Manager

Easily create, manage, and view leagues of any sport

## Setup
Download and install GoLang. Very highly recommend using JetBrains GoLand IDE.

### Database
1. Download postgres
2. Run `Backend/Database/createTables.sql` against postgres
3. Change the parameters in `Backend/conf.json` to the credentials you set in postgres install

### Go Libraries
```
go get ./...
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

### Server Tests
```
go test -v ./Backend/UnitTests/...
```

Or in GoLand, right click on the UnitTests folder in project view, and press Run->Go Test

#

This project is licensed under the terms of the MIT license