package poner

import (
	"fmt"
	"sort"
)

// CardPlay holds the ranking of a card play
type CardPlay struct {
	Card  Card
	Value int
}

// CardPlays is a group of CardPlays for sorting
type CardPlays []CardPlay

func (plays CardPlays) Len() int        { return len(plays) }
func (plays CardPlays) Swap(ii, jj int) { plays[ii], plays[jj] = plays[jj], plays[ii] }
func (plays CardPlays) Less(ii, jj int) bool {
	if plays[ii].Value > plays[jj].Value {
		return true
	}
	if plays[ii].Value == plays[jj].Value {
		return plays[ii].Card.Value > plays[jj].Card.Value
	}
	return false
}

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

// GetPlays gets all available plays
func (hand Hand) GetPlays(field Hand, nextPlayer Player) (plays CardPlays, cantPlay bool) {
	plays = CardPlays{}
	for _, card := range hand {
		if card.CanBePlayed(field) {
			plays = append(plays, CardPlay{Card: card})
		}
	}
	if len(plays) == 0 {
		cantPlay = true
	}

	for ii := range plays {
		play := &plays[ii]
		play.CalculateValue(field, nextPlayer)
	}
	sort.Sort(CardPlays(plays))

	return
}

// CalculateValue computes the value of a potential card play
func (play *CardPlay) CalculateValue(field Hand, nextPlayer Player) {
	playValue := 0

	// Bad moves
	newField := append(field, play.Card)
	for _, card := range nextPlayer.Discard.Played {
		scores := card.WouldScore(newField)
		if len(scores) > 0 {
			playValue--
		}
		if card.Value < 5 {
			card = Card{Name: names[4-card.Order], Value: 5 - card.Value, Order: 4 - card.Order}
		} else if card.Value > 5 {
			card = Card{Name: names[14-card.Order], Value: 15 - card.Value, Order: 14 - card.Order}
		}
		scores = card.WouldScore(newField)
		if len(scores) > 0 {
			playValue--
		}
	}
	total := play.Card.TotalWouldBe(field)
	if total == 5 || total == 10 || total == 21 {
		playValue--
	}

	// Good moves
	scores := play.Card.WouldScore(field)
	for _, score := range scores {
		playValue += score.Value
	}

	play.Value = playValue
}

// GetBestPlay gets the best card to play
func (hand Hand) GetBestPlay(field Hand, nextPlayer Player) (bestCard Card, cantPlay bool) {
	plays, cantPlay := hand.GetPlays(field, nextPlayer)
	if cantPlay {
		return
	}

	bestCard = plays[0].Card
	return
}
