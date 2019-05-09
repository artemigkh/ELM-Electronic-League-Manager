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
  team_Id         INT           NOT NULL         ,
  userId          INT           UNIQUE           ,
  gameIdentifier  VARCHAR(50)   NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  externalId      VARCHAR(50)                    ,
  mainRoster      BOOLEAN       NOT NULL         ,
  position        VARCHAR(20)
);
ALTER SEQUENCE player_id_seq OWNED BY players.id;

CREATE TABLE team (
  id              INT           PRIMARY KEY DEFAULT nextval('team_id_seq'),
  leagueId        INT           NOT NULL         ,
  name            VARCHAR(50)   NOT NULL         ,
  tag             VARCHAR(5)    NOT NULL         ,
  description     VARCHAR(500)                   ,
  wins            INT           NOT NULL         ,
  losses          INT           NOT NULL         ,
  iconSmall       VARCHAR(20)   NOT NULL         ,
  iconLarge       VARCHAR(20)   NOT NULL
);
ALTER SEQUENCE team_id_seq OWNED BY team.id;

CREATE TABLE league_permissions (
  userId          INT           NOT NULL         ,
  leagueId        INT           NOT NULL         ,
  administrator   BOOLEAN       NOT NULL         ,
  createTeams     BOOLEAN       NOT NULL         ,
  editTeams       BOOLEAN       NOT NULL         ,
  editGames       BOOLEAN       NOT NULL
);

CREATE TABLE team_permissions (
  userId          INT           NOT NULL         ,
  teamId          INT           NOT NULL         ,
  administrator   BOOLEAN       NOT NULL         ,
  information     BOOLEAN       NOT NULL         ,
  players         BOOLEAN       NOT NULL         ,
  reportResults   BOOLEAN       NOT NULL
);

CREATE TABLE game (
  id              INT           PRIMARY KEY DEFAULT nextval('game_id_seq'),
  externalId      VARCHAR(50)              NOT NULL      ,
  leagueId        INT                      NOT NULL      ,
  team1Id         INT                      NOT NULL      ,
  team2Id         INT                      NOT NULL      ,
  gametime        INT                      NOT NULL      ,
  complete        BOOLEAN                  NOT NULL      ,
  winnerId        INT                      NOT NULL      ,
  scoreteam1      INT                      NOT NULL      ,
  scoreteam2      INT                      NOT NULL
);

ALTER SEQUENCE game_id_seq OWNED BY game.id;
CREATE TABLE league_recurring_availability (
  id              INT           PRIMARY KEY DEFAULT nextval('league_recurring_availability_id_seq'),
  league_id       INT           NOT NULL REFERENCES league(id)   ,
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
  leagueId        INT           NOT NULL                ,
  start           INT                                   ,
  end             INT
);
ALTER SEQUENCE league_one_time_availability_id_seq OWNED BY leagueOneTimeAvailabilities.id;
