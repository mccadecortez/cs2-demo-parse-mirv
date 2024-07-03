//go:generate bash ../../scripts/build-mirv-script.sh
//go:generate go run aadsfadsf

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"mccadecortez.me/cs2-demo/v2/pkg/gamestate"
	"mccadecortez.me/cs2-demo/v2/pkg/toxml"
)

func main() {
	demoFile := flag.String("demo-file", os.Getenv("DEMO_FILE"), "Path to the demo file")
	outputXML := flag.String("output-xml", os.Getenv("OUTPUT_XML"), "Path to the output XML file")
	outputJSON := flag.String("output-json", os.Getenv("OUTPUT_JSON"), "Path to the output JSON file")
	shouldRecord := flag.Bool("should-record", true, "Tells the game to record the demo when replaying")

	flag.Parse()

	if *demoFile == "" {
		log.Fatal("missing --demo-file")
	}

	if *outputXML == "" {
		*outputXML = "output.xml"
	}

	if *outputJSON == "" {
		*outputJSON = "output.json"
	}

	f, err := os.Open(*demoFile)
	if err != nil {
		log.Panic("failed to open demo file: ", err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	var gs = gamestate.NewGameState()

	// MARK: Handlers

	var getCurrentTick = func() gamestate.Tick {
		return gamestate.Tick(p.GameState().IngameTick())
	}

	var handlerRoundEndOffical = func(e events.RoundEndOfficial) {
		gs.AddRoundEnd(getCurrentTick())
	}

	var handlerRoundFreezetimeEnd = func(e events.RoundFreezetimeEnd) {
		gs.AddFreezeTimeEnd(getCurrentTick())
	}

	var handlerKill = func(e events.Kill) {
		if e.Killer != nil {
			gs.AddKill(getCurrentTick(), gamestate.SteamID32(e.Killer.SteamID32()))
		}

		if e.Victim == nil {
			log.Panicln("e.Victim == nil, demo courrupted?")
		}

		gs.AddDeath(getCurrentTick(), gamestate.SteamID32(e.Victim.SteamID32()))
	}

	var handlerIsWarmupPeriodChanged = func(e events.IsWarmupPeriodChanged) {
		if e.NewIsWarmupPeriod /* new value */ && !gs.IsWarmupEnded() {
			gs.SetWarmupEnded(getCurrentTick())

			p.RegisterEventHandler(handlerRoundEndOffical)
			p.RegisterEventHandler(handlerRoundFreezetimeEnd)
			p.RegisterEventHandler(handlerKill)
		}
	}

	p.RegisterEventHandler(handlerIsWarmupPeriodChanged)

	// MARK: Parse Loop
	fmt.Println("Parsing demo", *demoFile)

	// assumes ParseNextFrame sets ok to `false` on err
	for ok := true; ok; ok, _ = p.ParseNextFrame() {
		players := p.GameState().Participants().Playing()

		gametick := getCurrentTick()
		// MARK: Button Pressed
		for _, player := range players {
			if !player.IsAlive() {
				continue
			}

			steamID := gamestate.SteamID32(player.SteamID32())
			m_nButtonDownMaskPrev := gamestate.ButtonMask(player.PlayerPawnEntity().PropertyValueMust("m_pMovementServices.m_nButtonDownMaskPrev").S2UInt64())
			player.PreviousFramePosition = r3.Vector(gs.GetPreviousFramePosition(steamID))

			gs.AddButton(gametick, steamID, m_nButtonDownMaskPrev)
			gs.AddVelocity(gametick, steamID, gamestate.VelocityOrPosition(player.Velocity()))
			gs.AddViewAngle(gametick, steamID, player.ViewDirectionX(), player.ViewDirectionY())

			gs.SetPreviousFramePosition(steamID, gamestate.VelocityOrPosition(player.Position()))
		}
	}

	gs.SetLastTick(getCurrentTick())

	// MARK: Sanity Check

	if len(gs.Rounds.RoundEnd) == 0 {
		log.Panic("len(gs.Rounds.RoundEnd) == 0")
	} else if len(gs.Rounds.FreezeTimeEnd) == 0 {
		log.Panic("gs.Rounds.FreezeTimeEnd) == 0")
	}

	// MARK: Write Files
	// MARK: JSON Gamestate
	s, err := json.MarshalIndent(gs, "", "\t")

	if err != nil {
		panic(err)
	}

	os.WriteFile(*outputJSON, s, 0664)
	log.Printf("Wrote json output to: %s", *outputJSON)

	// MARK: XML Demo Cam
	skip_to_first_round := gs.Rounds.FreezeTimeEnd[0]
	game_end := gs.LastTick

	for steamID, player := range gs.Players {
		var commands []toxml.Command = make([]toxml.Command, 0)
		var prev_goto_round int64 = int64(skip_to_first_round) // prevent a deadlock by doing `mirv_skip tick to` crossing paths with another `mirv_skip tick to`

		commands = append(commands, toxml.Command{Tick: "0", Command: fmt.Sprintf("echoln spectating %d", steamID)})
		commands = append(commands, toxml.Command{Tick: "0", Command: fmt.Sprintf("spec_lock_to_accountid %d", steamID)})
		commands = append(commands, toxml.Command{Tick: "0", Command: fmt.Sprintf("mirv_skip tick to %d", skip_to_first_round)})

		if *shouldRecord {
			commands = append(commands, toxml.Command{Tick: "0", Command: "mirv_streams record start"})
			commands = append(commands, toxml.Command{Tick: fmt.Sprintf("%d", game_end-1000), Command: "mirv_streams record end"})
		}

		commands = append(commands, toxml.Command{Tick: fmt.Sprintf("%d", game_end-100), Command: "quit"})

		iter := append(player.Deaths, gs.Rounds.RoundEnd...)
		slices.Sort(iter) // sort because of the if `tick >= round` check

		// XXX: Skip freezetime & time the player is not alive
		for _, tick := range iter {
			var goto_round int64 = 0
			for _, round := range gs.Rounds.FreezeTimeEnd {
				if tick >= round {
					continue
				}

				goto_round = int64(round)
				break
			}

			if goto_round != 0 && goto_round != prev_goto_round {
				prev_goto_round = goto_round

				commands = append(commands, toxml.Command{Tick: fmt.Sprintf("%d", tick), Command: fmt.Sprintf("mirv_skip tick to %d", goto_round)})
			}
		}

		out := fmt.Sprintf("%d_%s", steamID, *outputXML)
		os.WriteFile(out, []byte(toxml.ToXML(commands)), 0664)
		log.Printf("Wrote xml output to: %s", out)
	}
}
