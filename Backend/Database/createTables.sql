CREATE TABLE leagues (
  id              SERIAL        PRIMARY KEY,
  name            VARCHAR(50)   NOT NULL   ,
  tag             VARCHAR(5)    NOT NULL
);

CREATE TABLE users (
  id              SERIAL        PRIMARY KEY,
  name            VARCHAR(50)   NOT NULL
);

CREATE TABLE teams (
  id              SERIAL        PRIMARY KEY,
  name            VARCHAR(50)   NOT NULL   ,
  tag             VARCHAR(5)    NOT NULL   ,
  wins            INT                      ,
  losses          INT
);

CREATE TABLE permissions (
  userID          SERIAL                   ,
  leagueID        SERIAL                   ,
  teamID          SERIAL                   ,
  editTeams       BOOLEAN                  ,
  editUsers       BOOLEAN                  ,
  editSchedule    BOOLEAN                  ,
  editGames       BOOLEAN                  ,
  reportResult    BOOLEAN
);

CREATE TABLE games (
  id              SERIAL        PRIMARY KEY,
  team1           SERIAL                   ,
  team2           SERIAL                   ,
  gametime        TIMESTAMP WITH TIME ZONE ,
  complete        BOOLEAN                  ,
  winner          SERIAL                   ,
  scoreteam1      INT                      ,
  scoreteam2      INT
);


