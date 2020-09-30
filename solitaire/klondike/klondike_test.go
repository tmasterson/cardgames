package main

import (
	"testing"

	"github.com/gdamore/tcell"
	"github.com/tmasterson/cardgames/generic"
	"github.com/tmasterson/cardgames/solitaire"
)

func mkTestScreen(t *testing.T, charset string) tcell.SimulationScreen {
	s := tcell.NewSimulationScreen(charset)
	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

func TestPutString(t *testing.T) {
	s := mkTestScreen(t, "")
	defer s.Fini()
	putString(s, 1, 1, tcell.StyleDefault, "test")
	s.Show()
	txt, x, y := s.GetContents()
	if len(txt) != x*y {
		t.Errorf("Incorrect size of content: should be %d but was %d", x*y, len(txt))
	}
	var txt2 []byte
	for i := 81; i < 85; i++ {
		txt2 = append(txt2, txt[i].Bytes[0])
	}
	if string(txt2) != "test" {
		t.Errorf("Incorrect string should be test but was %s", txt2)
	}
}

func TestMakeBox(t *testing.T) {
	s := mkTestScreen(t, "")
	defer s.Fini()
	box, err := makeBox(s, "test", 1, 1, 6, 4, tcell.StyleDefault)
	if err == nil {
		t.Fatalf("Should get and error here")
	}
	box, err = makeBox(s, "test", 1, 1, 8, 4, tcell.StyleDefault)
	if err != nil {
		t.Fatalf("Should not get an error here %s", err)
	}
	s.Show()
	txt, x, y := s.GetContents()
	if len(txt) != x*y {
		t.Errorf("Incorrect size of content: should be %d but was %d", x*y, len(txt))
	}
	if box.leftX != 1 || box.rightX != 8 || box.topY != 1 || box.botY != 4 || box.style != tcell.StyleDefault || box.title != "test" || box.cardArea != 3 {
		t.Errorf("expected 1, 8, 1, 4, %v, test, 3 but was %d %d %d %d %v %s %d", tcell.StyleDefault, box.leftX, box.rightX, box.topY, box.botY, box.style, box.title, box.cardArea)
	}
	w, _ := s.Size()
	topline := box.topY * w
	centerpos := (box.rightX-box.leftX)/2 + topline + box.leftX
	titlepos := centerpos - len("test")/2
	if txt[titlepos].Runes[0] != 't' {
		t.Errorf("expected t but was %c at %d %d", txt[titlepos].Runes[0], titlepos, centerpos)
	}
	if txt[topline+box.leftX].Runes[0] != tcell.RuneULCorner {
		t.Errorf("Expected %c but got %c", tcell.RuneULCorner, txt[topline+box.leftX].Runes[0])
	}
	if txt[topline+box.leftX+box.rightX-1].Runes[0] != tcell.RuneURCorner {
		t.Errorf("Expected %c but got %c", tcell.RuneURCorner, txt[topline+box.leftX+box.rightX-1].Runes[0])
	}
	botline := box.botY * w
	if txt[botline+box.leftX].Runes[0] != tcell.RuneLLCorner {
		t.Errorf("Expected %c but got %c", tcell.RuneLLCorner, txt[botline+box.leftX].Runes[0])
	}
	if txt[botline+box.leftX+box.rightX-1].Runes[0] != tcell.RuneLRCorner {
		t.Errorf("Expected %c but got %c", tcell.RuneLRCorner, txt[botline+box.leftX+box.rightX-1].Runes[0])
	}
}

func TestDrawScreen(t *testing.T) {
	s := mkTestScreen(t, "")
	defer s.Fini()
	s.SetSize(30, 10)
	err := drawScreen(s, tcell.StyleDefault)
	if err == nil {
		t.Errorf("Expected an error here")
	}
	s.SetSize(80, 25)
	err = drawScreen(s, tcell.StyleDefault)
	if err != nil {
		t.Errorf("There should be no errors")
	}
}

func TestDealToWaste(t *testing.T) {
	deck := generic.NewDeck()
	stacks := make([]solitaire.Pile, 12)
	stacks[7].Cards = deck.Deal(3, 1)
	stacks[7].Firstfaceup = len(stacks[7].Cards) - 1
	pass := dealToWaste(stacks[:], &deck, 0)
	if len(stacks[7].Cards) != 6 {
		t.Errorf("Should have 6 cards on stack but have %d", len(stacks[7].Cards))
	}
	if pass != 0 {
		t.Errorf("Pass should be 0 but was %d", pass)
	}
	deck.AllDealt = true
	pass = dealToWaste(stacks[:], &deck, 0)
	if len(stacks[7].Cards) != 3 {
		t.Errorf("Should have 3 cards on stack but have %d", len(stacks[7].Cards))
	}
	if pass != 1 {
		t.Errorf("Pass should be 1 but was %d", pass)
	}
	if len(deck.Cards) != 6 {
		t.Errorf("deck should have 6 cards left but has %d", len(deck.Cards))
	}
}

func TestProcessKey(t *testing.T) {
	stacks := make([]solitaire.Pile, 12)
	deck := generic.NewDeck()
	stacks[0].Cards = append(stacks[0].Cards, generic.NewCard("Q", "S", "black", 12, 16, true))
	stacks[0].Ptype = 'T'
	stacks[1].Cards = append(stacks[1].Cards, generic.NewCard("4", "S", "black", 4, 16, true))
	stacks[1].Ptype = 'T'
	stacks[2].Cards = append(stacks[2].Cards, generic.NewCard("T", "S", "black", 10, 16, true))
	stacks[2].Ptype = 'T'
	stacks[3].Cards = append(stacks[3].Cards, generic.NewCard("8", "S", "black", 8, 16, true))
	stacks[3].Ptype = 'T'
	stacks[4].Cards = append(stacks[4].Cards, generic.NewCard("9", "S", "black", 9, 16, true))
	stacks[4].Ptype = 'T'
	stacks[5].Cards = append(stacks[5].Cards, generic.NewCard("7", "S", "black", 7, 16, true))
	stacks[5].Ptype = 'T'
	stacks[6].Cards = append(stacks[6].Cards, generic.NewCard("3", "D", "red", 3, 4, true))
	stacks[6].Ptype = 'T'
	mf, mt, pass := processKey('F', stacks[:], &deck, 0, -1)
	if mt != -1 || mf != -1 || pass != 0 {
		t.Errorf("Expected -1, -1, 0 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('Q', stacks[:], &deck, 0, -1)
	if mt != -1 || mf != -1 || pass != 3 {
		t.Errorf("Expected -1, -1, 3 but was %d, %d, %d", mt, mf, pass)
	}
	deck.AllDealt = true
	mf, mt, pass = processKey('D', stacks[:], &deck, 0, -1)
	if mt != -1 || mf != -1 || pass != 1 {
		t.Errorf("Expected -1, -1, 1 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('2', stacks[:], &deck, 0, -1)
	if mt != -1 || mf != 1 || pass != 0 {
		t.Errorf("Expected -1, 1, 0 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('6', stacks[:], &deck, 0, 1)
	if mt != 5 || mf != 1 || pass != 0 {
		t.Errorf("Expected 5, 1, 0 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('W', stacks[:], &deck, 0, -1)
	if mt != -1 || mf != 7 || pass != 0 {
		t.Errorf("Expected -1, 7, 0 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('A', stacks[:], &deck, 0, 1)
	if mt != 8 || mf != 1 || pass != 0 {
		t.Errorf("Expected 8, 1, 0 but was %d, %d, %d", mt, mf, pass)
	}
	mf, mt, pass = processKey('A', stacks[:], &deck, 0, 6)
	if mt != 10 || mf != 6 || pass != 0 {
		t.Errorf("Expected 10, 6, 0 but was %d, %d, %d", mt, mf, pass)
	}
	stacks[6].Cards[0].Suit = "H"
	mf, mt, pass = processKey('A', stacks[:], &deck, 0, 6)
	if mt != 9 || mf != 6 || pass != 0 {
		t.Errorf("Expected 9, 6, 0 but was %d, %d, %d", mt, mf, pass)
	}
	stacks[6].Cards[0].Suit = "C"
	mf, mt, pass = processKey('A', stacks[:], &deck, 0, 6)
	if mt != 11 || mf != 6 || pass != 0 {
		t.Errorf("Expected 11, 6, 0 but was %d, %d, %d", mt, mf, pass)
	}
}
