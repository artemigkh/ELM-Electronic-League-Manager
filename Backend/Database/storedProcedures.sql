CREATE OR REPLACE FUNCTION
create_league(
  name         VARCHAR(50),
  description  VARCHAR(500),
  public_view  BOOLEAN,
  public_join  BOOLEAN,
  signup_start INT,
  signup_end   INT,
  league_start INT,
  league_end   INT,
  game         VARCHAR(30),
  user_id      INT
)
RETURNS INT AS $$
  BEGIN
    INSERT INTO league(
      name,
      description,
      markdown_path,
      public_view,
      public_join,
      signup_start,
      signup_end,
      league_start,
      league_end,
      game
    )
    VALUES (
      name,
      description,
      '',
      public_view,
      public_join,
      signup_start,
      signup_end,
      league_start,
      league_end,
      game
    );
    INSERT INTO league_permissions(
      user_id,
      league_id,
			administrator,
			create_teams,
			edit_teams,
			edit_games
    )
    VALUES(
      user_id,
      currval('league_id_seq'),
      true,
      true,
      true,
      true
    );
    RETURN currval('league_id_seq');
  END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
create_team(
  league_id       INT,
  name            VARCHAR(50),
  tag             VARCHAR(5),
  description     VARCHAR(500) ,
  icon_small      VARCHAR(20),
  icon_large      VARCHAR(20),
  user_id         INT
)
RETURNS INT AS $$
  BEGIN
    INSERT INTO team(
      league_id,
      name,
      tag,
      description,
      wins,
      losses,
      icon_small,
      icon_large
    )
    VALUES (
      league_id,
      name,
      tag,
      description,
      0,
      0,
      icon_small,
      icon_large
    );
    INSERT INTO team_permissions(
      user_id,
      team_id,
			administrator,
			information,
			games
    )
    VALUES(
      user_id,
      currval('team_id_seq'),
      true,
      true,
      true
    );
    RETURN currval('team_id_seq');
  END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION
report_game(
  game_id         INT,
  winner_id       INT,
  loser_id        INT,
  score_team1     INT,
  score_team2     INT
)
RETURNS VOID AS $$
  DECLARE game_complete BOOLEAN;
  DECLARE old_winner_id INT;
  DECLARE old_loser_id INT;
  BEGIN
    SELECT game.complete, game.winner_id, game.loser_id INTO game_complete, old_winner_id, old_loser_id
      FROM game WHERE game.game_id = report_game.game_id;
    IF (game_complete = TRUE) THEN
      UPDATE team
        SET wins = wins - 1
      WHERE team_id = old_winner_id;

      UPDATE team
        SET losses = losses - 1
      WHERE team_id = old_loser_id;
    END IF;

    UPDATE game SET
      complete = TRUE,
      winner_id = report_game.winner_id,
      loser_id = report_game.loser_id,
      score_team1 = report_game.score_team1,
      score_team2 = report_game.score_team2
    WHERE game.game_id = report_game.game_id;

    UPDATE team
      SET wins = wins + 1
    WHERE team_id = winner_id;

    UPDATE team
      SET losses = losses + 1
    WHERE team_id = loser_id;
  END;
$$ LANGUAGE plpgsql;

CREATE TYPE league_game_ids AS (league_id INT, game_id INT);
CREATE OR REPLACE FUNCTION
report_game_by_external_id(
  external_id     VARCHAR(50),
  winner_id       INT,
  loser_id        INT,
  score_team1     INT,
  score_team2     INT
)
RETURNS SETOF league_game_ids AS $$
  DECLARE league_id INT;
  DECLARE game_id INT;
  DECLARE game_complete BOOLEAN;
  DECLARE old_winner_id INT;
  DECLARE old_loser_id INT;
  BEGIN
    SELECT game.league_id, game.game_id, game.complete, game.winner_id, game.loser_id INTO league_id, game_id, game_complete, old_winner_id, old_loser_id
      FROM game WHERE game.external_id = report_game_by_external_id.external_id;
    IF (game_complete = TRUE) THEN
      UPDATE team
        SET wins = wins - 1
      WHERE team_id = old_winner_id;

      UPDATE team
        SET losses = losses - 1
      WHERE team_id = old_loser_id;
    END IF;

    UPDATE game SET
      complete = TRUE,
      winner_id = report_game_by_external_id.winner_id,
      loser_id = report_game_by_external_id.loser_id,
      score_team1 = report_game_by_external_id.score_team1,
      score_team2 = report_game_by_external_id.score_team2
    WHERE game.external_id = report_game_by_external_id.external_id;

    UPDATE team
      SET wins = wins + 1
    WHERE team_id = winner_id;

    UPDATE team
      SET losses = losses + 1
    WHERE team_id = loser_id;

    RETURN QUERY (SELECT league_id, game_id);
  END;
$$ LANGUAGE plpgsql;

