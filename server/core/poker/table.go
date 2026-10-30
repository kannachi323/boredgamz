package poker

import (
	"errors"
	"fmt"
)

const (
	StartingStack = 1000
	SmallBlind    = 10
	BigBlind      = 20
)

type Seat struct {
	ID        string
	Name      string
	Stack     int
	Cards     []Card
	Folded    bool
	AllIn     bool
	RoundBet  int
	TotalBet  int
	Acted     bool
	Connected bool
	Discarded bool
}

type Table struct {
	Mode         string
	Seats        [2]*Seat
	Deck         []Card
	Community    []Card
	Phase        string
	Dealer       int
	Turn         int
	Pot          int
	AwardedPot   int
	CurrentBet   int
	MinimumRaise int
	Hand         int
	Winners      []int
	WinningHand  string
	Message      string
}

func NewTable(mode string, first, second *Seat) (*Table, error) {
	if mode != "holdem" && mode != "five-card-draw" {
		return nil, errors.New("unsupported poker mode")
	}
	for _, seat := range []*Seat{first, second} {
		if seat == nil || seat.ID == "" {
			return nil, errors.New("two identified players are required")
		}
		if seat.Stack <= 0 {
			seat.Stack = StartingStack
		}
		seat.Connected = true
	}
	table := &Table{Mode: mode, Seats: [2]*Seat{first, second}, Dealer: -1}
	if err := table.StartHand(); err != nil {
		return nil, err
	}
	return table, nil
}

func (t *Table) StartHand() error {
	if t.Seats[0] == nil || t.Seats[1] == nil {
		return errors.New("two players are required")
	}
	if t.Seats[0].Stack == 0 || t.Seats[1].Stack == 0 {
		return errors.New("both players need chips")
	}
	deck, err := SecureDeck()
	if err != nil {
		return err
	}
	t.Deck, t.Community, t.Pot, t.AwardedPot = deck, nil, 0, 0
	t.CurrentBet, t.MinimumRaise = 0, BigBlind
	t.Dealer = (t.Dealer + 1) % 2
	t.Hand++
	t.Winners, t.WinningHand, t.Message = nil, "", ""
	for _, seat := range t.Seats {
		seat.Cards, seat.Folded, seat.AllIn = nil, false, false
		seat.RoundBet, seat.TotalBet, seat.Acted, seat.Discarded = 0, 0, false, false
	}
	if t.Mode == "holdem" {
		t.startHoldem()
	} else {
		t.startFiveCardDraw()
	}
	return nil
}

func (t *Table) startHoldem() {
	for count := 0; count < 2; count++ {
		for seat := 0; seat < 2; seat++ {
			t.Seats[seat].Cards = append(t.Seats[seat].Cards, t.draw())
		}
	}
	t.Phase = "preflop"
	t.pay(t.Dealer, SmallBlind)
	t.pay(1-t.Dealer, BigBlind)
	t.CurrentBet, t.MinimumRaise, t.Turn = BigBlind, BigBlind, t.Dealer
	t.Message = "Blinds are posted. The button acts first."
	if t.Seats[0].AllIn || t.Seats[1].AllIn {
		t.runOut()
	}
}

func (t *Table) startFiveCardDraw() {
	for seat := 0; seat < 2; seat++ {
		t.pay(seat, SmallBlind)
	}
	for count := 0; count < 5; count++ {
		for seat := 0; seat < 2; seat++ {
			t.Seats[seat].Cards = append(t.Seats[seat].Cards, t.draw())
		}
	}
	t.collectBets()
	if t.Seats[0].AllIn || t.Seats[1].AllIn {
		t.Phase, t.Turn, t.Message = "draw", 1-t.Dealer, "Choose up to four cards to replace."
	} else {
		t.Phase, t.Turn, t.MinimumRaise = "draw-bet", 1-t.Dealer, BigBlind
		t.Message = "Opening betting round."
	}
}

func (t *Table) draw() Card {
	card := t.Deck[0]
	t.Deck = t.Deck[1:]
	return card
}

func (t *Table) pay(index, amount int) int {
	seat := t.Seats[index]
	paid := min(amount, seat.Stack)
	seat.Stack -= paid
	seat.RoundBet += paid
	seat.TotalBet += paid
	seat.AllIn = seat.Stack == 0
	return paid
}

