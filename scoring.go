package poner

import (
	"fmt"
	"sort"
)

// Score represents a single cribbage score
type Score struct {
	Name    string
	Value   int
	Pairing Hand
}

func (score Score) String() string {
	return fmt.Sprintf("%v for %v %v", score.Name, score.Value, score.Pairing)
}

// AddPairing adds a pairing to the score
func (score Score) AddPairing(pairing Hand) Score {
	score.Pairing = pairing
	return score
}

// Pairings are all the possible groupings of a hand
type Pairings []Hand

// The different types of scores
var (
	nobs            = Score{Name: "Nobs", Value: 1}
	fifteen         = Score{Name: "Fifteen", Value: 2}
	pair            = Score{Name: "Pair", Value: 2}
	pairRoyal       = Score{Name: "Pair Royal", Value: 6}
	doublePairRoyal = Score{Name: "Double Pair Royal", Value: 12}
	runOfThree      = Score{Name: "Run of Three", Value: 3}
	runOfFour       = Score{Name: "Run of Four", Value: 4}
	runOfFive       = Score{Name: "Run of Five", Value: 5}
	flushOfFour     = Score{Name: "Flush of Four", Value: 4}
	flushOfFive     = Score{Name: "Flush of Five", Value: 5}
	hisHeels        = Score{Name: "His Heels", Value: 2}
	goScore         = Score{Name: "Go", Value: 1}
	thirtyOne       = Score{Name: "Thirty One", Value: 2}
)

// Score scores a cribbage hand/crib
func (hand Hand) Score(starter Card, isCrib bool) (scores []Score, total int) {
	grossScores := []Score{}
	if len(hand) != 4 {
		return
	}
	sizedHand := append(hand[:4], starter)
	pairings := sizedHand.BuildPairings()

	grossScores = append(grossScores, hand[:4].NobsScore(starter))
	grossScores = append(grossScores, pairings.OfAKindScores()...)
	grossScores = append(grossScores, pairings.FifteenScores()...)
	grossScores = append(grossScores, pairings.RunScores()...)
	grossScores = append(grossScores, pairings.FlushScores(isCrib)...)

	scores = []Score{}
	for _, score := range grossScores {
		if score.Value > 0 {
			total += score.Value
			scores = append(scores, score)
		}
	}

	return
}

// BuildPairings builds all the possible card pairings for a hand
func (hand Hand) BuildPairings() (pairings Pairings) {
	pairings = Pairings{hand}
	// Quadruples
	for ii := 0; ii < len(hand)-3; ii++ {
		for jj := ii + 1; jj < len(hand)-2; jj++ {
			for kk := jj + 1; kk < len(hand)-1; kk++ {
				for ll := kk + 1; ll < len(hand); ll++ {
					pairing := Hand{hand[ii], hand[jj], hand[kk], hand[ll]}
					pairings = append(pairings, pairing)
				}
			}
		}
	}
	// Triples
	for ii := 0; ii < len(hand)-2; ii++ {
		for jj := ii + 1; jj < len(hand)-1; jj++ {
			for kk := jj + 1; kk < len(hand); kk++ {
				pairing := Hand{hand[ii], hand[jj], hand[kk]}
				pairings = append(pairings, pairing)
			}
		}
	}
	// Doubles
	for ii := 0; ii < len(hand)-1; ii++ {
		for jj := ii + 1; jj < len(hand); jj++ {
			pairing := Hand{hand[ii], hand[jj]}
			pairings = append(pairings, pairing)
		}
	}
	return
}

// GetTotal returns to total value of a hand
func (hand Hand) GetTotal() (total int) {
	for _, card := range hand {
		total += card.Value
	}
	return
}

// NobsScore finds the nob in a hand
func (hand Hand) NobsScore(starter Card) (score Score) {
	for _, card := range hand {
		if card.Order == 10 && card.Suit == starter.Suit {
			return nobs.AddPairing(Hand{card})
		}
	}
	return
}

// OfAKindScores finds all the pairs in pairings
func (pairings Pairings) OfAKindScores() (scores []Score) {
	scores = []Score{}
	foundPairs := []string{}
	for _, pairing := range pairings {
		name := pairing[0].Name
		// Skip already matched of a kinds
		alreadyFound := false
		for _, found := range foundPairs {
			if name == found {
				alreadyFound = true
				break
			}
		}
		if alreadyFound {
			continue
		}
		// Check if pairing is of a kind
		isOfAKind := true
		for _, card := range pairing {
			if card.Name != name {
				isOfAKind = false
				break
			}
		}
		if !isOfAKind {
			continue
		}
		// Add the of a kind
		foundPairs = append(foundPairs, name)
		switch len(pairing) {
		case 2:
			scores = append(scores, pair.AddPairing(pairing))
			break
		case 3:
			scores = append(scores, pairRoyal.AddPairing(pairing))
			break
		case 4:
			scores = append(scores, doublePairRoyal.AddPairing(pairing))
			break
		default:
			break
		}
	}
	return
}

// FifteenScores find all the 15s in pairings
func (pairings Pairings) FifteenScores() (scores []Score) {
	scores = []Score{}
	for _, pairing := range pairings {
		total := pairing.GetTotal()
		if total == 15 {
			scores = append(scores, fifteen.AddPairing(pairing))
		}
	}
	return
}

// RunScore checks a hand for a run
func (hand Hand) RunScore() (score Score) {
	// Check if hand is a run
	isRun := true
	sort.Sort(Hand(hand))
	for ii := 1; ii < len(hand); ii++ {
		if hand[ii].Order != hand[ii-1].Order+1 {
			isRun = false
			break
		}
	}
	if !isRun {
		return
	}
	// Score the run
	switch len(hand) {
	case 3:
		return runOfThree.AddPairing(hand)
	case 4:
		return runOfFour.AddPairing(hand)
	case 5:
		return runOfFive.AddPairing(hand)
	default:
		return
	}
}

// RunScores finds all the runs in pairings
func (pairings Pairings) RunScores() (scores []Score) {
	largestRun := 0
	scores = []Score{}
	for _, pairing := range pairings {
		if len(pairing) < 3 || len(pairing) < largestRun {
			return
		}
		score := pairing.RunScore()
		if score.Value > 0 {
			largestRun = len(pairing)
			scores = append(scores, score)
		}
	}
	return
}

// FlushScores finds all the flushes in pairings
func (pairings Pairings) FlushScores(isCrib bool) (scores []Score) {
	scores = []Score{}
	for _, pairing := range pairings {
		if len(pairing) < 4 || (isCrib && len(pairing) < 5) {
			return
		}
		// Check if pairing is a flush
		suit := pairing[0].Suit
		isFlush := true
		for _, card := range pairing {
			if card.Suit != suit {
				isFlush = false
				break
			}
		}
		if !isFlush {
			continue
		}
		switch len(pairing) {
		case 4:
			scores = append(scores, flushOfFour.AddPairing(pairing))
			return
		case 5:
			scores = append(scores, flushOfFive.AddPairing(pairing))
			return
		default:
			break
		}
	}
	return
}

// HisHeelsScore checks the starter for his heels (Jack)
func (card Card) HisHeelsScore() (score Score) {
	if card.Order == 10 {
		return hisHeels.AddPairing(Hand{card})
	}
	return
}

// ThirtyOneScore checks a hand for 31
func (hand Hand) ThirtyOneScore() (score Score) {
	total := hand.GetTotal()
	if total == 31 {
		return thirtyOne.AddPairing(hand)
	}
	return
}
