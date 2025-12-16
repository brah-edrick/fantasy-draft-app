package main

import "fmt"

func createTeamRoster(teamID string) FootballTeamRoster {
	qbCount := NFLRosterComposition["QB"]
	rbCount := NFLRosterComposition["RB"]
	wrCount := NFLRosterComposition["WR"]
	teCount := NFLRosterComposition["TE"]
	pkCount := NFLRosterComposition["PK"]

	// Create players with depth-based skill assignments
	qbPlayers := createPlayersWithDepthSkills(QB, teamID, qbCount)
	rbPlayers := createPlayersWithDepthSkills(RB, teamID, rbCount)
	wrPlayers := createPlayersWithDepthSkills(WR, teamID, wrCount)
	tePlayers := createPlayersWithDepthSkills(TE, teamID, teCount)
	pkPlayers := createPlayersWithDepthSkills(PK, teamID, pkCount)

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

func createPlayersWithDepthSkills(position Position, teamID string, count int) []Player {
	players := make([]Player, count)
	for depthIndex := range count {
		player := createNewPlayer(position, teamID)
		// Override the random skill with depth-based skill
		player.Skill = createSkillForDepthPosition(depthIndex, count)
		players[depthIndex] = player
	}
	return players
}
