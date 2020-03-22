package poner

import (
	"fmt"
)

// CanBePlayed returns whether a card can be played
func (card Card) CanBePlayed(field Hand) bool {
	total := field.GetTotal()
	return card.Value+total <= 31
}

// TotalWouldBe returns what the total would be if the card was played
func (card Card) TotalWouldBe(field Hand) int {
	return field.GetTotal() + card.Value
}

// WouldScore returns the scores that would occur if the card was played
func (card Card) WouldScore(field Hand) (scores []Score) {
	return append(field, card).FieldScore()
}

// CanPlay returns whether a hand has a playable card
func (hand Hand) CanPlay(field Hand) bool {
	for _, card := range hand {
		if card.CanBePlayed(field) {
			return true
		}
	}
	return false
}

// Play puts a card into the playfield
func (hand Hand) Play(card Card) (field Hand, scores []Score, err error) {
	field = hand
	if !card.CanBePlayed(hand) {
		err = fmt.Errorf("Play:: %v can't be played. Total would be %v", card, card.TotalWouldBe(hand))
	}

	field = append(field, card)
	scores = field.FieldScore()
	return
}

// BuildFieldPairings returns all the pairings from the playfield
func (hand Hand) BuildFieldPairings() (pairings Pairings) {
	pairings = Pairings{}
	if len(hand) < 2 {
		return
	}
	for ii := len(hand) - 2; ii >= 0; ii-- {
		pairing := Hand{}
		for jj := len(hand) - 1; jj >= ii; jj-- {
			pairing = append(pairing, hand[jj])
		}
		pairings = append(pairings, pairing)
	}
	return
}

// FieldScore returns the scores in the playfield
func (hand Hand) FieldScore() (scores []Score) {
	scores = []Score{}
	pairings := hand.BuildFieldPairings()
	scores = append(scores, pairings.OfAKindScores()...)
	scores = append(scores, pairings.RunScores()...)

	total := hand.GetTotal()
	if total == 15 {
		scores = append(scores, fifteen.AddPairing(hand))
	} else if total == 31 {
		scores = append(scores, thirtyOne.AddPairing(hand))
	}
	return
}

// GetBestPlay gets the best card to play
func (hand Hand) GetBestPlay(field Hand) (bestCard Card, cantPlay bool) {
	canPlay := Hand{}
	for _, card := range hand {
		if card.CanBePlayed(field) {
			canPlay = append(canPlay, card)
		}
	}
	if len(canPlay) == 0 {
		cantPlay = true
		return
	}

	nonOptimal := Hand{}
	optimal := Hand{}
	for _, card := range canPlay {
		total := card.TotalWouldBe(field)
		if total == 10 || total == 21 {
			nonOptimal = append(nonOptimal, card)
			continue
		}
		optimal = append(optimal, card)
	}
	if len(optimal) == 0 {
		optimal = nonOptimal
	}

	highestTotal := 0
	for _, card := range optimal {
		scores := card.WouldScore(field)
		total := 0
		for _, score := range scores {
			total += score.Value
		}
		if total >= highestTotal {
			bestCard = card
		}
	}

	return
}
