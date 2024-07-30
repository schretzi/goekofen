package main

// TODO: Refactor ofen/Oekofen struct of structs to map of structs to support multiple/other entities.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/schretzi/go_oekofen/src/env"
	"github.com/schretzi/go_oekofen/src/state"
	log "github.com/sirupsen/logrus"
)

type System struct {
	Ambient   int64 `json:"L_ambient,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Errors    int64 `json:"L_errors,string" oekofen:"number#None,int"`
	Usb_stick bool  `json:"L_usb_stick,string"`
}

type Weather struct {
	Temp            int64  `json:"L_temp,string"`
	Clouds          int64  `json:"L_clouds,string"`
	Forecast_temp   int64  `json:"L_forecast_temp,string"`
	Forecast_clouds int64  `json:"L_forecast_clouds,string"`
	Forecast_today  bool   `json:"L_forecast_today,string"`
	Starttime       int64  `json:"L_starttime,string"`
	Endtime         int64  `json:"L_endtime,string"`
	Source          string `json:"L_source"`
	Location        string `json:"L_location"`
	Cloud_limit     int64  `json:"cloud_limit,string"`
	Hysteresys      int64  `json:"hysteresys,string"`
	Offtemp         int64  `json:"offtemp,string"`
	Lead            int64  `json:"lead,string"`
	Refresh         bool   `json:"refresh,string"`
	Oekomode        int64  `json:"oekomode,string"`
}

type Forecast struct {
	W_0  string `json:"L_w_0"`
	W_1  string `json:"L_w_1"`
	W_2  string `json:"L_w_2"`
	W_3  string `json:"L_w_3"`
	W_4  string `json:"L_w_4"`
	W_5  string `json:"L_w_5"`
	W_6  string `json:"L_w_6"`
	W_7  string `json:"L_w_7"`
	W_8  string `json:"L_w_8"`
	W_9  string `json:"L_w_9"`
	W_10 string `json:"L_w_10"`
	W_11 string `json:"L_w_11"`
	W_12 string `json:"L_w_12"`
	W_13 string `json:"L_w_13"`
	W_14 string `json:"L_w_14"`
	W_15 string `json:"L_w_15"`
	W_16 string `json:"L_w_16"`
	W_17 string `json:"L_w_17"`
	W_18 string `json:"L_w_18"`
	W_19 string `json:"L_w_19"`
	W_20 string `json:"L_w_20"`
	W_21 string `json:"L_w_21"`
	W_22 string `json:"L_w_22"`
	W_23 string `json:"L_w_23"`
	W_24 string `json:"L_w_24"`
}

type Heating_curcuit struct {
	Roomtemp_act        int64  `json:"L_roomtemp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Roomtemp_set        int64  `json:"L_roomtemp_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Flowtemp_act        int64  `json:"L_flowtemp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Flowtemp_set        int64  `json:"L_flowtemp_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Comfort             int64  `json:"L_comfort,string" oekofen:"sensor#None,int"`
	State               int64  `json:"L_state,string" oekofen:"sensor#None,int"`
	Statetext           string `json:"L_statetext" oekofen:"text,string#Transform_Statetext"`
	Pump                bool   `json:"L_pump,string"`
	Remote_override     int64  `json:"remote_override,string" oekofen:"number#None,int"`
	Mode_auto           int64  `json:"mode_auto,string" oekofen:"number#None,int"`
	Time_prg            int64  `json:"time_prg,string" oekofen:"number#None,int"`
	Temp_setback        int64  `json:"temp_setback,string" oekofen:"number#temperature,float#Transform_deziCelsius"`
	Temp_heat           int64  `json:"temp_heat,string" oekofen:"number#temperature,float#Transform_deziCelsius"`
	Temp_vacation       int64  `json:"temp_vacation,string" oekofen:"number#temperature,float#Transform_deziCelsius"`
	Name                string `json:"name"`
	Oekomode            int64  `json:"oekomode,string" oekofen:"number#None,int"`
	Autocomfort         int64  `json:"autocomfort,string" oekofen:"number#None,int"`
	Autocomfort_sunset  int64  `json:"autocomfort_sunset,string" oekofen:"number#None,int"`
	Autocomfort_sunrise int64  `json:"autocomfort_sunrise,string" oekofen:"number#None,int"`
}

