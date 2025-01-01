// This version allows a user to play the game klondike
// It uses the ncurses tcell created by Garrett D'Amore which can be gotten by
// go get -u github.com/gdamore/tcell
//
// This is the 3 pass version of klondike in which you can go through the deck 3 times.
package main

import (
	"errors"
	"fmt"
	//"log"
	"os"
	"strings"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/tmasterson/cardgames/generic"
	"github.com/tmasterson/cardgames/solitaire"
)

// Box defines the coordinates of a box and the area where cards will
// be printed on the single card groups i.e waste and aces.
type box struct {
	title                     string
	leftX, rightX, topY, botY int
	cardArea                  int
	style                     tcell.Style
}

// Move is a structure to track moves.
// from:  Stack to mover cards from
// to:  Stack to move cards to
// pass:  Number of passes through the waste stack
// howmany:  Number of cards to move
type move struct {
    from int
    to int
    pass int
    howmany int
}

// This variable sets up logging to the file game.out which will show each move and is automatically truncated for each run.
//var (
//	f, err = os.OpenFile("game.out", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
//	logger = log.New(f, "", log.LstdFlags)
//)

// Define the boxes for easier manipulation
var wasteArea, ace1, ace2, ace3, ace4, playArea box

// Print a string in a given area of the screen in a given style
// s: The screen variable
// x, y: Cooridinates of the left end of the string.
// style: The string style.
// str: The string to be printed.
func putString(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		w := runewidth.RuneWidth(c)
		s.SetContent(x, y, c, nil, style)
		x += w
	}
}

// Makebox draws a box on the screen.
// s: The screen variable
// title: A string that defines a title for the box.  Can be blank or null
// leftX, topY, rightX, botY: The coordinates for the corners of the box
// style: The style for the box
//
// If this were to be made more generic it should be set so that you can do it without borders.
func makeBox(s tcell.Screen, title string, leftX, topY, rightX, botY int, style tcell.Style) (box, error) {
	center := (rightX - leftX) / 2
	var b box
	if len(title) > rightX-leftX-2 {
		return b, fmt.Errorf("title length %d must not exceed right-left-2 %d", len(title), rightX-leftX-2)
	}
	b.title = title
	b.leftX = leftX
	b.topY = topY
	b.rightX = rightX
	b.botY = botY
	b.cardArea = b.leftX + center - 1
	b.style = style
	// draw the box
	for col := b.leftX; col <= b.rightX; col++ {
		s.SetContent(col, b.topY, tcell.RuneHLine, nil, b.style)
		s.SetContent(col, b.botY, tcell.RuneHLine, nil, b.style)
	}
	for row := b.topY; row <= b.botY; row++ {
		s.SetContent(b.leftX, row, tcell.RuneVLine, nil, b.style)
		s.SetContent(b.rightX, row, tcell.RuneVLine, nil, b.style)
	}
	s.SetContent(b.leftX, b.topY, tcell.RuneULCorner, nil, style)
	s.SetContent(b.leftX, b.botY, tcell.RuneLLCorner, nil, style)
	s.SetContent(b.rightX, b.topY, tcell.RuneURCorner, nil, style)
	s.SetContent(b.rightX, b.botY, tcell.RuneLRCorner, nil, style)
	titlepos := b.leftX + center - len(title)/2
	if titlepos == b.leftX {
		titlepos++
	}
	putString(s, titlepos, b.topY, style, title)
	return b, nil
}

// DrawScreen draws the screen putting all the boxes in place
// s: The screen variable
// style: The style for the screen
//
// Returns: Returns an error if one occurs otherwise nil
//
func drawScreen(s tcell.Screen, style tcell.Style) error {
	w, h := s.Size()
	if w < 80 || h < 25 {
		return errors.New("Screen size must be at least 80 by 25")
	}
	center := w / 2
	title := "Klondike"
	putString(s, center-len(title)/2, 0, style, title)
	var err error
	wasteArea, err = makeBox(s, "Waste", 0, 2, 10, 4, style)
	if err != nil {
		return err
	}
	ace1, err = makeBox(s, "Ace", wasteArea.rightX+2, wasteArea.topY, wasteArea.rightX+8, wasteArea.topY+2, style)
	if err != nil {
		return err
	}
	ace2, err = makeBox(s, "Ace", ace1.rightX+1, wasteArea.topY, ace1.rightX+7, wasteArea.topY+2, style)
	if err != nil {
		return err
	}
	ace3, err = makeBox(s, "Ace", ace2.rightX+1, wasteArea.topY, ace2.rightX+7, wasteArea.topY+2, style)
	if err != nil {
		return err
	}
	ace4, err = makeBox(s, "Ace", ace3.rightX+1, wasteArea.topY, ace3.rightX+7, wasteArea.topY+2, style)
	if err != nil {
		return err
	}
	playArea, err = makeBox(s, "Tableau", 0, ace1.botY+2, 23, ace2.botY+16, style)
	if err != nil {
		return err
	}
	x := playArea.rightX + 5
	y := playArea.topY
	putString(s, (w-x)/2-3, y, style, "Moves:")
	y++
	putString(s, x, y, style, "All moves are a two character instruction.")
	y++
	putString(s, x, y, style, "AG will move from stack 1 to stack 7.")
	y++
	putString(s, x, y, style, "wA will move from waste to stack 1.")
	y++
	putString(s, x, y, style, "w<Enter will move from waste to an ace stack.")
	y++
	putString(s, x, y, style, "A<enter will move from stack 1 to an ace stack.")
	s.Show()
	return nil
}

