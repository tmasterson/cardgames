package main

import (
	"fmt"
	"github.com/tmasterson/cardgames/generic"
        "os"
        "log"
)

type Adder interface {
    Add(c []generic.Card)
}

type Reducer interface {
    Reduce(n int)
}

type SetFirstFaceUper interface {
    ChangeFirstFaceUp()
}

type Pile struct {
    Cards []generic.Card
    Firstfaceup int
    ptype byte
}

var (
    f, err = os.OpenFile("game.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    logger = log.New(f, "", log.LstdFlags)
)

// Determine if a move is legal
func isLegalMove(card1, card2 generic.Card) bool {
	if card1.Rvalue == card2.Rvalue-1 && card1.Color != card2.Color && card1.Faceup && card2.Faceup {
		return true
	}
        if card1.Rvalue == 13 && card2.Rvalue == 0 && card1.Faceup {
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
    p.Cards = p.Cards[:n]
}

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

func doMove(from, to *Pile, index int) {
    to.Add(from.Cards[index:])
    from.Reduce(index)
    from.ChangeFirstFaceUp()
}

// move 1 or more cards from pile1 to pile2 
// Return true if the move was made otherwise return false
func move(from, to *Pile, index int) bool {
    if len(from.Cards) == 0 { // Can not move empty pile
        return false
    }
    card1 := from.Cards[index]
    if index == 0 && card1.Rank == "K" {
        return false
    }
    // if pile2 is empty generate a dummy card for checking legal moves
    var card2 generic.Card
    if len(to.Cards) == 0 {
        card2 = generic.NewCard("", "", "", 0, 0, false)
    } else {
            card2 = to.Cards[len(to.Cards)-1]
        }
	if isLegalMove(card1, card2) {
            doMove(from, to, index)
            return true
	}
	return false
}

func showpile(pile []generic.Card) {
    str := ""
    for i := range pile {
        if pile[i].Faceup {
            fmt.Printf("%s%s ", pile[i].Rank, pile[i].Suit)
            str = str+" "+pile[i].Rank+pile[i].Suit
        }
    }
    fmt.Println()
    logger.Println(str)
}

func showTableau(tableau []Pile) {
    for i := range tableau {
        fmt.Printf("tableau %d: ", i)
        logger.Printf("tableau %d: Firstfaceup: %d", i, tableau[i].Firstfaceup)
        showpile(tableau[i].Cards)
    }
}

func showAces(aces []Pile) {
    for i := range aces {
        fmt.Printf("Aces %d: ", i)
        logger.Printf("Aces %d: ", i)
        showpile(aces[i].Cards)
    }
}

func showWaste(waste Pile) {
    fmt.Print("waste: ")
    logger.Print("waste: ")
    showpile(waste.Cards)
}

func moveAces(pile *Pile, aces []Pile) bool {
    index := 0
    lenpile := 0
    if len(pile.Cards) > 0 {
        lenpile = len(pile.Cards)-1
    }
    switch pile.Cards[lenpile].Suit {
    case "S":
        index = 0
    case "H":
        index = 1
    case "D":
        index = 2
    case "C":
        index = 3
    }
    move := false
    switch {
    case pile.Cards[lenpile].Rvalue == 1 && len(aces[index].Cards) == 0:
        move = true
    case len(aces[index].Cards) == 0:
        move = false
    case pile.Cards[lenpile].Rvalue-1 == aces[index].Cards[len(aces[index].Cards)-1].Rvalue:
        move = true
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
        for i := range piles {
            for j := range piles {
                if j == i || len(piles[j].Cards) == 0 {
                    continue
                }
                if moveAces(&piles[j], aces[:]) {
                    cnt++
                }
                if move(&piles[j], &piles[i], piles[j].Firstfaceup) {
                    cnt++
                }
            }
        }
        fmt.Printf("number of moves was %d\n", cnt)
        logger.Printf("number of moves was %d\n", cnt)
        showTableau(piles)
        showAces(aces)
        if cnt == 0 {
            done = true
        }
    }
}

func moveWaste(tableau, aces []Pile, waste *Pile, deck generic.Deck) {
    done := false
    for !done {
        for i := 0; i < len(tableau) && len(waste.Cards) > 0; i++ {
            if move(waste, &tableau[i], len(waste.Cards)-1) || moveAces(waste, aces[:]) {
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
        logger.Printf("Waste: %s%s", waste.Cards[len(waste.Cards)-1].Rank, waste.Cards[len(waste.Cards)-1].Suit)
    }
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
        moveWaste(tableau[:], aces[:], &waste, deck)
        logger.Printf("passes: %d\n", passes)
        if deck.AllDealt {
            passes++
            if len(waste.Cards) > 0 {
                waste.Cards[len(waste.Cards)-1].Turn()
            }
            copy(deck.Cards, waste.Cards)
            waste.Reduce(0)
            waste.Cards = deck.Deal(3, 1)
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
        if deck.Cards[i].Rvalue == 14 {  // change aces to 1 instead of 14
            deck.Cards[i].Rvalue = 1
        }
    }
    waste.Cards = deck.Deal(3,1)
    waste.ptype = 'W'
    for i := range tableau {
        tableau[i].Cards = deck.Deal(i+1, 1)
        tableau[i].Firstfaceup = len(tableau[i].Cards)-1
        tableau[i].ptype = 'T'
    }
    for i := range aces {
        aces[i].ptype = 'A'
    }
    playgame(tableau, aces, waste, deck)
}
