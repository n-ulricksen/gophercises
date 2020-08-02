package deck

type Suit uint8

const (
	Club Suit = iota + 1
	Spade
	Heart
	Diamond
)

type Rank uint8

type Card struct {
}
