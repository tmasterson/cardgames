package solitaire

import (
	"testing"

	"github.com/tmasterson/cardgames/generic"
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

func TestChangeFirstFaceUp(t *testing.T) {
	var p1 Pile
	p1.ChangeFirstFaceUp()
	if p1.Firstfaceup != 0 {
		t.Errorf("expected 0 but got %d", p1.Firstfaceup)
	}
	c1 := generic.NewCard("K", "S", "black", 13, 16, false)
	p1.Cards = append(p1.Cards, c1)
	c1 = generic.NewCard("Q", "S", "black", 12, 16, false)
	p1.Cards = append(p1.Cards, c1)
	c1 = generic.NewCard("K", "D", "Red", 13, 8, false)
	p1.Cards = append(p1.Cards, c1)
	c1 = generic.NewCard("Q", "C", "black", 12, 2, false)
	p1.Cards = append(p1.Cards, c1)
	p1.Firstfaceup = 4
	p1.ChangeFirstFaceUp()
	if p1.Firstfaceup != 3 {
		t.Errorf("expected 3 but was %d", p1.Firstfaceup)
	}
	if !p1.Cards[3].Faceup {
		t.Errorf("expected true but was false")
	}
	p1.Firstfaceup = 0
	p1.ChangeFirstFaceUp()
	if !p1.Cards[0].Faceup {
		t.Errorf("expected true but was false")
	}
	if p1.Firstfaceup != 0 {
		t.Errorf("expected 0 but was %d", p1.Firstfaceup)
	}
}

func TestCheckMove(t *testing.T) {
	var p1, p2 Pile
	if p1.CheckMove(&p2, 0) {
		t.Errorf("expected false but was true")
	}
	p1.Cards = append(p1.Cards, generic.NewCard("K", "S", "black", 13, 16, false))
	p1.Cards = append(p1.Cards, generic.NewCard("j", "S", "black", 11, 16, true))
	p2.Cards = append(p2.Cards, generic.NewCard("t", "H", "red", 10, 8, true))
	p1.Ptype = 'T'
	p2.Ptype = 'A'
	if p2.CheckMove(&p1, 0) {
		t.Errorf("expected false but was true")
	}
	p1.Cards[0].Faceup = true
	if p1.CheckMove(&p2, 0) {
		t.Errorf("Can not move a king at the bottom of a pile")
	}
    p1.Cards = p1.Cards[:0]
    p2.Cards = p2.Cards[:0]
	p1.Cards = append(p1.Cards, generic.NewCard("K", "S", "black", 13, 16, true))
	p2.Cards = append(p2.Cards, generic.NewCard("Q", "S", "black", 12, 16, true))
	if !p1.CheckMove(&p2, 0) {
		t.Errorf("Expected true moving %v %c, to %v %c", p1.Cards, p1.Ptype, p2.Cards, p2.Ptype)
	}
    p1.Cards = p1.Cards[:0]
    p2.Cards = p2.Cards[:0]
	p1.Cards = append(p1.Cards, generic.NewCard("K", "S", "black", 13, 16, true))
	p1.Cards = append(p1.Cards, generic.NewCard("j", "S", "black", 11, 16, true))
	p2.Cards = append(p2.Cards, generic.NewCard("t", "H", "red", 10, 8, true))
	p1.Ptype = 'T'
	p2.Ptype = 'T'
	if !p2.CheckMove(&p1, 0) {
		t.Errorf("Should have gotten true for move of %+v to %+v", p2.Cards[0], p1.Cards[len(p1.Cards)-1])
	}
	p1.Cards = p1.Cards[:0]
	p2.Cards = append(p2.Cards, generic.NewCard("K", "C", "black", 13, 2, true))
	if !p2.CheckMove(&p1, 1) {
		t.Errorf("SHould have been able to move %+v to %+v", p2.Cards, p1.Cards)
	}
	p1.Ptype = 'A'
	p2.Cards = append(p2.Cards, generic.NewCard("A", "C", "black", 1, 2, true))
	if !p2.CheckMove(&p1, 2) {
		t.Errorf("SHould have been able to move %+v to %+v", p2.Cards[2], p1.Cards)
	}
	p2.Ptype = 'A'
	p1.Ptype = 'T'
	p1.Cards = append(p1.Cards, generic.NewCard("2", "C", "black", 2, 2, true))
	if !p1.CheckMove(&p2, 0) {
		t.Errorf("SHould have been able to move %+v to %+v", p1.Cards[0], p2.Cards[2])
	}
}

func TestDoMove(t *testing.T) {
	var p1, p2 Pile
	p1.Cards = append(p1.Cards, generic.NewCard("j", "S", "black", 11, 16, true))
	p2.Cards = append(p2.Cards, generic.NewCard("t", "H", "red", 10, 8, true))
	p2.Cards = append(p2.Cards, generic.NewCard("9", "S", "black", 9, 16, true))
	p2.Cards = append(p2.Cards, generic.NewCard("8", "S", "black", 8, 16, true))
	p1.Ptype = 'T'
	p2.Ptype = 'T'
	p1.Firstfaceup = len(p1.Cards) - 1
	p2.Firstfaceup = len(p2.Cards) - 1
	p2.DoMove(&p1, p2.Firstfaceup)
	if len(p1.Cards) != 2 {
		t.Errorf("expected 2 but was %d", len(p1.Cards))
	}
	if len(p2.Cards) != 2 {
		t.Errorf("expected 2 but was %d", len(p2.Cards))
	}
	if p2.Firstfaceup != 1 {
		t.Errorf("First faceup should be 1 but is %d", p2.Firstfaceup)
	}
	p2.Cards = append(p1.Cards, generic.NewCard("8", "S", "black", 8, 16, true))
	p2.DoMove(&p1, 2)
	if p2.Firstfaceup != 1 {
		t.Errorf("First faceup should be 1 but is %d", p2.Firstfaceup)
	}
}
