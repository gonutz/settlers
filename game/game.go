package game

// Tiles: 5 resources: brick, lumber, wool, grain, ore
// or water: nothing, 5 2:1 harbors, 3:1 harbors
// or desert.
// Each Player has: 5 settlements, 4 cities, 15 streets.
// There is one robber.
// Player has cards: assets and special
// Special Cards: knight, victory point, road building, monopoly, 2 resources
// Special achievements: largest army, longest road

// TODO store rand seed here in Game? and use fixed rand function?
type Game struct {
	State            State
	Tiles            [37]Tile
	Players          [4]Player
	PlayerCount      int
	CurrentPlayer    int
	Robber           Robber
	DevelopmentCards [25]DevelopmentCard
	CardsDealt       int
	Dice             [2]int
	// seed is for random number generation
	rand *randomNumberGenerator
}

type State int

const (
	NotStarted State = iota
	BuildingFirstSettlement
	BuildingFirstRoad
	BuildingSecondSettlement
	BuildingSecondRoad
	ChoosingNextAction
	BuildingNewRoad
	BuildingNewSettlement
	BuildingNewCity
	RollingDice
)

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
type TilePosition point

type point struct{ X, Y int }

func (p point) IsValid() bool { return p.X != 0 || p.Y != 0 }

// TileCorners are numbered in zigzag lines, the top-most line has y=0 and x
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
type TileCorner point

// TileEdges are TODO comment and finish drawing
//
// x     x     x
//  \    /\    /\
//   \  /  \  /  \
//    \/    \/    \/
//    x     x     x
//    |     |     |
//    |     |     |
//    |     |     |
//    x     x     x
//    /\    /\    /\
//   /  \  /  \  /
//  /    \/    \/
// x     x     x
type TileEdge point

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

type Harbor struct {
	Kind      HarborKind
	Direction Direction
}

type HarborKind int

const (
	NoHarbor HarborKind = iota
	WoolHarbor
	LumberHarbor
	BrickHarbor
	OreHarbor
	GrainHarbor
	ThreeToOneHarbor
)

type Direction int

const (
	Right Direction = iota
	TopRight
	TopLeft
	Left
	BottomLeft
	BottomRight
)

// TODO cards
type Player struct {
	Color          Color
	Roads          [15]Road
	Settlements    [5]Settlement
	Cities         [4]City
	Resources      [ResourceCount]int
	HasLongestRoad bool
	HasLargestArmy bool
}

type Color int

const (
	Red Color = iota
	White
	Blue
	Orange
)

type Settlement struct{ Position TileCorner }

type City struct{ Position TileCorner }

type Road struct{ Position TileEdge }

type Resource int

const (
	Lumber Resource = iota
	Brick
	Wool
	Ore
	Grain
	Nothing // NOTE this has to come last
)
const ResourceCount = 5

type Robber struct {
	Position TilePosition
}

type DevelopmentCard struct {
	Kind DevelopmentCardKind
}

type DevelopmentCardKind int

const (
	Knight DevelopmentCardKind = iota // Knight can move the robber
	VictoryPoint
	Monopoly // Monopoly gets all other players' resources of a chosen type
	BuildTwoRoads
	TakeTwoResources
)

