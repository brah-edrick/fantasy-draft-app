package main

import "fmt"

func createTeamRoster(teamID string) FootballTeamRoster {
	qbCount := NFLRosterComposition["QB"]
	rbCount := NFLRosterComposition["RB"]
	wrCount := NFLRosterComposition["WR"]
	teCount := NFLRosterComposition["TE"]
	pkCount := NFLRosterComposition["PK"]

	qbPlayers := make([]Player, qbCount)
	for qbIndex := range qbCount {
		qbPlayers[qbIndex] = createNewPlayer(QB, teamID)
	}
	rbPlayers := make([]Player, rbCount)
	for rbIndex := range rbCount {
		rbPlayers[rbIndex] = createNewPlayer(RB, teamID)
	}
	wrPlayers := make([]Player, wrCount)
	for wrIndex := range wrCount {
		wrPlayers[wrIndex] = createNewPlayer(WR, teamID)
	}
	tePlayers := make([]Player, teCount)
	for teIndex := range teCount {
		tePlayers[teIndex] = createNewPlayer(TE, teamID)
	}
	pkPlayers := make([]Player, pkCount)
	for pkIndex := range pkCount {
		pkPlayers[pkIndex] = createNewPlayer(PK, teamID)
	}

	roster := FootballTeamRoster{
		QB: qbPlayers,
		RB: rbPlayers,
		WR: wrPlayers,
		TE: tePlayers,
		PK: pkPlayers,
	}

	fmt.Printf("Roster created: %+v\n", roster)
	return roster
}