// ShowStack prints the face up cards in each stack
// s: Screen variable
// stacks: A slice containing all the stacks or Piles of cards
// style: The style for the cards
func showStacks(s tcell.Screen, stacks []solitaire.Pile, style tcell.Style) {
	for i, pile := range stacks {
		switch pile.Ptype {
		case 'T':
			k := 1
			for _, card := range pile.Cards {
				if card.Faceup {
					putString(s, playArea.leftX+i*3+2, playArea.topY+k, style, card.Rank+card.Suit)
					k++
				}
			}
			for playArea.topY+k < playArea.botY-1 {
				putString(s, playArea.leftX+i*3+2, playArea.topY+k, style, "  ")
				k++
			}
		case 'W':
			if len(pile.Cards) > 0 {
				card := pile.Cards[pile.Firstfaceup]
				putString(s, wasteArea.cardArea, wasteArea.topY+1, style, card.Rank+card.Suit)
			} else {
				putString(s, wasteArea.cardArea, wasteArea.topY+1, style, "  ")
			}
		case 'A':
			if len(pile.Cards) > 0 {
				card := pile.Cards[len(pile.Cards)-1]
				switch card.Suit {
				case "S":
					putString(s, ace1.cardArea, ace1.topY+1, style, card.Rank+card.Suit)
				case "H":
					putString(s, ace2.cardArea, ace2.topY+1, style, card.Rank+card.Suit)
				case "D":
					putString(s, ace3.cardArea, ace3.topY+1, style, card.Rank+card.Suit)
				case "C":
					putString(s, ace4.cardArea, ace4.topY+1, style, card.Rank+card.Suit)
				}
			}
		}
	}
	s.Show()
}

// DealToWaste deals cards from the deck to the waste stack.
// If all cards have been dealt from the deck it turns the waste stack into the deck
// and increments the pass count.
//
// stacks: A slice containing all the stacks (could probably be just the wast stack)
// deck: A pointer to the deck
// pass: The pass count
//
// returns: The pass count
func dealToWaste(stacks []solitaire.Pile, deck *generic.Deck, pass int) int {
	if len(stacks[7].Cards) > 0 {
		stacks[7].Cards[stacks[7].Firstfaceup].Turn()
	}
	if deck.AllDealt {
		deck.Cards = deck.Cards[:0]
		deck.Cards = append(deck.Cards, stacks[7].Cards...)
		deck.LastDealt = 0
		deck.AllDealt = false
		stacks[7].Cards = stacks[7].Cards[:0]
		pass++
	}
	stacks[7].Cards = append(stacks[7].Cards, deck.Deal(3, 1)...)
	stacks[7].Firstfaceup = len(stacks[7].Cards) - 1
	return pass
}

// processKey handles the processing of key strokes
//
// ev:  The event that that contains the key.
// stacks:  The card stacks.  Passed primarily to handle dealing to the waste stack
// deck:  The deck of cards.  Passed to handle dealing to the waste pile.
// pass:  Number of passes so far reset by the deal to waste function.
// movefrom:  The stack to move cards from on initial call it is -1 and is passed again to set the move to stack.
//
// returns: a filled move structure 
//
func processKey(ev *tcell.EventKey, stacks []solitaire.Pile, deck *generic.Deck, cm move) move {
    var ret move
    switch ev.Key() {
    case tcell.KeyEnter:
        if cm.from != -1 {
			if len(stacks[cm.from].Cards) != 0 {
                switch stacks[cm.from].Cards[len(stacks[cm.from].Cards) - 1].Suit {
                case "S":
                    ret.from = cm.from
                    ret.to = 8
                    ret.pass = cm.pass
                    ret.howmany = 1
                case "H":
                    ret.from = cm.from
                    ret.to = 9
                    ret.pass = cm.pass
                    ret.howmany = 1
                case "D":
                    ret.from = cm.from
                    ret.to = 10
                    ret.pass = cm.pass
                    ret.howmany = 1
                case "C":
                    ret.from = cm.from
                    ret.to = 11
                    ret.pass = cm.pass
                    ret.howmany = 1
                }
			} else {
                ret.from = -1
                ret.to = -1
                ret.pass = cm.pass
			}
		}
    case tcell.KeyRune:
        switch unicode.ToUpper(ev.Rune()) {
        case 'Q':
            ret.from = -1
            ret.to = -1
            ret.pass = 3
        case ' ':
            ret.from = -1
            ret.to = -1
            ret.pass = dealToWaste(stacks[:], deck, cm.pass)
        case 'A', 'B', 'C', 'D', 'E', 'F', 'G':
            //		logger.Printf("in tableau")
            if cm.from == -1 {
                ret.from = int(unicode.ToUpper(ev.Rune())-'A')
                ret.to = -1
                ret.pass = cm.pass
                ret.howmany = cm.howmany
            } else {
                ret.from = cm.from
                ret.to = int(unicode.ToUpper(ev.Rune())-'A')
                ret.pass = cm.pass
                ret.howmany = cm.howmany
            }
        case 'W':
            //		logger.Printf("in waste")
            ret.from = 7
            ret.to = -1
            ret.pass = cm.pass
        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
            ret.from = -1
            ret.to = -1
            ret.pass = cm.pass
            ret.howmany = cm.howmany*10+int(ev.Rune()-'0')
            if ret.howmany > 12 {
                ret.howmany = 0
            }
        default:
            ret.from = -1
            ret.to = -1
            ret.pass = cm.pass
        }
    }
	return ret
}

