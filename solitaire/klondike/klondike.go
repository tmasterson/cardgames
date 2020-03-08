// This version allows either a user or the computer to play the game klondike
// It uses the ncurses tcell created by Garrett D'Amore which can be gotten by
// go get -u github.com/gdamore/tcell
//
// This is the 3 pass version of klondike in which you can go through the deck 3 times.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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

// This variable sets up logging to the file game.out which will show each move and is automatically truncated for each run.
var (
	f, err = os.OpenFile("game.out", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	logger = log.New(f, "", log.LstdFlags)
)

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
func makeBox(s tcell.Screen, title string, leftX, topY, rightX, botY int, style tcell.Style) box {
	center := (rightX - leftX) / 2
	var b box
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
	putString(s, b.leftX+center-len(title)/2, b.topY, style, title)
	return b
}

// DrawScreen draws the screen putting all the boxes in place
// s: The screen variable
// style: The style for the screen
func drawScreen(s tcell.Screen, style tcell.Style) {
	w, _ := s.Size()
	center := w / 2
	title := "Klondike"
	putString(s, center-len(title)/2, 0, style, title)
	wasteArea = makeBox(s, "Waste", 0, 2, 10, 4, style)
	ace1 = makeBox(s, "Ace", wasteArea.rightX+2, wasteArea.topY, wasteArea.rightX+8, wasteArea.topY+2, style)
	ace2 = makeBox(s, "Ace", ace1.rightX+1, wasteArea.topY, ace1.rightX+7, wasteArea.topY+2, style)
	ace3 = makeBox(s, "Ace", ace2.rightX+1, wasteArea.topY, ace2.rightX+7, wasteArea.topY+2, style)
	ace4 = makeBox(s, "Ace", ace3.rightX+1, wasteArea.topY, ace3.rightX+7, wasteArea.topY+2, style)
	playArea = makeBox(s, "Tableau", 0, ace1.botY+2, 23, ace2.botY+16, style)
	x := playArea.rightX + 5
	y := playArea.topY
	putString(s, (w-x)/2-3, y, style, "Moves:")
	y++
	putString(s, x, y, style, "All moves are a two character instruction.")
	y++
	putString(s, x, y, style, "17 will move from stack 1 to stack 7.")
	y++
	putString(s, x, y, style, "w1 will move from waste to stack 1.")
	y++
	putString(s, x, y, style, "wa will move from waste to an ace stack.")
	y++
	putString(s, x, y, style, "1a will move from stack 1 to an ace stack.")
	s.Show()
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
			for playArea.topY+k < playArea.botY-2 {
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
// If all cards have been dealt from the deck it turns the wast stack into the deck
// and increments the pass count.
//
// stacks: A slice containing all the stacks (could probably be just the wast stack)
// deck: Apointer to the deck
// pass: THe pass count
//
// returns: The pas count
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

// CalcMove is the primary function for computer play allowing the computer to figure out the next move
//
// stacks: A slice containing all stacks of cards
// deck: A pointer to the deck used in case we need to deal to the waste stack
// pass: The pass count
//
// Returns: THe index of the stack to move from, the index of the stack to move to, and the pass count.
// If the pass count goes to 3 it will return -1, -1, and 3 ending the game.
func calcMove(stacks []solitaire.Pile, deck *generic.Deck, pass int) (int, int, int) {
	done := false
	index := 0
	for !done {
		for i, stack1 := range stacks {
			for j, stack2 := range stacks {
				if i == j {
					continue
				}
				if j > 7 && stack1.Firstfaceup < len(stack1.Cards)-1 {
					index = len(stack1.Cards) - 1
				} else {
					index = stack1.Firstfaceup
				}
				if stack1.CheckMove(&stack2, index) {
					return i, j, pass
				}
			}
		}
		pass = dealToWaste(stacks[:], deck, pass)
		if pass == 3 {
			done = true
		}
	}
	return -1, -1, pass
}

// PlayGame if the main function that handles all aspects of the game.
//
// s: Screnn variable.
// style: The style for the screen.
func playGame(s tcell.Screen, style tcell.Style) {
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
	ch := ' '
	putString(s, 0, h-1, style, fmt.Sprintf(" Enter C for computer, U for user or Q to quit:"))
	s.Show()
	for ch != 'C' && ch != 'U' {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyRune:
				switch unicode.ToUpper(ev.Rune()) {
				case 'Q':
					return
				case 'C':
					ch = 'C'
				case 'U':
					ch = 'U'
				}
			}
		}
	}
	putString(s, 0, h-1, style, strings.Repeat(" ", w-1))
	s.Show()
	pass := 0
	movefrom := -1
	moveto := -1
	for pass < 3 {
		showStacks(s, stacks, style)
		putString(s, 0, h-1, style, fmt.Sprintf("Pass# %02d, Waste# %02d, Deck# %02d", pass, len(stacks[7].Cards), len(deck.Cards)-deck.LastDealt))
		s.Show()
		switch ch {
		case 'U':
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyRune:
					switch unicode.ToUpper(ev.Rune()) {
					case 'Q':
						return
					case 'D':
						pass = dealToWaste(stacks[:], &deck, pass)
					case '1', '2', '3', '4', '5', '6', '7':
						logger.Printf("in tableau")
						if movefrom == -1 {
							movefrom = int(ev.Rune()-'0') - 1
							moveto = -1
						} else {
							if moveto == -1 {
								moveto = int(ev.Rune()-'0') - 1
							}
						}
					case 'W':
						logger.Printf("in waste")
						movefrom = 7
					case 'A':
						logger.Printf("in aces")
						if movefrom != -1 {
							var index int
							if len(stacks[movefrom].Cards) == 0 {
								index = 0
							} else {
								index = len(stacks[movefrom].Cards) - 1
							}
							switch stacks[movefrom].Cards[index].Suit {
							case "S":
								moveto = 8
							case "H":
								moveto = 9
							case "D":
								moveto = 10
							case "C":
								moveto = 11
							}
						}
					}
				}
			}
		case 'C':
			time.Sleep(time.Second) // give time to see what is being done
			movefrom, moveto, pass = calcMove(stacks[:], &deck, pass)
		}
		if movefrom > -1 && moveto > -1 {
			logger.Printf("from: %d to: %d faceup %d cards, %d lento %d", movefrom, moveto, stacks[movefrom].Firstfaceup, len(stacks[movefrom].Cards)-1, len(stacks[moveto].Cards))
			if moveto != movefrom {
				var index int
				if moveto > 7 && stacks[movefrom].Firstfaceup < len(stacks[movefrom].Cards)-1 {
					index = len(stacks[movefrom].Cards) - 1
				} else {
					index = stacks[movefrom].Firstfaceup
				}
				if stacks[movefrom].CheckMove(&stacks[moveto], index) {
					stacks[movefrom].DoMove(&stacks[moveto], index)
					logger.Printf("Moved from: %+v, to: %+v", stacks[movefrom], stacks[moveto])
				}
			}
			movefrom = -1
			moveto = -1
		}
	}
}

func main() {
	if err != nil {
		log.Fatal("Error opening error log.\n")
	}
	defer f.Close()
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
	defer s.Fini()
	drawScreen(s, tcell.StyleDefault)
	playGame(s, tcell.StyleDefault)
}
