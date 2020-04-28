package poner_test

import (
	"testing"

	"github.com/blakecallens/poner"
)

func TestSim2Players(t *testing.T) {
	players := []poner.Player{
		{Name: "Bob", IsComputer: true, SkillLevel: 4},
		{Name: "Sue", IsComputer: false, SkillLevel: 4},
	}
	game := poner.Game{}
	game.New(players)

	// Wait for a winner
	for game.Winner == nil {
		err := playRound(&game)
		if err != nil {
			t.Errorf("Error simulating game: %v", err)
			return
		}
		if game.Winner != nil {
			break
		}
		scorePlayerHands(&game)
	}
}

func TestSim3Players(t *testing.T) {
	players := []poner.Player{
		{Name: "Bob", IsComputer: true, SkillLevel: 4},
		{Name: "Sue", IsComputer: true, SkillLevel: 4},
		{Name: "Dan", IsComputer: true, SkillLevel: 3},
	}
	game := poner.Game{}
	game.New(players)

	// Wait for a winner
	for game.Winner == nil {
		err := playRound(&game)
		if err != nil {
			t.Errorf("Error simulating game: %v", err)
			return
		}
		if game.Winner != nil {
			break
		}
		scorePlayerHands(&game)
	}
}

func TestSim4Players(t *testing.T) {
	players := []poner.Player{
		{Name: "Bob", IsComputer: true, SkillLevel: 4},
		{Name: "Sue", IsComputer: true, SkillLevel: 4},
		{Name: "Dan", IsComputer: true, SkillLevel: 3},
		{Name: "Joe", IsComputer: true, SkillLevel: 2},
	}
	game := poner.Game{}
	game.New(players)

	// Wait for a winner
	for game.Winner == nil {
		err := playRound(&game)
		if err != nil {
			t.Errorf("Error simulating game: %v", err)
			return
		}
		if game.Winner != nil {
			break
		}
		scorePlayerHands(&game)
	}
}

func TestSim5Players(t *testing.T) {
	players := []poner.Player{
		{Name: "Bob", IsComputer: true, SkillLevel: 4},
		{Name: "Sue", IsComputer: true, SkillLevel: 4},
		{Name: "Dan", IsComputer: true, SkillLevel: 3},
		{Name: "Joe", IsComputer: true, SkillLevel: 2},
		{Name: "Foo", IsComputer: true, SkillLevel: 1},
	}
	game := poner.Game{}
	game.New(players)

	err := playRound(&game)
	if err == nil {
		t.Errorf("Error simulating game, no error for 5 players")
	}
}

func playRound(game *poner.Game) (err error) {
	// Start a new round an get his heels, if drawn
	_, err = game.NextRound()
	if err != nil {
		return
	}
	for ii := range game.Players {
		if !game.Players[ii].IsComputer {
			game.Players[ii].Discard = game.Players[ii].DealtHand.GetBestDiscard(&game.Deck, game.Dealer == ii)
			playingHand := poner.Hand{}
			for _, card := range game.Players[ii].Discard.Held {
				playingHand = append(playingHand, card)
			}
			game.Players[ii].PlayingHand = playingHand
		}
	}
	// Wait for all players to be out of cards
	for !game.AllPlaysDone() {
		// Field is the current playfield that players put their cards into
		game.ResetField()
		for !game.AllPlayersGone() {
			// Go to the next player and get their play
			isHuman, _, _, err := game.NextPlayer()
			if err != nil {
				break
			}

			player := &game.Players[game.ActivePlayer]
			if isHuman {
				nextPlayer := game.ActivePlayer + 1
				if nextPlayer >= len(game.Players) {
					nextPlayer = 0
				}
				bestCard, cantPlay := player.PlayingHand.GetBestPlay(game.Field, game.Players[nextPlayer])
				if cantPlay {
					game.HumanPlayGone()
				} else {
					game.HumanPlayCard(bestCard)
				}
			} else if player.Gone {
				continue
			}
			// Stop the round if the player has reached 121
			if game.Winner != nil {
				break
			}
		}
		if err != nil || game.Winner != nil {
			break
		}
		// GoScore will return a score for go, if the total isn't 31
		game.GoScore()
	}
	return
}

func scorePlayerHands(game *poner.Game) {
	ii := game.Dealer + 1
	scoresCounted := 0
	var player *poner.Player
	// Start with the player left of the dealer and make your way around
	for scoresCounted < len(game.Players) {
		if ii >= len(game.Players) {
			ii = 0
		}
		player = &game.Players[ii]
		// Returns the individual scores and the total value of the hand
		game.ScoreHand(player, false)
		// Stop the round if the player has reached 121
		if game.Winner != nil {
			return
		}
		scoresCounted++
		ii++
	}
	// Count the crib
	game.ScoreHand(player, true)
}
