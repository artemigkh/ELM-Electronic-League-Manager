CREATE SEQUENCE league_id_seq;
CREATE SEQUENCE user_id_seq;
CREATE SEQUENCE player_id_seq;
CREATE SEQUENCE team_id_seq;
CREATE SEQUENCE game_id_seq;
CREATE SEQUENCE league_recurring_availability_id_seq;
CREATE SEQUENCE league_one_time_availability_id_seq;

CREATE TABLE league (
  id              INT           PRIMARY KEY DEFAULT nextval('league_id_seq'),
  name            VARCHAR(50)   UNIQUE NOT NULL         ,
  description     VARCHAR(500)                          ,
  markdown_path   VARCHAR(20)   UNIQUE NOT NULL         ,
  public_view     BOOLEAN       NOT NULL                ,
  public_join     BOOLEAN       NOT NULL                ,
  signup_start    INT           NOT NULL                ,
  signup_end      INT           NOT NULL                ,
  league_start    INT           NOT NULL                ,
  league_end      INT           NOT NULL                ,
  game            VARCHAR(30)   NOT NULL
);
ALTER SEQUENCE league_id_seq OWNED BY league.id;

CREATE TABLE user_ (
  id              INT           PRIMARY KEY DEFAULT nextval('user_id_seq'),
  email           VARCHAR(256)  UNIQUE NOT NULL  ,
  salt            CHAR(64)      NOT NULL         ,
  hash            CHAR(128)     NOT NULL
);
ALTER SEQUENCE user_id_seq OWNED BY user_.id;

CREATE TABLE player (
  id              INT           PRIMARY KEY DEFAULT nextval('player_id_seq'),
  team_id         INT           NOT NULL REFERENCES team(id),
  user_id         INT           UNIQUE           ,
  game_identifier VARCHAR(50)   NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  external_id     VARCHAR(50)                    ,
  main_roster     BOOLEAN       NOT NULL         ,
  position        VARCHAR(20)
);
ALTER SEQUENCE player_id_seq OWNED BY players.id;

CREATE TABLE team (
  id              INT           PRIMARY KEY DEFAULT nextval('team_id_seq'),
  league_id       INT           NOT NULL REFERENCES league(id),
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  description     VARCHAR(500)                   ,
  wins            INT           NOT NULL         ,
  losses          INT           NOT NULL         ,
  icon_small      VARCHAR(20)   NOT NULL         ,
  icon_large      VARCHAR(20)   NOT NULL
);
ALTER SEQUENCE team_id_seq OWNED BY team.id;

CREATE TABLE league_permissions (
  user_id         INT           NOT NULL REFERENCES user_(id),
  league_id       INT           NOT NULL REFERENCES league(id),
  administrator   BOOLEAN       NOT NULL         ,
  create_teams    BOOLEAN       NOT NULL         ,
  edit_teams      BOOLEAN       NOT NULL         ,
  edit_games      BOOLEAN       NOT NULL
);

CREATE TABLE team_permissions (
  user_id         INT           NOT NULL REFERENCES user_(id),
  team_id         INT           NOT NULL REFERENCES team(id),
  administrator   BOOLEAN       NOT NULL         ,
  information     BOOLEAN       NOT NULL         ,
  players         BOOLEAN       NOT NULL         ,
  report_results  BOOLEAN       NOT NULL
);

CREATE TABLE game (
  id              INT           PRIMARY KEY DEFAULT nextval('game_id_seq'),
  external_id     VARCHAR(50)              NOT NULL      ,
  league_id       INT                      NOT NULL REFERENCES league(id),
  team1_id        INT                      NOT NULL REFERENCES team(id),
  team2_id        INT                      NOT NULL REFERENCES team(id),
  game_time       INT                      NOT NULL      ,
  complete        BOOLEAN                  NOT NULL      ,
  winner_id       INT                      NOT NULL      ,
  score_team1     INT                      NOT NULL      ,
  score_team2     INT                      NOT NULL
);

ALTER SEQUENCE game_id_seq OWNED BY game.id;
CREATE TABLE league_recurring_availability (
  id              INT           PRIMARY KEY DEFAULT nextval('league_recurring_availability_id_seq'),
  league_id       INT           NOT NULL REFERENCES league(id),
  weekday         SMALLINT      NOT NULL                ,
  timezone        INT           NOT NULL                ,
  hour            SMALLINT      NOT NULL                ,
  minute          SMALLINT      NOT NULL                ,
  duration        SMALLINT      NOT NULL                ,
  constrained     BOOLEAN       NOT NULL                ,
  start_time      INT                                   ,
  end_time        INT
);
ALTER SEQUENCE league_recurring_availability_id_seq OWNED BY league_recurring_availability.id;

CREATE TABLE league_one_time_availability (
  id              INT           PRIMARY KEY DEFAULT nextval('league_one_time_availability_id_seq'),
  league_id       INT           NOT NULL REFERENCES league(id),
  start           INT                                   ,
  end             INT
);
ALTER SEQUENCE league_one_time_availability_id_seq OWNED BY leagueOneTimeAvailabilities.id;
