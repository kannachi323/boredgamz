package poker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testTable(t *testing.T, mode string) *Table {
	t.Helper()
	table, err := NewTable(mode, &Seat{ID: "a", Name: "Alice"}, &Seat{ID: "b", Name: "Bob"})
	require.NoError(t, err)
	return table
}

func TestHoldemPostsHeadsUpBlindsAndValidatesRaises(t *testing.T) {
	table := testTable(t, "holdem")
	require.Equal(t, 0, table.Dealer)
	require.Equal(t, 0, table.Turn)
	require.Equal(t, 10, table.Seats[0].RoundBet)
	require.Equal(t, 20, table.Seats[1].RoundBet)
	require.EqualError(t, table.Act(0, "check", 0), "cannot check while facing a bet")
	require.EqualError(t, table.Act(0, "raise", 30), "minimum raise is 40")
	require.NoError(t, table.Act(0, "call", 0))
	require.NoError(t, table.Act(1, "check", 0))
	require.Equal(t, "flop", table.Phase)
	require.Len(t, table.Community, 3)
	require.Equal(t, 1, table.Turn)
}

func TestHoldemFoldAwardsEveryCommittedChip(t *testing.T) {
	table := testTable(t, "holdem")
	require.NoError(t, table.Act(0, "fold", 0))
	require.Equal(t, "hand-finished", table.Phase)
	require.Equal(t, []int{1}, table.Winners)
	require.Equal(t, 1010, table.Seats[1].Stack)
	require.Equal(t, 990, table.Seats[0].Stack)
	require.Equal(t, 30, table.AwardedPot)
	require.Equal(t, 2000, table.Seats[0].Stack+table.Seats[1].Stack)
}

func TestHoldemAllInRunsBoardAndRefundsUnmatchedChips(t *testing.T) {
	table, err := NewTable("holdem", &Seat{ID: "short", Name: "Short", Stack: 100}, &Seat{ID: "deep", Name: "Deep", Stack: 1000})
	require.NoError(t, err)
	require.NoError(t, table.Act(0, "all-in", 0))
	require.NoError(t, table.Act(1, "call", 0))
	require.Equal(t, "showdown", table.Phase)
	require.Len(t, table.Community, 5)
	require.Equal(t, 1100, table.Seats[0].Stack+table.Seats[1].Stack)
	require.Zero(t, table.Pot)
	require.Equal(t, 200, table.AwardedPot)
}

func TestShortBlindAutomaticallyRunsOutWithoutDeadlock(t *testing.T) {
	table, err := NewTable("holdem", &Seat{ID: "tiny", Name: "Tiny", Stack: 5}, &Seat{ID: "deep-auto", Name: "Deep", Stack: 1000})
	require.NoError(t, err)
	require.Equal(t, "showdown", table.Phase)
	require.Len(t, table.Community, 5)
	require.Equal(t, 1005, table.Seats[0].Stack+table.Seats[1].Stack)
}

func TestFiveCardDrawReplacesSelectedCardsOnce(t *testing.T) {
	table := testTable(t, "five-card-draw")
	require.Equal(t, 1-table.Dealer, table.Turn)
	require.NoError(t, table.Act(table.Turn, "check", 0))
	require.NoError(t, table.Act(table.Turn, "check", 0))
	require.Equal(t, "draw", table.Phase)
	require.Equal(t, 1-table.Dealer, table.Turn)
	first := table.Turn
	before := append([]Card{}, table.Seats[first].Cards...)
	require.NoError(t, table.Discard(first, []int{0, 2, 4}))
	require.NotEqual(t, before[0], table.Seats[first].Cards[0])
	require.Equal(t, before[1], table.Seats[first].Cards[1])
	require.EqualError(t, table.Discard(first, nil), "it is not this player's draw")
	require.NoError(t, table.Discard(table.Turn, nil))
	require.Equal(t, "final-bet", table.Phase)
	require.Equal(t, 1-table.Dealer, table.Turn)
}

func TestFiveCardDrawAllInStillReceivesDrawChoice(t *testing.T) {
	table, err := NewTable("five-card-draw", &Seat{ID: "draw-short", Name: "Short", Stack: 10}, &Seat{ID: "draw-deep", Name: "Deep", Stack: 100})
	require.NoError(t, err)
	require.Equal(t, "draw", table.Phase)
	require.NoError(t, table.Discard(table.Turn, []int{0}))
	require.NoError(t, table.Discard(table.Turn, []int{1}))
	require.Equal(t, "showdown", table.Phase)
	require.Equal(t, 110, table.Seats[0].Stack+table.Seats[1].Stack)
}
