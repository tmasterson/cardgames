package main

import (
	"github.com/tmasterson/cardgames/generic"
	"testing"
        "encoding/json"
        "os"
)

func TestIsLegalMove(t *testing.T) {
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
	if isLegalMove(card1, card2) {
		t.Errorf("Both cards must be face up %+v %+v.", card1, card2)
	}
	card1.Faceup = true
	if !isLegalMove(card1, card2) {
		t.Errorf("unknown error as cards match %+v, %+v.", card1, card2)
	}
	if isLegalMove(card2, card1) {
		t.Errorf("First card rank value should be less than second %d, %d.", card2.Rvalue, card1.Rvalue)
	}
	if isLegalMove(card1, card3) {
		t.Errorf("Must be different colors %s, %s.", card1.Color, card2.Color)
	}
	if isLegalMove(card1, card4) {
		t.Errorf("Rank values must differ: %d, %d.", card1.Rvalue, card4.Rvalue)
	}
        card1.Rvalue = 13
        card1.Faceup = false
        card2.Rvalue = 0
        if isLegalMove(card1, card2) {
            t.Errorf("expected false as card1 is face down but got true.")
        }
        card1.Faceup = true
        if !isLegalMove(card1, card2) {
            t.Errorf("Expected true but got false.")
        }
}

func TestAdd(t *testing.T) {
    var p1 Pile
    c1 := make([]generic.Card, 3)
    c2 := make([]generic.Card, 2)
    c1[0] = generic.NewCard("K", "S", "black", 13, 16, false)
    c1[1] = generic.NewCard("Q", "S", "black", 12, 16, false)
    c1[2] = generic.NewCard("K", "D", "Red", 13, 8, true)
    p1.Add(c1)
    if len(p1.Cards) != 3 {
        t.Errorf("Expected 3 but was %d.", len(p1.Cards))
    }
    if p1.Cards[0] != c1[0] {
        t.Errorf("Expected %+v, but got %+v.", c1[0], p1.Cards[0])
    }
    if p1.Cards[1] != c1[1] {
        t.Errorf("Expected %+v, but got %+v.", c1[1], p1.Cards[1])
    }
    if p1.Cards[2] != c1[2] {
        t.Errorf("Expected %+v, but got %+v.", c1[2], p1.Cards[2])
    }
    c2[0] = generic.NewCard("Q", "C", "black", 12, 2, true)
    c2[1] = generic.NewCard("J", "H", "red", 11, 8, true)
    p1.Add(c2)
    if len(p1.Cards) != 5 {
        t.Errorf("Expected 5 but was %d.", len(p1.Cards))
    }
    if p1.Cards[3] != c2[0] {
        t.Errorf("Expected %+v, but got %+v.", c2[0], p1.Cards[3])
    }
    if p1.Cards[4] != c2[1] {
        t.Errorf("Expected %+v, but got %+v.", c2[1], p1.Cards[4])
    }
}

func TestReduce(t *testing.T) {
    var p1 Pile
    c1 := make([]generic.Card, 5)
    c1[0] = generic.NewCard("K", "S", "black", 13, 16, false)
    c1[1] = generic.NewCard("Q", "S", "black", 12, 16, false)
    c1[2] = generic.NewCard("K", "D", "Red", 13, 8, true)
    c1[3] = generic.NewCard("Q", "C", "black", 12, 2, true)
    c1[4] = generic.NewCard("J", "H", "red", 11, 8, true)
    p1.Cards = c1
    p1.Reduce(3)
    if len(p1.Cards) != 3 {
        t.Errorf("Expected 3 but was %d.", len(p1.Cards))
    }
    if p1.Cards[0] != c1[0] {
        t.Errorf("Expected %+v, but got %+v.", c1[0], p1.Cards[0])
    }
    if p1.Cards[1] != c1[1] {
        t.Errorf("Expected %+v, but got %+v.", c1[1], p1.Cards[1])
    }
    if p1.Cards[2] != c1[2] {
        t.Errorf("Expected %+v, but got %+v.", c1[2], p1.Cards[2])
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
    if moveAces(&piles[3], aces[:]) {
        t.Error("expected false but was true.")
    }
    piles[3].Cards[3].Rank = "A"
    piles[3].Cards[3].Rvalue = 1
    if !moveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
    if len(piles[3].Cards) != 3 {
        t.Errorf("expected 3 but got %d.", len(piles[3].Cards))
    }
    if len(aces[0].Cards) != 1 {
        t.Errorf("expected 1 but got %d.", len(aces[0].Cards))
    }
    if !moveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
    piles[3].Cards[1].Rank = "2"
    piles[3].Cards[1].Suit = "S"
    piles[3].Cards[1].Rvalue = 2
    piles[3].Cards[1].Svalue = 16
    piles[3].Cards[1].Color = "black"
    if !moveAces(&piles[3], aces[:]) {
        t.Error("expected true but got false.")
    }
}

func TestTableauMoves(t *testing.T) {
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
    tableauMoves(piles[:], aces[:])
    if len(piles[0].Cards) != 2 {
        t.Errorf("expected 2 but got %d.", len(piles[0].Cards))
    }
    if len(piles[2].Cards) != 2 && piles[2].Firstfaceup != 1 {
        t.Errorf("expected length of 2 but was %d and faceup to be 1 but was %d", len(piles[2].Cards), piles[2].Firstfaceup)
    }
    copy(piles, piles2)
    piles[3].Cards[3].Rank = "7"
    piles[3].Cards[3].Rvalue = 7
    tableauMoves(piles[:], aces[:])
    if len(piles[0].Cards) != 0 {
        t.Errorf("expected 0 but got %d.", len(piles[0].Cards))
    }
    piles[6].Cards[6].Rank = "K"
    piles[6].Cards[6].Rvalue = 13
    tableauMoves(piles[:], aces[:])
    if len(piles[0].Cards) != 1 {
        t.Errorf("expected 1 but got %d.", len(piles[0].Cards))
    }
}
