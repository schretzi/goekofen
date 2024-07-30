package main

// TODO: Implement MQTT listening for status changes on selected parameters
// TODO: Make Timing PAramter configurable (Delay between Reads, validity of MQTT messages (Adapt Validity of State according to MQTT messages (+10sec)))
// TODO: Transform Mode in Mode_Text for easier display in Home assistant

import (
	"context"
	"sync"

	"github.com/schretzi/go_oekofen/src/env"
	"github.com/schretzi/go_oekofen/src/state"
	log "github.com/sirupsen/logrus"
)

func main() {
	state.NewState()

	err := env.ReadConfig()
	if err != nil {
		log.Panic(err)
	}
	if env.CfgDebug() {
		log.SetLevel(log.DebugLevel)
	}
	log.Info("goekofen started, setting up Connections")
	cliCfg := mqtt_setup_connection()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqtt_start_connection(cliCfg, ctx)

	if env.CfgHA() {
		var wg sync.WaitGroup
		log.Info("Setup home assistant discovery via MQTT")
		ofen := read_oekofen()
		log.Info("  Got discovery info from ofen ")
		log.Info("  Wait 3 seconds between consecutive Calls ")

		wg.Add(2)
		go discovery_main(ctx, ofen, &wg)
		go sleep_sec(3, &wg)
		wg.Wait()

	}

	var wg sync.WaitGroup

	for {
		wg.Add(1)
		log.Info("Send status to home assistant via MQTT")
		ofen := read_oekofen()
		log.Info("  Got Status from Ofen")
		log.Info("  Send to MQTT and wait 1 minute")
		if env.CfgHA() {
			wg.Add(1)
			go oekofen_push_status(ctx, ofen, &wg)
		}

		go sleep_sec(10, &wg)
		wg.Wait()
	}

}