func (t *Table) collectBets() {
	for _, seat := range t.Seats {
		t.Pot += seat.RoundBet
		seat.RoundBet, seat.Acted = 0, false
	}
	t.CurrentBet, t.MinimumRaise = 0, BigBlind
}

func (t *Table) ToCall(index int) int { return max(0, t.CurrentBet-t.Seats[index].RoundBet) }

func (t *Table) LegalActions(index int) []string {
	if index != t.Turn || !t.bettingPhase() || t.Seats[index].Folded || t.Seats[index].AllIn {
		return nil
	}
	actions := []string{"fold"}
	if t.ToCall(index) == 0 {
		actions = append(actions, "check")
	} else {
		actions = append(actions, "call")
	}
	if !t.Seats[1-index].AllIn && t.Seats[index].RoundBet+t.Seats[index].Stack >= t.CurrentBet+t.MinimumRaise {
		actions = append(actions, "raise")
	}
	if !t.Seats[1-index].AllIn && t.Seats[index].Stack > 0 {
		actions = append(actions, "all-in")
	}
	return actions
}

func (t *Table) bettingPhase() bool {
	return t.Phase == "preflop" || t.Phase == "flop" || t.Phase == "turn" || t.Phase == "river" || t.Phase == "draw-bet" || t.Phase == "final-bet"
}

func (t *Table) Act(index int, action string, amount int) error {
	if index < 0 || index > 1 || index != t.Turn {
		return errors.New("it is not this player's turn")
	}
	if !t.bettingPhase() {
		return errors.New("betting is not active")
	}
	seat := t.Seats[index]
	if seat.Folded || seat.AllIn {
		return errors.New("player cannot act")
	}
	switch action {
	case "fold":
		seat.Folded, seat.Acted = true, true
		t.finishByFold(1 - index)
		return nil
	case "check":
		if t.ToCall(index) != 0 {
			return errors.New("cannot check while facing a bet")
		}
		seat.Acted = true
	case "call":
		if t.ToCall(index) == 0 {
			return errors.New("there is no bet to call")
		}
		t.pay(index, t.ToCall(index))
		seat.Acted = true
	case "raise":
		if err := t.raiseTo(index, amount, false); err != nil {
			return err
		}
	case "all-in":
		if seat.Stack == 0 {
			return errors.New("player has no chips")
		}
		target := seat.RoundBet + seat.Stack
		if target > t.CurrentBet {
			if err := t.raiseTo(index, target, true); err != nil {
				return err
			}
		} else {
			t.pay(index, seat.Stack)
			seat.Acted = true
		}
	default:
		return errors.New("unsupported poker action")
	}
	t.afterAction(index)
	return nil
}

func (t *Table) raiseTo(index, target int, allIn bool) error {
	seat := t.Seats[index]
	if t.Seats[1-index].AllIn {
		return errors.New("cannot raise an all-in player")
	}
	if target <= t.CurrentBet {
		return errors.New("raise must exceed the current bet")
	}
	needed := target - seat.RoundBet
	if needed > seat.Stack {
		return errors.New("raise exceeds the player's stack")
	}
	raiseSize := target - t.CurrentBet
	if raiseSize < t.MinimumRaise && !(allIn && needed == seat.Stack) {
		return fmt.Errorf("minimum raise is %d", t.CurrentBet+t.MinimumRaise)
	}
	t.pay(index, needed)
	if raiseSize >= t.MinimumRaise {
		t.MinimumRaise = raiseSize
	}
	t.CurrentBet = target
	seat.Acted = true
	t.Seats[1-index].Acted = false
	return nil
}

func (t *Table) afterAction(index int) {
	other := 1 - index
	if t.Seats[other].Folded {
		t.finishByFold(index)
		return
	}
	matched := t.Seats[0].RoundBet == t.Seats[1].RoundBet
	if (t.Seats[0].AllIn || t.Seats[1].AllIn) && (matched || t.Seats[other].AllIn) {
		if t.Mode == "five-card-draw" && t.Phase == "draw-bet" {
			t.advancePhase()
		} else {
			t.runOut()
		}
		return
	}
	if matched && t.Seats[0].Acted && t.Seats[1].Acted {
		t.advancePhase()
		return
	}
	t.Turn = other
	if t.Seats[other].AllIn {
		t.Turn = index
	}
}

