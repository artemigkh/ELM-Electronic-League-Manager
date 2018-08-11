CREATE SEQUENCE leaguesIdSeq;
CREATE SEQUENCE usersIdSeq;
CREATE SEQUENCE teamsIdSeq;
CREATE SEQUENCE gamesIdSeq;

CREATE TABLE leagues (
  id              INT           PRIMARY KEY DEFAULT nextval('leaguesIdSeq'),
  name            VARCHAR(50)   UNIQUE NOT NULL         ,
  publicView      BOOLEAN       NOT NULL                ,
  publicJoin      BOOLEAN       NOT NULL
);
ALTER SEQUENCE leaguesIdSeq OWNED BY leagues.id;

CREATE TABLE users (
  id              INT           PRIMARY KEY DEFAULT nextval('usersIdSeq'),
  email           VARCHAR(256)  UNIQUE NOT NULL  ,
  salt            CHAR(64)      NOT NULL         ,
  hash            CHAR(128)     NOT NULL
);
ALTER SEQUENCE usersIdSeq OWNED BY users.id;

CREATE TABLE teams (
  id              INT           PRIMARY KEY DEFAULT nextval('teamsIdSeq'),
  leagueId        INT           NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  wins            INT           NOT NULL         ,
  losses          INT           NOT NULL
);
ALTER SEQUENCE teamsIdSeq OWNED BY teams.id;

CREATE TABLE leaguePermissions (
  userId          INT           NOT NULL         ,
  leagueId        INT           NOT NULL         ,
  editPermissions BOOLEAN       NOT NULL         ,
  editTeams       BOOLEAN       NOT NULL         ,
  editUsers       BOOLEAN       NOT NULL         ,
  editSchedule    BOOLEAN       NOT NULL         ,
  editResults     BOOLEAN       NOT NULL
);

CREATE TABLE teamPermissions (
  userId          INT           NOT NULL         ,
  teamId          INT           NOT NULL         ,
  editPermissions BOOLEAN       NOT NULL         ,
  editTeamInfo    BOOLEAN       NOT NULL         ,
  editUsers       BOOLEAN       NOT NULL         ,
  reportResult    BOOLEAN       NOT NULL
);

CREATE TABLE games (
  id              INT           PRIMARY KEY DEFAULT nextval('gamesIdSeq'),
  leagueId        INT                      NOT NULL      ,
  team1Id         INT                      NOT NULL      ,
  team2Id         INT                      NOT NULL      ,
  gametime        INT                      NOT NULL      ,
  complete        BOOLEAN                  NOT NULL      ,
  winnerId        INT                      NOT NULL      ,
  scoreteam1      INT                      NOT NULL      ,
  scoreteam2      INT                      NOT NULL
);
ALTER SEQUENCE gamesIdSeq OWNED BY games.id;