func New(colors []Color, randomSeed int) *Game {
	var game Game
	game.rand = newRNG(randomSeed)

	rand := func(tiles *[]Tile) Tile {
		tile := (*tiles)[0]
		*tiles = (*tiles)[1:]
		return tile
	}

	shuffle := func(tiles []Tile) {
		count := len(tiles)
		for i := 0; i < count-1; i++ {
			j := i + game.rand.next()%(count-i)
			tiles[i], tiles[j] = tiles[j], tiles[i]
		}
	}

	harbors := &[]Tile{
		{Terrain: Water, Harbor: Harbor{Kind: LumberHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: WoolHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: BrickHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: OreHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: GrainHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: ThreeToOneHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: ThreeToOneHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: ThreeToOneHarbor}},
		{Terrain: Water, Harbor: Harbor{Kind: ThreeToOneHarbor}},
	}
	shuffle(*harbors)
	game.Tiles[0] = rand(harbors)
	game.Tiles[0].Harbor.Direction = BottomRight
	game.Tiles[2] = rand(harbors)
	game.Tiles[2].Harbor.Direction = BottomLeft
	game.Tiles[8] = rand(harbors)
	game.Tiles[8].Harbor.Direction = BottomLeft
	game.Tiles[9] = rand(harbors)
	game.Tiles[9].Harbor.Direction = Right
	game.Tiles[21] = rand(harbors)
	game.Tiles[21].Harbor.Direction = Left
	game.Tiles[22] = rand(harbors)
	game.Tiles[22].Harbor.Direction = Right
	game.Tiles[32] = rand(harbors)
	game.Tiles[32].Harbor.Direction = TopLeft
	game.Tiles[33] = rand(harbors)
	game.Tiles[33].Harbor.Direction = TopRight
	game.Tiles[35] = rand(harbors)
	game.Tiles[35].Harbor.Direction = TopLeft

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

	// set tile positions
	var tilePositions = []TilePosition{
		{3, 0}, {5, 0}, {7, 0}, {9, 0},
		{2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1},
		{1, 2}, {3, 2}, {5, 2}, {7, 2}, {9, 2}, {11, 2},
		{0, 3}, {2, 3}, {4, 3}, {6, 3}, {8, 3}, {10, 3}, {12, 3},
		{1, 4}, {3, 4}, {5, 4}, {7, 4}, {9, 4}, {11, 4},
		{2, 5}, {4, 5}, {6, 5}, {8, 5}, {10, 5},
		{3, 6}, {5, 6}, {7, 6}, {9, 6},
	}
	for i := range game.Tiles {
		game.Tiles[i].Position = tilePositions[i]
	}
	// find the desert and place the robber on it
	for _, tile := range game.Tiles {
		if tile.Terrain == Desert {
			game.Robber.Position = tile.Position
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

	game.randomizePlayerOrder()
	game.State = NotStarted

	return &game
}

func (g *Game) Start() {
	if g.State == NotStarted {
		g.State = BuildingFirstSettlement
	}
}

func (g *Game) randomizePlayerOrder() {
	players := g.GetPlayers()
	var order []Player
	for len(players) > 0 {
		index := g.rand.next() % len(players)
		order = append(order, players[index])
		players = append(players[:index], players[index+1:]...)
	}
	for i, p := range order {
		g.Players[i] = p
	}
}

// TODO create a way to get to the information who received what resources
// this is necessary to e.g. animate the resources flying to the players
func (g *Game) DealResources(dice int) {
	for _, tile := range g.Tiles {
		if tile.Number == dice && g.Robber.Position != tile.Position {
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

func (g *Game) RemainingSettlements() int {
	p := g.GetCurrentPlayer()
	return len(p.Settlements) - len(p.GetBuiltSettlements())
}

func (g *Game) RemainingCities() int {
	p := g.GetCurrentPlayer()
	return len(p.Cities) - len(p.GetBuiltCities())
}

func (g *Game) CanPlayerBuildCity() bool {
	return !g.Players[g.CurrentPlayer].Cities[3].isSet()
}

// BuiltCity assumes that you checked CanBuildCityAt right before calling it.
func (g *Game) BuildCity(c TileCorner) {
	player := g.currentPlayerPointer()

	// un-build the settlement that is to be replaced with a city
	settlementIndex := -1
	for i, s := range player.GetBuiltSettlements() {
		if s.Position == c {
			settlementIndex = i
			break
		}
	}
	lastIndex := len(player.Settlements) - 1
	for i := settlementIndex; i < lastIndex; i++ {
		player.Settlements[i].Position = player.Settlements[i+1].Position
	}
	player.Settlements[lastIndex].Position = TileCorner{}

	// now build the city
	for i := range player.Cities {
		if !player.Cities[i].isSet() {
			player.Cities[i].Position = TileCorner(c)
			break
		}
	}

	g.State = ChoosingNextAction
}

// BuildRoad assumes that you check CanBuildRoadAt first.
func (g *Game) BuildRoad(e TileEdge) {
	player := g.currentPlayerPointer()
	for i := range player.Roads {
		if !player.Roads[i].isSet() {
			player.Roads[i].Position = TileEdge(e)
			break
		}
	}
	if g.State == BuildingFirstRoad {
		g.State = BuildingFirstSettlement
		g.CurrentPlayer++
		if g.CurrentPlayer == g.PlayerCount {
			g.State = BuildingSecondSettlement
			g.CurrentPlayer--
		}
	} else if g.State == BuildingSecondRoad {
		g.State = BuildingSecondSettlement
		g.CurrentPlayer--
		if g.CurrentPlayer == -1 {
			g.State = RollingDice
			g.randomizePlayerOrder()
			g.CurrentPlayer = 0
		}
	} else if g.State == BuildingNewRoad {
		g.State = ChoosingNextAction
	}
}

func (g *Game) Size() (w, h int) { return 13, 7 }

type building interface {
	isSet() bool
	resourceCount() int
	borderingTiles() []int
}

func (g *Game) NextTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % g.PlayerCount
	g.State = RollingDice
}

// IsSet returns true if the settlement is currently placed on the game field.
func (s Settlement) isSet() bool {
	return point(s.Position).IsValid()
}

func (Settlement) resourceCount() int { return 1 }

// IsSet returns true if the settlement is currently placed on the game field.
func (c City) isSet() bool {
	return point(c.Position).IsValid()
}

func (r Road) isSet() bool {
	return point(r.Position).IsValid()
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

func AdjacentTilesToCorner(c TileCorner) [3]TilePosition {
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

	if (p.Y%2 == 0 && (p.X-3)%4 == 0) ||
		(p.Y%2 == 1 && (p.X-1)%4 == 0) {
		return [2]TilePosition{
			{p.X/2 - 1, p.Y - 1},
			{p.X / 2, p.Y},
		}
	}

	return [2]TilePosition{
		{p.X/2 - 1, p.Y},
		{p.X / 2, p.Y - 1},
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

func AdjacentCornersToCorner(p TileCorner) [3]TileCorner {
	if (p.X+p.Y)%2 == 0 {
		return [3]TileCorner{
			{p.X - 1, p.Y},
			{p.X, p.Y - 1},
			{p.X + 1, p.Y},
		}
	}
	return [3]TileCorner{
		{p.X - 1, p.Y},
		{p.X, p.Y + 1},
		{p.X + 1, p.Y},
	}
}

func AdjacentEdgesToCorner(p TileCorner) [3]TileEdge {
	if (p.X+p.Y)%2 == 0 {
		return [3]TileEdge{
			{p.X*2 - 1, p.Y},
			{p.X * 2, p.Y - 1},
			{p.X*2 + 1, p.Y},
		}
	}
	return [3]TileEdge{
		{p.X*2 - 1, p.Y},
		{p.X * 2, p.Y},
		{p.X*2 + 1, p.Y},
	}
}

func AdjacentCornersToEdge(p TileEdge) [2]TileCorner {
	if p.X%2 == 0 {
		return [2]TileCorner{
			{p.X / 2, p.Y},
			{p.X / 2, p.Y + 1},
		}
	}
	return [2]TileCorner{
		{p.X / 2, p.Y},
		{p.X/2 + 1, p.Y},
	}
}

func AdjacentEdgesToEdge(p TileEdge) [4]TileEdge {
	if p.X%2 == 0 {
		return [4]TileEdge{
			{p.X - 1, p.Y},
			{p.X - 1, p.Y + 1},
			{p.X + 1, p.Y},
			{p.X + 1, p.Y + 1},
		}
	}
	// the same as this => falling edge
	//if (p.Y%2 == 0 && (p.X-1)%4 == 0) || (p.Y%2 == 1 && (p.X-3)%4 == 0) {
	if (p.X-1-2*(p.Y&1))%4 == 0 {
		return [4]TileEdge{
			{p.X - 2, p.Y},
			{p.X - 1, p.Y - 1},
			{p.X + 1, p.Y},
			{p.X + 2, p.Y},
		}
	}
	return [4]TileEdge{
		{p.X - 2, p.Y},
		{p.X - 1, p.Y},
		{p.X + 1, p.Y - 1},
		{p.X + 2, p.Y},
	}
}

func AdjacentTilesToTile(p TilePosition) [6]TilePosition {
	return [6]TilePosition{
		{p.X - 2, p.Y},
		{p.X - 1, p.Y - 1},
		{p.X - 1, p.Y + 1},
		{p.X + 1, p.Y - 1},
		{p.X + 1, p.Y + 1},
		{p.X + 2, p.Y},
	}
}

func (g *Game) GetTileAt(p TilePosition) (Tile, bool) {
	for _, tile := range g.Tiles {
		if tile.Position == p {
			return tile, true
		}
	}
	return Tile{}, false
}

// BuildSettlement assumes that CanBuildSettlement returned true right before
// you call this
func (g *Game) BuildSettlement(c TileCorner) {
	player := g.currentPlayerPointer()
	for i := range player.Settlements {
		if !player.Settlements[i].isSet() {
			player.Settlements[i].Position = c
			break
		}
	}

	if g.State == BuildingFirstSettlement {
		g.State = BuildingFirstRoad
	} else if g.State == BuildingSecondSettlement {
		// TODO make dealt resources available to the UI
		player := g.currentPlayerPointer()
		// deal resources for this settlement
		tilePositions := AdjacentTilesToCorner(c)
		for _, tilePosition := range tilePositions {
			tile, valid := g.GetTileAt(tilePosition)
			if valid {
				if resource := tile.Resource(); resource != Nothing {
					player.Resources[resource]++
				}
			}
		}
		g.State = BuildingSecondRoad
	} else if g.State == BuildingNewSettlement {
		g.State = ChoosingNextAction
	}
}

func (g *Game) CanBuildSettlementAt(c TileCorner) bool {
	if g.RemainingSettlements() == 0 {
		return false
	}

	if !g.canBuildBuildingAt(c) {
		return false
	}

	// if we are in the first phase, this is all we need to know to build
	if g.State == BuildingFirstSettlement || g.State == BuildingSecondSettlement {
		return true
	}

	return g.hasRoadToCorner(c)
}

func (g *Game) canBuildBuildingAt(c TileCorner) bool {
	// check that at least one adjacent tile is not water, can't build in water!
	tilePositions := AdjacentTilesToCorner(c)
	atLeastOneTileIsLand := false
	for _, pos := range tilePositions {
		if tile, ok := g.GetTileAt(pos); ok && tile.Terrain != Water {
			atLeastOneTileIsLand = true
			break
		}
	}
	if !atLeastOneTileIsLand {
		return false
	}

	// check that all adjacent corners have no building on them
	cornerPositions := AdjacentCornersToCorner(c)
	for _, player := range g.GetPlayers() {
		if player.HasBuildingOnCorner(c) {
			// if there is a building on this corner
			return false
		}
		// or if there is a building only one corner away
		for _, corner := range cornerPositions {
			if player.HasBuildingOnCorner(corner) {
				return false
			}
		}
	}

	return true
}

func (g *Game) CanBuildCityAt(c TileCorner) bool {
	if g.RemainingCities() == 0 {
		return false
	}

	player := g.currentPlayerPointer()
	for _, s := range player.GetBuiltSettlements() {
		if s.Position == c {
			return true
		}
	}

	return false
}

func (g *Game) hasRoadToCorner(c TileCorner) bool {
	edgePositions := AdjacentEdgesToCorner(c)
	for _, edge := range edgePositions {
		if g.GetCurrentPlayer().HasRoadOnEdge(edge) {
			return true
		}
	}
	return false
}

func (p Player) HasBuildingOnCorner(corner TileCorner) bool {
	for _, s := range p.GetBuiltSettlements() {
		if s.Position == corner {
			return true
		}
	}
	for _, c := range p.GetBuiltCities() {
		if c.Position == corner {
			return true
		}
	}
	return false
}

func (g *Game) GetCurrentPlayer() Player {
	return g.Players[g.CurrentPlayer]
}

// TODO this function is not necessary anymore, now there is State
func (g *Game) isFirstGamePhase() bool {
	// make sure there are no cities and each player has AT MOST 2 roads and 2
	// settlements
	for _, p := range g.GetPlayers() {
		if len(p.GetBuiltCities()) > 0 ||
			len(p.GetBuiltSettlements()) > 2 ||
			len(p.GetBuiltRoads()) > 2 {
			return false
		}
	}
	// if now any player has less than two roads, that player has not placed her
	// first two settlements and roads and thus this is still the first phase
	for _, p := range g.GetPlayers() {
		if len(p.GetBuiltRoads()) < 2 {
			return true
		}
	}
	return false
}

func (p Player) HasRoadOnEdge(edge TileEdge) bool {
	for _, e := range p.GetBuiltRoads() {
		if e.Position == edge {
			return true
		}
	}
	return false
}

func (g *Game) CanBuildRoadAt(edge TileEdge) bool {
	if g.RemainingRoads() == 0 {
		return false
	}

	// can't build here if there already is a road on this edge
	for _, p := range g.GetPlayers() {
		if p.HasRoadOnEdge(edge) {
			return false
		}
	}

	// can only build if at least one adjacent tile is land
	tiles := AdjacentTilesToEdge(edge)
	if !(g.isLand(tiles[0]) || g.isLand(tiles[1])) {
		return false
	}

	// find a building of the current player next to the road
	hasBuildingNextToIt := false
	var buildingCorner TileCorner
	for _, corner := range AdjacentCornersToEdge(edge) {
		if g.GetCurrentPlayer().HasBuildingOnCorner(corner) {
			hasBuildingNextToIt = true
			buildingCorner = corner
			break
		}
	}

	if g.State == BuildingFirstRoad || g.State == BuildingSecondRoad {
		if !hasBuildingNextToIt {
			return false
		}
		// In the first game phase, you can only build a road next to your
		// settlement. When it is time for the second road, there are two
		// settlements already and you have to build the road next to the new
		// one. The new one is that without any adjacent roads so if this
		// settlement already has a road next to it we return false, the player
		// has to build adjacent to the other settlement.
		if g.State == BuildingSecondRoad {
			for _, edge := range AdjacentEdgesToCorner(buildingCorner) {
				if g.GetCurrentPlayer().HasRoadOnEdge(edge) {
					return false
				}
			}
		}
		return true
	} else {
		if hasBuildingNextToIt {
			return true
		}
		// can also build if there is a road adjacent to the edge
		for _, e := range AdjacentEdgesToEdge(edge) {
			if g.GetCurrentPlayer().HasRoadOnEdge(e) {
				return true
			}
		}
		return false
	}
}

func (g *Game) RemainingRoads() int {
	p := g.GetCurrentPlayer()
	return len(p.Roads) - len(p.GetBuiltRoads())
}

func (g *Game) isLand(p TilePosition) bool {
	tile, valid := g.GetTileAt(p)
	return valid && tile.Terrain != Water
}

func (g *Game) CanBuyRoad() bool {
	player := g.currentPlayerPointer()
	hasOneLeft := !player.Roads[len(player.Roads)-1].isSet()

	return hasOneLeft &&
		player.Resources[Lumber] >= 1 &&
		player.Resources[Brick] >= 1
}

func (g *Game) BuyRoad() {
	player := g.currentPlayerPointer()
	player.Resources[Lumber]--
	player.Resources[Brick]--
	g.State = BuildingNewRoad
}

func (g *Game) CanBuySettlement() bool {
	player := g.currentPlayerPointer()
	hasOneLeft := !player.Settlements[len(player.Settlements)-1].isSet()

	return hasOneLeft &&
		g.canBuildAtAnyRoad(player) &&
		player.Resources[Lumber] >= 1 &&
		player.Resources[Brick] >= 1 &&
		player.Resources[Wool] >= 1 &&
		player.Resources[Grain] >= 1
}

func (g *Game) canBuildAtAnyRoad(player *Player) bool {
	for _, road := range player.GetBuiltRoads() {
		corners := AdjacentCornersToEdge(road.Position)
		for _, corner := range corners {
			if g.CanBuildSettlementAt(corner) {
				return true
			}
		}
	}
	return false
}

func (g *Game) BuySettlement() {
	player := g.currentPlayerPointer()
	player.Resources[Lumber]--
	player.Resources[Brick]--
	player.Resources[Wool]--
	player.Resources[Grain]--
	g.State = BuildingNewSettlement
}

func (g *Game) CanBuyCity() bool {
	player := g.currentPlayerPointer()
	hasOneLeft := !player.Cities[len(player.Cities)-1].isSet()

	return hasOneLeft &&
		player.Resources[Grain] >= 2 &&
		player.Resources[Ore] >= 3 &&
		len(player.GetBuiltSettlements()) > 0
}

func (g *Game) BuyCity() {
	player := g.currentPlayerPointer()
	player.Resources[Grain] -= 2
	player.Resources[Ore] -= 3
	g.State = BuildingNewCity
}

func (g *Game) CanBuyDevelopmentCard() bool {
	player := g.currentPlayerPointer()
	return g.CardsDealt < len(g.DevelopmentCards) &&
		player.Resources[Wool] >= 1 &&
		player.Resources[Grain] >= 1 &&
		player.Resources[Ore] >= 1
}

func (g *Game) BuyDevelopmentCard() {
	player := g.currentPlayerPointer()
	player.Resources[Wool]--
	player.Resources[Grain]--
	player.Resources[Ore]--

	// TODO deal card
	_ = g.DevelopmentCards[g.CardsDealt]
	g.CardsDealt++

	g.State = ChoosingNextAction
}

func (g *Game) currentPlayerPointer() *Player {
	return &g.Players[g.CurrentPlayer]
}

func (g *Game) RollTheDice() {
	g.Dice[0] = 1 + g.rand.next()%6
	g.Dice[1] = 1 + g.rand.next()%6
	g.DealResources(g.Dice[0] + g.Dice[1])
	g.State = ChoosingNextAction
}
