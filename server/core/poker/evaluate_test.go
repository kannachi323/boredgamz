package poker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func cards(codes ...string) []Card {
	deck := NewDeck()
	byCode := make(map[string]Card, len(deck))
	for _, card := range deck {
		byCode[card.Code] = card
	}
	result := make([]Card, 0, len(codes))
	for _, code := range codes {
		result = append(result, byCode[code])
	}
	return result
}

func TestDeckHas52UniqueCards(t *testing.T) {
	deck, err := SecureDeck()
	require.NoError(t, err)
	require.Len(t, deck, 52)
	seen := make(map[string]bool)
	for _, card := range deck {
		require.False(t, seen[card.Code])
		seen[card.Code] = true
	}
}

func TestEvaluateBestSelectsRoyalFlushFromSeven(t *testing.T) {
	rank := EvaluateBest(cards("As", "Ks", "Qs", "Js", "Ts", "2h", "2c"))
	require.Equal(t, "Straight flush", rank.Name)
	require.Equal(t, []int{14}, rank.Kickers)
}

func TestWheelStraightAndFullHouseOrdering(t *testing.T) {
	wheel := EvaluateBest(cards("As", "2h", "3d", "4c", "5s"))
	require.Equal(t, []int{5}, wheel.Kickers)
	fullHouse := EvaluateBest(cards("Kh", "Kd", "Kc", "2s", "2d"))
	flush := EvaluateBest(cards("Ah", "Jh", "8h", "5h", "3h"))
	require.Positive(t, CompareRanks(fullHouse, flush))
}

func TestPairKickersBreakTies(t *testing.T) {
	ace := EvaluateBest(cards("Qs", "Qh", "Ad", "6c", "2s"))
	king := EvaluateBest(cards("Qd", "Qc", "Kd", "6h", "2h"))
	require.Positive(t, CompareRanks(ace, king))
}
