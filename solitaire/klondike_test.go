package main

import (
	"github.com/tmasterson/cardgames/generic"
	"testing"
        "encoding/json"
        "os"
)

func TestIslegalmove(t *testing.T) {
	card1 := generic.Card{
		Rank:   "5",
		Suit:   "S",
		Rvalue: 5,
		Svalue: 16,
		Color:  "black",
		Faceup: false,
	}
	card2 := generic.Card{
		Rank:   "6",
		Suit:   "D",
		Rvalue: 6,
		Svalue: 4,
		Color:  "red",
		Faceup: true,
	}
	card3 := generic.Card{
		Rank:   "6",
		Suit:   "C",
		Rvalue: 5,
		Svalue: 2,
		Color:  "black",
		Faceup: true,
	}
	card4 := generic.Card{
		Rank:   "5",
		Suit:   "d",
		Rvalue: 5,
		Svalue: 4,
		Color:  "red",
		Faceup: true,
	}
	if Islegalmove(card1, card2) {
		t.Errorf("Both cards must be face up %+v %+v.", card1, card2)
	}
	card1.Faceup = true
	if !Islegalmove(card1, card2) {
		t.Errorf("unknown error as cards match %+v, %+v.", card1, card2)
	}
	if Islegalmove(card2, card1) {
		t.Errorf("First card rank value should be less than second %d, %d.", card2.Rvalue, card1.Rvalue)
	}
	if Islegalmove(card1, card3) {
		t.Errorf("Must be different colors %s, %s.", card1.Color, card2.Color)
	}
	if Islegalmove(card1, card4) {
		t.Errorf("Rank values must differ: %d, %d.", card1.Rvalue, card4.Rvalue)
	}
        card1.Rvalue = 13
        card1.Faceup = false
        card2.Rvalue = 0
        if Islegalmove(card1, card2) {
            t.Errorf("expected false as card1 is face down but got true.")
        }
        card1.Faceup = true
        if !Islegalmove(card1, card2) {
            t.Errorf("Expected true but got false.")
        }
        card1.Rvalue = 1
        if !Islegalmove(card1, card2) {
            t.Errorf("Expected true but got false.")
        }
}

func TestMove(t *testing.T) {
    pile1 := Pile{Cards: []generic.Card{{
        Rank: "K",
        Suit: "s",
        Rvalue: 13,
        Svalue: 16,
        Color: "black",
        Faceup: false,
    },{
        Rank: "Q",
        Suit: "D",
        Rvalue: 12,
        Svalue: 4,
        Color: "red",
        Faceup: true,
    }}}
    pile2 := Pile{Cards: []generic.Card{}}
    pile3 := Pile{Cards: []generic.Card{{
        Rank: "J",
        Suit: "s",
        Rvalue: 11,
        Svalue: 16,
        Color: "black",
        Faceup: true,
    },{
        Rank: "T",
        Suit: "D",
        Rvalue: 10,
        Svalue: 4,
        Color: "red",
        Faceup: true,
    }}}
    pile4 := Pile{Cards: []generic.Card{{
        Rank: "8",
        Suit: "d",
        Rvalue: 8,
        Svalue: 4,
        Color: "red",
        Faceup: true,
    },{
        Rank: "9",
        Suit: "S",
        Rvalue: 9,
        Svalue: 16,
        Color: "black",
        Faceup: true,
    }}}
    if !move(&pile3, &pile1, 0) {
        t.Errorf("expected true but got false, move was %+v, %+v.", pile3, pile1)
    }
    if len(pile1.Cards) < 4 {
        t.Errorf("expected 4 but got %d.", len(pile1.Cards))
    }
    if len(pile3.Cards) > 0 {
        t.Errorf("expected 0 but got %d.", len(pile3.Cards))
    }
    if move(&pile1, &pile2, 0) {
        t.Errorf("expected false as first moved card is facedown.")
    }
    if !move(&pile4, &pile1, 1) {
        t.Errorf("expected true but was false.")
    }
}

func TestMoveAces(t *testing.T) {
    piles := make([]Pile, 7)
    piles2 := make([]Pile, 7)  // keep a copy of the original
    datafile, err := os.Open("cards.json")
    if err != nil {
        t.Error(err)
    }
    datadecoder := json.NewDecoder(datafile)
    datadecoder.Decode(&piles)
    datafile.Close()
    aces := make([]Pile, 4)
    copy(piles2, piles)
    if MoveAces(&piles[3], aces[:]) {
        t.Error("expected false but was true.")
    }
    piles[3].Cards[3].Rank = "A"
    piles[3].Cards[3].Rvalue = 1
    if !MoveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
    if len(piles[3].Cards) != 3 {
        t.Errorf("expected 3 but got %d.", len(piles[3].Cards))
    }
    if len(aces[0].Cards) != 1 {
        t.Errorf("expected 1 but got %d.", len(aces[0].Cards))
    }
    if !MoveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
    piles[3].Cards[1].Rank = "2"
    piles[3].Cards[1].Suit = "S"
    piles[3].Cards[1].Rvalue = 2
    piles[3].Cards[1].Svalue = 16
    piles[3].Cards[1].Color = "black"
    if !MoveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
}

func TestMakeMoves(t *testing.T) {
    piles := make([]Pile, 7)
    piles2 := make([]Pile, 7)  // keep a copy of the original
    datafile, err := os.Open("cards.json")
    if err != nil {
        t.Error(err)
    }
    datadecoder := json.NewDecoder(datafile)
    datadecoder.Decode(&piles)
    datafile.Close()
    aces := make([]Pile, 4)
    copy(piles2, piles)
    cnt := 0
    cnt = MakeMoves(piles[:], aces[:])
    if cnt != 3 {
        t.Errorf("expected 3 but got %d.", cnt)
    }
    if len(piles[0].Cards) != 2 {
        t.Errorf("expected 2 but got %d.", len(piles[0].Cards))
    }
    if len(piles[2].Cards) != 2 && piles[2].Firstfaceup != 1 {
        t.Errorf("expected length of 2 but was %d and faceup to be 1 but was %d", len(piles[2].Cards), piles[2].Firstfaceup)
    }
    copy(piles, piles2)
    piles[3].Cards[3].Rank = "7"
    piles[3].Cards[3].Rvalue = 7
    cnt = MakeMoves(piles[:], aces[:])
    if len(piles[0].Cards) != 0 {
        t.Errorf("expected 0 but got %d.", len(piles[0].Cards))
    }
    piles[6].Cards[6].Rank = "K"
    piles[6].Cards[6].Rvalue = 13
    cnt = MakeMoves(piles[:], aces[:])
    if len(piles[0].Cards) != 1 {
        t.Errorf("expected 1 but got %d.", len(piles[0].Cards))
    }
}
