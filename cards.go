package poner

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var names = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var values = [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}
var suits = [4]string{"♠", "♣", "♥", "♦"}
var suitAlts = [4]string{"S", "C", "H", "D"}

// Card represents the data of a single playing card
type Card struct {
	Name  string
	Value int
	Order int
	Suit  string
}

func (card Card) String() string {
	return fmt.Sprintf(card.Name + card.Suit)
}

// Hand represents a collection of cards held by a player
type Hand []Card

func (hand Hand) Len() int             { return len(hand) }
func (hand Hand) Swap(ii, jj int)      { hand[ii], hand[jj] = hand[jj], hand[ii] }
func (hand Hand) Less(ii, jj int) bool { return hand[ii].Order < hand[jj].Order }

// RemoveCard removes a card from a hand
func (hand Hand) RemoveCard(card Card) (newHand Hand) {
	for index, handCard := range hand {
		if card == handCard {
			hand = append(hand[:index], hand[index+1:]...)
			break
		}
	}

	return hand
}

// Frequency represents the number of cards of a type left in the deck
type Frequency struct {
	Name  string
	Value int
	Order int
	Count int
}

// Deck is 52 individual cards
type Deck struct {
	Cards       []Card
	Frequencies []Frequency
}

// New returns a brand new deck of 52 cards
func (deck Deck) New() (newDeck Deck) {
	newDeck = Deck{}

	for _, suit := range suits {
		for ii, name := range names {
			newDeck.Cards = append(newDeck.Cards, Card{
				Name:  name,
				Value: values[ii],
				Order: ii,
				Suit:  suit,
			})
		}
	}

	newDeck.GetFrequencies()
	return
}

// Shuffle shuffles the cards in a deck
func (deck *Deck) Shuffle() {
	shuffledCards := []Card{}
	for len(deck.Cards) > 0 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(deck.Cards))
		shuffledCards = append(shuffledCards, deck.Cards[index])
		deck.Cards = append(deck.Cards[:index], deck.Cards[index+1:]...)
	}
	deck.Cards = shuffledCards
}

// Cut cuts the cards in a deck
func (deck *Deck) Cut() {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(deck.Cards))
	deck.Cards = append(deck.Cards[index:], deck.Cards[:index]...)
}

// Deal deals x number of cards to y players
func (deck *Deck) Deal(cards int, players int) (hands []Hand, err error) {
	if len(deck.Cards) < cards*players {
		err = errors.New("Deal:: not enough cards to deal the requested amount")
		return
	}
	hands = []Hand{}
	for len(hands) < players {
		hands = append(hands, Hand{})
	}
	for cc := 0; cc < cards; cc++ {
		for pp := 0; pp < players; pp++ {
			hands[pp] = append(hands[pp], deck.Cards[0])
			deck.Cards = deck.Cards[1:]
		}
	}
	deck.GetFrequencies()
	return
}

// DealCribbage auto deals the correct amount of cribbage cards for players
func (deck *Deck) DealCribbage(players int) (hands []Hand, err error) {
	cards := 6
	if players > 2 {
		cards = 5
	}
	return deck.Deal(cards, players)
}

// PullFromTop deals one card from the deck
func (deck *Deck) PullFromTop() (card Card, err error) {
	if len(deck.Cards) == 0 {
		err = errors.New("PullFromTop:: no cards left in the deck")
		return
	}
	card = deck.Cards[0]
	deck.Cards = deck.Cards[1:]
	deck.GetFrequencies()
	return
}

// PullCard finds a card in the deck and pulls it
func (deck *Deck) PullCard(name string, suit string) (pulledCard Card, err error) {
	for index, alt := range suitAlts {
		if alt == strings.ToUpper(suit) {
			suit = suits[index]
			break
		}
	}
	for index, card := range deck.Cards {
		if card.Name == strings.ToUpper(name) && card.Suit == suit {
			pulledCard = card
			deck.Cards = append(deck.Cards[:index], deck.Cards[index+1:]...)
			deck.GetFrequencies()
			return
		}
	}

	err = fmt.Errorf("PullCard:: no %v%v card in the deck", name, suit)
	return
}

// PullCards finds multiple cards in the deck and pulls them
func (deck *Deck) PullCards(cardString string) (pulledCards Hand, err error) {
	pulledCards = Hand{}
	cardStrings := strings.Split(cardString, " ")
	for _, cardDesc := range cardStrings {
		descSize := len(cardDesc)
		if descSize != 2 && descSize != 3 {
			err = fmt.Errorf("PullCards:: invalid card %v", cardDesc)
			return
		}
		var pulledCard Card
		pulledCard, err = deck.PullCard(string(cardDesc[:descSize-1]), string(cardDesc[descSize-1]))
		if err != nil {
			break
		}
		pulledCards = append(pulledCards, pulledCard)
	}

	return
}

// GetFrequencies builds the frequencies of remaining cards in the deck
func (deck *Deck) GetFrequencies() (frequencies []Frequency) {
	frequencies = []Frequency{}
	for ii, name := range names {
		frequency := Frequency{
			Name:  name,
			Value: values[ii],
			Order: ii,
		}
		for _, card := range deck.Cards {
			if card.Name == frequency.Name {
				frequency.Count++
			}
		}
		frequencies = append(frequencies, frequency)
	}
	deck.Frequencies = frequencies
	return
}
