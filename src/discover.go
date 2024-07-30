package main

// TODO: Integrate Device
// TODO: Integrate Availability
// TODO: Improve Error Handling (unkown or missing tags, tags in wrong order)
// TODO:

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/schretzi/go_oekofen/src/env"
	"github.com/schretzi/go_oekofen/src/state"
	log "github.com/sirupsen/logrus"
)

type Device_object struct {
	Identifiers        []string `default:"Oekofen"`
	Manufacturer       string   `default:"Ökofen"`
	Model              string   `default:"Peletronic"`
	Name               string   `default:"Oekofen"`
	Device_class       string
	Enabled_by_default bool
	State_topic        string
	Value_template     string
}

type Sensor_object struct {
	Unique_id           string `json:"unique_id"`
	Object_id           string `json:"object_id"`
	Name                string `json:"name"`
	State_topic         string `json:"state_topic"`
	Device_class        string `json:"device_class"`
	Unit_of_measurement string `json:"unit_of_measurement"`
	//Availability_topic    string `json:"availability_topic"`
	//Availability_mode     string `json:"availability_mode"`
	Expire_after int16 `json:"expire_after"`
	//Payload_available     string `json:"payload_available"`
	//Payload_not_available string `json:"payload_not_available"`
	//Value_template        string `json:"value_template"`
}

type Number_object struct {
	Unique_id      string `json:"unique_id"`
	Object_id      string `json:"object_id"`
	Name           string `json:"name"`
	State_topic    string `json:"state_topic"`
	Command_topic  string `json:"command_topic"`
	Max            int32  `json:"max"`
	Min            int32  `json:"min"`
	Value_template string `json:"value_template"`
}

type Text_object struct {
	Unique_id     string `json:"unique_id"`
	Name          string `json:"name"`
	State_topic   string `json:"state_topic"`
	Command_topic string `json:"command_topic"`
}

type Switch_object struct {
	Unique_id          string
	Name               string
	State_topic        string
	Command_topic      string `default:""`
	Availability_topic string `default:""`
	Payload_on         string `default:"ON"`
	Payload_off        string `default:"OFF"`
	State_on           string `default:"ON"`
	State_off          string `default:"OFF"`
	Optimistic         bool
}

const oekofen_tagName = "oekofen"

func discovery_main(ctx context.Context, ofen Oekofen, wg *sync.WaitGroup) {
	defer wg.Done()
	discovery_create_configs(ctx, reflect.TypeOf(&ofen).Elem(), "okeofen")

}

func discovery_create_configs(ctx context.Context, t reflect.Type, parent string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if env.CfgDebug() {
			log.Trace("*" + parent + "_" + field.Name)
		}
		if field.Type.Kind() == reflect.Struct {
			discovery_create_configs(ctx, field.Type, field.Name)
			continue
		}

		oekofen_tag_string := field.Tag.Get(oekofen_tagName)
		if len(oekofen_tag_string) == 0 {
			continue
		}
		oekofen_tags := strings.Split(oekofen_tag_string, ",")
		oekofen_ha_component := strings.Split(oekofen_tags[0], "#")
		oekofen_transform := strings.Split(oekofen_tags[1], "#")
		var transform string
		if len(oekofen_transform) == 2 {
			transform = oekofen_transform[1]
		} else {
			transform = ""
		}
		config_topic := "homeassistant/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/config"

		if oekofen_transform[0] == "float" {
			state.InitStateFloat(parent+"_"+field.Name, transform)
		} else if oekofen_transform[0] == "int" {
			state.InitStateInteger(parent+"_"+field.Name, transform)
		} else if oekofen_transform[0] == "string" {
			state.InitStateText(parent+"_"+field.Name, transform)
		} else if oekofen_transform[0] == "bool" {
			state.InitStateBool(parent+"_"+field.Name, transform)
		} else {
			log.Error("unkown data type for value: " + parent + "_" + field.Name)
		}

		if oekofen_ha_component[0] == "sensor" {

			var sensor Sensor_object

			sensor.State_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/state"
			//sensor.Availability_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/available"
			//sensor.Availability_mode = "any"
			sensor.Name = parent + "_" + field.Name
			sensor.Unique_id = parent + "_" + field.Name
			sensor.Object_id = parent + "_" + field.Name

			//sensor.Availability_mode = "latest"
			if oekofen_ha_component[1] == "temperature" {
				sensor.Unit_of_measurement = "°C"
				sensor.Device_class = "temperature"
			} else {
				sensor.Device_class = oekofen_ha_component[1]
				sensor.Unit_of_measurement = "None"
			}
			sensor.Expire_after = 3660
			//sensor.Payload_available = "online"
			//sensor.Payload_not_available = "offline"

			mqtt_push_discovery(ctx, config_topic, sensor)
			//mqtt_push_available(ctx, sensor.Availability_topic, "online")
		}

		if oekofen_ha_component[0] == "number" {
			var number Number_object

			number.State_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/state"
			number.Command_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/command"
			number.Name = parent + "_" + field.Name
			number.Unique_id = parent + "_" + field.Name
			number.Object_id = parent + "_" + field.Name
			number.Min = -1
			number.Max = 33000
			//number.Value_template = "'{{ ((value_json.state )) }}'"
			mqtt_push_discovery(ctx, config_topic, number)
		}

		if oekofen_ha_component[0] == "text" {
			var text Text_object

			text.State_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/state"
			text.Command_topic = "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + field.Name + "/command"
			text.Name = parent + "_" + field.Name
			text.Unique_id = parent + "_" + field.Name
			mqtt_push_discovery(ctx, config_topic, text)
		}

	}
}