type Puffer struct {
	Temp_oben_act   int64  `json:"L_tpo_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Temp_oben_set   int64  `json:"L_tpo_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Temp_mitte_act  int64  `json:"L_tpm_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Temp_mitte_set  int64  `json:"L_tpm_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Pump_release    int64  `json:"L_pump_release,string"`
	Pump            int64  `json:"L_pump,string" oekofen:"number#None,int"`
	State           int64  `json:"L_state,string" oekofen:"number#None,int"`
	Statetext       string `json:"L_statetext" oekofen:"text,string#Transform_Statetext"`
	Mintemp_off     int64  `json:"mintemp_off,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Mintemp_on      int64  `json:"mintemp_on,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Ext_mintemp_off int64  `json:"ext_mintemp_off,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Ext_mintemp_on  int64  `json:"ext_mintemp_on,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
}

type Warm_water struct {
	Temp_set        int64  `json:"L_temp_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Ontemp_act      int64  `json:"L_ontemp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Offtemp_act     int64  `json:"L_offtemp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Pump            bool   `json:"L_pump,string"`
	State           int64  `json:"L_state,string" oekofen:"number#None,int"`
	Statetext       string `json:"L_statetext" oekofen:"text,string#Transform_Statetext"`
	Time_prg        int64  `json:"time_prg,string"`
	Sensor_on       int64  `json:"sensor_on,string"`
	Sensor_off      int64  `json:"sensor_off,string"`
	Mode_auto       int64  `json:"mode_auto,string" oekofen:"number#None,int"`
	Mode_dhw        int64  `json:"mode_dhw,string" oekofen:"number#None,int"`
	Heat_once       bool   `json:"heat_once,string"`
	Temp_min_set    int64  `json:"temp_min_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Temp_max_set    int64  `json:"temp_max_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Name            string `json:"name"`
	Smartstart      int64  `json:"smartstart,string" oekofen:"number#None,int"`
	Use_boiler_heat int64  `json:"use_boiler_heat,string" oekofen:"number#None,int"`
	Oekomode        int64  `json:"oekomode,string" oekofen:"number#None,int"`
}

