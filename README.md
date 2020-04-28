# poner [![Build Status](https://travis-ci.com/blakecallens/poner.svg?branch=master)](https://travis-ci.com/blakecallens/poner) [![Go Report Card](https://goreportcard.com/badge/github.com/blakecallens/poner)](https://goreportcard.com/report/github.com/blakecallens/poner) [![GoDoc](https://godoc.org/github.com/blakecallens/poner?status.svg)](https://godoc.org/github.com/blakecallens/poner)

### A Golang cribbage engine with computer play of varying skill level

#### Examples

How about a nice game of cribbage?

```golang
package main

import (
	"github.com/blakecallens/poner"
	log "github.com/sirupsen/logrus"
)

func main() {
	formatter := log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: ".0000",
	}
	log.SetFormatter(&formatter)

	// poner can handle up to four players
	players := []poner.Player{
		{Name: "Bob", IsComputer: true, SkillLevel: 4},
		{Name: "Sue", IsComputer: true, SkillLevel: 4},
		// {Name: "Dan", IsComputer: true, SkillLevel: 3},
		// {Name: "Joe", IsComputer: true, SkillLevel: 2},
	}
	game := poner.Game{}
	game.New(players)

	// Wait for a winner
	for game.Winner == nil {
		playRound(&game)
		if game.Winner != nil {
			break
		}
		scorePlayerHands(&game)
	}
	// Output final scores
	for _, player := range game.Players {
		log.Infof("%v: %v", player.Name, player.Score)
	}
}

func playRound(game *poner.Game) {
	log.Info("Starting a new round")
	// Start a new round an get his heels, if drawn
	score, _ := game.NextRound()
	log.Infof("Starter: %v", game.Starter)
	if score.Value > 0 {
		log.Info(score)
	}
	// Wait for all players to be out of cards
	for !game.AllPlaysDone() {
		// Field is the current playfield that players put their cards into
		game.ResetField()
		for !game.AllPlayersGone() {
			// Go to the next player and get their play
			_, card, scores, err := game.NextPlayer()
			if err != nil {
				log.Fatal(err)
			}

			player := &game.Players[game.ActivePlayer]
			// poner uses Gone instead of go, for obvious reasons
			if player.Gone {
				continue
			}
			log.Infof("%v plays %v", player.Name, card)
			log.Info(game.Field)
			if len(scores) > 0 {
				log.Info(scores)
			}
			// Stop the round if the player has reached 121
			if game.Winner != nil {
				return
			}
		}
		// GoScore will return a score for go, if the total isn't 31
		score := game.GoScore()
		if score.Value > 0 {
			log.Info(score)
		}
	}
}

func scorePlayerHands(game *poner.Game) {
	ii := game.Dealer + 1
	scoresCounted := 0
	var player *poner.Player
	// Start with the player left of the dealer and make your way around
	for scoresCounted < len(game.Players) {
		if ii >= len(game.Players) {
			ii = 0
		}
		player = &game.Players[ii]
		// Returns the individual scores and the total value of the hand
		scores, total := game.ScoreHand(player, false)
		log.Infof("%v's hand %v: %v %v", player.Name, player.Discard.Held, total, scores)
		// Stop the round if the player has reached 121
		if game.Winner != nil {
			return
		}
		scoresCounted++
		ii++
	}
	// Count the crib
	scores, total := game.ScoreHand(player, true)
	log.Infof("%v's crib%v: %v %v", player.Name, game.Crib, total, scores)
}
```

Output:

