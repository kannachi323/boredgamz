package poker

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Card struct {
	Rank int    `json:"rank"`
	Suit string `json:"suit"`
	Code string `json:"code"`
}

var suits = []string{"s", "h", "d", "c"}
var rankLabels = map[int]string{2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "T", 11: "J", 12: "Q", 13: "K", 14: "A"}

func NewDeck() []Card {
	deck := make([]Card, 0, 52)
	for _, suit := range suits {
		for rank := 2; rank <= 14; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit, Code: rankLabels[rank] + suit})
		}
	}
	return deck
}

func SecureDeck() ([]Card, error) {
	deck := NewDeck()
	for index := len(deck) - 1; index > 0; index-- {
		pick, err := rand.Int(rand.Reader, big.NewInt(int64(index+1)))
		if err != nil {
			return nil, errors.New("secure deck shuffle failed")
		}
		swap := int(pick.Int64())
		deck[index], deck[swap] = deck[swap], deck[index]
	}
	return deck, nil
}
