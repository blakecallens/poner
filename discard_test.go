package poner_test

import (
	"strings"
	"testing"

	"github.com/blakecallens/poner"
)

func TestDiscardString(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hands, err := deck.Deal(6, 1)
	if err != nil {
		t.Errorf("Error dealing cards from deck: %v", err)
		return
	}
	discard := hands[0].GetBestDiscard(&deck, true)
	discardString := discard.String()
	if !strings.Contains(discardString, "Discarded") {
		t.Error("Error formatting discard into string")
	}
}

func TestGetAverageScore(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("2c 3c 5c Jc 4d 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	avg := hand[:4].GetAverageScore(&deck)
	if avg != 11.565217 {
		t.Errorf("Error getting average score, got %v, want 11.565217", avg)
	}
}

func TestGetDiscards(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("2c 3c 5c Jc 4d 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	discards := hand.GetDiscards(&deck, true)
	if len(discards) != 15 {
		t.Errorf("Error getting possible discards, got %v discards, want 15", len(discards))
	}
	if discards[0].HeldAverage < discards[1].HeldAverage {
		t.Error("Error sorting discards, greatest score not first")
	}
}

func TestGetBestDiscard(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("2c 3c 5c Jc 4d 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	discard := hand.GetBestDiscard(&deck, true)
	if discard.HeldAverage != 11.086957 {
		t.Errorf("Error getting average score, got %v, want 11.086957", discard.HeldAverage)
	}
	discard = hand.GetBestDiscard(&deck, false)
	if discard.HeldAverage != 12.478261 {
		t.Errorf("Error getting average score, got %v, want 12.478261", discard.HeldAverage)
	}
}
