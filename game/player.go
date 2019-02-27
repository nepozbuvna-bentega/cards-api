package game

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Player is representing a game participant
type Player struct {
	id    string
	name  string
	score int
	cards []*WhiteCard
}

// NewPlayer creates new Player with given name and unique ID
func NewPlayer(name string) *Player {
	return &Player{
		id:   uuid.NewV4().String(),
		name: name,
	}
}

func (p *Player) String() string {
	return fmt.Sprintf("<player:%s:%s>", p.name, p.id)
}

// GetCard returns card with given ID if user has it, otherwise ni is returned
func (p *Player) GetCard(cardID string) *WhiteCard {
	for _, card := range p.cards {
		if card.id == cardID {
			return card
		}
	}
	return nil
}
