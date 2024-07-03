//go:generate bash ../../scripts/build-mirv-script.sh
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"mccadecortez.me/cs2-demo/v2/pkg/button"
	"mccadecortez.me/cs2-demo/v2/pkg/gamestate"
)

type GameEvent struct {
	EventName string          `json:"eventName"`
	Values    []GameEventData `json:"values"`
}

type GameEventData struct {
	ServerTick int `json:"server_tick"`
}

var ENV_STEAMID, _ = strconv.ParseInt(os.Getenv("STEAMID"), 0, 64)

func main() {
	// MARK: Flags
	DemoJson := flag.String("demo-json", os.Getenv("DEMO_JSON"), "Path to the parsed demo JSON")
	// FIXME: This is needed because mirv hardcodes the URL for the client/game
	Url := flag.String("url", os.Getenv("URL"), "Websocket URI to connect to, ws://<host>/mirv")
	SpectatingSteamID := flag.Int64("steam-id", ENV_STEAMID, "Value of the convar 'spec_lock_to_accountid', a player entry in the Json")
	MirvServerDir := flag.String("mirv-streams-dir", os.Getenv("MIRV_STREAMS_DIR"), "Directory of the submodule mirv-script to run 'npm run server'")

	flag.Parse()

	if *DemoJson == "" {
		log.Fatal("missing --demo-json or DEMO_JSON")
	}

	if *Url == "" {
		log.Fatal("missing --url or URL")
	} else if !strings.HasSuffix(*Url, "?user=1") {
		*Url = *Url + "?user=1"
	}

	if *SpectatingSteamID == 0 {
		log.Fatal("missing --steam-id or STEAMID")
	}

	if *MirvServerDir == "" {
		log.Fatal("missing --mirv-streams-dir or MIRV_STREAMS_DIR")
	}

	// XXX: ignoring error because of ptr
	*DemoJson, _ = filepath.Abs(*DemoJson)
	*MirvServerDir, _ = filepath.Abs(*MirvServerDir)

	fmt.Println(*DemoJson, *MirvServerDir)

	// MARK: Start Mirv Server (nodejs)
	os.Chdir(*MirvServerDir)
	cmd := exec.Command("node", "dist/node/server.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// XXX: Not waiting on cmd, ignoring if the command fails because we won't be able to connect later
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	fmt.Println(cmd.Process.Pid)
	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
	}()

	time.Sleep(time.Second * 10)

	// MARK: Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(*Url, nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server:", err) // XXX: log.Fatal skips defer

		return
	}
	defer conn.Close()

	log.Println("Connected to", *Url)

	// MARK: Parse JSON Gamestate
	var gs gamestate.GameState = gamestate.LoadGameStateFromJson(*DemoJson)
	var event GameEvent
	var player = gs.Players[gamestate.SteamID32(*SpectatingSteamID)]

	if player == nil {
		panic("player == nil")
	}

	// MARK: Websocket Loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Panic("Error reading message:", err)
		}

		if err := json.Unmarshal(message, &event); err != nil || len(event.Values) == 0 {
			continue
		}

		value := event.Values[0]
		serverTick := value.ServerTick

		if serverTick <= 0 {
			continue
		}

		t := gamestate.Tick(serverTick)

		button_at_tick, ok := player.Button[t]
		if !ok {
			continue
		}

		viewangle, ok := player.ViewAngle[t]
		if !ok {
			continue
		}

		velocity, ok := player.Velocity[t]
		if !ok {
			continue
		}

		pressed := button.ParseButtonMask(uint64(button_at_tick))

		fmt.Print("\033[H\033[2J")
		fmt.Printf("%#v\n", pressed)
		fmt.Printf("\n%#v\n", viewangle)
		fmt.Printf("\n%#v\n", velocity)
	}
}
