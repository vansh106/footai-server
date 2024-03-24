package DB

import (
	"fmt"
)

func GenSqlQuery(prompt string) string{
	return fmt.Sprintf(`%s 
							 %s
							 %s
							 This is the database schema, Return SQL query for the data that could be required for inference for below question.
							 NOTE: Column names in the query should be in inverted commas.
							 NOTE: Query should ignore all the NULL values.
							 NOTE: Player names might be in short form so use substring but don't use substring to reduce the name size, Return full row.
							 NOTE: No. of rows in output should be never be more than 20.
							 NOTE: Team names and player names should be case insensitive in query.
							 NOTE: Give out statistics for all the teams or players asked in the prompt.
							 %s`, createTableQuery, playerQuery, fixturesQuery, prompt)
}

func GenPrompt(prompt string, rows []map[string]interface{}) string{
	return fmt.Sprintf(`%s 
						Answer the above question using below data.\n
						`, prompt, rows)
}

var createTableQuery string= `
	CREATE TABLE IF NOT EXISTS pl_standings (
		"rank" INTEGER,
		"team__id" INTEGER PRIMARY KEY,
		"team__name" VARCHAR(100),
		"points" INTEGER,
		"goalsDiff" INTEGER,
		"form" VARCHAR(10),
		"description" VARCHAR(255),
		"all__played" INTEGER,
		"all__win" INTEGER,
		"all__draw" INTEGER,
		"all__lose" INTEGER,
		"all__goals__for" INTEGER,
		"all__goals__against" INTEGER,
		"home__played" INTEGER,
		"home__win" INTEGER,
		"home__draw" INTEGER,
		"home__lose" INTEGER,
		"home__goals__for" INTEGER,
		"home__goals__against" INTEGER,
		"away__played" INTEGER,
		"away__win" INTEGER,
		"away__draw" INTEGER,
		"away__lose" INTEGER,
		"away__goals__for" INTEGER,
		"away__goals__against" INTEGER
	);`
	
var	playerQuery string = `CREATE TABLE IF NOT EXISTS players (
		"Id" INTEGER,
		"Name" VARCHAR(255),
		"Age" INTEGER,
		"Nationality" VARCHAR(255),
		"Injured" BOOLEAN,
		"Team_Name" VARCHAR(255),
		"Team_Id" INTEGER,
		"Position" VARCHAR(255),
		"Games" INTEGER,
		"Minutes" INTEGER,
		"Accuracy_Passes" INTEGER,
		"Key_Passes" INTEGER,
		"Total_Passes" INTEGER,
		"Shots_On" INTEGER,
		"Shots_Total" INTEGER,
		"Dribbles_Attempts" INTEGER,
		"Dribbles_Past" INTEGER,
		"Dribbles_Success" INTEGER,
		"Fouls_Drawn" INTEGER,
		"Fouls_Committed" INTEGER,
		"Tackled_Block" INTEGER,
		"Tackled_Intercept" INTEGER,
		"Tackled_Total" INTEGER,
		"Duels_Won" INTEGER,
		"Duels_Total" INTEGER,
		"Goals_Assist" INTEGER,
		"Goals_Total" INTEGER,
		"Goals_Conceded" INTEGER,
		"Goals_Saves" INTEGER,
		"Rating" DECIMAL(5,2),
		"Yellow_Cards" INTEGER,
		"Red_Cards" INTEGER,
		"Yellow_Red_Cards" INTEGER,
		"Captain" BOOLEAN,
		"Weight_kg" INTEGER,
		"Height_cm" INTEGER
	);`

var	fixturesQuery string = `CREATE TABLE IF NOT EXISTS fixtures (
		fixture_id INTEGER PRIMARY KEY,
		fixture_referee VARCHAR(255),
		fixture_timezone VARCHAR(255),
		fixture_date DATE,
		fixture_timestamp INTEGER,
		fixture_periods_first INTEGER,
		fixture_periods_second INTEGER,
		fixture_venue_id INTEGER,
		fixture_venue_name VARCHAR(255),
		fixture_venue_city VARCHAR(255),
		fixture_status_long VARCHAR(255),
		fixture_status_short VARCHAR(255),
		fixture_status_elapsed INTEGER,
		teams_home_id INTEGER,
		teams_home_name VARCHAR(255),
		teams_home_winner BOOLEAN,
		teams_away_id INTEGER,
		teams_away_name VARCHAR(255),
		teams_away_winner BOOLEAN,
		goals_home INTEGER,
		goals_away INTEGER,
		score_halftime_home INTEGER,
		score_halftime_away INTEGER,
		score_fulltime_home INTEGER,
		score_fulltime_away INTEGER
	);
	`
	/* content := fmt.Sprintf(`
		You are an assistant that generates JSONL prompts based off of JSON data for fine tuning.
		Each response should be formatted as:

		Please generate 10 most important questions on the basis of the their postion based off of the JSON and provide it in JSON format without square brackets.
		Each response should come from the following JSON:
		%s
	  `, jsonData2) */

