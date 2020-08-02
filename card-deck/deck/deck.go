package deck

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	return &Deck{}
}
