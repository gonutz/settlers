package game

import "testing"

func TestAdjacentTilesToCorner(t *testing.T) {
	checkTiles(t, AdjacentTilesToCorner(TileCorner{3, 1}), 1, 0, 2, 1, 3, 0)
	checkTiles(t, AdjacentTilesToCorner(TileCorner{4, 1}), 2, 1, 3, 0, 4, 1)
	checkTiles(t, AdjacentTilesToCorner(TileCorner{3, 2}), 1, 2, 2, 1, 3, 2)
	checkTiles(t, AdjacentTilesToCorner(TileCorner{4, 2}), 2, 1, 3, 2, 4, 1)
}

func checkTiles(t *testing.T, tiles [3]TilePosition, xys ...int) {
	for i, tile := range tiles {
		if tile.X != xys[i*2] || tile.Y != xys[i*2+1] {
			t.Errorf("tile %v: wanted %v,%v got %v,%v",
				i, xys[i*2], xys[i*2+1], tile.X, tile.Y)
		}
	}
}

func TestAdjacentEdgesToCorner(t *testing.T) {
	edges := AdjacentEdgesToCorner(TileCorner{3, 1})
	checkEdges(t, edges[:], 5, 1, 6, 0, 7, 1)

	edges = AdjacentEdgesToCorner(TileCorner{4, 1})
	checkEdges(t, edges[:], 7, 1, 8, 1, 9, 1)
}

func checkEdges(t *testing.T, edges []TileEdge, xys ...int) {
	for i, edge := range edges {
		if edge.X != xys[i*2] || edge.Y != xys[i*2+1] {
			t.Errorf("edge %v: wanted %v,%v got %v,%v",
				i, xys[i*2], xys[i*2+1], edge.X, edge.Y)
		}
	}
}

func TestAdjacentCornersToEdge(t *testing.T) {
	checkCorners(t, AdjacentCornersToEdge(TileEdge{4, 1}), 2, 1, 2, 2)
	checkCorners(t, AdjacentCornersToEdge(TileEdge{5, 1}), 2, 1, 3, 1)
	checkCorners(t, AdjacentCornersToEdge(TileEdge{5, 2}), 2, 2, 3, 2)
	checkCorners(t, AdjacentCornersToEdge(TileEdge{7, 1}), 3, 1, 4, 1)
	checkCorners(t, AdjacentCornersToEdge(TileEdge{7, 2}), 3, 2, 4, 2)
	checkCorners(t, AdjacentCornersToEdge(TileEdge{8, 1}), 4, 1, 4, 2)
}

func checkCorners(t *testing.T, corners [2]TileCorner, xys ...int) {
	for i, corner := range corners {
		if corner.X != xys[i*2] || corner.Y != xys[i*2+1] {
			t.Errorf("corner %v: wanted %v,%v got %v,%v",
				i, xys[i*2], xys[i*2+1], corner.X, corner.Y)
		}
	}
}

func TestAdjacentEdgesToEdge(t *testing.T) {
	edges := AdjacentEdgesToEdge(TileEdge{4, 1})
	checkEdges(t, edges[:], 3, 1, 3, 2, 5, 1, 5, 2)

	edges = AdjacentEdgesToEdge(TileEdge{5, 1})
	checkEdges(t, edges[:], 3, 1, 4, 1, 6, 0, 7, 1)

	edges = AdjacentEdgesToEdge(TileEdge{5, 2})
	checkEdges(t, edges[:], 3, 2, 4, 1, 6, 2, 7, 2)

	edges = AdjacentEdgesToEdge(TileEdge{7, 1})
	checkEdges(t, edges[:], 5, 1, 6, 0, 8, 1, 9, 1)

	edges = AdjacentEdgesToEdge(TileEdge{7, 2})
	checkEdges(t, edges[:], 5, 2, 6, 2, 8, 1, 9, 2)

	edges = AdjacentEdgesToEdge(TileEdge{8, 1})
	checkEdges(t, edges[:], 7, 1, 7, 2, 9, 1, 9, 2)
}

func TestAdjacentTilesToEdge(t *testing.T) {
	checkAdjacentTilesToEdge(t, TileEdge{5, 1}, TilePosition{1, 0}, TilePosition{2, 1})
	checkAdjacentTilesToEdge(t, TileEdge{7, 2}, TilePosition{2, 1}, TilePosition{3, 2})
	checkAdjacentTilesToEdge(t, TileEdge{6, 0}, TilePosition{1, 0}, TilePosition{3, 0})
	checkAdjacentTilesToEdge(t, TileEdge{8, 1}, TilePosition{2, 1}, TilePosition{4, 1})
	checkAdjacentTilesToEdge(t, TileEdge{5, 2}, TilePosition{1, 2}, TilePosition{2, 1})
	checkAdjacentTilesToEdge(t, TileEdge{7, 1}, TilePosition{2, 1}, TilePosition{3, 0})
}

func checkAdjacentTilesToEdge(t *testing.T, edge TileEdge, t1, t2 TilePosition) {
	tiles := AdjacentTilesToEdge(edge)
	if tiles[0] != t1 || tiles[1] != t2 {
		t.Errorf("for edge %v expected %v %v but was %v %v", edge, t1, t2, tiles[0], tiles[1])
	}
}
