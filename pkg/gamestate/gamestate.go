// Class gamestate helps parses and stores the gamestate of the demo parser, store/read to & from json
package gamestate

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
)

type SteamID32 uint32

// MARK: Alias Types
// Demotick in-game
type Tick int64

// A mask of button held down of the spectating player
type ButtonMask uint64

type VelocityOrPosition r3.Vector

type ViewAngle r2.Point

// MARK: Children Structs
// Player events, actions
type Player struct {
	Button    map[Tick]ButtonMask         `json:"button_pressed_at_ticks"`
	Velocity  map[Tick]VelocityOrPosition `json:"velocity_at_ticks"`
	ViewAngle map[Tick]ViewAngle          `json:"viewangle_at_ticks"`

	Deaths []Tick `json:"deaths_ticks"`
	Kills  []Tick `json:"kills_ticks"`

	previousFramePosition VelocityOrPosition
}

type Rounds struct {
	RoundEnd      []Tick `json:"round_end_ticks"`
	FreezeTimeEnd []Tick `json:"freeze_end_ticks"`
}

// MARK: Class
type GameState struct {
	// The start of the demo we care about:
	WarmupEnded Tick `json:"warmup_ended_tick"`
	// The end of the demo
	LastTick Tick `json:"last_tick"`

	// Round information to skip to
	Rounds Rounds `json:"rounds"`
	// Player information per round / tick
	Players map[SteamID32]*Player `json:"players"`
}

// MARK: Class Functions
func (game *GameState) IsWarmupEnded() bool {
	return (game.WarmupEnded > 0)
}

func (game *GameState) SetWarmupEnded(tick Tick) {
	game.WarmupEnded = tick
}

func (game *GameState) SetLastTick(tick Tick) {
	game.LastTick = tick
}

func (game *GameState) AddRoundEnd(tick Tick) {
	game.Rounds.RoundEnd = append(game.Rounds.RoundEnd, tick)
}

func (game *GameState) AddFreezeTimeEnd(tick Tick) {
	game.Rounds.FreezeTimeEnd = append(game.Rounds.FreezeTimeEnd, tick)
}

// MARK: Subclass Player

func (game *GameState) CreatePlayerIfMissingOrCurrent(steamID SteamID32) *Player {
	var player *Player

	if player, exists := game.Players[steamID]; exists {
		return player
	}

	player = &Player{
		Button:    make(map[Tick]ButtonMask),
		Deaths:    make([]Tick, 0),
		Kills:     make([]Tick, 0),
		Velocity:  make(map[Tick]VelocityOrPosition),
		ViewAngle: make(map[Tick]ViewAngle),
	}

	game.Players[steamID] = player
	return player
}

func (game *GameState) AddButton(tick Tick, steamID SteamID32, button ButtonMask) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.Button[tick] = button
}

func (game *GameState) AddDeath(tick Tick, steamID SteamID32) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.Deaths = append(player.Deaths, tick)
}

func (game *GameState) AddKill(tick Tick, steamID SteamID32) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.Kills = append(player.Kills, tick)
}

func (game *GameState) AddVelocity(tick Tick, steamID SteamID32, velocity VelocityOrPosition) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.Velocity[tick] = velocity
}

func (game *GameState) SetPreviousFramePosition(steamID SteamID32, position VelocityOrPosition) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.previousFramePosition = position
}

func (game *GameState) GetPreviousFramePosition(steamID SteamID32) VelocityOrPosition {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	return player.previousFramePosition
}

func (game *GameState) AddViewAngle(tick Tick, steamID SteamID32, x, y float32) {
	player := game.CreatePlayerIfMissingOrCurrent(steamID)

	player.ViewAngle[tick] = ViewAngle{X: float64(x), Y: float64(y)}
}

// MARK: Class Initalizer
func NewGameState() GameState {
	return GameState{
		WarmupEnded: -1,
		LastTick:    -1,
		Rounds:      Rounds{},
		Players:     make(map[SteamID32]*Player),
	}
}

func LoadGameStateFromJson(file string) GameState {
	var gamestate GameState
	fd, err := os.Open(file)

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(fd)
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&gamestate)

	if err != nil {
		panic(err)
	}

	return gamestate
}
