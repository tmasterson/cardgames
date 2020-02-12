//  This has the computer play a game of Klondike solitaire.
// It creates a gile named game.out which all moves are logged to as well
// as printing mthe major moves to the console.  The game.out file is opened in append mode
// so you will need to clean it out on occasion, or more likely just delete it.
//
// Todo:  Make Pile a type with a slice of hands and a type.  Move hands to the generic package along with the show function.

package main

import (
	"fmt"
	"log"
	"os"

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
}

var (
	f, err = os.OpenFile("game.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger = log.New(f, "", log.LstdFlags)
)

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
	to.Add(p.Cards[index:])
	p.Reduce(index)
	p.ChangeFirstFaceUp()
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
	if index == 0 && card1.Rank == "K" && card1.Faceup { // don't move a king if is the bottom card on a pile
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

func (p *Pile) show() {
	str := ""
	for _, card := range p.Cards {
		if card.Faceup {
			fmt.Printf("%s%s ", card.Rank, card.Suit)
			str = str + " " + card.Rank + card.Suit
		}
	}
	fmt.Println()
	logger.Println(str)
}

func showTableau(tableau []Pile) {
	for i, pile := range tableau {
		fmt.Printf("tableau %d: ", i)
		logger.Printf("tableau %d: Firstfaceup: %d", i, pile.Firstfaceup)
		pile.show()
	}
}

func showAces(aces []Pile) {
	for i, pile := range aces {
		fmt.Printf("Aces %d: ", i)
		logger.Printf("Aces %d: ", i)
		pile.show()
	}
}

func showWaste(waste Pile) {
	fmt.Print("waste: ")
	logger.Printf("waste: %d", len(waste.Cards))
	waste.show()
}

func tableauMoves(piles, aces []Pile) {
	logger.Println("Entering TableauMoves")
	done := false
	for !done {
		cnt := 0
		for i := range piles {
			for j := range piles {
				if j == i || len(piles[j].Cards) == 0 {
					continue
				}
				for l := range aces {
					if piles[j].CheckMove(&aces[l], piles[j].Firstfaceup) {
						cnt++
						piles[j].DoMove(&aces[l], piles[j].Firstfaceup)
					}
				}
				if piles[j].CheckMove(&piles[i], piles[j].Firstfaceup) {
					cnt++
					piles[j].DoMove(&piles[i], piles[j].Firstfaceup)
				}
			}
			logger.Printf("i = %d\n", i)
			logger.Printf("number of moves was %d\n", cnt)
			showTableau(piles)
			showAces(aces)
		}
		fmt.Printf("number of moves was %d\n", cnt)
		logger.Printf("number of moves was %d\n", cnt)
		showTableau(piles)
		showAces(aces)
		if cnt == 0 {
			done = true
		}
	}
	logger.Println("Exiting TableauMoves")
}

func moveWaste(tableau, aces []Pile, waste *Pile, deck *generic.Deck) {
	passes := 0
	for {
		cnt := 0
		for l := range aces {
			if waste.CheckMove(&aces[l], waste.Firstfaceup) {
				waste.DoMove(&aces[l], waste.Firstfaceup)
				cnt++
			}
		}
		for i := range tableau {
			if waste.CheckMove(&tableau[i], waste.Firstfaceup) {
				waste.DoMove(&tableau[i], waste.Firstfaceup)
				cnt++
				break
			}
		}
		if passes > 0 && cnt == 0 {
			break
		}
		if cnt == 0 || len(waste.Cards) == 0 {
			if deck.AllDealt {
				break
			}
			waste.Cards[waste.Firstfaceup].Turn()
			waste.Cards = append(waste.Cards, deck.Deal(3, 1)...)
			waste.Firstfaceup = len(waste.Cards) - 1
			passes++
		}
		showWaste(*waste)
	}
	showWaste(*waste)
}

func playgame(tableau, aces []Pile, waste Pile, deck generic.Deck) {
	passes := 0
	cnt := 0
	canplay := true
	for canplay {
		showTableau(tableau)
		showWaste(waste)
		showAces(aces)
		tableauMoves(tableau[:], aces[:])
		moveWaste(tableau[:], aces[:], &waste, &deck)
		logger.Printf("passes: %d\n", passes)
		if deck.AllDealt {
			passes++
			if len(waste.Cards) > 0 {
				waste.Cards[len(waste.Cards)-1].Turn()
			}
			deck.Cards = deck.Cards[:0]
			deck.Cards = append(deck.Cards, waste.Cards...)
			deck.LastDealt = 0
			deck.AllDealt = false
			waste.Reduce(0)
			waste.Cards = deck.Deal(3, 1)
			waste.Firstfaceup = len(waste.Cards) - 1
		}
		cnt = 0
		for i := range aces {
			cnt += len(aces[i].Cards)
		}
		if cnt == 52 {
			fmt.Println("you won.")
			canplay = false
		}
		if passes >= 3 {
			fmt.Printf("you have made %d passes through the deck and lost\n", passes)
			canplay = false
		}
	}
}

func main() {
	if err != nil {
		log.Fatal("Error opening error log.\n")
	}
	defer f.Close()
	tableau := make([]Pile, 7)
	aces := make([]Pile, 4)
	var waste Pile
	deck := generic.NewDeck()
	deck.Shuffle()
	// Make aces low card instead of high card
	for i := range deck.Cards {
		if deck.Cards[i].Rvalue == 14 { // change aces to 1 instead of 14
			deck.Cards[i].Rvalue = 1
		}
	}
	waste.Cards = deck.Deal(3, 1)
	waste.Ptype = 'W'
	waste.Firstfaceup = len(waste.Cards) - 1
	for i := range tableau {
		tableau[i].Cards = deck.Deal(i+1, 1)
		tableau[i].Firstfaceup = len(tableau[i].Cards) - 1
		tableau[i].Ptype = 'T'
	}
	for i := range aces {
		aces[i].Ptype = 'A'
	}
	playgame(tableau, aces, waste, deck)
}
