package poner_test

import (
	"sort"
	"testing"

	"github.com/blakecallens/poner"
)

func TestCardString(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	card, err := deck.PullCard("A", "s")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	if card.String() != "A♠" {
		t.Errorf("Error stringing card, got %v, want A♠", card)
	}
}

func TestSortHand(t *testing.T) {
	deck := poner.Deck{}.New()
	hand, err := deck.PullCards("6h 5h 4h 3h 2h As")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	sort.Sort(poner.Hand(hand))
	if hand[0].Name != "A" {
		t.Errorf("Error pulling card from deck, got %v hands, want A♠", hand[0])
	}
}

func TestRemoveCard(t *testing.T) {
	deck := poner.Deck{}.New()
	hand, err := deck.PullCards("Ah 2h 3h 4h 5h 6h")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	hand = hand.RemoveCard(hand[2])
	if hand[2].Value == 3 {
		t.Errorf("Failed to remove card, got %v, want %v", hand[2], hand[3])
	}
}

func TestNewDeck(t *testing.T) {
	deck := poner.Deck{}.New()
	if len(deck.Cards) != 52 {
		t.Errorf("Failed to create a new deck, got %v cards, want 52", len(deck.Cards))
	}
	for _, frequency := range deck.Frequencies {
		if frequency.Count != 4 {
			t.Errorf("Failed to create deck, got %v frequency, want 4", frequency.Count)
			break
		}
	}
}

func TestShuffleDeck(t *testing.T) {
	unshuffled := poner.Deck{}.New()
	shuffled := poner.Deck{}.New()
	shuffled.Shuffle()
	for ii := range unshuffled.Cards {
		if unshuffled.Cards[ii] != shuffled.Cards[ii] {
			return
		}
	}
	t.Error("Failed to shuffle deck, cards are still unshuffled")
}

func TestCutDeck(t *testing.T) {
	deck := poner.Deck{}.New()
	for range deck.Cards {
		card := deck.Cards[0]
		deck.Cut()
		if card != deck.Cards[0] {
			return
		}
	}
	t.Error("Failed to cut deck, cards are always the same")
}

func TestDeal(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	_, err := deck.Deal(10, 10)
	if err == nil {
		t.Error("Error dealing cards from deck, did not get err for too many cards")
	}
	deck = poner.Deck{}.New()
	deck.Shuffle()
	hands, err := deck.Deal(7, 5)
	if err != nil {
		t.Errorf("Error dealing cards from deck: %v", err)
		return
	}
	if len(hands) != 5 {
		t.Errorf("Error dealing cards from deck, got %v hands, want 5", len(hands))
	}
	for _, hand := range hands {
		if len(hand) != 7 {
			t.Errorf("Error dealing cards from deck, got %v cards, want 7", len(hand))
			break
		}
	}
}

func TestDealCribbage(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hands, err := deck.DealCribbage(2)
	if err != nil {
		t.Errorf("Error dealing cards from deck: %v", err)
		return
	}
	if len(hands) != 2 {
		t.Errorf("Error dealing cards from deck, got %v hands, want 2", len(hands))
	}
	for _, hand := range hands {
		if len(hand) != 6 {
			t.Errorf("Error dealing cards from deck, got %v cards, want 6", len(hand))
			break
		}
	}

	deck = poner.Deck{}.New()
	deck.Shuffle()
	hands, err = deck.DealCribbage(3)
	if err != nil {
		t.Errorf("Error dealing cards from deck: %v", err)
		return
	}
	if len(hands) != 3 {
		t.Errorf("Error dealing cards from deck, got %v hands, want 3", len(hands))
	}
	for _, hand := range hands {
		if len(hand) != 5 {
			t.Errorf("Error dealing cards from deck, got %v cards, want 5", len(hand))
			break
		}
	}
}

func TestPullFromTop(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	card1 := deck.Cards[0]
	card2, err := deck.PullFromTop()
	if err != nil {
		t.Errorf("Error puling top card from deck: %v", err)
		return
	}
	if card1 != card2 {
		t.Errorf("Error pulling top card from deck, got %v hands, want %v", card2, card1)
	}
	deck.Cards = []poner.Card{}
	_, err = deck.PullFromTop()
	if err == nil {
		t.Error("Error pulling top card from deck, did not get err for no cards")
	}
}

func TestPullCard(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	_, err := deck.PullCard("F", "r")
	if err == nil {
		t.Error("Error pulling card from deck, did not get err for bad card")
	}
	card, err := deck.PullCard("A", "s")
	if err != nil {
		t.Errorf("Error pulling card from deck: %v", err)
		return
	}
	if card.Name != "A" || card.Suit != "♠" {
		t.Errorf("Error pulling card from deck, got %v hands, want A♠", card)
	}
}

func TestPullCards(t *testing.T) {
	deck := poner.Deck{}.New()
	deck.Shuffle()
	_, err := deck.PullCards("5s 4h 3c 2d 113s")
	if err == nil {
		t.Error("Error pulling cards from deck, did not get err for bad card")
	}
	deck = poner.Deck{}.New()
	deck.Shuffle()
	_, err = deck.PullCards("5s 4h 3c 2d Fr")
	if err == nil {
		t.Error("Error pulling cards from deck, did not get err for bad card")
	}
	deck = poner.Deck{}.New()
	deck.Shuffle()
	cards, err := deck.PullCards("5s 4h 3c 2d As")
	if err != nil {
		t.Errorf("Error pulling cards from deck: %v", err)
		return
	}
	if cards[4].Name != "A" || cards[4].Suit != "♠" {
		t.Errorf("Error pulling cards from deck, got %v hands, want A♠", cards[4])
	}
}
