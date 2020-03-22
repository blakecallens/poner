package poner

import "sort"

// Player holds the data for a player in the game
type Player struct {
	Name        string
	Score       int
	LastScore   int
	GamesWon    int
	DealtHand   Hand
	PlayingHand Hand
	Discard     Discard
	Gone        bool
	IsComputer  bool
}

// AddScore adds scores to the player's total
func (player *Player) AddScore(scores []Score) (total int) {
	for _, score := range scores {
		total += score.Value
	}
	if total == 0 {
		return
	}

	player.LastScore = player.Score
	player.Score += total
	return
}

// TakeDeal gives dealt cards to a player
func (player *Player) TakeDeal(hand Hand, deck *Deck, isDealer bool) {
	player.DealtHand = hand
	if player.IsComputer {
		player.SetDiscard(hand.GetBestDiscard(deck, isDealer))
	} else {
		player.Discard = Discard{}
		player.PlayingHand = Hand{}
	}
}

// SetDiscard set's the player's discard and playing hand
func (player *Player) SetDiscard(discard Discard) {
	playingHand := Hand{}
	for _, card := range discard.Held {
		playingHand = append(playingHand, card)
	}
	player.Discard = discard
	player.PlayingHand = playingHand
	sort.Sort(Hand(player.PlayingHand))
	player.Gone = false
}