func moveCards(stacks []solitaire.Pile, cm move) move {
    //logger.Printf("cm = %v", cm)
    if cm.from > -1 && cm.to > -1 {
        if cm.to != cm.from {
            if cm.to > 7 && stacks[cm.from].Firstfaceup < len(stacks[cm.from].Cards)-1 {
                cm.howmany = len(stacks[cm.from].Cards) - 1
            } else {
                if cm.howmany < 1 {
                    cm.howmany = stacks[cm.from].Firstfaceup
                } else {
                    cm.howmany = len(stacks[cm.from].Cards) - cm.howmany
                }
            }
            if stacks[cm.from].CheckMove(&stacks[cm.to], cm.howmany) {
                stacks[cm.from].DoMove(&stacks[cm.to], cm.howmany)
            }
        }
        cm.from = -1
        cm.to = -1
        cm.howmany = 0
    }
    return cm
}

// PlayGame is the main function that handles all aspects of the game.
//
// s: Screnn variable.
// style: The style for the screen.
func playGame(s tcell.Screen, style tcell.Style) int {
	stacks := make([]solitaire.Pile, 12)
	deck := generic.NewDeck()
	deck.Shuffle()
	// Make aces low card instead of high card
	for i := range deck.Cards {
		if deck.Cards[i].Rvalue == 14 { // change aces to 1 instead of 14
			deck.Cards[i].Rvalue = 1
		}
	}
	for i := range stacks {
		switch i {
		case 0, 1, 2, 3, 4, 5, 6: // tableau
			stacks[i].Cards = deck.Deal(i+1, 1)
			stacks[i].Firstfaceup = len(stacks[i].Cards) - 1
			stacks[i].Ptype = 'T'
		case 7: // waste
			stacks[i].Cards = deck.Deal(3, 1)
			stacks[i].Ptype = 'W'
			stacks[i].Firstfaceup = len(stacks[i].Cards) - 1
		case 8, 9, 10, 11: // aces
			stacks[i].Ptype = 'A'
		}
	}
	w, h := s.Size()
	putString(s, 0, h-1, style, strings.Repeat(" ", w-1))
	s.Show()
    cardmove := move{from: -1, to: -1, pass: 0, howmany: 0}
	for cardmove.pass < 3 {
		showStacks(s, stacks, style)
		putString(s, 0, h-1, style, fmt.Sprintf("Pass# %02d, Waste# %02d, Deck# %02d", cardmove.pass, len(stacks[7].Cards), len(deck.Cards)-deck.LastDealt))
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlL {
                s.Sync()
            } else {
				cardmove = moveCards(stacks[:], processKey(ev, stacks[:], &deck, cardmove))
                //logger.Printf("cardmove = %v", cardmove)
			}
		}
		total := 0
		for i := 8; i < 12; i++ { // total the number of cards in aces
			total += len(stacks[i].Cards)
		}
		if total == 52 { // if all cards are in aces stacks we are done
			return -1
		}
	}
	return cardmove.pass
}

func main() {
	//if err != nil {
	//		log.Fatal("Error opening error log.\n")
	//}
	//defer f.Close()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	s.HideCursor()
	if err := drawScreen(s, tcell.StyleDefault); err != nil {
		s.Fini()
		fmt.Println(err)
		os.Exit(1)
	}
	st := playGame(s, tcell.StyleDefault)
	s.Fini()
	if st == -1 {
		fmt.Println("Congratulations you won!")
	} else {
		fmt.Println("You either quit or lost. Better luck next time.")
	}
}

