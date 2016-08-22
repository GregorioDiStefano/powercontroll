package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (m *macAddress) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))

	if err != nil {
		return err
	}

	hwAddress, err := net.ParseMAC(unquoted)

	if err != nil {
		return err
	}

	*m = macAddress(hwAddress.String())
	return nil
}

var connectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Println("Connection lost")
}

var connectionHandler = func(client MQTT.Client) {
	fmt.Println("Connection established")
}

var announceHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	var announceData map[string]macAddress

	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSGID: %d, MSG: %s\n", msg.MessageID(), msg.Payload())

	if err := json.Unmarshal(msg.Payload(), &announceData); err != nil {
		panic(err)
	}

	if _, ok := announceData["mac_address"]; !ok {
		//do something smart
	}

	added, err := db.addDevice(string(announceData["mac_address"]))
	if added {
		fmt.Println("New incoming device: ", announceData["mac_address"])
	} else if err != nil {
		fmt.Printf("Error adding new device <%s>: %s", announceData["mac_address"], err.Error())
	}

	time.Sleep(1 * time.Second)
}

var pongHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	var pongData map[string]macAddress

	if err := json.Unmarshal(msg.Payload(), &pongData); err != nil {
		panic(err)
	}

	if _, ok := pongData["mac_address"]; !ok {
		fmt.Println("Invalid data, 'mac_address' not found in pong message.")
	}

	deviceMacAddress := string(pongData["mac_address"])

	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSGID: %d, MSG: %s\n", msg.MessageID(), msg.Payload())
	updateDeviceLastSeen(deviceMacAddress)
}

var statusAlertHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSGID: %d, MSG: %s\n", msg.MessageID(), msg.Payload())

	var alertData map[string]interface{}

	if err := json.Unmarshal(msg.Payload(), &alertData); err != nil {
		panic(err)
	}

	if _, ok := alertData["mac_address"]; !ok {
		fmt.Println("Invalid data, 'mac_address' not found in pong message.")
	}

	deviceMacAddress := string(alertData["mac_address"].(string))

	db.setDeviceStatus(deviceMacAddress, alertData["status"].(string))
}
