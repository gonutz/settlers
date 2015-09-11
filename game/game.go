package game

// Tiles: 5 resources: brick, lumber, wool, grain, ore
// or water: nothing, 5 2:1 havens, 3:1 haven
// or desert.
// Each Player has: 5 settlements, 4 cities, 15 streets.
// There is one robber.
// Player has cards: assets and special
// Special Cards: knight, victory point, road building, monopoly, 2 resources
// Special achievements: largest army, longest road

// TODO store rand seed here in Game? and use fixed rand function?
type Game struct {
	Tiles            [37]Tile
	Players          [4]Player
	PlayerCount      int
	CurrentPlayer    int
	Robber           Robber
	DevelopmentCards [25]DevelopmentCard
	CardsDealt       int
	// LongestRoad and LargestArmy give the index of the player currently
	// holding the achievement or -1 if it is not yet accomplished by anybody.
	LongestRoad int
	LargestArmy int
}

type Tile struct {
	Position TilePosition
	Terrain  Terrain
	Number   int
	Harbor   Harbor
}

// TilePosition's coordinates will always add up to an odd number. The top-most
// horizontal row has y=0 and the left-most, only half visible, tile has x=-1 so
// the first full visible tile in that row is x=1.
// The second row has y=1 and starts with x=0.
// This means that x increases in steps of two when going to the next tile to
// the right.
//     ____  ____
//    /    \/    \
//    | 10 || 30 |
//  __\____/\____/__
// /    \/    \/    \
// | 01 || 21 || 41 |
// \____/\____/\____/
//    /    \/    \
//    | 12 || 32 |
//    \____/\____/
//
type TilePosition Point

type Point struct{ X, Y int }

type Terrain int

const (
	Hills Terrain = iota
	Pasture
	Mountains
	Field
	Forest
	Desert
	Water
)

// TODO give harbor an edge or two corners for making clear where it is
type Harbor int

const (
	NoHarbor Harbor = iota
	WoolHarbor
	LumberHarbor
	BrickHarbor
	OreHarbor
	GrainHarbor
	ThreeToOneHarbor
)

type Player struct {
	Color       Color
	Roads       [15]Road
	Settlements [5]Settlement
	Cities      [4]City
	Resources   [ResourceCount]int // TODO or just five int fields?
}

type Color int

const (
	Red Color = iota
	White
	Blue
	Orange
)

type Settlement struct{ Position TileCorner }

// TileCorners are numbered in zgizag lines, the top-most line has y=0 and x
// is sequentiel, starting at 0 on the left-most tile (the one that is only half
// way visible horizonally).
//
// 00  20  40
//  \  /\  /\
//   \/  \/  \/
//   10  30  50
//   ||  ||  ||
//   ||  ||  ||
//   11  31  51
//   /\  /\  /\
//  /  \/  \/
// 01  21  41
//
type TileCorner Point

func (c TileCorner) IsValid() bool { return c.X != 0 || c.Y != 0 }

// TODO comment this one and show the drawing
type TileEdge Point

func (e TileEdge) IsValid() bool { return e.X != 0 || e.Y != 0 }

type City struct{ Position TileCorner }

type Road struct{ Position TileEdge }

type Resource int

const (
	Brick Resource = iota
	Wool
	Ore
	Grain
	Lumber
	Nothing // NOTE this has to come last
)
const ResourceCount = 5

type Robber struct {
	Position int
}

type DevelopmentCard struct {
	Kind DevelopmentCardKind
	// Owner is the player index for who owns this card
	// NOTE this is 0 by default, even when the card was not dealt yet
	Owner int
}

type DevelopmentCardKind int

const (
	Knight DevelopmentCardKind = iota // Knight can move the robber
	VictoryPoint
	Monopoly // Monopoly gets all other players' resources of a chosen type
	BuildTwoRoads
	TakeTwoResources
)

