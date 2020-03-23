package poner_test

import (
	"sort"
	"testing"

	"github.com/blakecallens/poner"
)

func TestCardPlaySort(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("Ac 2c 3s 4c")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	plays := poner.CardPlays{
		poner.CardPlay{Card: hand[0], Value: -1},
		poner.CardPlay{Card: hand[1], Value: 1},
		poner.CardPlay{Card: hand[2], Value: 1},
		poner.CardPlay{Card: hand[3], Value: -2},
	}
	sort.Sort(poner.CardPlays(plays))
	if plays[0].Card.Order != 2 {
		t.Errorf("Error sorting card plays, got %v, want 3♠", plays[0].Card)
	}
}

func TestCanBePlayed(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc Js 5d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	card, err := deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	if !card.CanBePlayed(field) {
		t.Errorf("Error getting playability, got %v, want true", card.CanBePlayed(field))
	}

	card, err = deck.PullCard("10", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	if card.CanBePlayed(field) {
		t.Errorf("Error getting playability, got %v, want false", card.CanBePlayed(field))
	}
}

func TestTotalWouldBe(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc Js 5d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	card, err := deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	if card.TotalWouldBe(field) != 30 {
		t.Errorf("Error getting total, got %v, want 30", card.TotalWouldBe(field))
	}
}

func TestWouldScore(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc Js 5d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	card, err := deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	scores := card.WouldScore(field)
	if len(scores) != 1 || scores[0].Name != "Pair" {
		t.Errorf("Error getting scores, got %v, want Pair", scores[0].Name)
	}
}

func TestCanPlay(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc Js 5d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	hand, err := deck.PullCards("10s 10h 10d 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	if !hand.CanPlay(field) {
		t.Errorf("Error getting playability, got %v, want true", hand.CanPlay(field))
	}

	hand, err = deck.PullCards("Qs Qh")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	if hand.CanPlay(field) {
		t.Errorf("Error getting playability, got %v, want false", hand.CanPlay(field))
	}
}

func TestPlay(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc Js 5d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	card, err := deck.PullCard("10", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	_, _, err = field.Play(card)
	if err == nil {
		t.Error("Error playing invalid card, no error returned")
	}

	card, err = deck.PullCard("5", "c")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	field, scores, err := field.Play(card)
	if err != nil {
		t.Errorf("Error playing card: %v", err)
		return
	}
	if len(field) != 4 {
		t.Errorf("Error playing card, got field of %v, want 4", len(field))
		return
	}
	if field[3].Order != 4 {
		t.Errorf("Error playing card, got %v, want 5♣", field[3])
		return
	}
	if len(scores) != 1 || scores[0].Name != "Pair" {
		t.Errorf("Error playing card, got %v, want Pair for 2 [5♦ 5♣]", scores[0])
	}
}

func TestBuildFieldPairings(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	pairings := field.BuildFieldPairings()
	if len(pairings) > 0 {
		t.Errorf("Error getting pairings, got %v pairings, want 0", len(pairings))
	}

	field, err = deck.PullCards("As 2s 3s 4s")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	pairings = field.BuildFieldPairings()
	if len(pairings) != 3 {
		t.Errorf("Error getting pairings, got %v pairings, want 3", len(pairings))
	}
}

func TestFieldScore(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	field, err := deck.PullCards("Jc")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	scores := field.FieldScore()
	if len(scores) > 0 {
		t.Errorf("Error getting pairings, got %v scores, want 0", len(scores))
	}

	field, err = deck.PullCards("10c 3c 2s")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	scores = field.FieldScore()
	if len(scores) != 1 || scores[0].Name != "Fifteen" {
		t.Errorf("Error getting pairings, got %v, want Fifteen", scores)
	}

	field, err = deck.PullCards("Qs Kd 6d 5h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	scores = field.FieldScore()
	if len(scores) != 1 || scores[0].Name != "Thirty One" {
		t.Errorf("Error getting pairings, got %v, want Thirty One", scores)
	}
}

func TestGetBestPlay(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hand, err := deck.PullCards("6s 7s 8s 9s 10s Js")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	player := poner.Player{Name: "Test"}
	player.Discard = hand.GetBestDiscard(&deck, false)
	field, err := deck.PullCards("Jh Qh Kh")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	_, cantPlay := hand.GetBestPlay(field, player)
	if !cantPlay {
		t.Error("Error getting best play, got canPlay, want cantPlay")
	}

	field, err = deck.PullCards("9d")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	bestCard, cantPlay := hand.GetBestPlay(field, player)
	if cantPlay {
		t.Error("Error getting best play, got cantPlay, want canPlay")
	}
	if bestCard.Order != 8 {
		t.Errorf("Error getting best play, got %v, want 9♠", bestCard)
	}
}
