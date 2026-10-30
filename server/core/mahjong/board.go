package mahjong

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Tile struct {
	ID    int    `json:"id"`
	Face  string `json:"face"`
	Group string `json:"group"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Z     int    `json:"z"`
}

type Board struct {
	Tiles    []Tile
	Solution [][2]int
}

type position struct{ x, y, z int }
type facePair struct{ leftFace, rightFace, group string }

func NewBoard(layout string) (*Board, error) {
	positions := fullPositions()
	if layout == "compact" {
		positions = compactPositions()
	}
	if len(positions)%2 != 0 {
		return nil, errors.New("mahjong layout needs an even tile count")
	}
	solution, err := removalOrder(positions, make([]bool, len(positions)))
	if err != nil {
		return nil, err
	}
	pairs := standardPairs()
	if err := secureShuffle(pairs); err != nil {
		return nil, err
	}
	tiles := make([]Tile, len(positions))
	for pairIndex, indexes := range solution {
		pair := pairs[pairIndex]
		for offset, index := range indexes {
			pos := positions[index]
			face := pair.leftFace
			if offset == 1 {
				face = pair.rightFace
			}
			tiles[index] = Tile{ID: index, Face: face, Group: pair.group, X: pos.x, Y: pos.y, Z: pos.z}
		}
	}
	return &Board{Tiles: tiles, Solution: solution}, nil
}

func (b *Board) IsFree(removed []bool, index int) bool {
	if index < 0 || index >= len(b.Tiles) || len(removed) != len(b.Tiles) || removed[index] {
		return false
	}
	tile := b.Tiles[index]
	for otherIndex, other := range b.Tiles {
		if otherIndex != index && !removed[otherIndex] && other.Z > tile.Z && other.X == tile.X && other.Y == tile.Y {
			return false
		}
	}
	leftBlocked, rightBlocked := false, false
	for otherIndex, other := range b.Tiles {
		if otherIndex == index || removed[otherIndex] || other.Z != tile.Z || other.Y != tile.Y {
			continue
		}
		if other.X == tile.X-1 {
			leftBlocked = true
		}
		if other.X == tile.X+1 {
			rightBlocked = true
		}
	}
	return !leftBlocked || !rightBlocked
}

func (b *Board) CanRemove(removed []bool, first, second int) bool {
	return first != second && b.IsFree(removed, first) && b.IsFree(removed, second) && b.Tiles[first].Group == b.Tiles[second].Group
}

func (b *Board) Remaining(removed []bool) int {
	count := 0
	for index := range b.Tiles {
		if index >= len(removed) || !removed[index] {
			count++
		}
	}
	return count
}

func (b *Board) AvailablePair(removed []bool) (int, int, bool) {
	for first := range b.Tiles {
		if !b.IsFree(removed, first) {
			continue
		}
		for second := first + 1; second < len(b.Tiles); second++ {
			if b.CanRemove(removed, first, second) {
				return first, second, true
			}
		}
	}
	return 0, 0, false
}

func (b *Board) Reshuffle(removed []bool) error {
	if len(removed) != len(b.Tiles) {
		return errors.New("invalid removed tile state")
	}
	positions := make([]position, len(b.Tiles))
	for index, tile := range b.Tiles {
		positions[index] = position{tile.X, tile.Y, tile.Z}
	}
	solution, err := removalOrder(positions, removed)
	if err != nil {
		return err
	}
	byGroup := map[string][]string{}
	for index, tile := range b.Tiles {
		if !removed[index] {
			byGroup[tile.Group] = append(byGroup[tile.Group], tile.Face)
		}
	}
	pairs := make([]facePair, 0, len(solution))
	for group, faces := range byGroup {
		if len(faces)%2 != 0 {
			return errors.New("unpaired mahjong group")
		}
		for index := 0; index < len(faces); index += 2 {
			pairs = append(pairs, facePair{faces[index], faces[index+1], group})
		}
	}
	if len(pairs) != len(solution) {
		return errors.New("remaining pair count mismatch")
	}
	if err := secureShuffle(pairs); err != nil {
		return err
	}
	for pairIndex, indexes := range solution {
		pair := pairs[pairIndex]
		b.Tiles[indexes[0]].Face, b.Tiles[indexes[0]].Group = pair.leftFace, pair.group
		b.Tiles[indexes[1]].Face, b.Tiles[indexes[1]].Group = pair.rightFace, pair.group
	}
	b.Solution = solution
	return nil
}

func removalOrder(positions []position, initiallyRemoved []bool) ([][2]int, error) {
	if len(initiallyRemoved) != len(positions) {
		return nil, errors.New("invalid removed tile state")
	}
	remaining := 0
	for _, removed := range initiallyRemoved {
		if !removed {
			remaining++
		}
	}
	if remaining%2 != 0 {
		return nil, errors.New("mahjong layout has an odd number of remaining tiles")
	}
	for attempt := 0; attempt < 256; attempt++ {
		removed := append([]bool{}, initiallyRemoved...)
		solution := make([][2]int, 0, remaining/2)
		for len(solution) < remaining/2 {
			free := freePositions(positions, removed)
			if len(free) < 2 {
				break
			}
			firstAt, err := secureIndex(len(free))
			if err != nil {
				return nil, err
			}
			first := free[firstAt]
			free = append(free[:firstAt], free[firstAt+1:]...)
			secondAt, err := secureIndex(len(free))
			if err != nil {
				return nil, err
			}
			second := free[secondAt]
			removed[first], removed[second] = true, true
			solution = append(solution, [2]int{first, second})
		}
		if len(solution) == remaining/2 {
			return solution, nil
		}
	}
	return nil, errors.New("mahjong layout is not reducible")
}

func freePositions(positions []position, removed []bool) []int {
	board := &Board{Tiles: make([]Tile, len(positions))}
	for index, pos := range positions {
		board.Tiles[index] = Tile{ID: index, X: pos.x, Y: pos.y, Z: pos.z}
	}
	free := []int{}
	for index := range positions {
		if board.IsFree(removed, index) {
			free = append(free, index)
		}
	}
	return free
}

func fullPositions() []position {
	result := []position{}
	for y := 0; y < 8; y++ {
		for x := 0; x < 12; x++ {
			result = append(result, position{x, y, 0})
		}
	}
	for y := 2; y < 6; y++ {
		for x := 1; x < 11; x++ {
			result = append(result, position{x, y, 1})
		}
	}
	for y := 3; y < 5; y++ {
		for x := 4; x < 8; x++ {
			result = append(result, position{x, y, 2})
		}
	}
	return result
}
func compactPositions() []position {
	result := []position{}
	for y := 0; y < 6; y++ {
		for x := 0; x < 8; x++ {
			result = append(result, position{x, y, 0})
		}
	}
	for y := 1; y < 4; y++ {
		for x := 1; x < 7; x++ {
			result = append(result, position{x, y, 1})
		}
	}
	for y := 1; y < 3; y++ {
		for x := 3; x < 6; x++ {
			result = append(result, position{x, y, 2})
		}
	}
	return result
}

func standardPairs() []facePair {
	pairs := make([]facePair, 0, 72)
	for _, suit := range []string{"萬", "筒", "索"} {
		for rank := 1; rank <= 9; rank++ {
			face := string(rune('0'+rank)) + suit
			group := face
			pairs = append(pairs, facePair{face, face, group}, facePair{face, face, group})
		}
	}
	for _, face := range []string{"東", "南", "西", "北", "中", "發", "白"} {
		pairs = append(pairs, facePair{face, face, face}, facePair{face, face, face})
	}
	pairs = append(pairs, facePair{"春", "夏", "season"}, facePair{"秋", "冬", "season"}, facePair{"梅", "蘭", "flower"}, facePair{"菊", "竹", "flower"})
	return pairs
}

func secureIndex(length int) (int, error) {
	if length <= 0 {
		return 0, errors.New("empty secure choice")
	}
	value, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
	if err != nil {
		return 0, errors.New("secure mahjong shuffle failed")
	}
	return int(value.Int64()), nil
}
func secureShuffle[T any](values []T) error {
	for index := len(values) - 1; index > 0; index-- {
		swap, err := secureIndex(index + 1)
		if err != nil {
			return err
		}
		values[index], values[swap] = values[swap], values[index]
	}
	return nil
}
