CREATE TABLE championStats (
  championId    INT           NOT NULL                ,
  leagueId      INT           NOT NULL                ,
  name          VARCHAR(16)   NOT NULL                ,
  picks         INT           NOT NULL                ,
  wins          INT           NOT NULL                ,
  bans          INT           NOT NULL
);

CREATE TABLE playerStats (
  playerId          INT           NOT NULL                ,
  leagueId          INT           NOT NULL                ,
  numGames          INT           NOT NULL                ,
  goldPerMinute     FLOAT         NOT NULL                ,
  csPerMinute       FLOAT         NOT NULL                ,
  damagePerMinute   FLOAT         NOT NULL                ,
  kills             FLOAT         NOT NULL                ,
  deaths            FLOAT         NOT NULL                ,
  assists           FLOAT         NOT NULL                ,
  visionWards       FLOAT         NOT NULL                ,
  controlWards      FLOAT         NOT NULL
);

CREATE TABLE teamStats (
  teamId            INT           NOT NULL                ,
  leagueId          INT           NOT NULL                ,
  numGames          INT           NOT NULL                ,
  teamKDA           FLOAT         NOT NULL                ,
  gameTime          FLOAT         NOT NULL                ,
  firstBloods       FLOAT         NOT NULL                ,
  firstTurrets      FLOAT         NOT NULL
);
