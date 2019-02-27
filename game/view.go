package game

import "github.com/pkg/errors"

// View is a player's view of the game
type View struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}

// NewView create a View of given game for particular player
func NewView(g *Game, playerID string) (*View, error) {
	if g.GetPlayer(playerID) == nil {
		return nil, errors.Errorf(`unknown player "%s"`, playerID)
	}
	return &View{
		ID:   g.id,
		Host: g.host.id,
	}, nil
}
