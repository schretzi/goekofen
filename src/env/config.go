package env

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	oekofenURL     *url.URL
	oekofenVersion string
	oekofenSerial  string

	mqttURL  *url.URL // MQTT server URL
	clientID string   // Client ID to use when connecting to server
	topic    string   // Topic to subscribe to
	qos      byte     // QOS to use when subscribing

	connectUsername string
	connectPassword []byte

	keepAlive            uint16        // seconds between keepalive packets
	connectRetryDelay    time.Duration // Period between connection attempts
	delayBetweenMessages time.Duration // Period between publishing messages

	printMessages bool // If true then published messages will be written to the console
	debug         bool // autopaho and paho debug output requested

	// Enable Backends
	homeassistant bool
	influxdb      bool
}

var cfg config

// ReadConfig - Retrieves the configuration from the environment
func ReadConfig() error {
	var err error

	oekofen_ip, err := stringFromEnv("OEKOFEN_IP")
	if err != nil {
		return err
	}
	oekofen_port, err := stringFromEnv("OEKOFEN_PORT")
	if err != nil {
		return err
	}
	oekofen_path, err := stringFromEnv("OEKOFEN_PATH")
	if err != nil {
		return err
	}
	oekofen_entity, err := stringFromEnv("OEKOFEN_ENTITY")
	if err != nil {
		oekofen_entity = "all"
	}
	var oekofen_url = "http" + "://" + oekofen_ip + ":" + oekofen_port + "/" + oekofen_path + "/" + oekofen_entity
	cfg.oekofenURL, err = url.Parse(oekofen_url)
	if err != nil {
		return fmt.Errorf("input must be a valid URL (%w)", err)
	}

	if cfg.oekofenSerial, err = stringFromEnv("OEKOFEN_SERIAL"); err != nil {
		cfg.oekofenSerial = "P00XXXXXX_XXXXXX"
	}

	if cfg.oekofenVersion, err = stringFromEnv("OEKOFEN_VERSION"); err != nil {
		cfg.oekofenVersion = "VX.XX_X"
	}

	mqtt_url, err := stringFromEnv("MQTT_URL")
	if err != nil {
		return err
	}

	cfg.mqttURL, err = url.Parse(mqtt_url)
	if err != nil {
		return fmt.Errorf("input must be a valid URL (%w)", err)
	}

	if cfg.connectUsername, err = stringFromEnv("MQTT_USERNAME"); err != nil {
		return err
	}
	cPW, err := stringFromEnv("MQTT_PASSWORD")
	if err != nil {
		return err
	}
	cfg.connectPassword = []byte(cPW)

	if cfg.clientID, err = stringFromEnv("MQTT_CLIENTID"); err != nil {
		return err
	}
	if cfg.topic, err = stringFromEnv("MQTT_STATUS_TOPIC"); err != nil {
		return err
	}

	iQos, err := intFromEnv("MQTT_QOS")
	if err != nil {
		return err
	}
	cfg.qos = byte(iQos)

	iKa, err := intFromEnv("MQTT_KEEPALIVE")
	if err != nil {
		return err
	}
	cfg.keepAlive = uint16(iKa)

	if cfg.connectRetryDelay, err = milliSecondsFromEnv("MQTT_CONNECT_RETRY_DELAY"); err != nil {
		return err
	}

	if cfg.delayBetweenMessages, err = milliSecondsFromEnv("MQTT_DELAY_BETWEEN_MESSAGES"); err != nil {
		return err
	}

	if cfg.printMessages, err = booleanFromEnv("PRINTMESSAGES"); err != nil {
		return err
	}
	if cfg.debug, err = booleanFromEnv("DEBUG"); err != nil {
		return err
	}

	if cfg.homeassistant, err = booleanFromEnv("HOMEASSISTANT"); err != nil {
		return err
	}
	if cfg.influxdb, err = booleanFromEnv("INFLUXDB"); err != nil {
		return err
	}

	return nil
}

// Getters - Get the
func CfgOekofenURL() *url.URL {
	return cfg.oekofenURL
}

func CfgOekofenSerial() string {
	return cfg.oekofenSerial
}

func CfgOekofenVersion() string {
	return cfg.oekofenVersion
}

func CfgMqttURL() *url.URL {
	return cfg.mqttURL
}

func CfgMqttUsername() string {
	return cfg.connectUsername
}

func CfgMqttPassword() []byte {
	return cfg.connectPassword
}

func CfgMqttClientID() string {
	return cfg.clientID
}

func CfgMqttTopic() string {
	return cfg.topic
}

func CfgMqttQos() byte {
	return cfg.qos
}

func CfgMqttKeepAlive() uint16 {
	return cfg.keepAlive
}

func CfgMqttRetryDelay() time.Duration {
	return cfg.connectRetryDelay
}

func CfgMqttDelayBetweenMessages() time.Duration {
	return cfg.delayBetweenMessages
}

func CfgPrintMessage() bool {
	return cfg.printMessages
}

func CfgDebug() bool {
	//return cfg.debug
	return false
}
func CfgHA() bool {
	return cfg.homeassistant
}

func CfgInfluxDB() bool {
	return cfg.influxdb
}

// Helper Functions for parsing Environment variables

// stringFromEnv - Retrieves a string from the environment and ensures it is not blank (ort non-existent)
func stringFromEnv(key string) (string, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return "", fmt.Errorf("environmental variable %s must not be blank", key)
	}
	return s, nil
}

// intFromEnv - Retrieves an integer from the environment (must be present and valid)
func intFromEnv(key string) (int, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return 0, fmt.Errorf("environmental variable %s must not be blank", key)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environmental variable %s must be an integer", key)
	}
	return i, nil
}

// milliSecondsFromEnv - Retrieves milliseconds (as time.Duration) from the environment (must be present and valid)
func milliSecondsFromEnv(key string) (time.Duration, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return 0, fmt.Errorf("environmental variable %s must not be blank", key)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environmental variable %s must be an integer", key)
	}
	return time.Duration(i) * time.Millisecond, nil
}

// booleanFromEnv - Retrieves boolean from the environment (must be present and valid)
func booleanFromEnv(key string) (bool, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return false, fmt.Errorf("environmental variable %s must not be blank", key)
	}
	switch strings.ToUpper(s) {
	case "TRUE", "T", "1":
		return true, nil
	case "FALSE", "F", "0":
		return false, nil
	default:
		return false, fmt.Errorf("environmental variable %s be a valid boolean option (is %s)", key, s)
	}
}
