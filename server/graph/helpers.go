package graph

import (
	"fantasy-draft/graph/model"
)

// playerRows interface for scanning player rows from database
type playerRows interface {
	Next() bool
	Scan(dest ...any) error
}

// scanPlayers scans rows from the database into Player models
func scanPlayers(rows playerRows) ([]*model.Player, error) {
	var players []*model.Player
	for rows.Next() {
		var p model.Player
		var pos string
		var status string

		if err := rows.Scan(
			&p.ID, &p.FirstName, &p.LastName, &pos, &p.TeamID,
			&p.Height, &p.Weight, &p.Age, &p.YearsOfExperience,
			&p.DraftYear, &p.JerseyNumber, &status, &p.Skill,
		); err != nil {
			return nil, err
		}

		// Convert string to enum
		p.Position = model.Position(pos)
		p.Status = model.PlayerStatus(status)

		players = append(players, &p)
	}
	return players, nil
}
