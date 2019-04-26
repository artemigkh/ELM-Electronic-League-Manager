CREATE SEQUENCE gameIdSeq;

CREATE TABLE championStats (
  leagueId      INT           NOT NULL                ,
  name          VARCHAR(16)   NOT NULL                ,
  picks         INT           NOT NULL                ,
  wins          INT           NOT NULL                ,
  bans          INT           NOT NULL
);

CREATE TABLE leagueGame (
  id                INT           PRIMARY KEY DEFAULT nextval('gameIdSeq'),
  gameId            INT           NOT NULL                ,
  winTeamId         INT           NOT NULL                ,
  loseTeamId        INT           NOT NULL                ,
  leagueId          INT           NOT NULL                ,
  timestamp         INT           NOT NULL                ,
  duration          FLOAT         NOT NULL
);

CREATE TABLE playerStats (
  id                VARCHAR(50)   NOT NULL                ,
  name              VARCHAR(16)   NOT NULL                ,
  gameId            INT           NOT NULL                ,
  teamId            INT           NOT NULL                ,
  leagueId          INT           NOT NULL                ,
  duration          FLOAT         NOT NULL                ,
  championPicked    VARCHAR(16)   NOT NULL                ,
  gold              FLOAT         NOT NULL                ,
  cs                FLOAT         NOT NULL                ,
  damage            FLOAT         NOT NULL                ,
  kills             FLOAT         NOT NULL                ,
  deaths            FLOAT         NOT NULL                ,
  assists           FLOAT         NOT NULL                ,
  wards             FLOAT         NOT NULL                ,
  win               BOOLEAN       NOT NULL
);

CREATE TABLE teamStats (
  teamId            INT           NOT NULL                ,
  gameId            INT           NOT NULL                ,
  leagueId          INT           NOT NULL                ,
  duration          FLOAT         NOT NULL                ,
  side              INT           NOT NULL                , -- 100 blue 200 red
  firstBlood        BOOLEAN       NOT NULL                ,
  firstTurret       BOOLEAN       NOT NULL                ,
  win               BOOLEAN       NOT NULL
);
