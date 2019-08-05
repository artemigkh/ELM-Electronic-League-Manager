DROP SEQUENCE IF EXISTS lol_tournament_id_seq CASCADE ;
CREATE SEQUENCE lol_tournament_id_seq;
DROP TABLE IF EXISTS lol_tournament CASCADE;
CREATE TABLE lol_tournament (
  lol_tournament_id INT           PRIMARY KEY DEFAULT nextval('lol_tournament_id_seq'),
  league_id         INT           NOT NULL REFERENCES league(league_id),
  provider_id       INT           NOT NULL,
  tournament_id     INT           UNIQUE NOT NULL,
  UNIQUE (lol_tournament_id, tournament_id),
  UNIQUE (league_id, tournament_id)
);

DROP TABLE IF EXISTS lol_champion_stats CASCADE;
CREATE TABLE lol_champion_stats (
  league_id     INT           NOT NULL REFERENCES league(league_id),
  name          VARCHAR(16)   NOT NULL                ,
  picks         INT           NOT NULL                ,
  wins          INT           NOT NULL                ,
  bans          INT           NOT NULL
);

DROP TABLE IF EXISTS lol_player_stats CASCADE;
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

DROP TABLE IF EXISTS lol_team_stats CASCADE;
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

DROP TABLE IF EXISTS lol_game CASCADE;
CREATE TABLE lol_game (
  game_id           INT           NOT NULL REFERENCES game(game_id),
  tournament_code   VARCHAR(64)   NOT NULL         , -- unique in production
  match_id          VARCHAR(64)
);