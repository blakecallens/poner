package poner

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

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
	SkillLevel  int
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
		discards := hand.GetDiscards(deck, isDealer)
		skillAdjust := player.GetSkillAdjust(len(discards))
		player.SetDiscard(discards[skillAdjust])
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

// GetSkillAdjust gets a random skill ajustment for player skill
func (player *Player) GetSkillAdjust(maxAdjust int) int {
	rand.Seed(time.Now().UnixNano())
	maxSkilllevel := math.Min(4, float64(player.SkillLevel))
	largestOffset := math.Min(5-maxSkilllevel, float64(maxAdjust))
	largestOffset = math.Max(largestOffset, 0)
	return rand.Intn(int(largestOffset))
}
