package poker

import "sort"

type HandRank struct {
	Category int    `json:"category"`
	Name     string `json:"name"`
	Kickers  []int  `json:"kickers"`
}

var categoryNames = []string{"High card", "One pair", "Two pair", "Three of a kind", "Straight", "Flush", "Full house", "Four of a kind", "Straight flush"}

func EvaluateBest(cards []Card) HandRank {
	if len(cards) < 5 {
		return HandRank{Name: "Incomplete"}
	}
	best := HandRank{Category: -1}
	for a := 0; a < len(cards)-4; a++ {
		for b := a + 1; b < len(cards)-3; b++ {
			for c := b + 1; c < len(cards)-2; c++ {
				for d := c + 1; d < len(cards)-1; d++ {
					for e := d + 1; e < len(cards); e++ {
						rank := evaluateFive([5]Card{cards[a], cards[b], cards[c], cards[d], cards[e]})
						if CompareRanks(rank, best) > 0 {
							best = rank
						}
					}
				}
			}
		}
	}
	return best
}

func CompareRanks(left, right HandRank) int {
	if left.Category != right.Category {
		return compareInt(left.Category, right.Category)
	}
	length := max(len(left.Kickers), len(right.Kickers))
	for index := 0; index < length; index++ {
		var a, b int
		if index < len(left.Kickers) {
			a = left.Kickers[index]
		}
		if index < len(right.Kickers) {
			b = right.Kickers[index]
		}
		if a != b {
			return compareInt(a, b)
		}
	}
	return 0
}

func evaluateFive(cards [5]Card) HandRank {
	counts := make(map[int]int)
	values := make([]int, 0, 5)
	flush := true
	for index, card := range cards {
		counts[card.Rank]++
		values = append(values, card.Rank)
		if index > 0 && card.Suit != cards[0].Suit {
			flush = false
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	straightHigh := 0
	unique := uniqueDescending(values)
	if len(unique) == 5 {
		if unique[0]-unique[4] == 4 {
			straightHigh = unique[0]
		}
		if equalInts(unique, []int{14, 5, 4, 3, 2}) {
			straightHigh = 5
		}
	}
	if flush && straightHigh > 0 {
		return ranked(8, straightHigh)
	}
	groups := grouped(counts)
	if groups[0].count == 4 {
		return ranked(7, groups[0].rank, groups[1].rank)
	}
	if groups[0].count == 3 && groups[1].count == 2 {
		return ranked(6, groups[0].rank, groups[1].rank)
	}
	if flush {
		return ranked(5, values...)
	}
	if straightHigh > 0 {
		return ranked(4, straightHigh)
	}
	if groups[0].count == 3 {
		return ranked(3, expandGroups(groups)...)
	}
	if groups[0].count == 2 && groups[1].count == 2 {
		return ranked(2, groups[0].rank, groups[1].rank, groups[2].rank)
	}
	if groups[0].count == 2 {
		return ranked(1, expandGroups(groups)...)
	}
	return ranked(0, values...)
}

type rankGroup struct{ rank, count int }

func grouped(counts map[int]int) []rankGroup {
	groups := make([]rankGroup, 0, len(counts))
	for rank, count := range counts {
		groups = append(groups, rankGroup{rank, count})
	}
	sort.Slice(groups, func(i, j int) bool {
		if groups[i].count != groups[j].count {
			return groups[i].count > groups[j].count
		}
		return groups[i].rank > groups[j].rank
	})
	return groups
}

func expandGroups(groups []rankGroup) []int {
	values := make([]int, 0, 5)
	for _, group := range groups {
		for count := 0; count < group.count; count++ {
			values = append(values, group.rank)
		}
	}
	return values
}

func uniqueDescending(values []int) []int {
	result := make([]int, 0, len(values))
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}

func equalInts(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
func compareInt(left, right int) int {
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}
func ranked(category int, kickers ...int) HandRank {
	return HandRank{Category: category, Name: categoryNames[category], Kickers: kickers}
}
