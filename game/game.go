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
	// Tiles have an implicit position, the order is from top to bottom, left
	// to right, look at this ASCII art:
	//     ____  ____
	//    /    \/    \
	//    |  0 ||  1 |
	//  __\____/\____/__
	// /    \/    \/    \
	// |  2 ||  3 ||  4 |
	// \____/\____/\____/
	//
	Tiles            [37]Tile
	Players          [4]Player
	PlayerCount      int
	NextPlayer       int
	Robber           Robber
	DevelopmentCards [25]DevelopmentCard
	CardsDealt       int
	// LongestRoad and LargestArmy give the index of the player currently
	// holding the achievement or -1 if it is not yet accomplished by anybody.
	LongestRoad int
	LargestArmy int
}

type Tile struct {
	Terrain Terrain
	Number  int
	Harbor  Harbor
}

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

type Settlement struct{ Position CornerPosition }

// CornerPosition is identified by the 3 adjecent tiles, these are the indices
// of those tiles into the Tiles array of Game.
type CornerPosition [3]int

type City struct{ Position CornerPosition }

type Road struct{ Position EdgePosition }

// EdgePosition is identified by the 2 adjecent tiles, these are the indices
// of those tiles into the Tiles array of Game.
type EdgePosition [2]int

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

// TODO these are some thoughts

func New(colors []Color, randomNumberGenerator func() int) Game {
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
	for i, tile := range game.Tiles {
		if tile.Terrain == Desert {
			game.Robber.Position = i
		}
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

	return game
}

// TODO create a way to get to the information who received what resources
// this is necessary to animate the cards flying to the players
func (g *Game) DealResources(dice int) {

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

type PositionTile struct {
	Tile
	Position  TilePosition
	HasRobber bool
}

type TilePosition struct{ X, Y int }

var tilePositions = []TilePosition{
	{3, 0}, {5, 0}, {7, 0}, {9, 0},
	{2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1},
	{1, 2}, {3, 2}, {5, 2}, {7, 2}, {9, 2}, {11, 2},
	{0, 3}, {2, 3}, {4, 3}, {6, 3}, {8, 3}, {10, 3}, {12, 3},
	{1, 4}, {3, 4}, {5, 4}, {7, 4}, {9, 4}, {11, 4},
	{2, 5}, {4, 5}, {6, 5}, {8, 5}, {10, 5},
	{3, 6}, {5, 6}, {7, 6}, {9, 6},
}

type building interface {
	isSet() bool
	resourceCount() int
	borderingTiles() []int
}

func (g *Game) NextTurn() {
	g.NextPlayer = (g.NextPlayer + 1) % g.PlayerCount
}

func (p *CornerPosition) Sort() {
	if (*p)[2] < (*p)[0] {
		(*p)[0], (*p)[2] = (*p)[2], (*p)[0]
	}
	if (*p)[1] < (*p)[0] {
		(*p)[0], (*p)[1] = (*p)[1], (*p)[0]
	}
	if (*p)[2] < (*p)[1] {
		(*p)[1], (*p)[2] = (*p)[2], (*p)[1]
	}
}

func (p *CornerPosition) IsValid() bool {
	return !p.invalid()
}

func (p *CornerPosition) invalid() bool {
	return (*p)[0] == 0 && (*p)[1] == 0
}

func (p *EdgePosition) Sort() {
	if (*p)[1] < (*p)[0] {
		(*p)[0], (*p)[1] = (*p)[1], (*p)[0]
	}
}

func (p *EdgePosition) IsValid() bool {
	return !p.invalid()
}

func (p *EdgePosition) invalid() bool {
	return (*p)[0] == 0 && (*p)[1] == 0
}

// IsSet returns true if the settlement is currently placed on the game field.
func (s Settlement) isSet() bool {
	return s.Position.IsValid()
}

func (s Settlement) borderingTiles() []int {
	return s.Position[:]
}

func (Settlement) resourceCount() int { return 1 }

// IsSet returns true if the settlement is currently placed on the game field.
func (c City) isSet() bool {
	return c.Position.IsValid()
}

func (c City) borderingTiles() []int {
	return c.Position[:]
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
