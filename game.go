package poner

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Game represents
type Game struct {
	Players       []Player
	Round         int
	ToWin         int
	Dealer        int
	ActivePlayer  int
	ComputerDelay time.Duration
	Deck          Deck
	Starter       Card
	Field         Hand
	Crib          Hand
	Winner        *Player
}

// New creates a new game
func (game *Game) New(players []Player) {
	rand.Seed(time.Now().UnixNano())
	game.Players = players
	game.Round = 0
	game.Dealer = rand.Intn(len(game.Players))
	game.Winner = nil
	if game.ToWin == 0 {
		game.ToWin = 121
	}
}

// NextRound starts a new game round
func (game *Game) NextRound() (score Score, err error) {
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

	err = game.BuildCrib()
	if err != nil {
		return
	}
	game.Starter, err = game.Deck.PullFromTop()
	if err != nil {
		return
	}

	score = game.Starter.HisHeelsScore()
	if score.Value > 0 {
		player := &game.Players[game.Dealer]
		player.AddScore([]Score{score})
		game.CheckForWinner(player)
	}
	return
}

// BuildCrib creates the crib from the discards
func (game *Game) BuildCrib() (err error) {
	crib := Hand{}
	discards := []Discard{}
	for _, player := range game.Players {
		discards = append(discards, player.Discard)
	}
	if len(discards) > 4 {
		err = errors.New("BuildCrib:: no more than 4 player discards allowed")
		return
	}
	for _, discard := range discards {
		crib = append(crib, discard.Discarded...)
	}
	for len(crib) < 4 {
		var card Card
		card, err = game.Deck.PullFromTop()
		if err != nil {
			return
		}
		crib = append(crib, card)
	}
	game.Crib = crib
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
		player := &game.Players[game.ActivePlayer]
		player.AddScore([]Score{score})
		game.CheckForWinner(player)
	}
	return score
}

func (game *Game) delayComputer(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(game.ComputerDelay)
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

	var wg sync.WaitGroup
	wg.Add(1)
	go game.delayComputer(&wg)
	wg.Wait()

	nextPlayer := game.ActivePlayer + 1
	if nextPlayer >= len(game.Players) {
		nextPlayer = 0
	}

	plays, cantPlay := player.PlayingHand.GetPlays(game.Field, game.Players[nextPlayer])
	if cantPlay {
		player.Gone = true
		return
	}
	skillAdjust := player.GetSkillAdjust(len(plays))
	scores, err = game.PutCardIntoField(plays[skillAdjust].Card, player)
	return
}

// HumanPlayCard acts upon a human selected card for the playfield
func (game *Game) HumanPlayCard(card Card) (scores []Score, err error) {
	player := &game.Players[game.ActivePlayer]
	if player.IsComputer {
		err = errors.New("HumanPlayCard:: the active player is not human")
		return
	}
	if !card.CanBePlayed(game.Field) {
		err = fmt.Errorf("HumanPlayCard:: %v cannot be played", card)
		return
	}

	scores, err = game.PutCardIntoField(card, player)
	return
}

// HumanPlayGone acts upon a human saying go
func (game *Game) HumanPlayGone() (scores []Score, err error) {
	player := &game.Players[game.ActivePlayer]
	if player.IsComputer {
		err = errors.New("HumanPlayGone:: the active player is not human")
		return
	}
	if !player.PlayingHand.CanPlay(game.Field) {
		player.Gone = true
		return
	}

	err = errors.New("HumanPlayGone:: invalid go attempt. Card(s) can be played")
	return
}

// PutCardIntoField puts a card into the playfield for a player
func (game *Game) PutCardIntoField(card Card, player *Player) (scores []Score, err error) {
	field, scores, err := game.Field.Play(card)
	if err != nil {
		return
	}
	game.Field = field

	player.AddScore(scores)
	game.CheckForWinner(player)
	player.PlayingHand = player.PlayingHand.RemoveCard(card)
	player.Discard.Played = append(player.Discard.Played, card)

	return
}

// ScoreHand scores a player's hand or crib
func (game *Game) ScoreHand(player *Player, isCrib bool) (scores []Score, total int) {
	if !isCrib {
		scores, total = player.Discard.Held.Score(game.Starter, isCrib)
	} else {
		scores, total = game.Crib.Score(game.Starter, isCrib)
	}
	player.AddScore(scores)
	game.CheckForWinner(player)
	return
}

// CheckForWinner returns if the supplied player has won the game
func (game *Game) CheckForWinner(player *Player) bool {
	if player.Score >= game.ToWin {
		game.Winner = player
		return true
	}
	return false
}
