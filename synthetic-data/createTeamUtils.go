package main

import "fmt"

func createTeamRoster(teamID string) NFLTeamRoster {
	qbCount := NFLRosterComposition["QB"]
	rbCount := NFLRosterComposition["RB"]
	wrCount := NFLRosterComposition["WR"]
	teCount := NFLRosterComposition["TE"]
	pkCount := NFLRosterComposition["PK"]

	qbPlayers := make([]Player, qbCount)
	for qbIndex := 0; qbIndex < qbCount; qbIndex++ {
		qbPlayers[qbIndex] = createNewPlayer(QB, teamID)
	}
	rbPlayers := make([]Player, rbCount)
	for rbIndex := 0; rbIndex < rbCount; rbIndex++ {
		rbPlayers[rbIndex] = createNewPlayer(RB, teamID)
	}
	wrPlayers := make([]Player, wrCount)
	for wrIndex := 0; wrIndex < wrCount; wrIndex++ {
		wrPlayers[wrIndex] = createNewPlayer(WR, teamID)
	}
	tePlayers := make([]Player, teCount)
	for teIndex := 0; teIndex < teCount; teIndex++ {
		tePlayers[teIndex] = createNewPlayer(TE, teamID)
	}
	pkPlayers := make([]Player, pkCount)
	for pkIndex := 0; pkIndex < pkCount; pkIndex++ {
		pkPlayers[pkIndex] = createNewPlayer(PK, teamID)
	}

	roster := NFLTeamRoster{
		QB: qbPlayers,
		RB: rbPlayers,
		WR: wrPlayers,
		TE: tePlayers,
		PK: pkPlayers,
	}

	fmt.Printf("Roster created: %+v\n", roster)
	return roster
}
