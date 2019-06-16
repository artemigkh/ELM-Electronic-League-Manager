DROP SEQUENCE IF EXISTS league_id_seq CASCADE ;
CREATE SEQUENCE league_id_seq;
DROP TABLE IF EXISTS league CASCADE;
CREATE TABLE league (
  league_id       INT           PRIMARY KEY DEFAULT nextval('league_id_seq'),
  name            VARCHAR(50)   UNIQUE NOT NULL         ,
  description     VARCHAR(500)                          ,
  markdown_path   VARCHAR(20)   NOT NULL                ,
  public_view     BOOLEAN       NOT NULL                ,
  public_join     BOOLEAN       NOT NULL                ,
  signup_start    INT           NOT NULL                ,
  signup_end      INT           NOT NULL                ,
  league_start    INT           NOT NULL                ,
  league_end      INT           NOT NULL                ,
  game            VARCHAR(30)   NOT NULL
);
ALTER SEQUENCE league_id_seq OWNED BY league.league_id;

DROP SEQUENCE IF EXISTS user_id_seq CASCADE;
CREATE SEQUENCE user_id_seq;
DROP TABLE IF EXISTS user_ CASCADE;
CREATE TABLE user_ (
  user_id         INT           PRIMARY KEY DEFAULT nextval('user_id_seq'),
  email           VARCHAR(256)  UNIQUE NOT NULL  ,
  salt            CHAR(64)      NOT NULL         ,
  hash            CHAR(128)     NOT NULL
);
ALTER SEQUENCE user_id_seq OWNED BY user_.user_id;

DROP SEQUENCE IF EXISTS team_id_seq CASCADE;
CREATE SEQUENCE team_id_seq;
DROP TABLE IF EXISTS team CASCADE;
CREATE TABLE team (
  team_id         INT           PRIMARY KEY DEFAULT nextval('team_id_seq'),
  league_id       INT           NOT NULL REFERENCES league(league_id),
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  description     VARCHAR(500)                   ,
  wins            INT           NOT NULL         ,
  losses          INT           NOT NULL         ,
  icon_small      VARCHAR(20)   NOT NULL         ,
  icon_large      VARCHAR(20)   NOT NULL         ,
  UNIQUE (league_id, name)                       ,
  UNIQUE (league_id, tag)
);
ALTER SEQUENCE team_id_seq OWNED BY team.team_id;

DROP SEQUENCE IF EXISTS player_id_seq CASCADE;
CREATE SEQUENCE player_id_seq;
DROP TABLE IF EXISTS player CASCADE;
CREATE TABLE player (
  player_id       INT           PRIMARY KEY DEFAULT nextval('player_id_seq'),
  team_id         INT           NOT NULL REFERENCES team(team_id),
  league_id       INT           NOT NULL REFERENCES league(league_id),
  user_id         INT           UNIQUE           ,
  game_identifier VARCHAR(50)   NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  external_id     VARCHAR(50)                    ,
  main_roster     BOOLEAN       NOT NULL         ,
  position        VARCHAR(20)                    ,
  UNIQUE (league_id, game_identifier)            ,
  UNIQUE (league_id, external_id)
);
ALTER SEQUENCE player_id_seq OWNED BY player.player_id;

DROP TABLE IF EXISTS league_permissions;
CREATE TABLE league_permissions (
  user_id         INT           NOT NULL REFERENCES user_(user_id),
  league_id       INT           NOT NULL REFERENCES league(league_id),
  administrator   BOOLEAN       NOT NULL         ,
  create_teams    BOOLEAN       NOT NULL         ,
  edit_teams      BOOLEAN       NOT NULL         ,
  edit_games      BOOLEAN       NOT NULL
);

DROP TABLE IF EXISTS team_permissions;
CREATE TABLE team_permissions (
  user_id         INT           NOT NULL REFERENCES user_(user_id),
  team_id         INT           NOT NULL REFERENCES team(team_id),
  administrator   BOOLEAN       NOT NULL         ,
  information     BOOLEAN       NOT NULL         ,
  games           BOOLEAN       NOT NULL
);

DROP SEQUENCE IF EXISTS game_id_seq CASCADE;
CREATE SEQUENCE game_id_seq;
DROP TABLE IF EXISTS game CASCADE;
CREATE TABLE game (
  game_id         INT           PRIMARY KEY DEFAULT nextval('game_id_seq'),
  external_id     VARCHAR(50)                           ,
  league_id       INT                      NOT NULL REFERENCES league(league_id),
  team1_id        INT                      NOT NULL REFERENCES team(team_id),
  team2_id        INT                      NOT NULL REFERENCES team(team_id),
  game_time       INT                      NOT NULL      ,
  complete        BOOLEAN                  NOT NULL      ,
  winner_id       INT                      NOT NULL      ,
  loser_id        INT                      NOT NULL      ,
  score_team1     INT                      NOT NULL      ,
  score_team2     INT                      NOT NULL      ,
  UNIQUE (league_id, external_id)
);
ALTER SEQUENCE game_id_seq OWNED BY game.game_id;

DROP SEQUENCE IF EXISTS availability_id_seq CASCADE;
CREATE SEQUENCE availability_id_seq;
DROP TABLE IF EXISTS availability CASCADE;
CREATE TABLE availability (
  availability_id           INT           PRIMARY KEY DEFAULT nextval('availability_id_seq'),
  league_id                 INT           NOT NULL REFERENCES league(league_id),
  start_time                INT           NOT NULL                ,
  end_time                  INT           NOT NULL                ,
  is_recurring_weekly       BOOLEAN       NOT NULL
);
ALTER SEQUENCE availability_id_seq OWNED BY availability.availability_id;

DROP TABLE IF EXISTS weekly_recurrence CASCADE;
CREATE TABLE weekly_recurrence (
  availability_id           INT           NOT NULL REFERENCES availability(availability_id),
  weekday                   VARCHAR(9)    NOT NULL                ,
  timezone                  INT           NOT NULL                ,
  hour                      SMALLINT      NOT NULL                ,
  minute                    SMALLINT      NOT NULL                ,
  duration                  SMALLINT      NOT NULL
);
