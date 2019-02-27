package game

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// CardsPerUser is a maximum count of cards user may hold
const CardsPerUser = 3

// Stage is a game's blocking state awaiting for soma Action to be performed
type Stage string

const (
	// StageNewGame – game created bot round has not started
	StageNewGame Stage = "new-game"
	// StageAwaitResponse — round started and we are waiting for all responses to be received from users
	StageAwaitResponse Stage = "awaiting-response"
	// StageAwaitWinner — all responses received, waiting for host to pick the winner
	StageAwaitWinner Stage = "awaiting-winner"
	// StageFinal – desired count of rounds has been played
	StageFinal Stage = "final"
)

// Game is a state of a game with mehods
type Game struct {
	id           string
	stage        Stage
	roundsTotal  int
	roundsPlayed int
	question     *BlackCard
	answers      map[*Player]*WhiteCard
	host         *Player
	players      []*Player
	blackDeck    []*BlackCard
	whiteDeck    []*WhiteCard
}

// New creates new Game with defult configuration
func New(numRounds int, blackDeck []*BlackCard, whiteDeck []*WhiteCard) *Game {
	return &Game{
		id:          uuid.NewV4().String(),
		stage:       StageNewGame,
		roundsTotal: numRounds,
		answers:     map[*Player]*WhiteCard{},
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("<game @ %s, round %d/%d, %d players>", g.stage, g.roundsPlayed+1, g.roundsTotal, len(g.players))
}

// GetPlayer checks if player wth given ID is participating the game
func (g *Game) GetPlayer(id string) *Player {
	for _, p := range g.players {
		if p.id == id {
			return p
		}
	}
	return nil
}

// nextHost iterates over players returning next to the host player
func (g *Game) nextHost() (*Player, error) {
	if len(g.players) == 0 {
		return nil, errors.Errorf("game has no players")
	}
	if g.host == nil {
		return g.players[0], nil
	}
	for i, p := range g.players {
		if p.id == g.host.id {
			return g.players[(i+1)%len(g.players)], nil
		}
	}
	return nil, errors.Errorf("game host is unsyned with players list")
}

// pullBlackCard picks random card from black deck, returns card and rest of the deck
func (g *Game) pullBlackCard() (selected *BlackCard, rest []*BlackCard) {
	r := rand.Intn(len(g.blackDeck))
	for i, card := range g.blackDeck {
		if i == r {
			selected = card
			continue
		}
		rest = append(rest, card)
	}
	return selected, rest
}

// pullWhiteCard picks random card from white deck, returns card and rest of the deck
func (g *Game) pullWhiteCard() (selected *WhiteCard, rest []*WhiteCard) {
	r := rand.Intn(len(g.whiteDeck))
	for i, card := range g.whiteDeck {
		if i == r {
			selected = card
			continue
		}
		rest = append(rest, card)
	}
	return selected, rest
}
