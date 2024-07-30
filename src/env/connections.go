package env

import "github.com/eclipse/paho.golang/autopaho"

type Connections struct {
	mqtt *autopaho.ConnectionManager
}

var connections = Connections{}

func SetConnMqtt(mqtt *autopaho.ConnectionManager) {
	connections.mqtt = mqtt
}

func GetConnMqtt() (mqtt *autopaho.ConnectionManager) {
	return connections.mqtt
}
