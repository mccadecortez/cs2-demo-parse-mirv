package gamestate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mccadecortez.me/cs2-demo/v2/pkg/gamestate"
)

var g_gamestate = gamestate.NewGameState()
var g_tickcounter gamestate.Tick = 1000
var player_foo gamestate.SteamID32 = 1234
var player_bar gamestate.SteamID32 = 6789

func TestIsWarmupEnded(t *testing.T) {
	assert := assert.New(t)

	assert.False(g_gamestate.IsWarmupEnded())

	g_tickcounter++
	g_gamestate.SetWarmupEnded(g_tickcounter)

	assert.True(g_gamestate.IsWarmupEnded())
}

func TestAddRoundEnd(t *testing.T) {
	assert := assert.New(t)

	g_tickcounter++
	g_gamestate.AddRoundEnd(g_tickcounter)

	assert.Contains(g_gamestate.Rounds.RoundEnd, g_tickcounter)
}

func TestAddFreezeTimeEnd(t *testing.T) {
	assert := assert.New(t)

	g_tickcounter++
	g_gamestate.AddFreezeTimeEnd(g_tickcounter)

	assert.Contains(g_gamestate.Rounds.FreezeTimeEnd, g_tickcounter)
}

func TestAddButton(t *testing.T) {
	assert := assert.New(t)
	var button gamestate.ButtonMask = 513

	g_tickcounter++
	g_gamestate.AddButton(g_tickcounter, player_foo, button)
	g_gamestate.AddButton(g_tickcounter, player_bar, button)

	assert.Contains(g_gamestate.Players[player_foo].Button, g_tickcounter)
	assert.Equal(g_gamestate.Players[player_foo].Button[g_tickcounter], button)
}

func TestAddDeath(t *testing.T) {
	assert := assert.New(t)
	g_tickcounter++
	g_gamestate.AddDeath(g_tickcounter, player_foo)
	g_gamestate.AddDeath(g_tickcounter, player_bar)

	assert.Contains(g_gamestate.Players[player_foo].Deaths, g_tickcounter)
}

func TestAddKill(t *testing.T) {
	assert := assert.New(t)
	g_tickcounter++
	g_gamestate.AddKill(g_tickcounter, player_foo)
	g_gamestate.AddKill(g_tickcounter, player_bar)

	assert.Contains(g_gamestate.Players[player_foo].Kills, g_tickcounter)
}

// TODO: Getters