func New(colors []Color, randomNumberGenerator func() int) *Game {
	var game Game

	rand := func(tiles *[]Tile) Tile {
		tile := (*tiles)[0]
		*tiles = (*tiles)[1:]
		return tile
	}

	shuffle := func(tiles []Tile) {
		count := len(tiles)
		for i := 0; i < count-1; i++ {
			j := i + randomNumberGenerator()%(count-i)
			tiles[i], tiles[j] = tiles[j], tiles[i]
		}
	}

	harbors := &[]Tile{
		{Terrain: Water, Harbor: LumberHarbor},
		{Terrain: Water, Harbor: WoolHarbor},
		{Terrain: Water, Harbor: BrickHarbor},
		{Terrain: Water, Harbor: OreHarbor},
		{Terrain: Water, Harbor: GrainHarbor},
		{Terrain: Water, Harbor: ThreeToOneHarbor},
		{Terrain: Water, Harbor: ThreeToOneHarbor},
		{Terrain: Water, Harbor: ThreeToOneHarbor},
		{Terrain: Water, Harbor: ThreeToOneHarbor},
	}
	shuffle(*harbors)
	game.Tiles[0] = rand(harbors)
	game.Tiles[2] = rand(harbors)
	game.Tiles[8] = rand(harbors)
	game.Tiles[9] = rand(harbors)
	game.Tiles[21] = rand(harbors)
	game.Tiles[22] = rand(harbors)
	game.Tiles[32] = rand(harbors)
	game.Tiles[33] = rand(harbors)
	game.Tiles[35] = rand(harbors)

	water := Tile{Terrain: Water}
	game.Tiles[1] = water
	game.Tiles[3] = water
	game.Tiles[4] = water
	game.Tiles[14] = water
	game.Tiles[15] = water
	game.Tiles[27] = water
	game.Tiles[28] = water
	game.Tiles[34] = water
	game.Tiles[36] = water

	terrains := &[]Tile{
		{Terrain: Desert},
		{Terrain: Hills},
		{Terrain: Hills},
		{Terrain: Hills},
		{Terrain: Mountains},
		{Terrain: Mountains},
		{Terrain: Mountains},
		{Terrain: Pasture},
		{Terrain: Pasture},
		{Terrain: Pasture},
		{Terrain: Pasture},
		{Terrain: Forest},
		{Terrain: Forest},
		{Terrain: Forest},
		{Terrain: Forest},
		{Terrain: Field},
		{Terrain: Field},
		{Terrain: Field},
		{Terrain: Field},
	}
	shuffle(*terrains)
	game.Tiles[5] = rand(terrains)
	game.Tiles[6] = rand(terrains)
	game.Tiles[7] = rand(terrains)
	game.Tiles[10] = rand(terrains)
	game.Tiles[11] = rand(terrains)
	game.Tiles[12] = rand(terrains)
	game.Tiles[13] = rand(terrains)
	game.Tiles[16] = rand(terrains)
	game.Tiles[17] = rand(terrains)
	game.Tiles[18] = rand(terrains)
	game.Tiles[19] = rand(terrains)
	game.Tiles[20] = rand(terrains)
	game.Tiles[23] = rand(terrains)
	game.Tiles[24] = rand(terrains)
	game.Tiles[25] = rand(terrains)
	game.Tiles[26] = rand(terrains)
	game.Tiles[29] = rand(terrains)
	game.Tiles[30] = rand(terrains)
	game.Tiles[31] = rand(terrains)
	// find the desert and place the robber on it
	for i, tile := range game.Tiles {
		if tile.Terrain == Desert {
			game.Robber.Position = i
		}
	}
	// set tile positions
	for i, tile := range game.Tiles {
		tile.Position = tilePositions[i]
	}

	numbers := []int{5, 2, 6, 3, 8, 10, 9, 12, 11, 4, 8, 10, 9, 4, 5, 6, 3, 11}
	tileOrder := []int{16, 23, 29, 30, 31, 26, 20, 13, 7, 6, 5, 10, 17, 24, 25, 19, 12, 11, 18}
	for _, t := range tileOrder {
		if game.Tiles[t].Terrain != Desert {
			game.Tiles[t].Number = numbers[0]
			numbers = numbers[1:]
		}
	}

	game.PlayerCount = len(colors)
	for i := range colors {
		game.Players[i].Color = colors[i]
	}

	cards := []DevelopmentCard{
		{Kind: VictoryPoint},
		{Kind: VictoryPoint},
		{Kind: VictoryPoint},
		{Kind: VictoryPoint},
		{Kind: VictoryPoint},
		{Kind: Monopoly},
		{Kind: Monopoly},
		{Kind: BuildTwoRoads},
		{Kind: BuildTwoRoads},
		{Kind: TakeTwoResources},
		{Kind: TakeTwoResources},
	}
	var knights [14]DevelopmentCard
	for i := range knights {
		knights[i].Kind = Knight
	}
	cards = append(cards, knights[:]...)

	game.LongestRoad = -1
	game.LargestArmy = -1

	return &game
}

