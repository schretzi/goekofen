package state

import (
	"fmt"
	"strconv"
	"time"

	"github.com/schretzi/go_oekofen/src/env"
	log "github.com/sirupsen/logrus"
)

/*
	 This packages hold the latest state of the ofen known by the program.

	   Reason:
	   - Ökofen delivers temperatur in dezi-degree Celcius, so we need to transform temperatur and some others (state-text because of encoding in german language) for better display in home-assistant
	   - I want to reduce number of messages sent to MQTT if there is no update on the value, but for Lifetime-reasons updates are sent at least once per hour, even if the value stays the same

	   When getting new data from the ofen the actual datapoint in original value is compared, so we don't run transformation for a unchanged value.
	   - If both datapoints are equal no further action
	   - If datapoint is updated
		 - transform original value to homeassistant value
		 - save updated values in state
		 - check if last_update was sent within one hour
		 - return true, so updates to MQTT (and Influx?) needs to be sent
*/

type StateInteger struct {
	oekofen_value       int64
	homeassistant_value int64
	transform           string
	last_update         time.Time
}

type StateFloat struct {
	oekofen_value       int64
	homeassistant_value float64
	transform           string
	last_update         time.Time
}

type StateText struct {
	oekofen_value       string
	homeassistant_value string
	transform           string
	last_update         time.Time
}

type StateBool struct {
	oekofen_value       bool
	homeassistant_value bool
	transform           string
	last_update         time.Time
}

// var oekofen_state = Oekofen_state{}
var stateText map[string]*StateText
var stateInteger map[string]*StateInteger
var stateBool map[string]*StateBool
var stateFloat map[string]*StateFloat

func NewState() {
	if env.CfgDebug() {
		log.SetLevel(log.DebugLevel)
	}
	stateText = make(map[string]*StateText)
	stateInteger = make(map[string]*StateInteger)
	stateBool = make(map[string]*StateBool)
	stateFloat = make(map[string]*StateFloat)
}

func InitStateText(key string, transform string) {
	early_times_string := "2020-01-01 00:00:00"
	early_times, _ := time.Parse("2006-01-02 03:04:05", early_times_string)
	tmp_object := new(StateText)
	tmp_object.last_update = early_times
	tmp_object.transform = transform
	stateText[key] = tmp_object
}

func InitStateInteger(key string, transform string) {
	early_times_string := "2020-01-01 00:00:00"
	early_times, _ := time.Parse("2006-01-02 03:04:05", early_times_string)
	tmp_object := new(StateInteger)
	tmp_object.last_update = early_times
	tmp_object.transform = transform
	stateInteger[key] = tmp_object
}

func InitStateFloat(key string, transform string) {
	early_times_string := "2020-01-01 00:00:00"
	early_times, _ := time.Parse("2006-01-02 03:04:05", early_times_string)
	tmp_object := new(StateFloat)
	tmp_object.last_update = early_times
	tmp_object.transform = transform
	stateFloat[key] = tmp_object
}

func InitStateBool(key string, transform string) {
	early_times_string := "2020-01-01 00:00:00"
	early_times, _ := time.Parse("2006-01-02 03:04:05", early_times_string)
	tmp_object := new(StateBool)
	tmp_object.last_update = early_times
	tmp_object.transform = transform
	stateBool[key] = tmp_object
}

func UpdateStateText(key string, value string) (bool, string) {
	log.Info("HA Value: ", value)
	state_object := stateText[key]
	if state_object.oekofen_value == value {
		if state_object.last_update.Before(time.Now().Add(-1 * time.Hour)) {
			state_object.last_update = time.Now()
			return true, state_object.homeassistant_value
		} else {
			return false, ""
		}
	} else {
		log.Debug("HA Value: ", value)
		ha_text_transformed := Invoke_Text(Functions{}, state_object.transform, value)
		log.Debug("Transformed: ", ha_text_transformed)
		log.Info("Changed Value: " + key + " / From: " + state_object.homeassistant_value + " To: " + value)
		state_object.oekofen_value = value
		state_object.homeassistant_value = ha_text_transformed
		state_object.last_update = time.Now()

		return true, ha_text_transformed
	}
}

func UpdateStateFloat(key string, value int64) (bool, float64) {

	state_object := stateFloat[key]
	if state_object.oekofen_value == value {
		if state_object.last_update.Before(time.Now().Add(-1 * time.Hour)) {
			state_object.last_update = time.Now()
			log.Debug("Changed because of Time: " + key)
			return true, state_object.homeassistant_value
		} else {
			return false, 0
		}
	} else {
		ha_value_unrounded := Invoke_Float(Functions{}, state_object.transform, value)
		ha_value := roundFloat(ha_value_unrounded, 2)
		state_object.oekofen_value = value
		state_object.homeassistant_value = ha_value
		state_object.last_update = time.Now()
		log.Info("Changed Value: " + key + " / From: " + strconv.FormatFloat(state_object.homeassistant_value, 'f', 2, 64) + " To: " + strconv.FormatFloat(ha_value, 'f', 2, 64))

		return true, ha_value
	}
}

func UpdateStateInt(key string, value int64) (bool, int64) {
	state_object := stateInteger[key]
	if state_object.oekofen_value == value {
		if state_object.last_update.Before(time.Now().Add(-1 * time.Hour)) {
			state_object.last_update = time.Now()
			log.Debug("Changed because of Time: " + key)
			return true, state_object.homeassistant_value
		} else {
			return false, 0
		}
	} else {
		var ha_value int64
		if len(state_object.transform) > 0 {
			ha_value = Invoke_Int(Functions{}, state_object.transform, value)
		} else {
			ha_value = value
		}
		log.Info("Changed Value: " + key + " / From: " + strconv.FormatInt(state_object.homeassistant_value, 10) + " To: " + strconv.FormatInt(ha_value, 10))
		state_object.oekofen_value = value
		state_object.homeassistant_value = ha_value
		state_object.last_update = time.Now()

		return true, ha_value
	}
}

func UpdateStateBool(key string, value bool) (bool, bool) {
	state_object := stateBool[key]
	if state_object.oekofen_value == value {
		if state_object.last_update.Before(time.Now().Add(-1 * time.Hour)) {
			state_object.last_update = time.Now()
			log.Debug("Changed because of Time: " + key)
			return true, value
		} else {
			return false, false
		}
	} else {
		state_object.oekofen_value = value
		log.Info("Changed Value: " + key + " / From: " + strconv.FormatBool(state_object.homeassistant_value) + " To: " + strconv.FormatBool(value))
		state_object.homeassistant_value = value
		state_object.last_update = time.Now()

		return true, value
	}
}

func ShowState() {
	fmt.Printf("Text: %#v\n", stateText)
	fmt.Printf("Integer: %#v\n", stateInteger)
	fmt.Printf("Float: %#v\n", stateFloat)
	fmt.Printf("Bool: %#v\n", stateBool)

}
