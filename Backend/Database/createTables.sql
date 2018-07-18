CREATE SEQUENCE leaguesIDSeq;
CREATE SEQUENCE usersIDSeq;
CREATE SEQUENCE teamsIDSeq;
CREATE SEQUENCE gamesIDSeq;

CREATE TABLE leagues (
  id              INT           PRIMARY KEY DEFAULT nextval('leaguesIDSeq'),
  name            VARCHAR(50)   UNIQUE NOT NULL         ,
  publicView      BOOLEAN       NOT NULL                ,
  publicJoin      BOOLEAN       NOT NULL
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
  leagueID        INT           NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  wins            INT           NOT NULL         ,
  losses          INT           NOT NULL
);
ALTER SEQUENCE teamsIDSeq OWNED BY teams.id;

CREATE TABLE leaguePermissions (
  userID          INT           NOT NULL         ,
  leagueID        INT           NOT NULL         ,
  editPermissions BOOLEAN       NOT NULL         ,
  editTeams       BOOLEAN       NOT NULL         ,
  editUsers       BOOLEAN       NOT NULL         ,
  editSchedule    BOOLEAN       NOT NULL         ,
  editResults     BOOLEAN       NOT NULL
);

CREATE TABLE teamPermissions (
  userID          INT           NOT NULL         ,
  teamID          INT           NOT NULL         ,
  editPermissions BOOLEAN       NOT NULL         ,
  editTeamInfo    BOOLEAN       NOT NULL         ,
  editUsers       BOOLEAN       NOT NULL         ,
  reportResult    BOOLEAN       NOT NULL
);

CREATE TABLE games (
  id              INT           PRIMARY KEY DEFAULT nextval('gamesIDSeq'),
  team1ID         INT                      NOT NULL      ,
  team2ID         INT                      NOT NULL      ,
  gametime        TIMESTAMP WITH TIME ZONE NOT NULL      ,
  complete        BOOLEAN                  NOT NULL      ,
  winnerID        INT                      NOT NULL      ,
  scoreteam1      INT                      NOT NULL      ,
  scoreteam2      INT                      NOT NULL
);
ALTER SEQUENCE gamesIDSeq OWNED BY games.id;

