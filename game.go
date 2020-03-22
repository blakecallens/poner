package poner

import (
	"math/rand"
	"time"
)

// Game represents
type Game struct {
	Players      []Player
	Round        int
	ToWin        int
	Dealer       int
	ActivePlayer int
	Deck         Deck
	Starter      Card
	Field        Hand
	Crib         Hand
}

// New creates a new game
func (game *Game) New(players []Player) {
	rand.Seed(time.Now().UnixNano())
	game.Players = players
	game.Round = 0
	game.Dealer = rand.Intn(len(game.Players))
	if game.ToWin == 0 {
		game.ToWin = 121
	}
}

// NextRound starts a new game round
func (game *Game) NextRound() (err error) {
	game.Round++
	game.Dealer++
	if game.Dealer >= len(game.Players) {
		game.Dealer = 0
	}
	game.ActivePlayer = game.Dealer

	game.Field = Hand{}
	game.Deck = Deck{}.New()
	game.Deck.Shuffle()
	game.Deck.Cut()

	hands, _ := game.Deck.DealCribbage(len(game.Players))
	for index, hand := range hands {
		game.Players[index].TakeDeal(hand, &game.Deck, index == game.Dealer)
	}

	discards := []Discard{}
	for _, player := range game.Players {
		discards = append(discards, player.Discard)
	}
	game.Crib, err = BuildCrib(discards, &game.Deck)
	game.Starter, _ = game.Deck.PullFromTop()
	return
}

// AllPlayersGone returns whether all players have called go
func (game *Game) AllPlayersGone() bool {
	for _, player := range game.Players {
		if !player.Gone {
			return false
		}
	}
	return true
}

// AllPlaysDone returns whether all players have emptied their hands
func (game *Game) AllPlaysDone() bool {
	for _, player := range game.Players {
		if len(player.PlayingHand) > 0 {
			return false
		}
	}
	return true
}

// ResetField resets the playing field
func (game *Game) ResetField() {
	game.Field = Hand{}
	for ii := range game.Players {
		player := &game.Players[ii]
		player.Gone = false
	}
}

// GoScore calculates whether a finished field is a go
func (game *Game) GoScore() (score Score) {
	if game.Field.GetTotal() != 31 {
		score = goScore.AddPairing(game.Field)
		game.Players[game.ActivePlayer].AddScore([]Score{score})
	}
	return score
}

// NextPlayer runs the next player turn
func (game *Game) NextPlayer() (isHuman bool, card Card, scores []Score, err error) {
	game.ActivePlayer++
	if game.ActivePlayer >= len(game.Players) {
		game.ActivePlayer = 0
	}
	player := &game.Players[game.ActivePlayer]
	isHuman = !player.IsComputer
	if isHuman {
		return
	}

	card, cantPlay := player.PlayingHand.GetBestPlay(game.Field)
	if cantPlay {
		player.Gone = true
		return
	}
	field, scores, err := game.Field.Play(card)
	if err != nil {
		return
	}
	game.Field = field
	player.AddScore(scores)
	player.PlayingHand = player.PlayingHand.RemoveCard(card)

	return
}
