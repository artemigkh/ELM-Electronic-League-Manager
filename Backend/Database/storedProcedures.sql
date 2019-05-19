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
			players,
			report_results
    )
    VALUES(
      user_id,
      currval('team_id_seq'),
      true,
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
  BEGIN
    UPDATE game SET
      complete = TRUE,
      winner_id = winner_id,
      loser_id = loser_id,
      score_team1 = score_team1,
      score_team2 = score_team2
    WHERE game_id = game_id;

    UPDATE team
      SET wins = wins + 1
    WHERE team_id = winner_id;

    UPDATE team
      SET losses = losses + 1
    WHERE team_id = loser_id;
  END;
$$ LANGUAGE plpgsql;


SELECT create_team(
  8,
  'test team name',
  'tag',
  'test description',
  'smallabc',
  'largeabc',
  1
);


SELECT create_league(
  'new_league_name_2',
  'description',
  true,
  true,
  0,
  0,
  0,
  0,
  'genericsport'
);