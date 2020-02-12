// Package generic is a set of types and meghods used by all cardgames.
package generic

import (
	"math/rand"
	"time"
)

// Dealer is an interface for dealing cards
type Dealer interface {
	Deal(n, nfaceup int) []Card
}

// Shuffler is the interface for shuffling cards
type Shuffler interface {
	Shuffle()
}

// Card holds the card suits and types in the deck.
// The rank and suit are held both in numeric and string format.
// Ace is the high card in a suit by default.
// It also contains a color and whether the card is faceup or not.
type Card struct {
	Rank   string
	Suit   string
	Rvalue int    // Numeric value of a cards rank
	Svalue int    // Numeric value of cards suit
	Color  string // Color of card, red or black
	Faceup bool
}

// Deck holds the cards in the deck to be shuffled and the last one dealt.
// AllDealt tells if all the cards in the deck have been dealt as we don't actually remove the cards we have dealt.
type Deck struct {
	Cards     []Card
	LastDealt int
	AllDealt  bool
}

// NewCard returns a new card with the given criteria.
// Rank, suit, and color are strings
// Rv and sv are the rank and suit numeric value respectively.
// Up is a bollean telling if the card is face up.
//
// Note that color, rv, and sv could be set automatecally.
func NewCard(rank, suit, color string, rv, sv int, up bool) Card {
	return Card{
		Rank:   rank,
		Suit:   suit,
		Color:  color,
		Rvalue: rv,
		Svalue: sv,
		Faceup: up,
	}
}

// NewDeck returns a deck of cards to be used.
func NewDeck() (deck Deck) {

	// Valid types include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	rvalues := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	// Valid suits include Heart, Diamond, Club & Spade
	suits := []string{"H", "D", "C", "S"}
	svalues := []int{8, 4, 2, 16} // In order clubs, diamonds, hearts, spades
	colors := []string{"red", "red", "black", "black"}

	// Loop over each type and suit appending to the deck
	for i, v1 := range ranks {
		for n, v2 := range suits {
			card := NewCard(v1, v2, colors[n], rvalues[i], svalues[n], false)
			deck.Cards = append(deck.Cards, card)
		}
	}
	deck.LastDealt = 0
	deck.AllDealt = false
	return
}

// Shuffle the deck.
func (d *Deck) Shuffle() {
	for i := range d.Cards {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d.Cards[r], d.Cards[i] = d.Cards[i], d.Cards[r]
		}
	}
}

// Deal a specified number of cards and set a specified number face up.
func (d *Deck) Deal(n, nfaceup int) (hand []Card) {
	for i := d.LastDealt; i < d.LastDealt+n && i < len(d.Cards); i++ {
		card := d.Cards[i]
		hand = append(hand, card)
	}
	d.LastDealt = d.LastDealt + n
	if d.LastDealt >= len(d.Cards) {
		d.AllDealt = true
	}
	if nfaceup > len(hand) {
		nfaceup = len(hand)
	}
	for i := len(hand) - nfaceup; i < len(hand); i++ {
		hand[i].Faceup = true
	}
	return hand
}

// Turn a card over so if it is face down turn it face up and vice versa
func (c *Card) Turn() {
	if c.Faceup {
		c.Faceup = false
	} else {
		c.Faceup = true
	}
}

// Seed our randomness with the current time
func init() {
	rand.Seed(time.Now().UnixNano())
}