type Pellets struct {
	Temp_act               int64  `json:"L_temp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Temp_set               int64  `json:"L_temp_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Ext_temp               int64  `json:"L_ext_temp,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Frt_temp_act           int64  `json:"L_frt_temp_act,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Frt_temp_set           int64  `json:"L_frt_temp_set,string" oekofen:"sensor#temperature,float#Transform_deziCelsius"`
	Br                     bool   `json:"L_br,string"`
	Ak                     bool   `json:"L_ak,string"`
	Not                    bool   `json:"L_not,string"`
	Stb                    bool   `json:"L_stb,string"`
	Modulation             int64  `json:"L_modulation,string" oekofen:"number#None,int"`
	Uw_speed               int64  `json:"L_uw_speed,string" oekofen:"number#None,int"`
	State                  int64  `json:"L_state,string" oekofen:"number#None,int"`
	Statetext              string `json:"L_statetext" oekofen:"text,string#Transform_Statetext"`
	Pe_type                int64  `json:"L_type,string" oekofen:"number#None,int"`
	Starts                 int64  `json:"L_starts,string" oekofen:"number#None,int"`
	Runtime                int64  `json:"L_runtime,string" oekofen:"number#None,int"`
	Avg_runtime            int64  `json:"L_avg_runtime,string" oekofen:"number#None,int"`
	Uw_release             int64  `json:"L_uw_release,string" oekofen:"number#None,int"`
	Uw                     int64  `json:"L_uw,string" oekofen:"number#None,int"`
	Storage_fill           int64  `json:"L_storage_fill,string"`
	Storage_min            int64  `json:"L_storage_min,string"`
	Storage_max            int64  `json:"L_storage_max,string"`
	Storage_popper         int64  `json:"L_storage_popper,string"`
	Storage_fill_today     int64  `json:"storage_fill_today,string"`
	Storage_fill_yesterday int64  `json:"storage_fill_yesterday,string"`
	Mode                   int64  `json:"mode,string" oekofen:"number#None,int"`
}

type Ofen_error struct {
	Ofen_error string
}

type Oekofen struct {
	System     System          `json:"system"`
	Weather    Weather         `json:"weather"`
	Forecast   Forecast        `json:"forecast"`
	HK1        Heating_curcuit `json:"hk1"`
	HK2        Heating_curcuit `json:"hk2"`
	PU1        Puffer          `json:"pu1"`
	WW1        Warm_water      `json:"ww1"`
	PE1        Pellets         `json:"pe1"`
	Ofen_error Ofen_error      `json:"error"`
}

func read_oekofen() Oekofen {
	var ofen Oekofen

	req, err := http.NewRequest(http.MethodGet, env.CfgOekofenURL().String(), nil)
	if err != nil {
		log.Panicf("could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("error making http request: %s\n", err)
	}

	log.Debugf("client: got response!\n")
	log.Debugf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicf("client: could not read response body: %s\n", err)
	}

	fmt.Printf("client: response body: %s\n", resBody)

	j_err := json.Unmarshal([]byte(resBody), &ofen)
	if j_err != nil {
		log.Panicf("An error occured: %v\n", j_err)
	}
	log.Debugf("Oekofen: %v\n", ofen)

	return ofen
}

func oekofen_push_status(ctx context.Context, ofen Oekofen, wg *sync.WaitGroup) {
	defer wg.Done()
	oekofen_state(ctx, ofen, "okeofen")

}

func oekofen_state(ctx context.Context, val interface{}, parent string) {
	log.Debug("Started Parsing of State")
	if reflect.ValueOf(val).Kind() == reflect.Struct {
		v := reflect.ValueOf(val)
		t := reflect.TypeOf(val)
		for i := 0; i < v.NumField(); i++ {
			log.Debug(t.Field(i).Name)
			if t.Field(i).Type.Kind() == reflect.Struct {
				oekofen_state(ctx, v.Field(i).Interface(), t.Field(i).Name)
			} else {

				oekofen_tag_string := t.Field(i).Tag.Get(oekofen_tagName)
				oekofen_tags := strings.Split(oekofen_tag_string, ",")
				oekofen_ha_component := strings.Split(oekofen_tags[0], "#")

				if len(oekofen_tag_string) == 0 {
					continue
				}
				state_topic := "oekofen/" + oekofen_ha_component[0] + "/" + parent + "_" + t.Field(i).Name + "/state"
				var state_value any
				if len(oekofen_tags) > 1 {
					oekofen_transform_options := strings.Split(oekofen_tags[1], "#")
					log.Debug(oekofen_transform_options[0])

					if oekofen_transform_options[0] == "float" {
						log.Debug(v.Field(i).Int())
						changed, changed_value := state.UpdateStateFloat(parent+"_"+t.Field(i).Name, v.Field(i).Int())
						if changed {
							state_value = changed_value
							mqtt_push_status(ctx, state_topic, state_value)
						}
					} else if oekofen_transform_options[0] == "int" {
						log.Debug(v.Field(i).Int())
						changed, changed_value := state.UpdateStateInt(parent+"_"+t.Field(i).Name, v.Field(i).Int())
						if changed {
							state_value = changed_value
							mqtt_push_status(ctx, state_topic, state_value)
						}
					} else if oekofen_transform_options[0] == "string" {
						log.Debug(v.Field(i).String())
						changed, changed_value := state.UpdateStateText(parent+"_"+t.Field(i).Name, v.Field(i).String())
						if changed {
							state_value = changed_value
							mqtt_push_status(ctx, state_topic, state_value)
						}
					}
				} else {
					state_value = v.Field(i).String()
				}

				log.Tracef("Field: %s \t type: %T \t value: %v \t  tags: %s \t state: %s \n",
					parent+"/"+t.Field(i).Name, v.Field(i), v.Field(i), oekofen_tags, state_value)
				log.Trace(state_topic)
				log.Trace(state_value)

			}
		}
	}

}
