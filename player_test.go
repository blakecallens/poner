package poner_test

import (
	"testing"

	"github.com/blakecallens/poner"
)

func TestAddScore(t *testing.T) {
	player := poner.Player{Name: "Test"}
	total := player.AddScore([]poner.Score{})
	if total != 0 {
		t.Errorf("Error adding score to player, got %v, want 0", total)
	}
	total = player.AddScore([]poner.Score{poner.Score{Name: "Fifteen", Value: 2}})
	if player.Score != 2 {
		t.Errorf("Error adding score to player, got %v, want 2", player.Score)
	}
	if total != player.Score {
		t.Errorf("Error adding score to player, got %v, want %v", total, player.Score)
	}
}

func TestTakeDeal(t *testing.T) {
	// Computer
	player := poner.Player{Name: "Test", IsComputer: true}
	deck := poner.Deck{}.New()
	deck.Shuffle()
	hands, err := deck.Deal(6, 1)
	if err != nil {
		t.Errorf("Error dealing cards from deck: %v", err)
		return
	}
	player.TakeDeal(hands[0], &deck, true)
	if len(player.Discard.Discarded) != 2 {
		t.Errorf("Error with computer taking deal, got %v discarded, want 2", len(player.Discard.Discarded))
	}
	// Human
	player.IsComputer = false
	player.TakeDeal(hands[0], &deck, false)
	if len(player.Discard.Discarded) != 0 {
		t.Errorf("Error with computer taking deal, got %v discarded, want 0", len(player.Discard.Discarded))
	}
}

func TestSkillAdjust(t *testing.T) {
	player := poner.Player{Name: "Test", IsComputer: true}
	skillAdjust := player.GetSkillAdjust(5)
	if skillAdjust < 0 || skillAdjust > 5 {
		t.Errorf("Error getting skill adjust, got %v, want 0-4", skillAdjust)
	}
	player.SkillLevel = 4
	skillAdjust = player.GetSkillAdjust(5)
	if skillAdjust > 0 {
		t.Errorf("Error getting skill adjust, got %v, want 0", skillAdjust)
	}
}