func (t *Table) advancePhase() {
	t.collectBets()
	if t.Mode == "five-card-draw" {
		if t.Phase == "draw-bet" {
			t.Phase, t.Turn = "draw", 1-t.Dealer
			t.Message = "Choose up to four cards to replace."
			return
		}
		t.showdown()
		return
	}
	switch t.Phase {
	case "preflop":
		t.burn()
		t.Community = append(t.Community, t.draw(), t.draw(), t.draw())
		t.Phase = "flop"
	case "flop":
		t.burn()
		t.Community = append(t.Community, t.draw())
		t.Phase = "turn"
	case "turn":
		t.burn()
		t.Community = append(t.Community, t.draw())
		t.Phase = "river"
	case "river":
		t.showdown()
		return
	}
	t.Turn = 1 - t.Dealer
	t.Message = "Betting continues on the " + t.Phase + "."
}

func (t *Table) Discard(index int, positions []int) error {
	if t.Mode != "five-card-draw" || t.Phase != "draw" {
		return errors.New("cards cannot be discarded now")
	}
	if index != t.Turn || t.Seats[index].Discarded {
		return errors.New("it is not this player's draw")
	}
	if len(positions) > 4 {
		return errors.New("at most four cards may be replaced")
	}
	seen := map[int]bool{}
	for _, position := range positions {
		if position < 0 || position >= 5 || seen[position] {
			return errors.New("invalid discard selection")
		}
		seen[position] = true
	}
	for position := range seen {
		t.Seats[index].Cards[position] = t.draw()
	}
	t.Seats[index].Discarded = true
	other := 1 - index
	if t.Seats[other].Discarded {
		if t.Seats[0].AllIn || t.Seats[1].AllIn {
			t.showdown()
		} else {
			t.Phase, t.Turn, t.Message = "final-bet", 1-t.Dealer, "Final betting round."
		}
		return nil
	}
	t.Turn = other
	return nil
}

func (t *Table) runOut() {
	t.collectBets()
	if t.Mode == "five-card-draw" {
		t.showdown()
		return
	}
	if len(t.Community) == 0 {
		t.burn()
		t.Community = append(t.Community, t.draw(), t.draw(), t.draw())
	}
	for len(t.Community) < 5 {
		t.burn()
		t.Community = append(t.Community, t.draw())
	}
	t.showdown()
}

func (t *Table) burn() {
	_ = t.draw()
}

func (t *Table) finishByFold(winner int) {
	t.collectBets()
	t.AwardedPot = t.Pot
	t.Seats[winner].Stack += t.Pot
	t.Winners, t.Phase, t.Message = []int{winner}, "hand-finished", t.Seats[winner].Name+" wins the uncontested pot."
	t.Pot = 0
}

func (t *Table) showdown() {
	t.collectBets()
	left, right := EvaluateBest(append(append([]Card{}, t.Seats[0].Cards...), t.Community...)), EvaluateBest(append(append([]Card{}, t.Seats[1].Cards...), t.Community...))
	comparison := CompareRanks(left, right)
	if comparison > 0 {
		t.Winners, t.WinningHand = []int{0}, left.Name
	} else if comparison < 0 {
		t.Winners, t.WinningHand = []int{1}, right.Name
	} else {
		t.Winners, t.WinningHand = []int{0, 1}, left.Name
	}
	matched := min(t.Seats[0].TotalBet, t.Seats[1].TotalBet)
	for index := range t.Seats {
		excess := t.Seats[index].TotalBet - matched
		if excess > 0 {
			t.Seats[index].Stack += excess
			t.Pot -= excess
		}
	}
	share := t.Pot / len(t.Winners)
	t.AwardedPot = t.Pot
	remainder := t.Pot % len(t.Winners)
	for winnerIndex, seatIndex := range t.Winners {
		t.Seats[seatIndex].Stack += share
		if winnerIndex == 0 {
			t.Seats[seatIndex].Stack += remainder
		}
	}
	if len(t.Winners) == 1 {
		t.Message = t.Seats[t.Winners[0]].Name + " wins with " + t.WinningHand + "."
	} else {
		t.Message = "The pot is split with " + t.WinningHand + "."
	}
	t.Phase, t.Pot = "showdown", 0
}
