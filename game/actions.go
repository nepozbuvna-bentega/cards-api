package game

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ErrBadAction means action can't be applied to the game at current stage
var ErrBadAction = fmt.Errorf("action is not applicable")

// Action is making changes to the game state.
// Action failes if not applicable to game.stage or other game state params
type Action func(g *Game) error

// ActionJoin adds new Player by name
func ActionJoin(name string) Action {
	return func(g *Game) error {
		if g.stage != StageNewGame {
			return errors.Wrapf(ErrBadAction, "can't join the game @ %s", g.stage)
		}

		name = strings.TrimSpace(name)
		for _, p := range g.players {
			if strings.EqualFold(p.name, name) {
				return errors.Errorf(`user with name "%s" already exists`, name)
			}
		}

		g.players = append(g.players, NewPlayer(name))
		return nil
	}
}

// ActionStartRound starts new round:
// - pulls the black card
// - spawns white cards to players
// - switches the game to StageAwaitResponse
func ActionStartRound(callerID string) Action {
	return func(g *Game) error {
		if g.stage != StageNewGame && g.stage != StageFinal && g.stage != StageAwaitWinner {
			return errors.Wrapf(ErrBadAction, "can't start a round @ %s", g.stage)
		}

		player := g.GetPlayer(callerID)
		if player == nil {
			return errors.Errorf(`player "%s" is not participating the game`, callerID)
		}
		if callerID != g.host.id {
			return errors.Errorf(`player "%s" is not a host and can't start the game`, player)
		}

		g.question, g.blackDeck = g.pullBlackCard()
		if g.question == nil {
			return errors.Errorf("out of black cards")
		}

		var donePulling bool
		for !donePulling {
			donePulling = true
			for _, p := range g.players {
				if len(p.cards) == CardsPerUser {
					continue
				}

				card, rest := g.pullWhiteCard()
				if card == nil {
					donePulling = true
					break
				}

				p.cards = append(p.cards, card)
				g.whiteDeck = rest

				if len(p.cards) < CardsPerUser {
					donePulling = false
				}
			}
		}

		g.stage = StageAwaitResponse

		return nil
	}
}

// ActionAnswer is providing player's answer.
// Once all players answered game switches to StageAwaitWinner
func ActionAnswer(callerID string, cardID string) Action {
	return func(g *Game) error {
		if g.stage != StageAwaitResponse {
			return errors.Wrapf(ErrBadAction, "can't answer @ %s", g.stage)
		}

		player := g.GetPlayer(callerID)
		if player == nil {
			return errors.Errorf(`player "%s" is not participating the game`, callerID)
		}
		if callerID == g.host.id {
			return errors.Errorf("host is not allowed to answer")
		}

		card := player.GetCard(cardID)
		if card == nil {
			return errors.Errorf(`player %s doesn't have card "%s"`, player, cardID)
		}

		g.answers[player] = card

		if len(g.answers) == len(g.players) {
			g.stage = StageAwaitWinner
		}

		return nil
	}
}

// ActionAnnounceWinner is setting round winner
func ActionAnnounceWinner(callerID, playerID string) Action {
	return func(g *Game) error {
		if g.stage != StageAwaitWinner {
			return errors.Wrapf(ErrBadAction, "can't announe round winner @ %s", g.stage)
		}
		if callerID != g.host.id {
			return errors.Errorf("only host is allowed to to announce the winner")
		}
		if callerID != playerID {
			return errors.Errorf("host can't be a winner")
		}

		player := g.GetPlayer(callerID)
		if player == nil {
			return errors.Errorf(`player "%s" is not participating the game`, callerID)
		}

		player.score++
		return nil
	}
}
