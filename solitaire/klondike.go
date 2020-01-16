package main

import (
	"fmt"
	"github.com/tmasterson/cardgames/generic"
)

type Adder interface {
    Add(c []generic.Card)
}

type Reducer interface {
    Reduce(n int)
}

type Pile struct {
    Cards []generic.Card
    Firstfaceup int
}

var tableau = make([]Pile, 7)
var aces = make([]Pile, 4)
var waste Pile
var deck = generic.NewDeck()

// Determine if a move is legal
func isLegalMove(card1, card2 generic.Card) bool {
	if card1.Rvalue == card2.Rvalue-1 && card1.Color != card2.Color && card1.Faceup && card2.Faceup {
		return true
	}
        if (card1.Rvalue == 13 || card1.Rvalue == 1) && card2.Rvalue == 0 && card1.Faceup {
            return true
        }
	return false
}

// Add a card or cards to a pile
func (p *Pile) Add(c []generic.Card) {
    p.Cards = append(p.Cards, c...)
}

// Remove one or more cards from a pile
func (p *Pile) Reduce(n int) {
    for i := range p.Cards[n:] {
        p.Cards[i] = generic.NewCard("", "", "", 0, 0, false)
    }
    p.Cards = p.Cards[:n]
}

func doMove(pile1, pile2 *Pile, index int) {
    pile2.Add(pile1.Cards[index:])
    pile1.Reduce(index)
}

// move 1 or more cards from pile1 to pile2 
// Return true if the move was made otherwise return false
func move(pile1, pile2 *Pile, index int) bool {
    if len(pile1.Cards) == 0 { // Can not move empty pile
        return false
    }
    card1 := pile1.Cards[index]
    // if pile2 is empty generate a dummy card for checking legal moves
    var card2 generic.Card
    if len(pile2.Cards) == 0 {
        card2 = generic.Card{
            Rank: "",
            Suit: "",
            Rvalue: 0,
            Svalue: 0,
            Color: "",
            Faceup: true,
        }
    } else {
            card2 = pile2.Cards[len(pile2.Cards)-1]
        }
	if isLegalMove(card1, card2) {
            doMove(pile1, pile2, index)
            return true
	}
	return false
}

func showpile(pile []generic.Card) {
	for i := len(pile) - 1; i >= 0; i-- {
		if pile[i].Faceup {
			fmt.Printf("%s%s ", pile[i].Rank, pile[i].Suit)
		}
	}
	fmt.Println()
}

func showTableau() {
    for i := 0; i < 7; i++ {
        fmt.Printf("tableau %d: ", i)
        showpile(tableau[i].Cards)
    }
}

func showAces() {
    for i := 0; i < 4; i++ {
        fmt.Printf("Aces %d: ", i)
        showpile(aces[i].Cards)
    }
}

func showWaste() {
    fmt.Print("waste: ")
    showpile(waste.Cards)
}

func moveAces(pile *Pile, aces []Pile) bool {
    index := 0
    lenpile := 0
    if len(pile.Cards) > 0 {
        lenpile = len(pile.Cards)-1
    }
    if pile.Cards[lenpile].Suit == "S" {
        index = 0
    } else {
        if pile.Cards[lenpile].Suit == "H" {
            index = 1
        } else {
            if pile.Cards[lenpile].Suit == "D" {
                index = 2
            } else {
                if pile.Cards[lenpile].Suit == "C" {
                    index = 3
                }
            }
        }
    }
    move := false
    if pile.Cards[lenpile].Rvalue == 1 && len(aces[index].Cards) == 0 {
        move = true
    } else {
        if len(aces[index].Cards) == 0 {
            move = false
        } else {
            if pile.Cards[lenpile].Rvalue-1 == aces[index].Cards[len(aces[index].Cards)-1].Rvalue {
                move = true
            }
        }
    }
    if move {
        doMove(pile, &aces[index], lenpile)
    }
    return move
}

func tableauMoves(piles, aces []Pile) {
    done := false
    for !done {
        cnt := 0
        for i := 0; i < len(piles); i++ {
            for j := 0; j < len(piles); j++ {
                if j == i || len(piles[j].Cards) == 0 {
                    continue
                }
                if moveAces(&piles[j], aces[:]) {
                    cnt++
                    if len(piles[j].Cards) > 0 {
                        piles[j].Firstfaceup--
                        piles[j].Cards[piles[j].Firstfaceup].Turn()
                    }
                }
                if move(&piles[j], &piles[i], piles[j].Firstfaceup) {
                    cnt++
                    if piles[j].Firstfaceup > 0 {
                        piles[j].Firstfaceup--
                    }
                    if len(piles[j].Cards) > 0 {
                        piles[j].Cards[piles[j].Firstfaceup].Turn()
                    }
                }
            }
        }
        fmt.Printf("number of moves was %d\n", cnt)
        showTableau()
        showWaste()
        showAces()
        if cnt == 0 {
            done = true
        }
    }
}

func initTableau() {
    for i := 0; i < len(deck.Cards); i++ {
        if deck.Cards[i].Rvalue == 14 {  // change aces to 1 instead of 14
            deck.Cards[i].Rvalue = 1
        }
    }
    for i := 1; i < 8; i++ {
        tableau[i-1].Cards = deck.Deal(i, 1)
        tableau[i-1].Firstfaceup = len(tableau[i-1].Cards)-1
    }
}

func initWaste() {
    waste.Cards = deck.Deal(3,1)
}

func initgame() {
    deck.Shuffle()
    initTableau()
    initWaste()
}

func moveWaste(tableau []Pile, waste *Pile ) {
    done := false
    for !done {
        for i := 0; i < 7 && len(waste.Cards) > 0; i++ {
            if move(waste, &tableau[i], len(waste.Cards)-1) {
                done = true
                if len(waste.Cards) > 0 {
                    waste.Cards[len(waste.Cards)-1].Turn()
                } else {
                    waste.Cards = deck.Deal(3, 1)
                }
            }
        }
        if !done {
            if deck.AllDealt {
                done = true
            } else {
                waste.Cards = deck.Deal(3,1)
            }
        }
    }
}

func playgame() {
    passes := 0
    cnt := 0
    canplay := true
    for canplay {
        showTableau()
        showWaste()
        showAces()
        tableauMoves(tableau[:], aces[:])
        moveWaste(tableau[:], &waste)
        if deck.AllDealt {
            passes++
            if len(waste.Cards) > 0 {
                waste.Cards[len(waste.Cards)-1].Turn()
            }
            copy(deck.Cards, waste.Cards)
            waste.Cards = waste.Cards[:0]
            waste.Cards = deck.Deal(3, 1)
        }
        cnt = 0
        for i := 0; i < 4; i++ {
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
    initgame()
    playgame()
}
