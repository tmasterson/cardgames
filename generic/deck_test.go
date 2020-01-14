package generic

import "testing"

func TestNewDeck(t *testing.T) {
    deck := NewDeck()
    if len(deck.Cards) != 52 {
        t.Errorf("expected length of 52 got %d.", len(deck.Cards))
    }
    if deck.LastDealt != 0 {
        t.Errorf("LastDealt should be 0 was %d.", deck.LastDealt)
    }
}

func TestShuffle(t *testing.T) {
    deck := NewDeck()
    var cards1 [52]string
    for i := 0; i < len(deck.Cards); i++ {
        cards1[i] = deck.Cards[i].Rank+deck.Cards[i].Suit
    }
    deck.Shuffle()
    matches := 0
    for i := 0; i < len(deck.Cards); i++ {
        if cards1[i] == deck.Cards[i].Rank+deck.Cards[i].Suit {
            matches++
        }
    }
    if matches == len(deck.Cards) {
        t.Errorf("Cards were not shuffled.")
    }
}

func TestDeal(t *testing.T) {
    deck := NewDeck()
    deck.Shuffle()
    hand := deck.Deal(5, 2)
    if len(hand) < 5 {
        t.Errorf("expected a hand of 5 but got %d.", len(hand))
    }
    if !hand[3].Faceup && !hand[4].Faceup {
        t.Errorf("should have 2 faceup cards but did not.")
    }
    var dup bool
    dup = false
    for i := 0; i < len(hand); i++ {
        for j := 0; j < len(hand); j++ {
            if hand[i].Rank == hand[j].Rank && hand[j].Suit == hand[i].Suit && i != j {
                dup = true
            }
        }
    }
    if dup {
        t.Errorf("There should be no duplicates in a hand.")
    }
    deck.LastDealt = 49
    hand2 := deck.Deal(5, 1)
    if len(hand2) == 5 {
        t.Errorf("There should be less than 5 cards in the hand.")
    }
    if !hand2[len(hand2)-1].Faceup {
        t.Errorf("last card in hand should be face up but was not.")
    }
    deck.LastDealt = 5
    hand3 := deck.Deal(5, 6)
    cnt := 0
    for i := 0; i < len(hand3); i++ {
        if hand3[i].Faceup {
            cnt++
        }
    }
    if cnt < 5 {
        t.Errorf("Expected 5 but was %d.", cnt)
    }
}

func TestTurn(t *testing.T) {
    c := Card{"A", "S", 14, 16, "black", true}
    c.Turn()
    if c.Faceup {
        t.Error("expected false but was true")
    }
    c.Turn()
    if !c.Faceup {
        t.Error("expected true but was false")
    }
}
