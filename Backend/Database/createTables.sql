CREATE SEQUENCE leaguesIDSeq;
CREATE SEQUENCE usersIDSeq;
CREATE SEQUENCE teamsIDSeq;
CREATE SEQUENCE gamesIDSeq;

CREATE TABLE leagues (
  id              INT           PRIMARY KEY DEFAULT nextval('leaguesIDSeq'),
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  publicView      BOOLEAN                        ,
  publicJoin      BOOLEAN
);
ALTER SEQUENCE leaguesIDSeq OWNED BY leagues.id;

CREATE TABLE users (
  id              INT           PRIMARY KEY DEFAULT nextval('usersIDSeq'),
  email           VARCHAR(256)  UNIQUE NOT NULL  ,
  salt            CHAR(64)      NOT NULL         ,
  hash            CHAR(128)     NOT NULL
);
ALTER SEQUENCE usersIDSeq OWNED BY users.id;

CREATE TABLE teams (
  id              INT           PRIMARY KEY DEFAULT nextval('teamsIDSeq'),
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  wins            INT                            ,
  losses          INT
);
ALTER SEQUENCE teamsIDSeq OWNED BY teams.id;

CREATE TABLE permissions (
  userID          INT                            ,
  leagueID        INT                            ,
  teamID          INT                            ,
  editTeams       BOOLEAN                        ,
  editUsers       BOOLEAN                        ,
  editSchedule    BOOLEAN                        ,
  editGames       BOOLEAN                        ,
  reportResult    BOOLEAN
);


CREATE TABLE games (
  id              INT           PRIMARY KEY DEFAULT nextval('gamesIDSeq'),
  team1           INT                            ,
  team2           INT                            ,
  gametime        TIMESTAMP WITH TIME ZONE       ,
  complete        BOOLEAN                        ,
  winner          INT                            ,
  scoreteam1      INT                            ,
  scoreteam2      INT
);
ALTER SEQUENCE gamesIDSeq OWNED BY games.id;

