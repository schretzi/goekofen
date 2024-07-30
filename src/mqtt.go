package main

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/schretzi/go_oekofen/src/env"
	log "github.com/sirupsen/logrus"
)

func mqtt_setup_connection() autopaho.ClientConfig {
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{env.CfgMqttURL()},
		KeepAlive:         env.CfgMqttKeepAlive(),
		ConnectRetryDelay: env.CfgMqttRetryDelay(),
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { log.Info("mqtt connection up") },
		OnConnectError:    func(err error) { log.Errorf("error whilst attempting connection: %s\n", err) },
		Debug:             paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:      env.CfgMqttClientID(),
			OnClientError: func(err error) { log.Errorf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Errorf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Errorf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	cliCfg.SetUsernamePassword(env.CfgMqttUsername(), env.CfgMqttPassword())

	cliCfg.Debug = logger{prefix: "autoPaho"}
	cliCfg.PahoDebug = logger{prefix: "paho"}

	return cliCfg
}

func mqtt_start_connection(cliCfg autopaho.ClientConfig, ctx context.Context) {
	// Connect to the broker - this will return immediately after initiating the connection process
	cm, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		log.Panic(err)
	}

	// AwaitConnection will return immediately if connection is up; adding this call stops publication whilst
	// connection is unavailable.
	err = cm.AwaitConnection(ctx)
	if err != nil { // Should only happen when context is cancelled
		log.Debugf("publisher done (AwaitConnection: %s)\n", err)
		return
	}
	env.SetConnMqtt(cm)

}

func mqtt_push_discovery(ctx context.Context, discovery_topic string, ha_entity any) {
	msg, err := json.Marshal(ha_entity)
	if err != nil {
		log.Panic(err)
	}
	cm := env.GetConnMqtt()
	pr, err := cm.Publish(ctx, &paho.Publish{
		PacketID:   0,
		QoS:        env.CfgMqttQos(),
		Retain:     true,
		Topic:      discovery_topic,
		Properties: &paho.PublishProperties{},
		Payload:    msg,
	})
	if err != nil {
		log.Errorf("error publishing: %s\n", err)
	} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
		log.Errorf("reason code %d received\n", pr.ReasonCode)
	} else if env.CfgPrintMessage() {
		log.Tracef("sent message Discovery: %s\n", msg)
	}
}

func mqtt_push_status(ctx context.Context, status_topic string, status_value any) {
	msg, err := json.Marshal(status_value)
	if err != nil {
		log.Panic(err)
	}
	cm := env.GetConnMqtt()

	pr, err := cm.Publish(ctx, &paho.Publish{
		PacketID:   0,
		QoS:        env.CfgMqttQos(),
		Retain:     true,
		Topic:      status_topic,
		Properties: &paho.PublishProperties{},
		Payload:    msg,
	})
	if err != nil {
		log.Errorf("error publishing: %s\n", err)
	} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
		log.Errorf("reason code %d received\n", pr.ReasonCode)
	} else if env.CfgPrintMessage() {
		log.Tracef("sent status message: %s\n", msg)
	}
}

func mqtt_push_available(ctx context.Context, availability_topic string, availability_value string) {
	msg, err := json.Marshal(availability_value)
	if err != nil {
		log.Panic(err)
	}
	cm := env.GetConnMqtt()

	pr, err := cm.Publish(ctx, &paho.Publish{
		PacketID:   0,
		QoS:        env.CfgMqttQos(),
		Retain:     true,
		Topic:      availability_topic,
		Properties: &paho.PublishProperties{},
		Payload:    msg,
	})
	if err != nil {
		log.Errorf("error publishing: %s\n", err)
	} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
		log.Errorf("reason code %d received\n", pr.ReasonCode)
	} else if env.CfgPrintMessage() {
		log.Tracef("sent status message: %s\n", msg)
	}
}
