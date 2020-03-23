package poner_test

import (
	"testing"

	"github.com/blakecallens/poner"
)

func TestScoreString(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("As Ac")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	score := poner.Score{Name: "Pair", Value: 2, Pairing: hand}
	if score.String() != "Pair for 2 [A♠ A♣]" {
		t.Errorf("Error getting score string, got %v, want Pair for 2 [A♠ A♣]", score.String())
	}
}

func TestScore(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("Ac 2c 3c")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	starter, err := deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	_, total := hand.Score(starter, false)
	if total > 0 {
		t.Errorf("Error getting score on bad hand, got  %v, want 0", total)
	}

	deck = poner.Deck{}.New()
	deck.Shuffle()
	hand, err = deck.PullCards("Ac 2c 3c 4c")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	starter, err = deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	score, total := hand.Score(starter, false)
	if len(score) != 3 || total != 12 {
		t.Errorf("Error getting score, got %v scores for %v, want 3 for 12", len(score), total)
	}

	deck = poner.Deck{}.New()
	deck.Shuffle()
	hand, err = deck.PullCards("Jc Js Jh Jd")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	starter, err = deck.PullCard("3", "d")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	score, total = hand.Score(starter, false)
	if len(score) != 2 || total != 13 {
		t.Errorf("Error getting score, got %v scores for %v, want 1 for 12", len(score), total)
	}
}

func TestHisHeels(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	starter, err := deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	score := starter.HisHeelsScore()
	if score.Value != 0 {
		t.Errorf("Error getting score, got %v, want 0", score.Value)
	}

	starter, err = deck.PullCard("J", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	score = starter.HisHeelsScore()
	if score.Value != 2 {
		t.Errorf("Error getting score, got %v, want 2", score.Value)
	}
}

func TestThirtyOne(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("Jc 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	score := hand.ThirtyOneScore()
	if score.Value != 0 {
		t.Errorf("Error getting score, got %v, want 0", score.Value)
	}

	hand, err = deck.PullCards("Js Qh Kd As")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	score = hand.ThirtyOneScore()
	if score.Value != 2 {
		t.Errorf("Error getting score, got %v, want 2", score.Value)
	}
}
