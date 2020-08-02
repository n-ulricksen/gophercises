//go:generate stringer -type=Suit,Rank

package deck

import "fmt"

type Suit uint8

const (
	Club Suit = iota
	Spade
	Heart
	Diamond
	Joker
)

type Rank uint8

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}