```
INFO[.5684] Starting a new round                         
INFO[.5721] Starter: 3♥                                  
INFO[.5721] Bob plays 8♣                                 
INFO[.5721] [8♣]                                         
INFO[.5721] Sue plays Q♦                                 
INFO[.5721] [8♣ Q♦]                                      
INFO[.5722] Bob plays 8♠                                 
INFO[.5722] [8♣ Q♦ 8♠]                                   
INFO[.5722] Bob plays 4♦                                 
INFO[.5722] [8♣ Q♦ 8♠ 4♦]                                
INFO[.5722] Go for 1 [8♣ Q♦ 8♠ 4♦]                       
INFO[.5722] Sue plays Q♥                                 
INFO[.5722] [Q♥]                                         
INFO[.5722] Bob plays 7♣                                 
INFO[.5722] [Q♥ 7♣]                                      
INFO[.5722] Sue plays J♣                                 
INFO[.5723] [Q♥ 7♣ J♣]                                   
INFO[.5723] Go for 1 [Q♥ 7♣ J♣]                          
INFO[.5723] Sue plays J♠                                 
INFO[.5723] [J♠]                                         
INFO[.5723] Go for 1 [J♠]                                
INFO[.5723] Bob's hand [8♠ 4♦ 8♣ 7♣]: 10 [Pair for 2 [8♠ 8♣] Fifteen for 2 [3♥ 4♦ 8♠] Fifteen for 2 [3♥ 4♦ 8♣] Fifteen for 2 [8♠ 7♣] Fifteen for 2 [8♣ 7♣]] 
INFO[.5724] Sue's hand [J♠ Q♥ J♣ Q♦]: 4 [Pair for 2 [J♠ J♣] Pair for 2 [Q♥ Q♦]] 
INFO[.5724] Sue's crib[A♦ 10♦ 7♠ 4♥]: 4 [Fifteen for 2 [A♦ 3♥ 4♥ 7♠] Fifteen for 2 [A♦ 4♥ 10♦]] 
INFO[.5724] Starting a new round                         
INFO[.5761] Starter: J♣                                  
INFO[.5761] His Heels for 2 [J♣]                         
INFO[.5761] Sue plays 8♣                                 
INFO[.5762] [8♣]                                         
INFO[.5762] Bob plays K♦                                 
INFO[.5762] [8♣ K♦]                                      
INFO[.5762] Sue plays 7♦                                 
INFO[.5762] [8♣ K♦ 7♦]                                   
INFO[.5762] Sue plays 6♠                                 
INFO[.5762] [8♣ K♦ 7♦ 6♠]                                
INFO[.5762] [Thirty One for 2 [8♣ K♦ 7♦ 6♠]]             
INFO[.5762] Bob plays 9♦                                 
INFO[.5762] [9♦]                                         
INFO[.5762] Sue plays 5♠                                 
INFO[.5763] [9♦ 5♠]                                      
INFO[.5763] Bob plays 10♠                                
INFO[.5763] [9♦ 5♠ 10♠]                                  
INFO[.5763] Go for 1 [9♦ 5♠ 10♠]                         
INFO[.5763] Bob plays 8♥                                 
INFO[.5763] [8♥]                                         
INFO[.5763] Go for 1 [8♥]                                
INFO[.5764] Sue's hand [7♦ 5♠ 8♣ 6♠]: 8 [Fifteen for 2 [7♦ 8♣] Fifteen for 2 [5♠ J♣] Run of Four for 4 [5♠ 6♠ 7♦ 8♣]] 
INFO[.5764] Bob's hand [9♦ K♦ 10♠ 8♥]: 4 [Run of Four for 4 [8♥ 9♦ 10♠ J♣]] 
INFO[.5765] Bob's crib[2♥ 3♦ 3♠ Q♥]: 10 [Pair for 2 [3♦ 3♠] Fifteen for 2 [2♥ 3♦ Q♥] Fifteen for 2 [2♥ 3♦ J♣] Fifteen for 2 [2♥ 3♠ Q♥] Fifteen for 2 [2♥ 3♠ J♣]] 
INFO[.5765] Bob: 130                                     
INFO[.5765] Sue: 103
```

Get the best discards for the hand 2♣ 3♣ 4<span style="color: darkred">♦</span> 5<span style="color: darkred">♥</span> 5♣ J♣:

```golang
package poner

import (
	"github.com/blakecallens/poner"
	log "github.com/sirupsen/logrus"
)

func main() {
	deck := poner.Deck{}.New()
	hand, _ := deck.PullCards("2c 3c 4d 5h 5c Jc")
	log.Info("Best discard for your crib")
	log.Info(hand.GetBestDiscard(&deck, true))
	log.Info("Best discard for opponent's crib")
	log.Info(hand.GetBestDiscard(&deck, false))
}
```

Output:

```
INFO[0000] Best discard for your crib                   
INFO[0000] Held: [2♣ 3♣ 5♣ J♣], Discarded: [4♦ 5♥], HeldAvg: 11.086957, DiscardedAvg: 6.48 
INFO[0000] Best discard for opponent's crib             
INFO[0000] Held: [3♣ 4♦ 5♥ 5♣], Discarded: [2♣ J♣], HeldAvg: 12.478261, DiscardedAvg: 4.8
```