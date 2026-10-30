package mahjong

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullBoardUsesStandard144TilesAndIsSolvable(t *testing.T) {
	board, err := NewBoard("full")
	require.NoError(t, err)
	require.Len(t, board.Tiles, 144)
	require.Len(t, board.Solution, 72)
	removed := make([]bool, len(board.Tiles))
	for _, pair := range board.Solution {
		require.True(t, board.CanRemove(removed, pair[0], pair[1]))
		removed[pair[0]], removed[pair[1]] = true, true
	}
	require.Zero(t, board.Remaining(removed))
}

func TestCompactBoardIsSolvableAndHasUniqueInstances(t *testing.T) {
	board, err := NewBoard("compact")
	require.NoError(t, err)
	require.Len(t, board.Tiles, 72)
	seen := map[int]bool{}
	for _, tile := range board.Tiles {
		require.False(t, seen[tile.ID])
		seen[tile.ID] = true
	}
	removed := make([]bool, 72)
	for _, pair := range board.Solution {
		require.True(t, board.CanRemove(removed, pair[0], pair[1]))
		removed[pair[0]], removed[pair[1]] = true, true
	}
}

func TestCoveredAndInteriorTilesAreBlocked(t *testing.T) {
	board, err := NewBoard("full")
	require.NoError(t, err)
	removed := make([]bool, len(board.Tiles))
	baseCenter := -1
	for index, tile := range board.Tiles {
		if tile.X == 4 && tile.Y == 3 && tile.Z == 0 {
			baseCenter = index
			break
		}
	}
	require.NotEqual(t, -1, baseCenter)
	require.False(t, board.IsFree(removed, baseCenter))
	require.True(t, board.IsFree(removed, 0))
}

func TestSeasonAndFlowerTilesMatchTheirFamilies(t *testing.T) {
	board := &Board{Tiles: []Tile{{ID: 0, Face: "春", Group: "season", X: 0}, {ID: 1, Face: "冬", Group: "season", X: 2}, {ID: 2, Face: "梅", Group: "flower", X: 4}}}
	removed := make([]bool, 3)
	require.True(t, board.CanRemove(removed, 0, 1))
	require.False(t, board.CanRemove(removed, 0, 2))
}

func TestReshufflePreservesRemovedTilesAndRestoresSolvableOrder(t *testing.T) {
	board, err := NewBoard("compact")
	require.NoError(t, err)
	removed := make([]bool, len(board.Tiles))
	for _, pair := range board.Solution[:5] {
		removed[pair[0]], removed[pair[1]] = true, true
	}
	require.NoError(t, board.Reshuffle(removed))
	require.Len(t, board.Solution, 31)
	for _, pair := range board.Solution {
		require.True(t, board.CanRemove(removed, pair[0], pair[1]))
		removed[pair[0]], removed[pair[1]] = true, true
	}
	require.Zero(t, board.Remaining(removed))
}