// TODO create a way to get to the information who received what resources
// this is necessary to animate the cards flying to the players
func (g *Game) DealResources(dice int) {
	for _, tile := range g.Tiles {
		if tile.Number == dice {
			corners := AdjacentCornersToTile(tile.Position)
			for _, corner := range corners {
				for playerIndex, p := range g.GetPlayers() {
					for _, s := range p.GetBuiltSettlements() {
						if s.Position == corner {
							g.Players[playerIndex].Resources[tile.Resource()]++
						}
					}
					for _, c := range p.GetBuiltCities() {
						if c.Position == corner {
							g.Players[playerIndex].Resources[tile.Resource()] += 2
						}
					}
				}
			}
		}
	}
	// TODO re-create this functionality
	/*
		for p := 0; p < g.PlayerCount; p++ {
			player := &g.Players[p]

			var buildings []building
			for _, settlement := range player.Settlements {
				buildings = append(buildings, settlement)
			}
			for _, city := range player.Cities {
				buildings = append(buildings, city)
			}

			for _, b := range buildings {
				if b.isSet() {
					for _, bordering := range b.borderingTiles() {
						tile := g.Tiles[bordering]
						if g.Robber.Position != bordering && tile.Number == dice {
							// the tile's resource cannot be Nothing because it as a
							// number that was rolled by the dice (!= 0)
							player.Resources[tile.Resource()] += b.resourceCount()
						}
					}
				}
			}
		}
	*/
}

func (g *Game) GetPlayers() []Player {
	return g.Players[:g.PlayerCount]
}

func (p Player) GetBuiltSettlements() []Settlement {
	last := 0
	for last < len(p.Settlements) && p.Settlements[last].isSet() {
		last++
	}
	return p.Settlements[:last]
}

func (p Player) GetBuiltCities() []City {
	last := 0
	for last < len(p.Cities) && p.Cities[last].isSet() {
		last++
	}
	return p.Cities[:last]
}

func (p Player) GetBuiltRoads() []Road {
	last := 0
	for last < len(p.Roads) && p.Roads[last].isSet() {
		last++
	}
	return p.Roads[:last]
}

func (g *Game) CanPlayerBuildSettlement() bool {
	// the player can build as long as the last settlement has not been placed
	return !g.Players[g.CurrentPlayer].Settlements[4].isSet()
}

func (g *Game) CanPlayerBuildCity() bool {
	return !g.Players[g.CurrentPlayer].Cities[3].isSet()
}

func (g *Game) CanPlayerBuildRoad() bool {
	return !g.Players[g.CurrentPlayer].Roads[14].isSet()
}

func (g *Game) BuildSettlement(p Point) {
	player := &g.Players[g.CurrentPlayer]
	for i := range player.Settlements {
		if !player.Settlements[i].isSet() {
			player.Settlements[i].Position = TileCorner(p)
			return
		}
	}
}

func (g *Game) BuildCity(p Point) {
	player := &g.Players[g.CurrentPlayer]
	for i := range player.Cities {
		if !player.Cities[i].isSet() {
			player.Cities[i].Position = TileCorner(p)
			return
		}
	}
}

func (g *Game) BuildRoad(p Point) {
	player := &g.Players[g.CurrentPlayer]
	for i := range player.Roads {
		if !player.Roads[i].isSet() {
			player.Roads[i].Position = TileEdge(p)
			return
		}
	}
}

