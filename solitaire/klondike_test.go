package main

import (
	"encoding/json"
	"github.com/tmasterson/cardgames/generic"
	"log"
	"os"
	"testing"
)

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

func setup() (tabs, aces []Pile) {
	tabs = make([]Pile, 7)
	datafile, err := os.Open("cards.json")
	if err != nil {
		log.Fatal(err)
	}
	datadecoder := json.NewDecoder(datafile)
	datadecoder.Decode(&tabs)
	defer datafile.Close()
	aces = make([]Pile, 4)
	for i := range aces {
		aces[i].Ptype = 'A'
	}
	return tabs, aces
}

func TestTableauMoves(t *testing.T) {
	piles, aces := setup()
	tableauMoves(piles[:], aces[:])
	if len(piles[0].Cards) != 2 {
		t.Errorf("expected 2 but got %d.", len(piles[0].Cards))
	}
	if len(piles[2].Cards) != 2 && piles[2].Firstfaceup != 1 {
		t.Errorf("expected length of 2 but was %d and faceup to be 1 but was %d", len(piles[2].Cards), piles[2].Firstfaceup)
	}
	piles, aces = setup()
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
