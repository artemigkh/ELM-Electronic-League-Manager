CREATE SEQUENCE leaguesIdSeq;
CREATE SEQUENCE usersIdSeq;
CREATE SEQUENCE playersIdSeq;
CREATE SEQUENCE teamsIdSeq;
CREATE SEQUENCE gamesIdSeq;

CREATE TABLE leagues (
  id              INT           PRIMARY KEY DEFAULT nextval('leaguesIdSeq'),
  name            VARCHAR(50)   UNIQUE NOT NULL         ,
  description     VARCHAR(500)                          ,
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

CREATE TABLE players (
  id              INT           PRIMARY KEY DEFAULT nextval('playersIdSeq'),
  teamId          INT           NOT NULL         ,
  userId          INT           UNIQUE           ,
  gameIdentifier  VARCHAR(50)   NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  mainRoster      BOOLEAN       NOT NULL
);
ALTER SEQUENCE playersIdSeq OWNED BY players.id;

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
  administrator   BOOLEAN       NOT NULL         ,
  createTeams     BOOLEAN       NOT NULL         ,
  editTeams       BOOLEAN       NOT NULL         ,
  editGames       BOOLEAN       NOT NULL
);

-- TODO: if efficiency a problem, add leagueID for faster filter
CREATE TABLE teamPermissions (
  userId          INT           NOT NULL         ,
  teamId          INT           NOT NULL         ,
  administrator   BOOLEAN       NOT NULL         ,
  information     BOOLEAN       NOT NULL         ,
  players         BOOLEAN       NOT NULL         ,
  reportResults   BOOLEAN       NOT NULL
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

CREATE TABLE teamManagerInviteCodes (
  code            CHAR(16)      UNIQUE NOT NULL   ,
  creationTime    INT           NOT NULL          ,
  leagueId        INT           NOT NULL          ,
  teamId          INT           NOT NULL          ,
  administrator   BOOLEAN       NOT NULL          ,
  information     BOOLEAN       NOT NULL          ,
  players         BOOLEAN       NOT NULL          ,
  reportResults   BOOLEAN       NOT NULL
);

CREATE TABLE leagueManagerInviteCodes (
  code            CHAR(16)      UNIQUE NOT NULL  ,
  creationTime    INT           NOT NULL         ,
  leagueId        INT           NOT NULL         ,
  editPermissions BOOLEAN       NOT NULL         ,
  createTeams     BOOLEAN       NOT NULL         ,
  editTeams       BOOLEAN       NOT NULL         ,
  editUsers       BOOLEAN       NOT NULL         ,
  editSchedule    BOOLEAN       NOT NULL         ,
  editResults     BOOLEAN       NOT NULL
);

CREATE TABLE teamManagerJoinRequests (
  userId          INT           NOT NULL          ,
  teamId          INT           NOT NULL
);

CREATE TABLE leagueManagerJoinRequests (
  userId          INT           NOT NULL          ,
  leagueId        INT           NOT NULL
);