CREATE TABLE lol_champion_stats (
  league_id     INT           NOT NULL REFERENCES league(league_id),
  name          VARCHAR(16)   NOT NULL                ,
  picks         INT           NOT NULL                ,
  wins          INT           NOT NULL                ,
  bans          INT           NOT NULL
);

CREATE TABLE lol_player_stats (
  id                VARCHAR(50)   NOT NULL                ,
  name              VARCHAR(16)   NOT NULL                ,
  game_id           INT           NOT NULL REFERENCES game(game_id),
  team_id           INT           NOT NULL REFERENCES team(team_id),
  league_id         INT           NOT NULL REFERENCES league(league_id),
  duration          FLOAT         NOT NULL                ,
  champion_picked   VARCHAR(16)   NOT NULL                ,
  gold              FLOAT         NOT NULL                ,
  cs                FLOAT         NOT NULL                ,
  damage            FLOAT         NOT NULL                ,
  kills             FLOAT         NOT NULL                ,
  deaths            FLOAT         NOT NULL                ,
  assists           FLOAT         NOT NULL                ,
  wards             FLOAT         NOT NULL                ,
  win               BOOLEAN       NOT NULL
);

CREATE TABLE lol_team_stats (
  team_id           INT           NOT NULL REFERENCES team(team_id),
  game_id           INT           NOT NULL REFERENCES game(game_id),
  league_id         INT           NOT NULL REFERENCES league(league_id),
  duration          FLOAT         NOT NULL                ,
  side              INT           NOT NULL                , -- 100 blue 200 red
  first_blood       BOOLEAN       NOT NULL                ,
  first_turret      BOOLEAN       NOT NULL                ,
  win               BOOLEAN       NOT NULL
);
