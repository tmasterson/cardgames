// This package is a set of types and functions useful for all solitaire games.
// A possible exception is something like spider.
package solitaire

import (
	"sync"

	"github.com/tmasterson/cardgames/generic"
)

// Adder is an iterface primarily designed to allow for isolation while testing
type Adder interface {
	Add(c []generic.Card)
}

// Reducer is an interface desgnied primarl=ily for isolation during testing
type Reducer interface {
	Reduce(n int)
}

// SetFirstFaceUper  Interface to allow for testing.  See method for full docs.
type SetFirstFaceUper interface {
	ChangeFirstFaceUp()
}

// CheckMover  An interfae for testing moves.
type CheckMover interface {
	CheckMove(to *Pile, index int) bool
}

// Pile is the base for all card stacks.
// Could possibly be made into a generic hand and moved to the generic package
type Pile struct {
	Cards       []generic.Card
	Firstfaceup int
	Ptype       rune
	mux         sync.Mutex // set up locking for safety
}

// Add adds a card or cards to a pile
func (p *Pile) Add(c []generic.Card) {
	p.Cards = append(p.Cards, c...)
}

// Reduce Removes one or more cards from a pile.
func (p *Pile) Reduce(n int) {
	p.Cards = p.Cards[:n]
}

// ChangeFirstFaceUp  Sets the first faceup card in a Pile.
func (p *Pile) ChangeFirstFaceUp() {
	if len(p.Cards) > 0 {
		if p.Firstfaceup > 0 {
			p.Firstfaceup--
		}
		p.Cards[p.Firstfaceup].Turn()
	} else {
		p.Firstfaceup = 0
	}
}

// DoMove  Does the actual moving of cards from one Pile to another.
func (p *Pile) DoMove(to *Pile, index int) {
	p.mux.Lock() // lock the piles for safety
	to.mux.Lock()
	to.Add(p.Cards[index:])
	p.Reduce(index)
    if index <= p.Firstfaceup {
        p.ChangeFirstFaceUp()
    }
	to.mux.Unlock() // unlock the piles
	p.mux.Unlock()
}

// CheckMove  Checks to make sure that a move is valid
// to is the Pile you are moving cards onto.
// index is the position of the first card in the stack to be moved.
// Returns true if the move is valid otherwise false.
func (p *Pile) CheckMove(to *Pile, index int) bool {
	if len(p.Cards) == 0 { // Can not move empty pile
		return false
	}
	if p.Ptype == 'A' { // can't move form aces
		return false
	}
	card1 := p.Cards[index]
	if index == 0 && card1.Rank == "K" && card1.Faceup && to.Ptype != 'A' { // don't move a king if is the bottom card on a pile unless moving it to an ace pile
		return false
	}
	// if to is empty generate a dummy card for checking legal moves
	var card2 generic.Card
	if len(to.Cards) == 0 {
		card2 = generic.NewCard("", "", "", 0, 0, false)
	} else {
		card2 = to.Cards[len(to.Cards)-1]
	}
	switch to.Ptype {
	case 'T':
		if card1.Rvalue == card2.Rvalue-1 && card1.Color != card2.Color && card1.Faceup && card2.Faceup {
			return true
		}
		if card1.Rvalue == 13 && card2.Rvalue == 0 && card1.Faceup {
			return true
		}
	case 'A':
		switch {
		case card1.Rvalue == 1 && card2.Rvalue == 0:
			return true
		case card1.Rvalue-1 == card2.Rvalue && card1.Suit == card2.Suit:
			return true
		}
	}
	return false
}