func (g *Game) GetTiles() []PositionTile {
	// TODO buffer these
	tiles := make([]PositionTile, len(tilePositions))
	for i := range tilePositions /*TODO g.Tiles*/ {
		tiles[i].Tile = g.Tiles[i]
		tiles[i].Position = tilePositions[i]
		tiles[i].HasRobber = g.Robber.Position == i
	}
	return tiles
}

func (g *Game) GetTilePosition(tile int) (x, y int) {
	return tilePositions[tile].X, tilePositions[tile].Y
}

type PositionTile struct {
	Tile
	Position  TilePosition
	HasRobber bool
}

var tilePositions = []TilePosition{
	{3, 0}, {5, 0}, {7, 0}, {9, 0},
	{2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1},
	{1, 2}, {3, 2}, {5, 2}, {7, 2}, {9, 2}, {11, 2},
	{0, 3}, {2, 3}, {4, 3}, {6, 3}, {8, 3}, {10, 3}, {12, 3},
	{1, 4}, {3, 4}, {5, 4}, {7, 4}, {9, 4}, {11, 4},
	{2, 5}, {4, 5}, {6, 5}, {8, 5}, {10, 5},
	{3, 6}, {5, 6}, {7, 6}, {9, 6},
}

func (g *Game) Size() (w, h int) { return 13, 7 }

type building interface {
	isSet() bool
	resourceCount() int
	borderingTiles() []int
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % g.PlayerCount
}

// IsSet returns true if the settlement is currently placed on the game field.
func (s Settlement) isSet() bool {
	return s.Position.IsValid()
}

func (Settlement) resourceCount() int { return 1 }

// IsSet returns true if the settlement is currently placed on the game field.
func (c City) isSet() bool {
	return c.Position.IsValid()
}

func (r Road) isSet() bool {
	return r.Position.IsValid()
}

func (City) resourceCount() int { return 2 }

func (t Tile) Resource() Resource {
	switch t.Terrain {
	case Hills:
		return Brick
	case Pasture:
		return Wool
	case Mountains:
		return Ore
	case Field:
		return Grain
	case Forest:
		return Lumber
	default:
		return Nothing
	}
}

func AdjacentTileToCorner(c TileCorner) [3]TilePosition {
	if (c.X+c.Y)%2 == 0 {
		// two on top, one below
		return [3]TilePosition{
			{c.X - 2, c.Y - 1},
			{c.X - 1, c.Y},
			{c.X - 0, c.Y - 1},
		}
	}
	// one on top, two below
	return [3]TilePosition{
		{c.X - 2, c.Y},
		{c.X - 1, c.Y - 1},
		{c.X - 0, c.Y},
	}
}

func AdjacentCornersToTile(p TilePosition) [6]TileCorner {
	return [6]TileCorner{
		{p.X + 0, p.Y},
		{p.X + 0, p.Y + 1},
		{p.X + 1, p.Y},
		{p.X + 1, p.Y + 1},
		{p.X + 2, p.Y},
		{p.X + 2, p.Y + 1},
	}
}

func AdjacentTilesToEdge(p TileEdge) [2]TilePosition {
	if p.X%2 == 0 {
		// vertical edge
		return [2]TilePosition{
			{p.X/2 - 2, p.Y},
			{p.X / 2, p.Y},
		}
	}
	// now it's either rising or falling edge
	if (p.Y%2 == 0 && (p.X-1)%4 == 0) || (p.Y%2 == 1 && (p.X-3)%4 == 0) {
		// falling edge
		return [2]TilePosition{
			{p.X / 4 * 2, p.Y},
			{p.X/4*2 + 1, p.Y - 1},
		}
	}
	// rising edge
	return [2]TilePosition{
		{p.X/4*2 - 1, p.Y - 1},
		{p.X / 4 * 2, p.Y},
	}
}

func AdjacentEdgesToTile(p TilePosition) [6]TileEdge {
	return [6]TileEdge{
		{p.X * 2, p.Y},       // left
		{p.X*2 + 1, p.Y},     // top-left
		{p.X*2 + 3, p.Y},     // top-right
		{p.X*2 + 4, p.Y},     // right
		{p.X*2 + 3, p.Y + 1}, // bottom-right
		{p.X*2 + 1, p.Y + 1}, // bottom-left
	}
}
